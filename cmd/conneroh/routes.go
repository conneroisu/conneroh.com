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
	h.Handle(
		"GET /favicon.ico",
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write(static.Favicon)
		}),
	)

	h.Handle(
		"GET /posts",
		postsHandler(layouts.Page))
	h.Handle(
		"GET /morph/posts",
		postsHandler(layouts.Morpher))
	h.Handle(
		"GET /projects",
		projectsHandler(layouts.Page))
	h.Handle(
		"GET /morph/projects",
		projectsHandler(layouts.Morpher))
	h.Handle(
		"GET /tags",
		tagsHandler(layouts.Page))
	h.Handle(
		"GET /morph/tags",
		tagsHandler(layouts.Morpher))

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
