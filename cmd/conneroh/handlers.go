package conneroh

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
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
	allPosts *[]*gen.Post,
	allProjects *[]*gen.Project,
	allTags *[]*gen.Tag,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(layouts.Page(posts)).ServeHTTP(w, r)
	}
}
