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
	numWorkers        = 20
	taskBufferInt     = 1000
	fullAssetLoc      = assets.AssetsLoc
	fullPostLoc       = assets.PostsLoc
	fullProjectLoc    = assets.ProjectsLoc
	fullTagLoc        = assets.TagsLoc
	fullEmploymentLoc = assets.EmploymentsLoc
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

	err := run(ctx, os.Getenv)
	if err != nil {
		panic(err)
	}
}

// run executes the main application logic.
func run(
	ctx context.Context,
	getenv func(string) string,
) error {
	var (
		relFns []assets.RelationshipFn
		items  []assets.DirMatchItem
	)
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
	ti, err := assets.NewTigris(getenv)
	if err != nil {
		return eris.Wrap(err, "failed to create Tigris client")
	}
	bucketName := getenv("BUCKET_NAME")
	if bucketName == "" {
		return eris.New("BUCKET_NAME environment variable is not set")
	}
	md := assets.NewMD(fs)

	// Assets
	items, err = assets.HashDirMatch(ctx, fs, assets.AssetsLoc, db)
	if err != nil {
		return err
	}
	for _, item := range items {
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

	// Posts
	items, err = assets.HashDirMatch(ctx, fs, assets.PostsLoc, db)
	if err != nil {
		return eris.Wrap(err, "failed to hash posts")
	}
	for _, item := range items {
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

	// Projects
	items, err = assets.HashDirMatch(ctx, fs, assets.ProjectsLoc, db)
	if err != nil {
		return eris.Wrap(err, "failed to hash projects")
	}
	for _, item := range items {
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

	// Tags
	items, err = assets.HashDirMatch(ctx, fs, assets.TagsLoc, db)
	if err != nil {
		return eris.Wrap(err, "failed to hash tags")
	}
	for _, item := range items {
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

	todoEmployments, err := assets.HashDirMatch(ctx, fs, fullEmploymentLoc, db)
	if err != nil {
		return eris.Wrap(err, "failed to hash employments")
	}
	for _, item := range todoEmployments {
		slog.Info("processing employment", "path", item.Path)
		var (
			employment assets.Employment
			doc        *assets.Doc
			relFn      assets.RelationshipFn
		)
		doc, err = assets.ParseMarkdown(md, item)
		if err != nil {
			return eris.Wrap(err, "failed to parse markdown")
		}
		copygen.ToEmployment(&employment, doc)
		relFn, err = assets.UpsertEmployment(ctx, db, &employment)
		if err != nil {
			return eris.Wrap(err, "failed to upsert employment")
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
