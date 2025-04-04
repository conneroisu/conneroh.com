package conneroh

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/routing"
)

var posts = views.List(
	routing.PluralTargetPost,
	&gen.AllPosts,
	&gen.AllProjects,
	&gen.AllTags,
	[]string{},
	[]string{},
	[]string{},
)

func postsHandler(
	renderFn func(templ.Component) templ.Component,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(renderFn(posts)).ServeHTTP(w, r)
	}
}
