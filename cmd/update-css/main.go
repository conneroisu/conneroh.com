// Package main updates the CSS.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/conneroisu/conneroh.com/cmd/conneroh/components"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/conneroisu/twerge"
	"github.com/rotisserie/eris"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"

	_ "modernc.org/sqlite"
)

var cwd = flag.String("cwd", "", "current working directory")

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("(update-css) Done in %s.\n", elapsed)
	}()
	flag.Parse()
	if *cwd != "" {
		err := os.Chdir(*cwd)
		if err != nil {
			panic(err)
		}
	}
	sqlDB, err := sql.Open("sqlite", assets.DBName())
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqlDB, sqlitedialect.New())
	assets.RegisterModels(db)
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	var (
		allPosts    []*assets.Post
		allProjects []*assets.Project
		allTags     []*assets.Tag
	)

	err = db.NewSelect().
		Model(&allPosts).
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Scan(ctx)
	if err != nil {
		return eris.Wrapf(err, "(update-css) failed to get posts")
	}
	err = db.NewSelect().
		Model(&allProjects).
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Scan(ctx)
	if err != nil {
		return eris.Wrapf(err, "(update-css) failed to get projects")
	}
	err = db.NewSelect().
		Model(&allTags).
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Scan(ctx)
	if err != nil {
		return eris.Wrapf(err, "(update-css) failed to get tags")
	}

	return twerge.CodeGen(
		twerge.Default(),
		"internal/data/css/classes.go",
		"input.css",
		"internal/data/css/classes.html",
		layouts.Page(views.Home(
			&allPosts,
			&allProjects,
			&allTags,
		)),
		layouts.Page(views.List(
			routing.PostPluralPath,
			&allPosts,
			&allProjects,
			&allTags,
			"",
			1,
			10,
		)),
		layouts.Page(views.Code500()),
		layouts.Page(views.List(
			routing.ProjectPluralPath,
			&allPosts,
			&allProjects,
			&allTags,
			"",
			1,
			10,
		)),
		layouts.Page(views.List(
			routing.TagsPluralPath,
			&allPosts,
			&allProjects,
			&allTags,
			"",
			1,
			10,
		)),
		components.TagControl(
			&assets.Tag{},
		),
		layouts.Page(views.Post(
			allPosts[0],
		)),
		layouts.Page(views.Project(
			allProjects[0],
		)),
		layouts.Page(views.Tag(
			allTags[0],
		)),
		views.Post(
			allPosts[0],
		),
		layouts.Layout("hello"),
		components.ThankYou(),
	)
}
