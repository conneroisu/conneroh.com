package conneroh

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/routing"
)

var (
	posts = views.List(
		routing.PluralTargetPost,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
		[]string{},
		[]string{},
		[]string{},
	)
	home = views.Home(
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
	)
	projects = views.List(
		routing.PluralTargetProject,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
		[]string{},
		[]string{},
		[]string{},
	)
	tags = views.List(
		routing.PluralTargetTag,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
		[]string{},
		[]string{},
		[]string{},
	)
)

func postsHandler(
	renderFn func(templ.Component) templ.Component,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(renderFn(posts)).ServeHTTP(w, r)
	}
}

func projectsHandler(
	renderFn func(templ.Component) templ.Component,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(renderFn(projects)).ServeHTTP(w, r)
	}
}

func tagsHandler(
	renderFn func(templ.Component) templ.Component,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(renderFn(tags)).ServeHTTP(w, r)
	}
}
