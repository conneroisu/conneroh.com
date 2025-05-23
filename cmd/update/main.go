// Package main is the entry point for the application
package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/copygen"
	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
	_ "modernc.org/sqlite"
)

const (
	numWorkers     = 20
	taskBufferInt  = 1000
	fullAssetLoc   = assets.AssetsLoc
	fullPostLoc    = assets.PostsLoc
	fullProjectLoc = assets.ProjectsLoc
	fullTagLoc     = assets.TagsLoc
)

var (
	workers    = flag.Int("workers", numWorkers, "number of parallel workers")
	taskBuffer = flag.Int("buffer", taskBufferInt, "size of task buffer")
)

func main() {
	flag.Parse()
	slog.SetDefault(logger.DefaultLogger)

	// Create context that will be canceled on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP)
	defer stop()

	err := UpdateDB(ctx, os.Getenv, *workers, *taskBuffer)
	if err != nil {
		panic(err)
	}
}

// UpdateDB executes the main application logic.
func UpdateDB(
	ctx context.Context,
	getenv func(string) string,
	_, _ int,
) error {
	var relFns []assets.RelationshipFn
	sqldb, err := sql.Open("sqlite", assets.DBName())
	if err != nil {
		return eris.Wrap(err, "failed to open database")
	}
	defer sqldb.Close()
	db := bun.NewDB(sqldb, sqlitedialect.New())
	if os.Getenv("DEBUG") == "true" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	err = assets.InitDB(ctx, db)
	if err != nil {
		return eris.Wrap(err, "failed to initialize database")
	}
	fs := afero.NewBasePathFs(afero.NewOsFs(), assets.VaultLoc)
	// ol, err := llama.NewOllamaClient(getenv)
	// if err != nil {
	// 	return eris.Wrap(err, "failed to create Ollama client")
	// }
	ti, err := assets.NewTigris(getenv)
	if err != nil {
		return eris.Wrap(err, "failed to create Tigris client")
	}
	bucketName := getenv("BUCKET_NAME")
	if bucketName == "" {
		return eris.New("BUCKET_NAME environment variable is not set")
	}
	md := assets.NewMD(fs)

	todoAssets, err := assets.HashDirMatch(ctx, fs, fullAssetLoc, db)
	if err != nil {
		return err
	}
	for _, item := range todoAssets {
		slog.Info("uploading to S3", "path", item.Path)
		err = assets.UploadToS3(
			ctx,
			ti,
			bucketName,
			item.Path,
			[]byte(item.Content),
		)
		if err != nil {
			return eris.Wrap(err, "failed to upload to S3")
		}
	}

	todoPosts, err := assets.HashDirMatch(ctx, fs, fullPostLoc, db)
	if err != nil {
		return eris.Wrap(err, "failed to hash posts")
	}
	for _, item := range todoPosts {
		slog.Info("processing post", "path", item.Path)
		var (
			post  assets.Post
			doc   *assets.Doc
			relFn assets.RelationshipFn
		)
		doc, err = assets.ParseMarkdown(md, item)
		if err != nil {
			return eris.Wrap(err, "failed to parse markdown")
		}
		copygen.ToPost(&post, doc)
		relFn, err = assets.UpsertPost(ctx, db, &post)
		if err != nil {
			return eris.Wrap(err, "failed to upsert post")
		}
		relFns = append(relFns, relFn)
	}

	todoProjects, err := assets.HashDirMatch(ctx, fs, fullProjectLoc, db)
	if err != nil {
		return eris.Wrap(err, "failed to hash projects")
	}
	for _, item := range todoProjects {
		slog.Info("processing project", "path", item.Path)
		var (
			project assets.Project
			doc     *assets.Doc
			relFn   assets.RelationshipFn
		)
		doc, err = assets.ParseMarkdown(md, item)
		if err != nil {
			return eris.Wrap(err, "failed to parse markdown")
		}
		copygen.ToProject(&project, doc)
		relFn, err = assets.UpsertProject(ctx, db, &project)
		if err != nil {
			return eris.Wrap(err, "failed to upsert project")
		}
		relFns = append(relFns, relFn)
	}

	todoTags, err := assets.HashDirMatch(ctx, fs, fullTagLoc, db)
	if err != nil {
		return eris.Wrap(err, "failed to hash tags")
	}
	for _, item := range todoTags {
		slog.Info("processing tag", "path", item.Path)
		var (
			tag   assets.Tag
			doc   *assets.Doc
			relFn assets.RelationshipFn
		)
		doc, err = assets.ParseMarkdown(md, item)
		if err != nil {
			return eris.Wrap(err, "failed to parse markdown")
		}
		copygen.ToTag(&tag, doc)
		relFn, err = assets.UpsertTag(ctx, db, &tag)
		if err != nil {
			return eris.Wrap(err, "failed to upsert tag")
		}
		relFns = append(relFns, relFn)
	}

	slog.Info("upserting relationships")
	for _, fn := range relFns {
		err := fn(ctx)
		if err != nil {
			return eris.Wrap(err, "failed to run relationship function")
		}
	}

	return nil
}
