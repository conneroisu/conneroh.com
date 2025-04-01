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

var (
	home = views.Home(
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
	)
	posts = views.List(
		routing.PluralTargetPost,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
	)
	projects = views.List(
		routing.PluralTargetProject,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
	)
	tags = views.List(
		routing.PluralTargetTag,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
	)
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
		templ.Handler(layouts.Morpher(home)),
	)
	h.Handle(
		"GET /dist/",
		http.FileServer(http.FS(static.Dist)),
	)
	h.Handle(
		"GET /favicon.ico",
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write(static.Favicon)
			if err != nil {
				return
			}
		}),
	)

	h.Handle(
		"GET /posts",
		templ.Handler(layouts.Page(posts)))
	h.Handle(
		"GET /morph/posts",
		templ.Handler(layouts.Morpher(posts)))
	h.Handle(
		"GET /projects",
		templ.Handler(layouts.Page(projects)))
	h.Handle(
		"GET /morph/projects",
		templ.Handler(layouts.Morpher(projects)))
	h.Handle(
		"GET /tags",
		templ.Handler(layouts.Page(tags)))
	h.Handle(
		"GET /morph/tags",
		templ.Handler(layouts.Morpher(tags)))

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
			"GET /morphs/post/%s",
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
			fmt.Sprintf("GET /morphs/project/%s", p.Slug),
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
			fmt.Sprintf("GET /morphs/tag/%s", t.Slug),
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
