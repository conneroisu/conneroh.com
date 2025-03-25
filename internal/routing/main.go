package routing

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/master"
)

// FullFn is a function that handles a full request.
type FullFn func(
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostsSlugMap *map[string]master.FullPost,
	fullProjectsSlugMap *map[string]master.FullProject,
	fullTagsSlugMap *map[string]master.FullTag,
) templ.Component

// SingleFn returns a fullFn for the single view.
type SingleFn func(
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostSlugMap *map[string]master.FullPost,
	fullProjectSlugMap *map[string]master.FullProject,
	fullTagSlugMap *map[string]master.FullTag,
) func(target SingleTarget, id string) templ.Component

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

// SingleTarget is the target of a single view.
// string
type SingleTarget string

const (
	// SingleTargetPost is the target of a single post view.
	SingleTargetPost SingleTarget = "post"
	// SingleTargetProject is the target of a single project view.
	SingleTargetProject SingleTarget = "project"
	// SingleTargetTag is the target of a single tag view.
	SingleTargetTag SingleTarget = "tag"
)

// PluralTarget is the target of a plural view.
// string
type PluralTarget string

const (
	// PluralTargetPost is the target of a plural post view.
	PluralTargetPost PluralTarget = "posts"
	// PluralTargetProject is the target of a plural project view.
	PluralTargetProject PluralTarget = "projects"
	// PluralTargetTag is the target of a plural tag view.
	PluralTargetTag PluralTarget = "tags"
)
