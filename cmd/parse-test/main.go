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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/cache"
	"github.com/conneroisu/conneroh.com/internal/credited"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/tigris"
	mathjax "github.com/litao91/goldmark-mathjax"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/quail-ink/goldmark-enclave/core"
	"github.com/rotisserie/eris"
	"github.com/sourcegraph/conc/pool"
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
	if err := run(
		context.Background(),
		os.Getenv,
		*workers,
	); err != nil {
		panic(err)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	workers int,
) error {
	var (
		fs   = afero.NewBasePathFs(afero.NewOsFs(), vaultLoc)
		m    = NewMathjax(ctx, fs)
		err  error
		pool = pool.New().WithMaxGoroutines(workers).WithContext(ctx)

		docPaths   []string
		assetPaths []string

		parsedPosts    []*gen.Post
		parsedProjects []*gen.Project
		parsedTags     []*gen.Tag
		cacheObj       *cache.Cache
	)
	cacheObj, err = cache.LoadCache(hashFile)
	if err != nil {
		return err
	}
	tigris, err := tigris.New(getenv)
	if err != nil {
		return err
	}
	ollama, err := credited.NewOllamaClient(getenv)
	if err != nil {
		return err
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
			ext := filepath.Ext(fPath)
			if ext == ".md" {
				docPaths = append(docPaths, fPath)
				return nil
			}
			assetPaths = append(assetPaths, fPath)
			return nil
		})
	if err != nil {
		return err
	}

	for _, path := range docPaths {
		pool.Go(func(_ context.Context) error {
			slog.Info("parsing", "path", path)
			defer slog.Info("parsed", "path", path)

			pCtx := parser.NewContext()
			content, derr := afero.ReadFile(fs, path)
			if derr != nil {
				return err
			}
			if cacheObj.OldDoc(path, content) {
				emb, ok := cacheObj.GetDoc(path)
				if !ok {
					return eris.Errorf("failed to get doc from cache")
				}
				if strings.HasPrefix(path, postsLoc) {
					parsedPosts = append(parsedPosts, gen.New[gen.Post](&emb))
				} else if strings.HasPrefix(path, projectsLoc) {
					parsedProjects = append(parsedProjects, gen.New[gen.Project](&emb))
				} else if strings.HasPrefix(path, tagsLoc) {
					parsedTags = append(parsedTags, gen.New[gen.Tag](&emb))
				}
				return nil
			}
			emb, derr := Convert(pCtx, m, path, content)
			if derr != nil {
				return err
			}

			err = gen.Defaults(&emb)
			if err != nil {
				return err
			}

			err = gen.Validate(&emb)
			if err != nil {
				return err
			}

			err = ollama.Embeddings(ctx, emb.RawContent, &emb)
			if err != nil {
				return err
			}

			if strings.HasPrefix(path, postsLoc) {
				parsedPosts = append(parsedPosts, gen.New[gen.Post](&emb))
			} else if strings.HasPrefix(path, projectsLoc) {
				parsedProjects = append(parsedProjects, gen.New[gen.Project](&emb))
			} else if strings.HasPrefix(path, tagsLoc) {
				parsedTags = append(parsedTags, gen.New[gen.Tag](&emb))
			}
			cacheObj.SetDoc(path, emb)
			return nil
		})
	}
	for _, path := range assetPaths {
		path := path
		pool.Go(func(poolCtx context.Context) error {
			slog.Debug("updating asset", "path", path)
			defer slog.Debug("asset updated", "path", path)
			data, perr := afero.ReadFile(fs, path)
			if perr != nil {
				return perr
			}
			hash := cache.Hash(data)
			known, ok := cacheObj.Hashes[path]
			if ok && known == hash {
				return nil
			}
			cacheObj.Set(path, hash)
			err = Upload(poolCtx, tigris, path, data)
			if err != nil {
				return err
			}
			return nil
		})
	}

	err = pool.Wait()
	if err != nil {
		return err
	}

	err = cacheObj.Close()
	if err != nil {
		slog.Error("failed to close cache", "error", err)
	}
	return nil
}

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
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("conneroh"),
		Key:         aws.String(path),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		// Log error but let errgroup handle it
		slog.Error("asset upload failed", "path", path, "error", err)
		return fmt.Errorf("failed to upload asset %s: %w", path, err)
	}

	return nil
}

// NewMathjax creates a new markdown parser.
func NewMathjax(
	ctx context.Context,
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
				Resolver: newResolver(ctx, fs),
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
) (gen.Embedded, error) {
	var (
		buf      = bytes.NewBufferString("")
		metadata *frontmatter.Data
		emb      gen.Embedded
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
	Ctx context.Context
	fs  afero.Fs
}

// newResolver creates a new wikilink resolver.
func newResolver(
	ctx context.Context,
	fs afero.Fs,
) *resolver {
	return &resolver{
		fs:  fs,
		Ctx: ctx,
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
