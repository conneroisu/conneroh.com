package conneroh

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	static "github.com/conneroisu/conneroh.com/cmd/conneroh/_static"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/components"
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
		templ.Handler(layouts.Page(views.Home(
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		))),
	)
	h.Handle(
		"GET /hateoas/morph/home",
		templ.Handler(components.Morpher(views.Home(
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		))),
	)
	h.Handle("GET /dist/", http.FileServer(http.FS(static.Dist)))

	h.Handle(
		"GET /posts",
		templ.Handler(layouts.Page(views.List(
			routing.PluralTargetPost,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		))))
	h.Handle(
		"GET /hateoas/morph/posts",
		templ.Handler(components.Morpher(views.List(
			routing.PluralTargetPost,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		))))
	h.Handle(
		"GET /projects",
		templ.Handler(layouts.Page(views.List(
			routing.PluralTargetProject,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		))))
	h.Handle(
		"GET /hateoas/morph/projects",
		templ.Handler(components.Morpher(views.List(
			routing.PluralTargetProject,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		))),
	)
	h.Handle(
		"GET /tags",
		templ.Handler(layouts.Page(views.List(
			routing.PluralTargetTag,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		))),
	)
	h.Handle(
		"GET /hateoas/morph/tags",
		templ.Handler(components.Morpher(views.List(
			routing.PluralTargetTag,
			&gen.AllPosts,
			&gen.AllProjects,
			&gen.AllTags,
		))),
	)

	for _, p := range gen.AllPosts {
		h.Handle(
			fmt.Sprintf("GET /post/%s", p.Slug),
			templ.Handler(layouts.Page(views.Post(p, &gen.AllPosts, &gen.AllProjects, &gen.AllTags))),
		)
		h.Handle(fmt.Sprintf(
			"GET /hateoas/morphs/post/%s",
			p.Slug,
		), templ.Handler(components.Morpher(views.Post(p, &gen.AllPosts, &gen.AllProjects, &gen.AllTags))))
	}
	for _, p := range gen.AllProjects {
		h.Handle(
			fmt.Sprintf("GET /project/%s", p.Slug),
			templ.Handler(layouts.Page(views.Project(p, &gen.AllPosts, &gen.AllProjects, &gen.AllTags))),
		)
		h.Handle(
			fmt.Sprintf("GET /hateoas/morphs/project/%s", p.Slug),
			templ.Handler(components.Morpher(views.Project(p, &gen.AllPosts, &gen.AllProjects, &gen.AllTags))),
		)
	}
	for _, t := range gen.AllTags {
		h.Handle(
			fmt.Sprintf("GET /tag/%s", t.Slug),
			templ.Handler(layouts.Page(views.Tag(t, &gen.AllPosts, &gen.AllProjects, &gen.AllTags))),
		)
		h.Handle(
			fmt.Sprintf("GET /hateoas/morphs/tag/%s", t.Slug),
			templ.Handler(components.Morpher(views.Tag(t, &gen.AllPosts, &gen.AllProjects, &gen.AllTags))),
		)
	}

	return nil
}
