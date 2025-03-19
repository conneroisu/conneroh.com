package routing

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/master"
)

// APIFn is a function that handles an API request.
type APIFn func(http.ResponseWriter, *http.Request) error

// APIHandler is a function that returns an APIFn.
type APIHandler func(
	ctx context.Context,
	db *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	postsSlugMap *map[string]master.FullPost,
	projectsSlugMap *map[string]master.FullProject,
	tagsSlugMap *map[string]master.FullTag,
) (APIFn, error)

// APIMap is a map of API functions.
type APIMap map[string]APIHandler

// AddRoutes adds all routes to the router.
func (m APIMap) AddRoutes(
	ctx context.Context,
	mux *http.ServeMux,
	db *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	postsSlugMap *map[string]master.FullPost,
	projectsSlugMap *map[string]master.FullProject,
	tagsSlugMap *map[string]master.FullTag,
) error {
	for path, fn := range m {
		h, err := fn(
			ctx,
			db,
			fullPosts,
			fullProjects,
			fullTags,
			postsSlugMap,
			projectsSlugMap,
			tagsSlugMap,
		)
		if err != nil {
			return err
		}
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if err := h(w, r); err != nil {
				slog.Error("error handling request", "err", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}

	return nil
}
