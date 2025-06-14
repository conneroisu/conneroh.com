// Package main updates the CSS.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/a-h/templ"
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
	if os.Getenv("DEBUG") != "" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	var (
		allPosts       []*assets.Post
		allProjects    []*assets.Project
		allTags        []*assets.Tag
		allEmployments []*assets.Employment
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
	err = db.NewSelect().
		Model(&allEmployments).
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Relation("Employments").
		Scan(ctx)
	if err != nil {
		return eris.Wrapf(err, "(update-css) failed to get employments")
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

	return twerge.CodeGen(
		twerge.Default(),
		"cmd/conneroh/classes/classes.go",
		"input.css",
		"cmd/conneroh/classes/classes.html",
		comps...,
	)
}
