package conneroh

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	static "github.com/conneroisu/conneroh.com/cmd/conneroh/_static"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/routing"
)

// AddRoutes adds all routes to the router.
func AddRoutes(
	_ context.Context,
	h *http.ServeMux,
) error {
	slog.Info("adding routes")
	defer slog.Info("added routes")

	h.Handle(
		"/{$}",
		routing.MorphableHandler(
			layouts.Page(home),
			home,
		),
	)
	h.Handle(
		"GET /dist/",
		http.FileServer(http.FS(static.Dist)))
	h.HandleFunc(
		"GET /favicon.ico",
		func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write(static.Favicon)
		})
	h.HandleFunc(
		"GET /search/all",
		globalSearchHandler(gen.AllPosts, gen.AllProjects, gen.AllTags))

	h.Handle(
		"GET /posts",
		routing.MorphableHandler(layouts.Page(posts), posts))
	h.Handle(
		"GET /search/posts",
		listHandler(routing.PostPluralPath))

	h.Handle(
		"GET /projects",
		routing.MorphableHandler(layouts.Page(projects), projects))
	h.Handle(
		"GET /search/projects",
		listHandler(routing.ProjectPluralPath))

	h.Handle(
		"GET /tags",
		routing.MorphableHandler(layouts.Page(tags), tags))
	h.Handle(
		"GET /search/tags",
		listHandler(routing.TagsPluralPath))

	for _, p := range gen.AllPosts {
		h.Handle(
			fmt.Sprintf("GET /post/%s", p.Slug),
			routing.MorphableHandler(
				layouts.Page(views.Post(
					p,
					&gen.AllPosts,
					&gen.AllProjects,
					&gen.AllTags,
				)),
				views.Post(
					p,
					&gen.AllPosts,
					&gen.AllProjects,
					&gen.AllTags,
				),
			),
		)
	}
	for _, p := range gen.AllProjects {
		h.Handle(
			fmt.Sprintf("GET /project/%s", p.Slug),
			routing.MorphableHandler(
				layouts.Page(views.Project(
					p,
					&gen.AllPosts,
					&gen.AllProjects,
					&gen.AllTags,
				)),
				views.Project(
					p,
					&gen.AllPosts,
					&gen.AllProjects,
					&gen.AllTags,
				),
			),
		)
	}
	for _, t := range gen.AllTags {
		h.Handle(
			fmt.Sprintf("GET /tag/%s", t.Slug),
			routing.MorphableHandler(
				layouts.Page(views.Tag(
					t,
					&gen.AllPosts,
					&gen.AllProjects,
					&gen.AllTags,
				)),
				views.Tag(
					t,
					&gen.AllPosts,
					&gen.AllProjects,
					&gen.AllTags,
				),
			),
		)
	}

	return nil
}
