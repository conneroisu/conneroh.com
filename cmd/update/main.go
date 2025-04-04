// Package main updates the database with new vault content.
package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"math"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/ollama/ollama/api"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/quail-ink/goldmark-enclave/core"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/cache"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
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
	vaultLoc       = "internal/data/docs/"
	assetsLoc      = "internal/data/docs/assets/"
	postsLoc       = "internal/data/docs/posts/"
	tagsLoc        = "internal/data/docs/tags/"
	projectsLoc    = "internal/data/docs/projects/"
)

var (
	workers = flag.Int("jobs", 10, "number of parallel uploads")
	cwd     = flag.String("cwd", "", "current working directory")
	client  = s3.NewFromConfig(aws.Config{
		Region:       "auto",
		BaseEndpoint: aws.String("https://fly.storage.tigris.dev"),
		Credentials: &CredHandler{
			ID:  os.Getenv("AWS_ACCESS_KEY_ID"),
			Key: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	})
	mdParser = newMDParser(client)
	llama    = &api.Client{}
)

func main() {
	flag.Parse()
	ctx := context.Background()
	llamaURL, err := url.Parse(os.Getenv("OLLAMA_URL"))
	noError(err)
	llama = api.NewClient(llamaURL, http.DefaultClient)

	if *cwd != "" {
		noError(os.Chdir(*cwd))
	}
	if err := run(ctx); err != nil {
		slog.Error("error", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	eg := errgroup.Group{}
	eg.SetLimit(*workers)

	cache, err := cache.LoadCache(hashFile)
	noError(err)
	defer noError(cache.Close())

	assets, _, err := parse(cache, assetsLoc)
	noError(err)

	err = actualizeAssets(ctx, client, assets)
	noError(err)

	parsedPosts, err := actualize[gen.Post](ctx, cache, postsLoc)
	if err != nil {
		return err
	}

	parsedProjects, err := actualize[gen.Project](ctx, cache, projectsLoc)
	if err != nil {
		return err
	}

	parsedTags, err := actualize[gen.Tag](ctx, cache, tagsLoc)
	if err != nil {
		return err
	}

	// postGen, err := genstruct.NewGenerator(genstruct.Config{
	// 	PackageName: "gen",
	// 	OutputFile:  "internal/data/gen/generated_data.go",
	// }, parsedPosts, parsedTags, parsedProjects)
	// if err != nil {
	// 	return err
	// }
	//
	// return postGen.Generate()
	var result ResultLog
	result.Posts = parsedPosts
	result.Tags = parsedTags
	result.Projects = parsedProjects
	return json.NewEncoder(os.Stdout).Encode(result)
}

// ResultLog is the result of the update command.
type ResultLog struct {
	Posts    []*gen.Post    `json:"posts"`
	Tags     []*gen.Tag     `json:"tags"`
	Projects []*gen.Project `json:"projects"`
}

func actualize[T gen.Post | gen.Tag | gen.Project](
	ctx context.Context,
	cache *cache.Cache,
	loc string,
) ([]*T, error) {
	type result struct {
		item  *T
		index int
		err   error
	}
	type input struct {
		asset Asset
		index int
	}

	// Get files to process
	contents, ignored, err := parse(cache, loc)
	if err != nil {
		return nil, err
	}
	amount := len(contents)
	defer slog.Info("actualization complete", slog.String("loc", loc), slog.Int("ignored", len(ignored)))
	slog.Info("actualizing", slog.String("loc", loc), slog.Int("amount", amount))
	var (
		resultCh = make(chan *result, *workers) // safely collect results from goroutines
		inputCh  = make(chan input, *workers)   // safely collect inputs from goroutines
		parsed   = make([]*T, amount)
	)

	for wi := range *workers {
		go func(wi int) {
			for i := range inputCh {
				slog.Info("processing", "worker", wi, "path", i.asset.Path)
				realized, workErr := realizeMD[T](ctx, parser.NewContext(), mdParser, i.asset)
				resultCh <- &result{item: realized, err: workErr, index: i.index}
			}
		}(wi)
	}

	// Process each asset concurrently
	go func() {
		for i, content := range contents {
			inputCh <- input{asset: content, index: i}
		}
	}()

	j := 1
	for res := range resultCh {
		if res.err != nil {
			return nil, res.err
		}
		parsed[res.index] = res.item
		if j == amount {
			break
		}
		j++
	}
	// Wait for all processing to complete
	close(resultCh) // Close channel after all goroutines are done
	close(inputCh)  // Close channel after all goroutines are done

	// Include previously parsed items that were ignored this time
	remembered := rememberMD[T](ignored)
	parsed = append(parsed, remembered...)

	return parsed, err
}
func hashFileContent(content []byte) string {
	sum := md5.Sum(content)
	return hex.EncodeToString(sum[:])
}
func rememberMD[T gen.Post | gen.Project | gen.Tag](ignored []string) []*T {
	ignoredMap := make(map[string]struct{}, len(ignored))
	for _, slug := range ignored {
		ignoredMap[slug] = struct{}{}
	}

	var typeExample T
	switch any(typeExample).(type) {
	case gen.Post:
		var posts []*gen.Post
		for _, post := range gen.AllPosts {
			if _, ok := ignoredMap[post.Slug]; ok {
				posts = append(posts, post)
			}
		}
		return any(posts).([]*T)
	case gen.Project:
		var projects []*gen.Project
		for _, project := range gen.AllProjects {
			if _, ok := ignoredMap[project.Slug]; ok {
				projects = append(projects, project)
			}
		}
		return any(projects).([]*T)
	case gen.Tag:
		var tags []*gen.Tag
		for _, tag := range gen.AllTags {
			if _, ok := ignoredMap[tag.Slug]; ok {
				tags = append(tags, tag)
			}
		}
		return any(tags).([]*T)
	default:
		return nil
	}
}
func realizeMD[T gen.Post | gen.Project | gen.Tag](
	ctx context.Context,
	pCtx parser.Context,
	md goldmark.Markdown,
	parsed Asset,
) (*T, error) {
	var (
		fm       gen.Embedded
		buf      = bytes.NewBufferString("")
		metadata = frontmatter.Get(pCtx)
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
	resp, err := llama.Embeddings(ctx, &api.EmbeddingRequest{
		Model:  embeddingModel,
		Prompt: fm.RawContent,
	})
	if err != nil {
		return nil, err
	}
	fm.X, fm.Y, fm.Z = projectTo3D(
		resp.Embedding,
		projectionMatrixCreate(gen.EmbedLength, 3),
	)
	copy(fm.Vec[:], resp.Embedding[:gen.EmbedLength])
	return gen.New[T](&fm), nil
}

func actualizeAssets(
	ctx context.Context,
	client *s3.Client,
	assets []Asset,
) (err error) {
	var (
		asset   Asset
		amount  = len(assets)
		inputCh = make(chan Asset, *workers)
		errCh   = make(chan error)
	)

	for wi := range *workers {
		go func(wi int) {
			for asset = range inputCh {
				contentType := mime.TypeByExtension(filepath.Ext(asset.Path))
				slog.Info("uploading", "worker", wi, "path", asset.Path, "contentType", contentType)
				_, err = client.PutObject(ctx, &s3.PutObjectInput{
					Bucket:      aws.String("conneroh"),
					Key:         aws.String(asset.Path),
					Body:        bytes.NewReader(asset.Data),
					ContentType: aws.String(contentType),
				})
				errCh <- err
			}
		}(wi)
	}

	j := 0
	for _, asset = range assets {
		select {
		case err = <-errCh:
			return err
		case inputCh <- asset:
			if j == amount-1 {
				break
			}
			j++
		}
	}
	close(inputCh)
	return nil
}

// Generate a deterministic projection matrix for dimensionality reduction
func projectionMatrixCreate(inputDim, outputDim int) *mat.Dense {
	// Create a matrix with dimensions [outputDim x inputDim]
	// This will multiply directly with the input vector without transposition
	matrix := mat.NewDense(outputDim, inputDim, nil)

	for i := range outputDim { // For each output dimension
		// center point for this dimension's projection focus
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

var exts = []goldmark.Option{
	goldmark.WithExtensions(
		extension.GFM, extension.Footnote,
		extension.Strikethrough, extension.Table,
		extension.TaskList, extension.DefinitionList,
		mathjax.MathJax,
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
}

// newMDParser creates a new markdown parser.
func newMDParser(
	client *s3.Client,
) goldmark.Markdown {
	exts = append(exts, goldmark.WithExtensions(
		&wikilink.Extender{
			Resolver: NewCustomResolver(client),
		},
	))
	return goldmark.New(exts...)
}

// Asset is a struct for embedding assets.
type Asset struct {
	Slug string
	Path string
	Data []byte
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
	ID  string
	Key string
}

// Retrieve returns the credentials for the bucket.
func (b *CredHandler) Retrieve(
	_ context.Context,
) (aws.Credentials, error) {
	return aws.Credentials{AccessKeyID: b.ID, SecretAccessKey: b.Key}, nil
}

// CustomResolver is a wikilink.Resolver that resolves pages and media referenced by
// wikilinks to their destinations.
type CustomResolver struct{ client *s3.Client }

// NewCustomResolver creates a new wikilink resolver.
func NewCustomResolver(client *s3.Client) *CustomResolver {
	return &CustomResolver{client: client}
}

// ResolveWikilink returns the address of the page that the provided
// wikilink points to. The destination will be URL-escaped before
// being placed into a link.
func (c *CustomResolver) ResolveWikilink(n *wikilink.Node) (destination []byte, err error) {
	targetStr := string(n.Target)
	_, err = c.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("conneroh"),
		Key:    aws.String(targetStr),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object(%s): %w", targetStr, err)
	}
	return fmt.Appendf(nil,
		"https://conneroh.fly.storage.tigris.dev/%s",
		targetStr,
	), nil
}
func parse(
	cache *cache.Cache,
	loc string,
) (parsedAssets []Asset, ignoredSlugs []string, err error) {
	err = filepath.WalkDir(
		loc,
		func(fPath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			var (
				hash string
				path string
				slug string
			)
			if d.IsDir() {
				return nil
			}
			body, err := os.ReadFile(fPath)
			if err != nil {
				return err
			}
			hash = hashFileContent(body)
			path = pathify(fPath)
			slug = slugify(fPath)
			if cache.Hashes[fPath] == hash {
				cache.Hashes[fPath] = hash
				ignoredSlugs = append(ignoredSlugs, slug)
				return nil
			}
			cache.Hashes[fPath] = hash
			parsedAssets = append(parsedAssets, Asset{
				Path: path,
				Data: body,
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
	var path string
	var ok bool
	path, ok = strings.CutPrefix(s, assetsLoc)
	if ok {
		return path
	}
	return strings.TrimSuffix(pathify(s), filepath.Ext(path))
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
func noError(err error) {
	if err != nil {
		panic(err)
	}
}
