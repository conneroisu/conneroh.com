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
	fullPostsSlugMap map[string]master.FullPost,
	fullProjectsSlugMap map[string]master.FullProject,
	fullTagsSlugMap map[string]master.FullTag,
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
	}, nil

}

// Projects is the projects handler.
func Projects(
	ctx context.Context,
	db *data.Database[master.Queries],
	_ *[]master.FullPost,
	fullProjects *[]master.FullProject,
	_ *[]master.FullTag,
	_ *map[string]master.FullPost,
	_ *map[string]master.FullProject,
	_ *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		templ.Handler(
			views.Page(views.Projects(fullProjects)),
		).ServeHTTP(w, r)
		return nil
	}, nil
}

// Project is the project handler.
func Project(
	ctx context.Context,
	db *data.Database[master.Queries],
	_ *[]master.FullPost,
	_ *[]master.FullProject,
	_ *[]master.FullTag,
	_ *map[string]master.FullPost,
	projects *map[string]master.FullProject,
	_ *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := r.PathValue("id")
		if id == "" {
			return routing.ErrMissingParam{}
		}
		proj, ok := (*projects)[id]
		if !ok {
			return routing.ErrNotFound{URL: r.URL}
		}
		templ.Handler(
			views.Page(views.Project(&proj)),
		).ServeHTTP(w, r)
		return nil
	}, nil
}

// Posts is the posts handler.
func Posts(
	_ context.Context,
	_ *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	_ *[]master.FullProject,
	_ *[]master.FullTag,
	_ *map[string]master.FullPost,
	_ *map[string]master.FullProject,
	_ *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		templ.Handler(
			views.Page(views.Posts(fullPosts)),
		).ServeHTTP(w, r)
		return nil
	}, nil
}

// Post is the post handler.
func Post(
	_ context.Context,
	_ *data.Database[master.Queries],
	_ *[]master.FullPost,
	_ *[]master.FullProject,
	_ *[]master.FullTag,
	posts *map[string]master.FullPost,
	_ *map[string]master.FullProject,
	_ *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := r.PathValue("id")
		if id == "" {
			return routing.ErrMissingParam{ID: id, View: "Post"}
		}
		post, ok := (*posts)[id]
		if !ok {
			return routing.ErrNotFound{URL: r.URL}
		}
		templ.Handler(
			views.Page(views.Post(&post)),
		).ServeHTTP(w, r)
		return nil
	}, nil
}

// Tags is the tags handler.
func Tags(
	_ context.Context,
	_ *data.Database[master.Queries],
	_ *[]master.FullPost,
	_ *[]master.FullProject,
	fullTags *[]master.FullTag,
	_ *map[string]master.FullPost,
	_ *map[string]master.FullProject,
	_ *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		templ.Handler(
			views.Page(views.Tags(fullTags)),
		).ServeHTTP(w, r)
		return nil
	}, nil
}

// Tag is the tag handler.
func Tag(
	_ context.Context,
	_ *data.Database[master.Queries],
	_ *[]master.FullPost,
	_ *[]master.FullProject,
	_ *[]master.FullTag,
	_ *map[string]master.FullPost,
	_ *map[string]master.FullProject,
	tags *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := r.PathValue("id")
		if id == "" {
			return routing.ErrMissingParam{ID: id, View: "Tag"}
		}
		tag, ok := (*tags)[id]
		if !ok {
			return routing.ErrNotFound{URL: r.URL}
		}
		templ.Handler(
			views.Page(views.Tag(&tag)),
		).ServeHTTP(w, r)
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
	_ *map[string]master.FullPost,
	_ *map[string]master.FullProject,
	_ *map[string]master.FullTag,
) (routing.APIFn, error) {
	var morphMap = map[string]templ.Component{
		"projects": views.Projects(fullProjects),
		"posts":    views.Posts(fullPosts),
		"tags":     views.Tags(fullTags),
		"home":     views.Home(),
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		view := r.PathValue("view")
		val, ok := morphMap[view]
		if !ok {
			return fmt.Errorf("unknown view: %s", view)
		}
		morphed := views.Morpher(val)
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
	_ *[]master.FullPost,
	_ *[]master.FullProject,
	_ *[]master.FullTag,
	posts *map[string]master.FullPost,
	projects *map[string]master.FullProject,
	tags *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			view = r.PathValue("view")
			id   = r.PathValue("id")
		)
		switch view {
		case "project":
			proj, ok := (*projects)[id]
			if !ok {
				return routing.ErrNotFound{URL: r.URL}
			}
			morphed := views.Morpher(views.Project(&proj))
			err := morphed.Render(r.Context(), w)
			if err != nil {
				return err
			}
		case "post":
			post, ok := (*posts)[id]
			if !ok {
				return routing.ErrNotFound{URL: r.URL}
			}
			morphed := views.Morpher(views.Post(&post))
			err := morphed.Render(r.Context(), w)
			if err != nil {
				return err
			}
		case "tag":
			tag, ok := (*tags)[id]
			if !ok {
				return routing.ErrNotFound{URL: r.URL}
			}
			morphed := views.Morpher(views.Tag(&tag))
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
