package routing

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
)

// SingleFn returns a fullFn for the single view.
type SingleFn func(target SingleTarget, id string) templ.Component

// APIFn is a function that handles an API request.
type APIFn func(http.ResponseWriter, *http.Request) error

// APIHandler is a function that returns an APIFn.
type APIHandler func(
	fullPosts *[]gen.Post,
	fullProjects *[]gen.Project,
	fullTags *[]gen.Tag,
	postsSlugMap *map[string]gen.Post,
	projectsSlugMap *map[string]gen.Project,
	tagsSlugMap *map[string]gen.Tag,
) (APIFn, error)

// APIMap is a map of API functions.
type APIMap map[string]APIHandler

// AddRoutes adds all routes to the router.
func (m APIMap) AddRoutes(
	mux *http.ServeMux,
	fullPosts *[]gen.Post,
	fullProjects *[]gen.Project,
	fullTags *[]gen.Tag,
	postsSlugMap *map[string]gen.Post,
	projectsSlugMap *map[string]gen.Project,
	tagsSlugMap *map[string]gen.Tag,
) error {
	for path, fn := range m {
		h, err := fn(
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
