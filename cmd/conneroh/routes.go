package conneroh

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
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
		templ.Handler(layouts.Page(home)),
	)
	h.Handle(
		"GET /morph/home",
		templ.Handler(layouts.Morpher(home)))
	h.Handle(
		"GET /dist/",
		http.FileServer(http.FS(static.Dist)))
	h.HandleFunc(
		"GET /favicon.ico",
		func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write(static.Favicon)
		},
	)

	h.Handle(
		"GET /posts",
		templ.Handler(layouts.Page(posts)))
	h.Handle(
		"GET /morph/posts",
		templ.Handler(layouts.Morpher(posts)))
	h.Handle(
		"GET /search/posts",
		listHandler(routing.PostPluralPath))

	h.Handle(
		"GET /projects",
		templ.Handler(layouts.Page(projects)))
	h.Handle(
		"GET /morph/projects",
		templ.Handler(layouts.Morpher(projects)))
	h.Handle(
		"GET /search/projects",
		listHandler(routing.ProjectPluralPath))

	h.Handle(
		"GET /tags",
		templ.Handler(layouts.Page(tags)))
	h.Handle(
		"GET /morph/tags",
		templ.Handler(layouts.Morpher(tags)))
	h.Handle(
		"GET /search/tags",
		listHandler(routing.TagsPluralPath))

	for _, p := range gen.AllPosts {
		h.Handle(
			fmt.Sprintf("GET /post/%s", p.Slug),
			templ.Handler(layouts.Page(views.Post(
				p,
				&gen.AllPosts,
				&gen.AllProjects,
				&gen.AllTags,
			))),
		)
		h.Handle(fmt.Sprintf(
			"GET /morph/post/%s",
			p.Slug,
		), templ.Handler(layouts.Morpher(views.Post(
			p,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		))))
	}
	for _, p := range gen.AllProjects {
		h.Handle(
			fmt.Sprintf("GET /project/%s", p.Slug),
			templ.Handler(layouts.Page(views.Project(
				p,
				&gen.AllPosts,
				&gen.AllProjects,
				&gen.AllTags,
			))),
		)
		h.Handle(
			fmt.Sprintf("GET /morph/project/%s", p.Slug),
			templ.Handler(layouts.Morpher(views.Project(
				p,
				&gen.AllPosts,
				&gen.AllProjects,
				&gen.AllTags,
			))),
		)
	}
	for _, t := range gen.AllTags {
		h.Handle(
			fmt.Sprintf("GET /tag/%s", t.Slug),
			templ.Handler(layouts.Page(views.Tag(
				t,
				&gen.AllPosts,
				&gen.AllProjects,
				&gen.AllTags,
			))),
		)
		h.Handle(
			fmt.Sprintf("GET /morph/tag/%s", t.Slug),
			templ.Handler(layouts.Morpher(views.Tag(
				t,
				&gen.AllPosts,
				&gen.AllProjects,
				&gen.AllTags,
			))),
		)
	}

	return nil
}
