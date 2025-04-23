// Package main is the main package for the updated the generated code.
// TODO: Remove afero dependency.
// TODO: Action Queue
package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/copygen"
	"github.com/conneroisu/conneroh.com/internal/credited"
	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/conneroisu/conneroh.com/internal/tigris"
	mathjax "github.com/litao91/goldmark-mathjax"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/quail-ink/goldmark-enclave/core"
	"github.com/rotisserie/eris"
	"github.com/sourcegraph/conc/iter"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
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
	_ "modernc.org/sqlite"
)

var (
	workers = flag.Int("jobs", 80, "number of parallel workers")
	cwd     = flag.String("cwd", "", "current working directory")
)

func main() {
	flag.Parse()
	slog.SetDefault(logger.DefaultLogger)
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

func run(
	ctx context.Context,
	getenv func(string) string,
	workers int,
) error {
	var (
		fs       = afero.NewBasePathFs(afero.NewOsFs(), assets.VaultLoc)
		renderer = NewMD(fs)
		paths    []string
	)
	tigris, err := tigris.New(getenv)
	if err != nil {
		return err
	}
	ollama, err := credited.NewOllamaClient(getenv)
	if err != nil {
		return err
	}
	iterator := iter.Iterator[string]{MaxGoroutines: workers}
	sqldb, err := sql.Open("sqlite", assets.DBName())
	if err != nil {
		return err
	}
	defer sqldb.Close()
	db := bun.NewDB(sqldb, sqlitedialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	err = assets.InitDB(ctx, db)
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
			if assets.IsAllowedAsset(fPath) {
				paths = append(paths, fPath)
			}
			return nil
		})
	if err != nil {
		return err
	}

	iterator.ForEach(paths, func(path *string) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic", "error", err, "handing", *path)
				log.Printf("panic: %v", err)
				panic(err)
			}
		}()
		slog.Info("parsing", "path", *path)
		defer slog.Info("parsed", "path", *path)
		hErr := handle(ctx, db, *path, fs, renderer, ollama, tigris)
		if hErr != nil {
			slog.Error("failed to handle doc", "path", *path, "error", hErr)
			panic(hErr)
		}
	})
	return nil
}

// UploadManager manages concurrent uploads using multiple mutexes
type UploadManager struct {
	mutexes    []sync.Mutex
	mutexCount int
	inUse      []bool
	lock       sync.Mutex // Meta-mutex to protect the inUse array
}

// NewUploadManager creates a new UploadManager with the specified number of mutexes
func NewUploadManager(count int) *UploadManager {
	if count <= 0 {
		count = 10 // Default to 10 mutexes if invalid count provided
	}

	manager := &UploadManager{
		mutexes:    make([]sync.Mutex, count),
		mutexCount: count,
		inUse:      make([]bool, count),
	}

	return manager
}

// AcquireMutex finds an available mutex, locks it, and returns its index
// This function will block until a mutex is available
func (um *UploadManager) AcquireMutex() int {
	for {
		um.lock.Lock()
		for i := range um.mutexCount {
			if !um.inUse[i] {
				// Mark this mutex as in use
				um.inUse[i] = true
				um.lock.Unlock()

				// Lock the mutex
				um.mutexes[i].Lock()
				return i
			}
		}
		// No mutex available, release meta-lock and wait a bit before trying again
		um.lock.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
}

// ReleaseMutex unlocks the mutex at the specified index and marks it as available
func (um *UploadManager) ReleaseMutex(index int) {
	if index < 0 || index >= um.mutexCount {
		return
	}

	// First unlock the actual mutex
	um.mutexes[index].Unlock()

	// Then mark it as not in use
	um.lock.Lock()
	um.inUse[index] = false
	um.lock.Unlock()
}

// Global upload manager instance
var uploadManager = NewUploadManager(20) // Create 20 mutexes for more concurrency

// Upload uploads the provided asset to the specified bucket.
// This version supports concurrent uploads using multiple mutexes
func Upload(
	ctx context.Context,
	client tigris.Client,
	path string,
	data []byte,
) error {
	slog.Debug(
		"asset changed uploading...",
		"path", path,
	)
	extension := filepath.Ext(path)
	if extension == "" {
		return fmt.Errorf("failed to get extension for %s", path)
	}
	contentType := mime.TypeByExtension(extension)

	// Acquire a mutex for this upload - this will block until one is available
	mutexIndex := uploadManager.AcquireMutex()

	// Ensure we release the mutex when done
	defer uploadManager.ReleaseMutex(mutexIndex)

	// Check if context is done before proceeding with upload
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Continue with upload
	}

	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("conneroh"),
		Key:         aws.String(path),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		// Log error but let errgroup handle it
		return fmt.Errorf("failed to upload asset %s: %w", path, err)
	}

	slog.Info("asset upload successful", "path", path, "mutexIndex", mutexIndex)
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
	m goldmark.Markdown,
	path string,
	data []byte,
) (assets.Doc, error) {
	var (
		metadata *frontmatter.Data
		doc      assets.Doc
		err      error
	)
	pCtx := parser.NewContext()
	buf := bytes.NewBufferString("")

	err = m.Convert(data, buf, parser.WithContext(pCtx))
	if err != nil {
		return doc, eris.Wrapf(
			err,
			"failed to convert %s's markdown to HTML",
			path,
		)
	}

	// Get frontmatter
	metadata = frontmatter.Get(pCtx)
	if metadata == nil {
		return doc, eris.Errorf(
			"frontmatter is nil for %s",
			path,
		)
	}

	// Decode frontmatter
	err = metadata.Decode(&doc)
	if err != nil {
		return doc, eris.Wrapf(
			err,
			"failed to decode frontmatter of %s",
			path,
		)
	}

	// Set slug and content
	doc.Slug = assets.Slugify(path)
	doc.Content = buf.String()
	doc.RawContent = data

	// Get frontmatter
	metadata = frontmatter.Get(pCtx)
	if metadata == nil {
		return doc, eris.Errorf(
			"frontmatter is nil for %s",
			path,
		)
	}

	// Decode frontmatter
	err = metadata.Decode(&doc)
	if err != nil {
		return doc, eris.Wrapf(
			err,
			"failed to decode frontmatter of %s",
			path,
		)
	}

	return doc, nil
}

// resolver is a wikilink.Resolver that resolves pages and media referenced by
// wikilinks to their destinations.
type resolver struct{ fs afero.Fs }

// newResolver creates a new wikilink resolver.
func newResolver(fs afero.Fs) *resolver {
	return &resolver{fs: fs}
}

// ResolveWikilink returns the address of the page that the provided
// wikilink points to. The destination will be URL-escaped before
// being placed into a link.
func (r *resolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	var (
		err       error
		targetStr = string(n.Target)
	)
	// first check if it a url
	if strings.HasPrefix(targetStr, "http") {
		return []byte(targetStr), nil
	}
	_, err = afero.ReadFile(r.fs, fmt.Sprintf("/assets/%s", targetStr))
	if err != nil {
		return nil, err
	}

	return fmt.Appendf(nil,
		"https://conneroh.fly.storage.tigris.dev/%s",
		targetStr,
	), nil
}

func match(
	ctx context.Context,
	db *bun.DB,
	doc *assets.Doc,
	ollama *credited.OllamaClient,
	tigris tigris.Client,
) error {
	pathified := assets.Pathify(doc.Path)
	switch {
	case strings.HasPrefix(doc.Path, assets.PostsLoc) &&
		strings.HasSuffix(doc.Path, ".md"):
		var post assets.Post
		_, err := db.NewSelect().Model(assets.EmpPost).
			Where("path = ?", doc.Path).Exec(ctx, &post)
		if errors.Is(err, sql.ErrNoRows) {
			err = ollama.Embeddings(ctx, doc.RawContent, doc)
			if err != nil {
				return err
			}
			copygen.ToPost(&post, doc)
			return assets.SavePost(ctx, db, &post)
		}
		if err != nil {
			return err
		}
		if post.Hash != doc.Hash {
			err = ollama.Embeddings(ctx, doc.RawContent, doc)
			if err != nil {
				return err
			}
			copygen.ToPost(&post, doc)
			return assets.SavePost(ctx, db, &post)
		}
	case strings.HasPrefix(doc.Path, assets.ProjectsLoc) && strings.HasSuffix(doc.Path, ".md"):
		var project assets.Project
		_, err := db.NewSelect().Model(assets.EmpProject).
			Where("path = ?", doc.Path).Exec(ctx, &project)
		if errors.Is(err, sql.ErrNoRows) {
			err = ollama.Embeddings(ctx, doc.RawContent, doc)
			if err != nil {
				return err
			}
			copygen.ToProject(&project, doc)
			return assets.SaveProject(ctx, db, &project)
		}
		if err != nil {
			return err
		}
		if project.Hash != doc.Hash {
			err = ollama.Embeddings(ctx, doc.RawContent, doc)
			if err != nil {
				return err
			}
			copygen.ToProject(&project, doc)
			err = assets.SaveProject(ctx, db, &project)
			if err != nil {
				return err
			}
		}
	case strings.HasPrefix(doc.Path, assets.TagsLoc) && strings.HasSuffix(doc.Path, ".md"):
		var tag assets.Tag
		_, err := db.NewSelect().Model(assets.EmpTag).
			Where("path = ?", doc.Path).Exec(ctx, &tag)
		if errors.Is(err, sql.ErrNoRows) {
			err = ollama.Embeddings(ctx, doc.RawContent, doc)
			if err != nil {
				return err
			}
			copygen.ToTag(&tag, doc)
			return assets.SaveTag(ctx, db, &tag)
		}
		if err != nil {
			return err
		}
		if tag.Hash != doc.Hash {
			err = ollama.Embeddings(ctx, doc.RawContent, doc)
			if err != nil {
				return err
			}
			copygen.ToTag(&tag, doc)
			return assets.SaveTag(ctx, db, &tag)
		}
	case strings.HasPrefix(doc.Path, "assets"):
		var asset assets.Asset
		_, err := db.NewSelect().
			Model(assets.EmpAsset).
			Where("path = ?", doc.Path).
			Exec(ctx, &asset)
		if errors.Is(err, sql.ErrNoRows) {
			err = Upload(ctx, tigris, pathified, []byte(doc.RawContent))
			if err != nil {
				return eris.Wrapf(err, "failed to upload asset %s", doc.Path)
			}
			copygen.ToAsset(&asset, doc)
			_, err = db.NewInsert().Model(&assets.Asset{
				Path: doc.Path,
				Hash: doc.Hash,
			}).Exec(ctx)
			if err != nil {
				return eris.Wrapf(err, "failed to insert asset %s", doc.Path)
			}
			return nil
		}
		if err != nil {
			return err
		}
		if asset.Hash != doc.Hash {
			_, err = db.NewUpdate().Model(&asset).Where("path = ?", doc.Path).Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to update asset after changes were made %s",
					doc.Path,
				)
			}
			err = Upload(ctx, tigris, pathified, []byte(doc.RawContent))
			if err != nil {
				return eris.Wrapf(err, "failed to upload asset %s", doc.Path)
			}
			return nil
		}
		return nil
	}

	return nil
}

func handle(
	ctx context.Context,
	db *bun.DB,
	path string,
	fs afero.Fs,
	m goldmark.Markdown,
	ollama *credited.OllamaClient, // TODO: take interface
	tigris tigris.Client,
) error {
	var doc assets.Doc
	content, err := afero.ReadFile(fs, path)
	if err != nil {
		return err
	}
	doc.RawContent = content
	// doc.Path = path
	if strings.HasSuffix(path, ".md") {
		doc, err = Convert(m, path, content)
		if err != nil {
			return err
		}
		doc.Hash = assets.Hash([]byte(doc.RawContent))

		err = assets.Defaults(&doc)
		if err != nil {
			return err
		}

		err = assets.Validate(path, &doc)
		if err != nil {
			return err
		}

		err = match(ctx, db, &doc, ollama, tigris)
		if err != nil {
			return err
		}
		return nil
	}

	// doc.Hash = assets.Hash([]byte(doc.RawContent))
	// doc.Slug = assets.Slugify(path)
	// doc.Path = path
	// err = match(ctx, db, &doc, ollama, tigris)
	// if err != nil {
	// 	return err
	// }
	return nil
}
