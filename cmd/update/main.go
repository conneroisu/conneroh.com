// Package main is the main package for the updated the generated code.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/cache"
	"github.com/conneroisu/conneroh.com/internal/credited"
	"github.com/conneroisu/conneroh.com/internal/tigris"
	"github.com/conneroisu/genstruct"
	mathjax "github.com/litao91/goldmark-mathjax"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/quail-ink/goldmark-enclave/core"
	"github.com/rotisserie/eris"
	"github.com/sourcegraph/conc/iter"
	"github.com/spf13/afero"
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
	workers = flag.Int("jobs", 20, "number of parallel workers")
	cwd     = flag.String("cwd", "", "current working directory")
)

const (
	hashFile    = ".config.json"
	vaultLoc    = "internal/data/docs/"
	assetsLoc   = "assets/"
	postsLoc    = "posts/"
	tagsLoc     = "tags/"
	projectsLoc = "projects/"
)

func main() {
	flag.Parse()
	// slog.SetDefault(logger.DefaultLogger)
	if *cwd == "" {
		err := os.Chdir(*cwd)
		if err != nil {
			panic(err)
		}
	}
	if err := run(
		context.Background(),
		os.Getenv,
		*workers,
	); err != nil {
		panic(err)
	}
}

var (
	parsedPosts    []*assets.Post
	parsedProjects []*assets.Project
	parsedTags     []*assets.Tag
)

func run(
	ctx context.Context,
	getenv func(string) string,
	workers int,
) error {
	var (
		fs  = afero.NewBasePathFs(afero.NewOsFs(), vaultLoc)
		m   = NewMD(fs)
		err error

		paths []string

		cacheObj *cache.Cache
		ollama   *credited.OllamaClient
	)
	cacheObj, err = cache.LoadCache(hashFile)
	if err != nil {
		return err
	}
	tigris, err := tigris.New(getenv)
	if err != nil {
		return err
	}
	ollama, err = credited.NewOllamaClient(getenv)
	if err != nil {
		return err
	}
	iterator := iter.Iterator[string]{
		MaxGoroutines: workers,
	}

	err = afero.Walk(
		fs,
		".",
		func(fPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			paths = append(paths, fPath)
			return nil
		})
	if err != nil {
		return err
	}

	iterator.ForEach(paths, func(path *string) {
		slog.Info("parsing", "path", *path)
		defer slog.Info("parsed", "path", *path)

		if strings.HasPrefix(*path, "assets") {
			hErr := handleAsset(ctx, *path, fs, cacheObj, tigris)
			if hErr != nil {
				panic(hErr)
			}
			return
		}
		if strings.HasSuffix(*path, ".md") {
			hErr := handleDoc(ctx, *path, fs, m, cacheObj, ollama)
			if hErr != nil {
				panic(hErr)
			}
			return
		}

		slog.Info("failed to find handler for", "path", *path)
	})

	err = genstruct.NewGenerator(genstruct.WithOutputFile("internal/data/gen/generated_data.go")).
		Generate(parsedPosts, parsedTags, parsedProjects)
	if err != nil {
		return err
	}

	err = cacheObj.Close()
	if err != nil {
		slog.Error("failed to close cache", "error", err)
	}
	return nil
}

var mu sync.Mutex

// Upload uploads the provided asset to the specified bucket.
func Upload(
	ctx context.Context,
	client tigris.Client,
	path string,
	data []byte,
) error {
	extension := filepath.Ext(path)
	if extension == "" {
		return fmt.Errorf("failed to get extension for %s", path)
	}

	contentType := mime.TypeByExtension(extension)
	mu.Lock()
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("conneroh"),
		Key:         aws.String(path),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	mu.Unlock()

	if err != nil {
		// Log error but let errgroup handle it
		slog.Error("asset upload failed", "path", path, "error", err)
		return fmt.Errorf("failed to upload asset %s: %w", path, err)
	}

	return nil
}

// NewMD creates a new markdown parser.
func NewMD(
	fs afero.Fs,
) goldmark.Markdown {
	return goldmark.New(goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
		parser.WithAttribute(),
	),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			extension.WithFootnoteBacklinkClass("footnote-backref"),
			extension.WithFootnoteLinkClass("footnote-ref"),
		),
		goldmark.WithExtensions(
			&wikilink.Extender{
				Resolver: newResolver(fs),
			},
			extension.GFM,
			extension.Footnote,
			extension.Strikethrough,
			extension.Table,
			extension.TaskList,
			extension.DefinitionList,
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
		))
}

// Convert converts the provided markdown to HTML.
func Convert(
	pCtx parser.Context,
	m goldmark.Markdown,
	path string,
	data []byte,
) (assets.Embedded, error) {
	var (
		buf      = bytes.NewBufferString("")
		metadata *frontmatter.Data
		emb      assets.Embedded
		err      error
	)
	err = m.Convert(data, buf, parser.WithContext(pCtx))
	if err != nil {
		return emb, eris.Wrapf(
			err,
			"failed to convert %s's markdown to HTML",
			path,
		)
	}

	// Get frontmatter
	metadata = frontmatter.Get(pCtx)
	if metadata == nil {
		return emb, eris.Errorf(
			"frontmatter is nil for %s",
			path,
		)
	}

	// Decode frontmatter
	err = metadata.Decode(&emb)
	if err != nil {
		return emb, eris.Wrapf(
			err,
			"failed to decode frontmatter of %s",
			path,
		)
	}

	// Set slug and content
	emb.Slug = slugify(path)
	emb.Content = buf.String()
	emb.RawContent = string(data)

	// Get frontmatter
	metadata = frontmatter.Get(pCtx)
	if metadata == nil {
		return emb, eris.Errorf(
			"frontmatter is nil for %s",
			path,
		)
	}

	// Decode frontmatter
	err = metadata.Decode(&emb)
	if err != nil {
		return emb, eris.Wrapf(
			err,
			"failed to decode frontmatter of %s",
			path,
		)
	}
	return emb, nil
}

// resolver is a wikilink.Resolver that resolves pages and media referenced by
// wikilinks to their destinations.
type resolver struct {
	fs afero.Fs
}

// newResolver creates a new wikilink resolver.
func newResolver(
	fs afero.Fs,
) *resolver {
	return &resolver{
		fs: fs,
	}
}

// ResolveWikilink returns the address of the page that the provided
// wikilink points to. The destination will be URL-escaped before
// being placed into a link.
func (r *resolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	var (
		err       error
		targetStr = string(n.Target)
	)
	_, err = afero.ReadFile(r.fs, fmt.Sprintf("/assets/%s", targetStr))
	if err != nil {
		return nil, err
	}
	return fmt.Appendf(nil,
		"https://conneroh.fly.storage.tigris.dev/%s",
		targetStr,
	), nil
}

func slugify(s string) string {
	var path string
	var ok bool
	path, ok = strings.CutPrefix(s, assetsLoc)
	if ok {
		return path
	}
	return strings.TrimSuffix(pathify(s), filepath.Ext(s))
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
	panic(fmt.Errorf("failed to pathify %s", s))
}

func match(path string, emb *assets.Embedded) {
	if strings.HasPrefix(path, postsLoc) {
		parsedPosts = append(parsedPosts, assets.New[assets.Post](emb))
	} else if strings.HasPrefix(path, projectsLoc) {
		parsedProjects = append(parsedProjects, assets.New[assets.Project](emb))
	} else if strings.HasPrefix(path, tagsLoc) {
		parsedTags = append(parsedTags, assets.New[assets.Tag](emb))
	}
}

func handleDoc(
	ctx context.Context,
	path string,
	fs afero.Fs,
	m goldmark.Markdown,
	cacheObj *cache.Cache,
	ollama *credited.OllamaClient,
) error {
	var (
		err error
		ok  bool
		emb assets.Embedded
	)

	pCtx := parser.NewContext()
	content, err := afero.ReadFile(fs, path)
	if err != nil {
		return err
	}
	emb, err = Convert(pCtx, m, path, content)
	if err != nil {
		return err
	}

	err = assets.Defaults(&emb)
	if err != nil {
		return err
	}

	err = assets.Validate(&emb)
	if err != nil {
		return err
	}

	hash := cache.Hash([]byte(emb.RawContent))
	known, ok := cacheObj.Get(path)
	if ok && known == hash {
		slog.Info("found cached version", "path", path)
		minV, ok := cacheObj.GetDoc(path)
		if ok {
			emb.X = minV.X
			emb.Y = minV.Y
			emb.Z = minV.Z
		} else {
			slog.Error("failed to get min version", "path", path)
		}
	} else {
		err = ollama.Embeddings(ctx, emb.RawContent, &emb)
		if err != nil {
			return err
		}
	}

	cacheObj.SetDoc(path, emb)
	cacheObj.Set(path, hash)
	match(path, &emb)
	return nil
}

func handleAsset(
	ctx context.Context,
	path string,
	fs afero.Fs,
	cacheObj *cache.Cache,
	tigris tigris.Client,
) error {
	var (
		err error
	)

	slog.Debug("updating asset", "path", path)
	defer slog.Debug("asset updated", "path", path)

	data, perr := afero.ReadFile(fs, path)
	if perr != nil {
		return perr
	}
	hash := cache.Hash(data)
	known, ok := cacheObj.Get(path)
	if ok && known == hash {
		slog.Debug("asset unchanged", "path", path)
		return nil
	}
	cacheObj.Set(path, hash)
	slog.Debug("asset changed uploading...", "path", path, "pathified", pathify(path))
	err = Upload(ctx, tigris, pathify(path), data)
	if err != nil {
		return eris.Wrapf(err, "failed to upload asset %s", path)
	}
	return nil
}
