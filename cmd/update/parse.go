package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/conneroisu/conneroh.com/internal/data/docs"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
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
)

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
](fsPath string, embedFs embed.FS) ([]T, error) {
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
			parsed, err := parse[T](path, embedFs)
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
](fsPath string, embedFs embed.FS) (*T, error) {
	if filepath.Ext(fsPath) != ".md" {
		return nil, nil
	}

	var (
		ctx      = parser.NewContext()
		fm       gen.Embedded
		metadata = frontmatter.Get(ctx)
		body     []byte
		buf      = bytes.NewBufferString("")
		err      error
	)

	body, err = embedFs.ReadFile(fsPath)
	if err != nil {
		return nil, err
	}
	err = md.Convert(body, buf, parser.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	metadata = frontmatter.Get(ctx)
	if metadata == nil {
		return nil, fmt.Errorf("frontmatter is nil for %s", fsPath)
	}
	err = metadata.Decode(&fm)
	if err != nil {
		return nil, err
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
	fm.Vec, fm.X, fm.Y, fm.Z, err = embedIt(context.TODO(), fm.RawContent)
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
	proj := generateProjectionMatrix(embeddingSize, 3)
	x, y, z := projectTo3D(resp.Embedding, proj)
	embs := [gen.EmbedLength]float64{}
	for i := range embs {
		embs[i] = resp.Embedding[i]
	}
	return embs, x, y, z, nil
}
