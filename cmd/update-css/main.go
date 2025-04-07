// Package main updates the CSS.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/conneroisu/twerge"
)

var cwd = flag.String("cwd", "", "current working directory")

func genCSS(ctx context.Context) error {
	var (
		_ = layouts.Page(views.Home(
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)).Render(ctx, io.Discard)
		_ = layouts.Page(views.List(
			routing.PluralTargetPost,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
			"",
		)).Render(ctx, io.Discard)
		_ = layouts.Page(views.List(
			routing.PluralTargetProject,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
			"",
		)).Render(ctx, io.Discard)
		_ = layouts.Page(views.List(
			routing.PluralTargetTag,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
			"",
		)).Render(ctx, io.Discard)
		_ = views.TagControl(
			&gen.Tag{},
			"#list-project",
		).Render(ctx, io.Discard)
		_ = layouts.Page(views.Post(
			gen.AllPosts[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)).Render(ctx, io.Discard)
		_ = layouts.Page(views.Project(
			gen.AllProjects[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)).Render(ctx, io.Discard)
		_ = layouts.Page(views.Tag(
			gen.AllTags[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)).Render(ctx, io.Discard)
		_ = views.TagControl(
			&gen.Tag{},
			"#list-project",
		).Render(ctx, io.Discard)
		_ = layouts.Morpher(views.Post(
			gen.AllPosts[0],
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		)).Render(ctx, io.Discard)
		_ = layouts.Layout("hello").Render(ctx, io.Discard)
	)
	content := twerge.GenerateClassMapCode("css")
	f, err := os.Create("internal/data/css/classes.go")
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	err = twerge.GenerateTailwind("input.css")
	if err != nil {
		return err
	}
	err = twerge.GenerateTempl("internal/data/css/classes.templ")
	if err != nil {
		return err
	}
	println("Generated classes.go.")
	return nil
}

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
	if err := genCSS(context.Background()); err != nil {
		panic(err)
	}
}
