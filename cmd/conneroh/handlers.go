package conneroh

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	static "github.com/conneroisu/conneroh.com/cmd/conneroh/_static"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/master"
	"github.com/conneroisu/conneroh.com/internal/routing"
)

const tagsParamContextKey contextKey = "tagsParam"
const currentURLContextKey contextKey = "currentURL"

type (
	fullFn func(
		fullPosts *[]master.FullPost,
		fullProjects *[]master.FullProject,
		fullTags *[]master.FullTag,
		fullPostsSlugMap *map[string]master.FullPost,
		fullProjectsSlugMap *map[string]master.FullProject,
		fullTagsSlugMap *map[string]master.FullTag,
	) templ.Component

	// Context keys for passing data to templates
	contextKey string
)

// Dist is the dist handler for serving/distributing static files.
func Dist(
	_ context.Context,
	_ *data.Database[master.Queries],
	_ *[]master.FullPost,
	_ *[]master.FullProject,
	_ *[]master.FullTag,
	_ *map[string]master.FullPost,
	_ *map[string]master.FullProject,
	_ *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		http.FileServer(http.FS(static.Dist)).ServeHTTP(w, r)
		return nil
	}, nil
}

// Home is the home page handler.
func Home(
	ctx context.Context,
	db *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostsSlugMap *map[string]master.FullPost,
	fullProjectsSlugMap *map[string]master.FullProject,
	fullTagsSlugMap *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		templ.Handler(views.Page(views.Home(
			fullPosts,
			fullProjects,
			fullTags,
			fullPostsSlugMap,
			fullProjectsSlugMap,
			fullTagsSlugMap,
		))).ServeHTTP(w, r)
		return nil
	}, nil

}

// Morph renders a morphed view.
func Morph(
	ctx context.Context,
	db *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostSlugMap *map[string]master.FullPost,
	fullProjectSlugMap *map[string]master.FullProject,
	fullTagSlugMap *map[string]master.FullTag,
) (routing.APIFn, error) {
	var morphMap = map[string]fullFn{
		"projects": views.Projects,
		"posts":    views.Posts,
		"tags":     views.Tags,
		"home":     views.Home,
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		view := r.PathValue("view")
		val, ok := morphMap[view]
		if !ok {
			return fmt.Errorf("unknown view: %s", view)
		}
		morphed := views.Morpher(val(
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		))
		err := morphed.Render(r.Context(), w)
		if err != nil {
			return err
		}
		return nil
	}, nil
}

// Morphs renders a morphed view.
func Morphs(
	_ context.Context,
	_ *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostSlugMap *map[string]master.FullPost,
	fullProjectSlugMap *map[string]master.FullProject,
	fullTagSlugMap *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			view = r.PathValue("view")
			id   = r.PathValue("id")
		)
		switch view {
		case "project":
			proj, ok := (*fullProjectSlugMap)[id]
			if !ok {
				return routing.ErrNotFound{URL: r.URL}
			}
			morphed := views.Morpher(views.Project(
				&proj,
				fullPosts,
				fullProjects,
				fullTags,
				fullPostSlugMap,
				fullProjectSlugMap,
				fullTagSlugMap,
			))
			err := morphed.Render(r.Context(), w)
			if err != nil {
				return err
			}
		case "post":
			post, ok := (*fullPostSlugMap)[id]
			if !ok {
				return routing.ErrNotFound{URL: r.URL}
			}
			morphed := views.Morpher(views.Post(
				&post,
				fullPosts,
				fullProjects,
				fullTags,
				fullPostSlugMap,
				fullProjectSlugMap,
				fullTagSlugMap,
			))
			err := morphed.Render(r.Context(), w)
			if err != nil {
				return err
			}
		case "tag":
			tag, ok := (*fullTagSlugMap)[id]
			if !ok {
				return routing.ErrNotFound{URL: r.URL}
			}
			morphed := views.Morpher(views.Tag(
				&tag,
				fullPosts,
				fullProjects,
				fullTags,
				fullPostSlugMap,
				fullProjectSlugMap,
				fullTagSlugMap,
			))
			err := morphed.Render(r.Context(), w)
			if err != nil {
				return err
			}
		default:
			return routing.ErrNotFound{URL: r.URL}
		}
		return nil
	}, nil
}
