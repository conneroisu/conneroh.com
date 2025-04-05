package main

import (
	"context"
	"fmt"
	"io"

	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/conneroisu/twerge"
	"github.com/rotisserie/eris"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	//
	// for k, v := range css.ClassMapStr {
	// 	for k2, v2 := range css.ClassMapStr {
	// 		if v == v2 {
	// 			return eris.Errorf(`
	// 				'%s' : %s
	// 				'%s' : %s
	// 				Duplicate class name
	// 			`, k, k2, v, v2)
	// 		}
	// 	}
	// }

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
			nil,
			nil,
			nil,
		)).Render(ctx, io.Discard)
		_ = layouts.Page(views.List(
			routing.PluralTargetProject,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
			nil,
			nil,
			nil,
		)).Render(ctx, io.Discard)
		_ = layouts.Page(views.List(
			routing.PluralTargetTag,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
			nil,
			nil,
			nil,
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

	for gen, genMerged := range twerge.GenClassMergeStr {
		fmt.Printf("%s : %s\n", gen, genMerged)
	}
	for orig, class := range twerge.ClassMapStr {
		fmt.Printf("%s : %s\n", orig, class)
	}

	for gen, genMerged := range twerge.GenClassMergeStr {
		for orig, class := range twerge.ClassMapStr {
			if gen == class {
				if genMerged != orig {
					return eris.Errorf("\n\n '%s' has been merged to '%s'", orig, genMerged)
				}
			}
		}
	}
	return nil
}
