// Package main updates the database with new vault content.
package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"math"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	mathjax "github.com/litao91/goldmark-mathjax"
	ollama "github.com/prathyushnallamothu/ollamago"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/quail-ink/goldmark-enclave/core"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/genstruct"
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
	"golang.org/x/sync/errgroup"
	"gonum.org/v1/gonum/mat"
)

const (
	embeddingModel = "nomic-embed-text"
	hashFile       = "config.json"
	assetsLoc      = "internal/data/docs/assets/"
	postsLoc       = "internal/data/docs/posts/"
	tagsLoc        = "internal/data/docs/tags/"
	projectsLoc    = "internal/data/docs/projects/"
)

var (
	uploadJobs = flag.Int("jobs", 10, "number of parallel uploads")
	cwd        = flag.String("cwd", "", "current working directory")
	logger     = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: false,
		Level:           log.DebugLevel,
	})
)

// Cache is the configuration for the updater.
type Cache struct {
	AssetsHash string `json:"hash"`
	AssetsMap  map[string]string
}

// Close writes the config to disk and closes the file.
func (c *Cache) Close() error {
	// write the config to disk
	body, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(hashFile, body, 0644)
}

// CompareMap returns true if the two maps are equal.
func (c *Cache) CompareMap(m map[string]string) bool {
	for k, v := range c.AssetsMap {
		if m[k] != v {
			return false
		}
	}
	return true
}

func main() {
	flag.Parse()

	if *cwd != "" {
		err := os.Chdir(*cwd)
		if err != nil {
			panic(err)
		}
	}
	if err := Run(
		context.Background(),
		s3.NewFromConfig(aws.Config{
			Region:       "auto",
			BaseEndpoint: aws.String("https://fly.storage.tigris.dev"),
			Credentials: &CredHandler{
				Name: "conneroh",
				ID:   os.Getenv("AWS_ACCESS_KEY_ID"),
				Key:  os.Getenv("AWS_SECRET_ACCESS_KEY"),
			},
		}),
		ollama.NewClient(ollama.WithTimeout(time.Minute*5)),
	); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadConfig() (Cache, error) {
	var cache Cache
	data, err := os.ReadFile(hashFile)
	if err != nil {
		return cache, err
	}
	if err := json.Unmarshal(data, &cache); err != nil {
		return cache, err
	}
	return cache, nil
}

// Run parses all markdown files in the database.
func Run(
	ctx context.Context,
	client *s3.Client,
	llama *ollama.Client,
) error {
	cache, err := loadConfig()
	if err != nil {
		return err
	}
	defer cache.Close()
	assetHash, err := hashDir(assetsLoc)
	if err != nil {
		return err
	}
	assets, hashMap, err := assetsParse(cache, assetsLoc)
	if err != nil {
		return err
	}
	if cache.AssetsHash != assetHash || !cache.CompareMap(hashMap) {
		cache.AssetsHash = assetHash
		cache.AssetsMap = hashMap
		err = uploadAssets(ctx, client, assets, *uploadJobs)
		if err != nil {
			return err
		}
	}
	logger.Info("assets are up to date")
	logger.Info("parsing tags")
	parsedTags, err := pathParse[gen.Tag](ctx, &assets, tagsLoc, llama)
	if err != nil {
		return err
	}
	logger.Info("parsed tags")
	logger.Info("parsing posts")
	parsedPosts, err := pathParse[gen.Post](ctx, &assets, postsLoc, llama)
	if err != nil {
		return err
	}
	logger.Info("parsed posts")
	logger.Info("parsing projects")
	parsedProjects, err := pathParse[gen.Project](ctx, &assets, projectsLoc, llama)
	if err != nil {
		return fmt.Errorf("failed to parse projects: %v", err)
	}
	logger.Info("parsed projects")

	postGen, err := genstruct.NewGenerator(genstruct.Config{
		PackageName: "gen",
		OutputFile:  "internal/data/gen/generated_data.go",
	}, parsedPosts, parsedTags, parsedProjects)
	if err != nil {
		return err
	}

	return postGen.Generate()
}

func uploadAssets(
	ctx context.Context,
	client *s3.Client,
	assets []Asset,
	uploadJobs int,
) error {
	eg := errgroup.Group{}
	eg.SetLimit(uploadJobs)
	for _, asset := range assets {
		eg.Go(func() error {
			_, err := client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String("conneroh.com"),
				Key:    aws.String(asset.Path),
				Body:   bytes.NewReader(asset.Data),
			})
			return err
		})
	}
	return eg.Wait()
}

// pathParse parses the markdown files in the given path.
func pathParse[T gen.Post | gen.Project | gen.Tag](
	ctx context.Context,
	assets *[]Asset,
	fsPath string,
	llama *ollama.Client,
) ([]T, error) {
	var parseds []T
	err := filepath.WalkDir(
		fsPath,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf(
					"failed to walk fsPath (%s): %w",
					fsPath,
					err,
				)
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			parsed, err := Parse[T](
				ctx,
				parser.NewContext(),
				path,
				NewParser(assets),
				llama,
			)
			if err != nil {
				return err
			}
			if parsed == nil {
				return nil
			}
			parseds = append(parseds, *parsed)
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to walk fsPath (%s): %w", fsPath, err)
	}
	return parseds, nil
}

func assetsParse(
	cfg Cache,
	assetsLoc string,
) ([]Asset, map[string]string, error) {
	var assets []Asset
	var contentsMap = make(map[string]string)
	err := filepath.WalkDir(
		assetsLoc,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf(
					"failed to walk fsPath (%s): %w",
					"assets",
					err,
				)
			}
			if d.IsDir() {
				return nil
			}
			asset, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			hash, err := hashFileContent(string(asset))
			if err != nil {
				return err
			}
			realPath := strings.TrimPrefix(path, "assets/")
			if cfg.AssetsMap[realPath] == hash {
				return nil
			}
			assets = append(assets, Asset{
				Path: realPath,
				Data: asset,
			})
			contentsMap[realPath] = hash
			return nil
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return assets, contentsMap, nil
}

func hashDir(dir string) (string, error) {
	hash := md5.New()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(hash, file); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func hashFileContent(content string) (string, error) {
	hash := md5.New()
	_, err := io.WriteString(hash, content)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Parse parses the markdown file at the given path.
func Parse[T gen.Post | gen.Project | gen.Tag](
	ctx context.Context,
	pCtx parser.Context,
	path string,
	md goldmark.Markdown,
	llama *ollama.Client,
) (*T, error) {
	var (
		fm       gen.Embedded
		metadata = frontmatter.Get(pCtx)
		buf      = bytes.NewBufferString("")
		err      error
	)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = md.Convert(content, buf, parser.WithContext(pCtx))
	if err != nil {
		return nil, err
	}
	metadata = frontmatter.Get(pCtx)
	if metadata == nil {
		return nil, fmt.Errorf("frontmatter is nil for %s", path)
	}
	err = metadata.Decode(&fm)
	if err != nil {
		return nil, fmt.Errorf("failed to decode frontmatter: %w", err)
	}

	path = strings.Replace(path, postsLoc, "", 1)
	path = strings.Replace(path, projectsLoc, "", 1)
	path = strings.Replace(path, tagsLoc, "", 1)
	path = strings.TrimSuffix(path, filepath.Ext(path))
	logger.Debugf("path: %s", path)
	fm.Slug = path
	fm.Content = buf.String()
	if fm.Description == "" {
		return nil, fmt.Errorf("description is empty for %s", path)
	}
	fm.RawContent = string(content)

	if fm.Icon == "" {
		fm.Icon = "tag"
	}

	if fm.Content == "" {
		return nil, fmt.Errorf("content is empty for %s", path)
	}
	fm.Vec, fm.X, fm.Y, fm.Z, err = TextEmbeddingCreate(
		ctx,
		llama,
		fm.RawContent,
	)
	if err != nil {
		return nil, err
	}

	return gen.New[T](fm), nil
}

// TextEmbeddingCreate creates an embedding for the given text.
func TextEmbeddingCreate(
	ctx context.Context,
	client *ollama.Client,
	input string,
) ([gen.EmbedLength]float64, float64, float64, float64, error) {
	resp, err := client.Embeddings(ctx, ollama.EmbeddingsRequest{
		Model:  embeddingModel,
		Prompt: input,
	})
	if err != nil {
		return [gen.EmbedLength]float64{}, 0, 0, 0, err
	}
	proj := projectionMatrixCreate(gen.EmbedLength, 3)
	x, y, z := projectTo3D(resp.Embedding, proj)
	embs := [gen.EmbedLength]float64{}
	for i := range embs {
		embs[i] = resp.Embedding[i]
	}
	return embs, x, y, z, nil
}

// Generate a deterministic projection matrix for dimensionality reduction
func projectionMatrixCreate(inputDim, outputDim int) *mat.Dense {
	// Create a matrix with dimensions [outputDim x inputDim]
	// This will multiply directly with the input vector without transposition
	matrix := mat.NewDense(outputDim, inputDim, nil)

	// For each output dimension
	for i := range outputDim {
		// Calculate a center point for this dimension's projection focus
		center := (i * inputDim) / outputDim

		for j := range inputDim {
			// Calculate a weight based on distance from the center point
			// Use a Gaussian-like function for smooth projection
			dist := float64(j - center)
			weight := math.Exp(-0.5 * (dist * dist) / float64(inputDim/10))

			// Set the weight in the projection matrix
			matrix.Set(i, j, weight)
		}

		// Normalize the row to unit length
		var rowSum float64
		for j := range inputDim {
			rowSum += matrix.At(i, j) * matrix.At(i, j)
		}

		rowSum = math.Sqrt(rowSum)
		if rowSum > 0 {
			for j := range inputDim {
				matrix.Set(i, j, matrix.At(i, j)/rowSum)
			}
		}
	}

	return matrix
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

// NewParser creates a new markdown parser.
func NewParser(
	assets *[]Asset,
) goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, extension.Footnote,
			extension.Strikethrough, extension.Table,
			extension.TaskList, extension.DefinitionList,
			mathjax.MathJax,
			&wikilink.Extender{
				Resolver: NewCustomResolver(assets),
			},
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
}

// Asset is a struct for embedding assets.
type Asset struct {
	Path string
	Data []byte
}

// Upload uploads the asset to the s3 bucket.
func (a *Asset) Upload(ctx context.Context, client *s3.Client) error {
	mimeType := mime.TypeByExtension(filepath.Ext(a.Path))
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("BUCKET_NAME")),
		Key:         aws.String(a.Path),
		Body:        bytes.NewReader(a.Data),
		ContentType: aws.String(mimeType),
	})

	return err
}

// URL returns the url of the asset.
func (a *Asset) URL() string {
	return fmt.Sprintf(
		"https://conneroh.fly.storage.tigris.dev/%s",
		a.Path,
	)
}

// Filename returns the filename of the asset.
func (a *Asset) Filename() string {
	str := a.Path
	str = str[strings.LastIndex(str, "/")+1:]
	return str
}

// Resolvable returns true if the asset is resolvable.
func (a *Asset) Resolvable(path string) bool {
	return strings.HasPrefix(path, a.Path)
}

// CredHandler is the bucket for the api security.
// It implements the aws.Authenitcator
type CredHandler struct {
	Name string
	ID   string
	Key  string
}

// Retrieve returns the credentials for the bucket.
func (b *CredHandler) Retrieve(
	_ context.Context,
) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     b.ID,
		SecretAccessKey: b.Key,
	}, nil
}

// CustomResolver is a wikilink.Resolver that uses the default wikilink
// resolver but also resolves the "custom" wikilink type.
//
// A Resolver resolves pages referenced by wikilinks to their destinations.
//
// It implements the wikilink.Resolver interface.
//
//	type Resolver interface {
//		// ResolveWikilink returns the address of the page that the provided
//		// wikilink points to. The destination will be URL-escaped before
//		// being placed into a link.
//		//
//		// If ResolveWikilink returns a non-nil error, rendering will be
//		// halted.
//		//
//		// If ResolveWikilink returns a nil destination and error, the
//		// Renderer will omit the link and render its contents as a regular
//		// string.
//		ResolveWikilink(*Node) (destination []byte, err error)
//	}
//
// ```
type CustomResolver struct {
	Assets *[]Asset
}

// NewCustomResolver creates a new wikilink resolver.
func NewCustomResolver(assets *[]Asset) *CustomResolver {
	return &CustomResolver{Assets: assets}
}

// ResolveWikilink returns the address of the page that the provided
// wikilink points to. The destination will be URL-escaped before
// being placed into a link.
//
// If ResolveWikilink returns a non-nil error, rendering will be
// halted.
//
// If ResolveWikilink returns a nil destination and error, the
// Renderer will omit the link and render its contents as a regular
// string.
func (c *CustomResolver) ResolveWikilink(n *wikilink.Node) (destination []byte, err error) {
	targetStr := string(n.Target)
	for _, asset := range *c.Assets {
		if targetStr == asset.Path || strings.HasSuffix(targetStr, asset.Filename()) {
			return []byte(asset.URL()), nil
		}
	}

	return nil, nil
}
