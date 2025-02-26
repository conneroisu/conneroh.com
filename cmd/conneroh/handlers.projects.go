package conneroh

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/master"
	"github.com/conneroisu/conneroh.com/internal/routing"
)

// Projects is the projects handler.
func Projects(
	ctx context.Context,
	db *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostSlugMap *map[string]master.FullPost,
	fullProjectSlugMap *map[string]master.FullProject,
	fullTagSlugMap *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		templ.Handler(
			views.Page(views.Projects(
				fullPosts,
				fullProjects,
				fullTags,
				fullPostSlugMap,
				fullProjectSlugMap,
				fullTagSlugMap,
			)),
		).ServeHTTP(w, r)
		return nil
	}, nil
}

// Project is the project handler.
func Project(
	ctx context.Context,
	db *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostSlugMap *map[string]master.FullPost,
	fullProjectSlugMap *map[string]master.FullProject,
	fullTagSlugMap *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := r.PathValue("id")
		if id == "" {
			return routing.ErrMissingParam{}
		}
		proj, ok := (*fullProjectSlugMap)[id]
		if !ok {
			return routing.ErrNotFound{URL: r.URL}
		}
		templ.Handler(
			views.Page(views.Project(&proj,
				fullPosts,
				fullProjects,
				fullTags,
				fullPostSlugMap,
				fullProjectSlugMap,
				fullTagSlugMap,
			)),
		).ServeHTTP(w, r)
		return nil
	}, nil
}
