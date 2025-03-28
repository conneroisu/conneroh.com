// Package main updates the database with new vault content.
package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/conneroisu/conneroh.com/internal/data/docs"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/genstruct"
	mathjax "github.com/litao91/goldmark-mathjax"
	ollama "github.com/prathyushnallamothu/ollamago"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/quail-ink/goldmark-enclave/core"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/wikilink"
	"gonum.org/v1/gonum/mat"
)

var client = ollama.NewClient(ollama.WithTimeout(time.Minute * 5))

func main() {
	if err := Run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Run parses all markdown files in the database.
func Run(ctx context.Context) error {
	parsedTags, err := pathParse[gen.Tag](ctx, "tags", docs.Tags)
	if err != nil {
		return err
	}

	parsedPosts, err := pathParse[gen.Post](ctx, "posts", docs.Posts)
	if err != nil {
		return err
	}

	parsedProjects, err := pathParse[gen.Project](ctx, "projects", docs.Projects)
	if err != nil {
		return fmt.Errorf("failed to parse projects: %v", err)
	}
	// Create a generator for people
	postGen, err := genstruct.NewGenerator(genstruct.Config{
		PackageName: "gen",
		OutputFile:  fmt.Sprintf("internal/data/gen/generated_%s.go", "data"),
	}, parsedPosts, parsedTags, parsedProjects)
	if err != nil {
		return err
	}

	return postGen.Generate()
}

// Project a single embedding to 3D using a projection matrix
func projectTo3D(embedding []float64, projectionMatrix *mat.Dense) (x, y, z float64) {
	// Create a vector from the embedding
	embedVec := mat.NewVecDense(len(embedding), embedding)

	// Project to 3D - no need to transpose
	result := mat.NewVecDense(3, nil)
	result.MulVec(projectionMatrix, embedVec)

	// Extract x, y, z coordinates
	x = result.AtVec(0)
	y = result.AtVec(1)
	z = result.AtVec(2)

	return
}

// Generate a deterministic projection matrix for dimensionality reduction
func generateProjectionMatrix(inputDim, outputDim int) *mat.Dense {
	// Create a matrix with dimensions [outputDim x inputDim]
	// This will multiply directly with the input vector without transposition
	matrix := mat.NewDense(outputDim, inputDim, nil)

	// For each output dimension
	for i := 0; i < outputDim; i++ {
		// Calculate a center point for this dimension's projection focus
		center := (i * inputDim) / outputDim

		for j := 0; j < inputDim; j++ {
			// Calculate a weight based on distance from the center point
			// Use a Gaussian-like function for smooth projection
			dist := float64(j - center)
			weight := math.Exp(-0.5 * (dist * dist) / float64(inputDim/10))

			// Set the weight in the projection matrix
			matrix.Set(i, j, weight)
		}

		// Normalize the row to unit length
		var rowSum float64
		for j := 0; j < inputDim; j++ {
			rowSum += matrix.At(i, j) * matrix.At(i, j)
		}

		rowSum = math.Sqrt(rowSum)
		if rowSum > 0 {
			for j := 0; j < inputDim; j++ {
				matrix.Set(i, j, matrix.At(i, j)/rowSum)
			}
		}
	}

	return matrix
}

var (
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, extension.Footnote,
			extension.Strikethrough, extension.Table,
			extension.TaskList, extension.DefinitionList,
			mathjax.MathJax, &wikilink.Extender{},
			extension.NewTypographer(
				extension.WithTypographicSubstitutions(
					extension.TypographicSubstitutions{
						extension.Apostrophe: []byte("'"),
					}),
			),
			enclave.New(&core.Config{DefaultImageAltPrefix: "caption: "}),
			extension.NewFootnote(
				extension.WithFootnoteIDPrefix("fn"),
			),
			&anchor.Extender{
				Position: anchor.Before,
				Texter:   anchor.Text("#"),
				Attributer: anchor.Attributes{
					"class": "anchor permalink p-4",
				},
			},
			&mermaid.Extender{
				RenderMode: mermaid.RenderModeClient,
			},
			&frontmatter.Extender{
				Formats: []frontmatter.Format{frontmatter.YAML},
			},
			&hashtag.Extender{
				Variant: hashtag.ObsidianVariant,
			},
			highlighting.NewHighlighting(highlighting.WithStyle("monokai")),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			extension.WithFootnoteBacklinkClass("footnote-backref"),
			extension.WithFootnoteLinkClass("footnote-ref"),
		),
	)
)

// pathParse parses the markdown files in the given path.
func pathParse[
	T gen.Post | gen.Project | gen.Tag,
](ctx context.Context, fsPath string, embedFs embed.FS) ([]T, error) {
	var parseds []T
	err := fs.WalkDir(
		embedFs,
		fsPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("failed to walk projects: %w", err)
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			parsed, err := parse[T](ctx, path, embedFs)
			if err != nil {
				return fmt.Errorf("failed to parse project %s: %w", path, err)
			}
			if parsed == nil {
				return nil
			}
			parseds = append(parseds, *parsed)
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update projects: %v", err)
	}
	return parseds, nil
}

// parse parses the markdown file at the given path.
func parse[
	T gen.Post | gen.Project | gen.Tag,
](ctx context.Context, fsPath string, embedFs embed.FS) (*T, error) {
	var (
		pCtx     = parser.NewContext()
		fm       gen.Embedded
		metadata = frontmatter.Get(pCtx)
		body     []byte
		buf      = bytes.NewBufferString("")
		err      error
	)

	body, err = embedFs.ReadFile(fsPath)
	if err != nil {
		return nil, err
	}
	err = md.Convert(body, buf, parser.WithContext(pCtx))
	if err != nil {
		return nil, err
	}
	metadata = frontmatter.Get(pCtx)
	if metadata == nil {
		return nil, fmt.Errorf("frontmatter is nil for %s", fsPath)
	}
	err = metadata.Decode(&fm)
	if err != nil {
		return nil, fmt.Errorf("failed to decode frontmatter: %w", err)
	}

	switch embedFs {
	case docs.Posts:
		fsPath = strings.Replace(fsPath, "posts/", "", 1)
	case docs.Tags:
		fsPath = strings.Replace(fsPath, "tags/", "", 1)
	case docs.Projects:
		fsPath = strings.Replace(fsPath, "projects/", "", 1)
	default:
		return nil, fmt.Errorf("unknown embedFs %v", embedFs)
	}
	fsPath = strings.TrimSuffix(fsPath, filepath.Ext(fsPath))
	fm.Slug = fsPath
	fm.Content = buf.String()
	if fm.Description == "" {
		return nil, fmt.Errorf("description is empty for %s", fsPath)
	}
	fm.RawContent = string(body)

	if fm.Icon == "" {
		fm.Icon = "tag"
	}

	if fm.Content == "" {
		return nil, fmt.Errorf("content is empty for %s", fsPath)
	}
	fm.Vec, fm.X, fm.Y, fm.Z, err = embedIt(ctx, fm.RawContent)
	if err != nil {
		return nil, err
	}

	return gen.New[T](fm), nil
}

func embedIt(
	ctx context.Context,
	input string,
) ([gen.EmbedLength]float64, float64, float64, float64, error) {
	resp, err := client.Embeddings(ctx, ollama.EmbeddingsRequest{
		Model:  "nomic-embed-text",
		Prompt: input,
	})
	if err != nil {
		return [gen.EmbedLength]float64{}, 0, 0, 0, err
	}
	proj := generateProjectionMatrix(gen.EmbedLength, 3)
	x, y, z := projectTo3D(resp.Embedding, proj)
	embs := [gen.EmbedLength]float64{}
	for i := range embs {
		embs[i] = resp.Embedding[i]
	}
	return embs, x, y, z, nil
}
