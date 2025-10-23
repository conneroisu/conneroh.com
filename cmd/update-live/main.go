// Package main provides live update functionality for assets and CSS generation.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/components"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/copygen"
	"github.com/conneroisu/conneroh.com/internal/logger"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/conneroisu/twerge"
	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"

	_ "modernc.org/sqlite"
)

var cwd = flag.String("cwd", "", "current working directory")

func main() {
	flag.Parse()
	if *cwd != "" {
		if err := os.Chdir(*cwd); err != nil {
			fmt.Fprintln(os.Stderr, eris.Wrap(err, "failed to change directory"))
			os.Exit(1)
		}
	}

	slog.SetDefault(logger.DefaultLogger)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
	)
	defer stop()

	if err := run(ctx, os.Getenv); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string) error {
	start := time.Now()
	defer func() {
		fmt.Printf("(update-live) Done in %s.\n", time.Since(start))
	}()

	if err := updateDatabase(ctx, getenv); err != nil {
		return err
	}

	if err := regenerateCSS(ctx); err != nil {
		return err
	}

	return nil
}

func updateDatabase(ctx context.Context, getenv func(string) string) error {
	sqldb, err := sql.Open("sqlite", assets.DBName())
	if err != nil {
		return eris.Wrap(err, "failed to open database")
	}
	defer sqldb.Close()

	db := bun.NewDB(sqldb, sqlitedialect.New())
	if getenv("DEBUG") == "true" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	if err = assets.InitDB(ctx, db); err != nil {
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
	relFns := make([]assets.RelationshipFn, 0)

	// Assets
	items, err := assets.HashDirMatch(ctx, fs, assets.AssetsLoc, db)
	if err != nil {
		return err
	}
	for _, item := range items {
		slog.Info("uploading to S3", "path", item.Path)
		if err = assets.UploadToS3(
			ctx,
			ti,
			bucketName,
			item.Path,
			[]byte(item.Content),
		); err != nil {
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
		doc, docErr := assets.ParseMarkdown(md, item)
		if docErr != nil {
			return eris.Wrap(docErr, "failed to parse markdown")
		}

		var post assets.Post
		copygen.ToPost(&post, doc)
		relFn, relErr := assets.UpsertPost(ctx, db, &post)
		if relErr != nil {
			return eris.Wrap(relErr, "failed to upsert post")
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
		doc, docErr := assets.ParseMarkdown(md, item)
		if docErr != nil {
			return eris.Wrap(docErr, "failed to parse markdown")
		}

		var project assets.Project
		copygen.ToProject(&project, doc)
		relFn, relErr := assets.UpsertProject(ctx, db, &project)
		if relErr != nil {
			return eris.Wrap(relErr, "failed to upsert project")
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
		doc, docErr := assets.ParseMarkdown(md, item)
		if docErr != nil {
			return eris.Wrap(docErr, "failed to parse markdown")
		}

		var tag assets.Tag
		copygen.ToTag(&tag, doc)
		relFn, relErr := assets.UpsertTag(ctx, db, &tag)
		if relErr != nil {
			return eris.Wrap(relErr, "failed to upsert tag")
		}
		relFns = append(relFns, relFn)
	}

	// Employments
	items, err = assets.HashDirMatch(ctx, fs, assets.EmploymentsLoc, db)
	if err != nil {
		return eris.Wrap(err, "failed to hash employments")
	}
	for _, item := range items {
		slog.Info("processing employment", "path", item.Path)
		doc, docErr := assets.ParseMarkdown(md, item)
		if docErr != nil {
			return eris.Wrap(docErr, "failed to parse markdown")
		}

		var employment assets.Employment
		copygen.ToEmployment(&employment, doc)
		relFn, relErr := assets.UpsertEmployment(ctx, db, &employment)
		if relErr != nil {
			return eris.Wrap(relErr, "failed to upsert employment")
		}
		relFns = append(relFns, relFn)
	}

	slog.Info("upserting relationships")
	for _, fn := range relFns {
		if err := fn(ctx); err != nil {
			return eris.Wrap(err, "failed to run relationship function")
		}
	}

	return nil
}

func regenerateCSS(ctx context.Context) error {
	sqlDB, err := sql.Open("sqlite", assets.DBName())
	if err != nil {
		return eris.Wrap(err, "failed to open database")
	}
	defer sqlDB.Close()

	db := bun.NewDB(sqlDB, sqlitedialect.New())
	assets.RegisterModels(db)
	if os.Getenv("DEBUG") != "" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	var (
		allPosts       []*assets.Post
		allProjects    []*assets.Project
		allTags        []*assets.Tag
		allEmployments []*assets.Employment
	)

	if err = db.NewSelect().
		Model(&allPosts).
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Scan(ctx); err != nil {
		return eris.Wrapf(err, "(update-live) failed to get posts")
	}

	if err = db.NewSelect().
		Model(&allProjects).
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Scan(ctx); err != nil {
		return eris.Wrapf(err, "(update-live) failed to get projects")
	}

	if err = db.NewSelect().
		Model(&allTags).
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Scan(ctx); err != nil {
		return eris.Wrapf(err, "(update-live) failed to get tags")
	}

	if err = db.NewSelect().
		Model(&allEmployments).
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Relation("Employments").
		Scan(ctx); err != nil {
		return eris.Wrapf(err, "(update-live) failed to get employments")
	}

	comps := []templ.Component{}
	for _, post := range allPosts {
		comps = append(comps, views.Post(post))
	}
	for _, project := range allProjects {
		comps = append(comps, views.Project(project))
	}
	for _, tag := range allTags {
		comps = append(comps, views.Tag(tag))
	}
	for _, employment := range allEmployments {
		comps = append(comps, views.Employment(employment))
	}

	comps = append(comps, views.List(
		routing.ProjectPluralPath,
		&allPosts,
		&allProjects,
		&allTags,
		&allEmployments,
		"",
		1,
		10,
	))
	comps = append(comps, views.List(
		routing.TagsPluralPath,
		&allPosts,
		&allProjects,
		&allTags,
		&allEmployments,
		"",
		1,
		10,
	))
	comps = append(comps, views.List(
		routing.PostPluralPath,
		&allPosts,
		&allProjects,
		&allTags,
		&allEmployments,
		"",
		1,
		10,
	))
	comps = append(comps, layouts.Page(views.Home(
		&allPosts,
		&allProjects,
		&allTags,
		&allEmployments,
	)))
	comps = append(comps, views.Code500())
	comps = append(comps, components.ThankYou())

	return twerge.CodeGen(
		twerge.Default(),
		"cmd/conneroh/classes/classes.go",
		"input.css",
		"cmd/conneroh/classes/classes.html",
		comps...,
	)
}
