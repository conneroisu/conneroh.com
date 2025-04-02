package main

import (
	"context"
	"io"
	"os"

	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/conneroisu/twerge"
)

func genCSS(ctx context.Context) error {
	var (
		_ = views.Home(
			&[]*gen.Post{},
			&[]*gen.Project{},
			&[]*gen.Tag{},
		).Render(ctx, io.Discard)
		_ = views.List(
			routing.PluralTargetPost,
			&[]*gen.Post{},
			&[]*gen.Project{},
			&[]*gen.Tag{},
		).Render(ctx, io.Discard)
		_ = views.List(
			routing.PluralTargetProject,
			&[]*gen.Post{},
			&[]*gen.Project{},
			&[]*gen.Tag{},
		).Render(ctx, io.Discard)
		_ = views.List(
			routing.PluralTargetTag,
			&[]*gen.Post{},
			&[]*gen.Project{},
			&[]*gen.Tag{},
		).Render(ctx, io.Discard)
	)
	content := twerge.GenerateClassMapCode("css")
	f, err := os.Create("internal/data/css/classes.go")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	err = twerge.GenerateTailwind("input.css", "input.css", twerge.ClassMapStr)
	if err != nil {
		return err
	}
	println("Generated classes.go.")
	return nil
}

func main() {
	if err := genCSS(context.Background()); err != nil {
		panic(err)
	}
}
