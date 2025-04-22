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
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	_ "modernc.org/sqlite"
)

var cwd = flag.String("cwd", "", "current working directory")

func main() {
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
	sqlDB, err := sql.Open("sqlite", "file:test.db?cache=shared&mode=rwc")
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqlDB, sqlitedialect.New())

	var (
		allPosts    []*assets.Post
		allProjects []*assets.Project
		allTags     []*assets.Tag
	)
	_, err = db.NewSelect().Model(assets.EmpPost).Exec(context.Background(), &allPosts)
	if err != nil {
		panic(err)
	}
	_, err = db.NewSelect().Model(assets.EmpProject).Exec(context.Background(), &allProjects)
	if err != nil {
		panic(err)
	}
	_, err = db.NewSelect().Model(assets.EmpTag).Exec(context.Background(), &allTags)
	if err != nil {
		panic(err)
	}
	if err := twerge.CodeGen(
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
	); err != nil {
		panic(err)
	}
}
