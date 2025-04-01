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
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/log"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/ollama/ollama/api"
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
	workers = flag.Int("jobs", 10, "number of parallel uploads")
	cwd     = flag.String("cwd", "", "current working directory")
	logger  = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: false,
		Level:           log.DebugLevel,
	})

	client = s3.NewFromConfig(aws.Config{
		Region:       "auto",
		BaseEndpoint: aws.String("https://fly.storage.tigris.dev"),
		Credentials: &CredHandler{
			Name: "conneroh",
			ID:   os.Getenv("AWS_ACCESS_KEY_ID"),
			Key:  os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	})
	llama = &api.Client{}
)

// Cache is the storage of previous hashes.
type Cache struct {
	Hashes map[string]string
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

func main() {
	flag.Parse()
	llamaURL, err := url.Parse(os.Getenv("OLLAMA_URL"))
	if err != nil {
		panic(err)
	}
	llama = api.NewClient(llamaURL, http.DefaultClient)

	if *cwd != "" {
		err := os.Chdir(*cwd)
		if err != nil {
			panic(err)
		}
	}
	pos, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting at (PWD=%s)", pos)
	if err := Run(
		context.Background(),
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
	if cache.Hashes == nil {
		cache.Hashes = make(map[string]string)
	}
	return cache, nil
}

// Run parses all markdown files in the database.
func Run(
	ctx context.Context,
) error {
	var (
		parsedTags     []gen.Tag
		parsedProjects []gen.Project
		parsedPosts    []gen.Post
	)
	eg := errgroup.Group{}
	eg.SetLimit(*workers)
	cache, err := loadConfig()
	if err != nil {
		return err
	}
	defer cache.Close()

	logger.Info("parsing assets")
	assets, _, err := parse(cache, assetsLoc)
	if err != nil {
		return err
	}
	for _, asset := range assets {
		eg.Go(func() error {
			return uploadAsset(ctx, client, asset)
		})
	}
	if eg.Wait() != nil {
		return err
	}
	logger.Info("assets are up to date")

	mdParser := NewMDParser(assets)

	logger.Info("actualizing posts", "dest", postsLoc)
	parsedPosts, err = actualize[gen.Post](ctx, cache, mdParser, postsLoc)
	if err != nil {
		return err
	}
	logger.Info("actualized posts", "len", len(parsedPosts))

	logger.Info("actualizing projects", "dest", projectsLoc)
	parsedProjects, err = actualize[gen.Project](ctx, cache, mdParser, projectsLoc)
	if err != nil {
		return err
	}
	logger.Info("actualized projects", "len", len(parsedProjects))

	logger.Info("actualizing tags", "dest", tagsLoc)
	parsedTags, err = actualize[gen.Tag](ctx, cache, mdParser, tagsLoc)
	if err != nil {
		return err
	}
	logger.Info("actualized tags", "len", len(parsedTags))

	postGen, err := genstruct.NewGenerator(genstruct.Config{
		PackageName: "gen",
		OutputFile:  "internal/data/gen/generated_data.go",
	}, parsedPosts, parsedTags, parsedProjects)
	if err != nil {
		return err
	}

	return postGen.Generate()
}

func actualize[T gen.Post | gen.Tag | gen.Project](
	ctx context.Context,
	cache Cache,
	mdParser goldmark.Markdown,
	loc string,
) ([]T, error) {
	// Create a channel to safely collect results from goroutines
	type result struct {
		item T
		err  error
	}
	resultCh := make(chan result)

	eg := errgroup.Group{}
	eg.SetLimit(*workers)

	// Get files to process
	assets, ignored, err := parse(cache, loc)
	if err != nil {
		return nil, err
	}

	// Start a goroutine to collect results
	var parsed []T
	var parseWg sync.WaitGroup
	parseWg.Add(1)
	go func() {
		defer parseWg.Done()
		for res := range resultCh {
			if res.err == nil {
				parsed = append(parsed, res.item)
			}
		}
	}()

	// Process each asset concurrently
	for _, asset := range assets {
		assetCopy := asset // Create a copy for the goroutine
		eg.Go(func() error {
			realized, err := realizeMD[T](ctx, mdParser, assetCopy)
			if err == nil && realized != nil {
				resultCh <- result{item: *realized, err: nil}
			}
			return err
		})
	}

	// Wait for all processing to complete
	err = eg.Wait()
	close(resultCh) // Close channel after all goroutines are done

	// Wait for collection goroutine to finish
	parseWg.Wait()

	// Include previously parsed items that were ignored this time
	remembered := rememberMD[T](ignored)
	parsed = append(parsed, remembered...)

	logger.Debugf("Parsed %d assets. Ignored %d", len(parsed), len(ignored))
	return parsed, err
}

func uploadAsset(
	ctx context.Context,
	client *s3.Client,
	asset Asset,
) error {
	logger.Debugf("Uploading asset %s", asset.Path)
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("conneroh"),
		Key:    aws.String(asset.Path),
		Body:   bytes.NewReader(asset.Data),
	})
	return err
}

func hashFileContent(content string) (string, error) {
	hash := md5.New()
	_, err := io.WriteString(hash, content)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func rememberMD[T gen.Post | gen.Project | gen.Tag](ignored []string) []T {
	// Use type assertions with a type variable to determine what we're working with
	var typeExample T
	switch any(typeExample).(type) {
	case gen.Post:
		var posts []gen.Post
		for _, post := range gen.AllPosts {
			for _, ignored := range ignored {
				if post.Slug == ignored {
					posts = append(posts, *post)
				}
			}
		}
		return any(posts).([]T) // Type assertion to convert the concrete type to []T
	case gen.Project:
		var projects []gen.Project
		for _, project := range gen.AllProjects {
			for _, ignored := range ignored {
				if project.Slug == ignored {
					projects = append(projects, *project)
				}
			}
		}
		return any(projects).([]T) // Type assertion to convert the concrete type to []T
	case gen.Tag:
		var tags []gen.Tag
		for _, tag := range gen.AllTags {
			for _, ignored := range ignored {
				if tag.Slug == ignored {
					tags = append(tags, *tag)
				}
			}
		}
		return any(tags).([]T) // Type assertion to convert the concrete type to []T
	default:
		panic("unknown type")
	}
}

// realizeMD parses the markdown file at the given path.
func realizeMD[T gen.Post | gen.Project | gen.Tag](
	ctx context.Context,
	md goldmark.Markdown,
	parsed Asset,
) (*T, error) {
	var (
		fm       gen.Embedded
		pCtx     = parser.NewContext()
		metadata = frontmatter.Get(pCtx)
		buf      = bytes.NewBufferString("")
		err      error
	)
	err = md.Convert(parsed.Data, buf, parser.WithContext(pCtx))
	if err != nil {
		return nil, err
	}
	metadata = frontmatter.Get(pCtx)
	if metadata == nil {
		return nil, fmt.Errorf("frontmatter is nil for %s", parsed.Path)
	}
	err = metadata.Decode(&fm)
	if err != nil {
		return nil, fmt.Errorf("failed to decode frontmatter: %w", err)
	}

	fm.Slug = parsed.Slug
	fm.Content = buf.String()
	if fm.Description == "" {
		return nil, fmt.Errorf("description is empty for %s", parsed.Path)
	}
	fm.RawContent = string(parsed.Data)

	if fm.Icon == "" {
		fm.Icon = "tag"
	}

	if fm.Content == "" {
		return nil, fmt.Errorf("content is empty for %s", parsed.Path)
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
	client *api.Client,
	input string,
) ([gen.EmbedLength]float64, float64, float64, float64, error) {
	resp, err := client.Embeddings(ctx, &api.EmbeddingRequest{
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

// NewMDParser creates a new markdown parser.
func NewMDParser(
	assets []Asset,
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
	Slug string
	Path string
	Data []byte
}

// Upload uploads the asset to the s3 bucket.
func (a *Asset) Upload(ctx context.Context, client *s3.Client) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("BUCKET_NAME")),
		Key:         aws.String(a.Path),
		Body:        bytes.NewReader(a.Data),
		ContentType: aws.String(mime.TypeByExtension(filepath.Ext(a.Path))),
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

// CustomResolver is a wikilink.Resolver that resolves pages referenced by
// wikilinks to their destinations.
type CustomResolver struct {
	Assets []Asset
}

// NewCustomResolver creates a new wikilink resolver.
func NewCustomResolver(assets []Asset) *CustomResolver {
	return &CustomResolver{Assets: assets}
}

// ResolveWikilink returns the address of the page that the provided
// wikilink points to. The destination will be URL-escaped before
// being placed into a link.
func (c *CustomResolver) ResolveWikilink(n *wikilink.Node) (destination []byte, err error) {
	targetStr := string(n.Target)
	for _, asset := range c.Assets {
		if targetStr == asset.Path || strings.HasSuffix(targetStr, asset.Path[strings.LastIndex(asset.Path, "/")+1:]) {
			return []byte(asset.URL()), nil
		}
	}

	return nil, nil
}

// parse parses the markdown files in the given path.
//
// ignored is a list of slugs that were ignored as they did not change.
func parse(
	cache Cache,
	loc string,
) (parseds []Asset, ignored []string, err error) {
	err = filepath.WalkDir(
		loc,
		func(fPath string, d fs.DirEntry, err error) error {
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
			asset, err := os.ReadFile(fPath)
			if err != nil {
				return err
			}
			hash, err := hashFileContent(string(asset))
			if err != nil {
				return err
			}
			path := pathify(fPath)
			slug := slugify(fPath)
			if cache.Hashes[fPath] == hash {
				cache.Hashes[fPath] = hash
				ignored = append(ignored, slug)
				return nil
			}
			cache.Hashes[fPath] = hash
			parseds = append(parseds, Asset{
				Path: path,
				Data: asset,
				Slug: slug,
			})
			return nil
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to walk fsPath (%s): %w", loc, err)
	}
	return
}

func slugify(s string) string {
	path := strings.Replace(s, postsLoc, "", 1)
	path = strings.Replace(path, projectsLoc, "", 1)
	path = strings.Replace(path, tagsLoc, "", 1)
	path = strings.TrimSuffix(path, filepath.Ext(path))
	return path
}
func pathify(s string) string {
	var path string
	var ok bool
	path, ok = strings.CutPrefix(s, postsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, projectsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, tagsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, assetsLoc)
	if ok {
		return path
	}
	panic("unknown path")
}
