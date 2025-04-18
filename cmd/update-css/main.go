// Package main updates the CSS.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/conneroisu/conneroh.com/cmd/conneroh/components"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/conneroisu/twerge"
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
	if err := twerge.GenCSS(
		"internal/data/css/classes.go",
		"input.css",
		"internal/data/css/classes.html",
		layouts.Page(views.Home(
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)),
		layouts.Page(views.List(
			routing.PostPluralPath,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
			"",
			1,
			10,
		)),
		layouts.Page(views.List(
			routing.ProjectPluralPath,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
			"",
			1,
			10,
		)),
		layouts.Page(views.List(
			routing.TagsPluralPath,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
			"",
			1,
			10,
		)),
		components.TagControl(
			&assets.Tag{},
		),
		layouts.Page(views.Post(
			gen.AllPosts[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)),
		layouts.Page(views.Project(
			gen.AllProjects[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)),
		layouts.Page(views.Tag(
			gen.AllTags[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)),
		views.Post(
			gen.AllPosts[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		),
		layouts.Layout("hello"),
		components.ThankYou(),
	); err != nil {
		panic(err)
	}
}
