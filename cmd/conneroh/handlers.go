package conneroh

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	static "github.com/conneroisu/conneroh.com/cmd/conneroh/_static"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/components"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
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

// Favicon is the favicon handler.
func Favicon(
	_ context.Context,
	_ *data.Database[master.Queries],
	_ *[]master.FullPost,
	_ *[]master.FullProject,
	_ *[]master.FullTag,
	_ *map[string]master.FullPost,
	_ *map[string]master.FullProject,
	_ *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, _ *http.Request) error {
		w.Header().Set("Content-Type", "image/x-icon")
		_, err := w.Write(static.Favicon)
		if err != nil {
			return err
		}
		return nil
	}, nil
}

// Home is the home page handler.
func Home(
	_ context.Context,
	_ *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostsSlugMap *map[string]master.FullPost,
	fullProjectsSlugMap *map[string]master.FullProject,
	fullTagsSlugMap *map[string]master.FullTag,
) (routing.APIFn, error) {
	return func(w http.ResponseWriter, r *http.Request) error {
		templ.Handler(layouts.Page(views.Home(
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

// ListMorph renders a morphed view.
func ListMorph(
	_ context.Context,
	_ *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostSlugMap *map[string]master.FullPost,
	fullProjectSlugMap *map[string]master.FullProject,
	fullTagSlugMap *map[string]master.FullTag,
) (routing.APIFn, error) {
	var morphMap = map[string]templ.Component{
		routing.PluralTargetProject: views.List(
			views.ListTargetProjects,
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		),
		routing.PluralTargetPost: views.List(
			views.ListTargetPosts,
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		),
		routing.PluralTargetTag: views.List(
			views.ListTargetTags,
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		),
		"home": views.Home(
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		),
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		view := r.PathValue("view")
		val, ok := morphMap[view]
		if !ok {
			return fmt.Errorf("unknown view: %s", view)
		}
		morphed := components.Morpher(val)
		err := morphed.Render(r.Context(), w)
		if err != nil {
			return err
		}
		return nil
	}, nil
}

// SingleMorph renders a morphed view.
func SingleMorph(
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
			templ.Handler(
				components.Morpher(views.Project(
					proj,
					fullPosts,
					fullProjects,
					fullTags,
					fullPostSlugMap,
					fullProjectSlugMap,
					fullTagSlugMap,
				))).ServeHTTP(w, r)
		case "post":
			post, ok := (*fullPostSlugMap)[id]
			if !ok {
				return routing.ErrNotFound{URL: r.URL}
			}
			templ.Handler(components.Morpher(views.Post(
				post,
				fullPosts,
				fullProjects,
				fullTags,
				fullPostSlugMap,
				fullProjectSlugMap,
				fullTagSlugMap,
			))).ServeHTTP(w, r)
		case "tag":
			tag, ok := (*fullTagSlugMap)[id]
			if !ok {
				return routing.ErrNotFound{URL: r.URL}
			}
			morphed := components.Morpher(views.Tag(
				tag,
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

// List handles the GET /list/{targets} endpoint.
func List(
	_ context.Context,
	_ *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostSlugMap *map[string]master.FullPost,
	fullProjectSlugMap *map[string]master.FullProject,
	fullTagSlugMap *map[string]master.FullTag,
) (routing.APIFn, error) {
	targetMap := map[views.ListTarget]templ.Component{
		views.ListTargetPosts: views.List(
			views.ListTargetPosts,
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		),
		views.ListTargetProjects: views.List(
			views.ListTargetProjects,
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		),
		views.ListTargetTags: views.List(
			views.ListTargetTags,
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		),
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		target := r.PathValue("targets")
		if target == "" {
			return routing.ErrMissingParam{ID: target, View: "List"}
		}
		comp, ok := targetMap[target]
		if !ok {
			return routing.ErrNotFound{URL: r.URL}
		}
		templ.Handler(layouts.Page(comp)).ServeHTTP(w, r)
		return nil
	}, nil
}

// Single handles the GET /{target}/{id...} endpoint.
func Single(
	_ context.Context,
	_ *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostSlugMap *map[string]master.FullPost,
	fullProjectSlugMap *map[string]master.FullProject,
	fullTagSlugMap *map[string]master.FullTag,
) (routing.APIFn, error) {
	singleMap := map[routing.SingleTarget]routing.SingleFn{
		routing.SingleTargetPost: views.Single(
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		),
		routing.SingleTargetProject: views.Single(
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		),
		routing.SingleTargetTag: views.Single(
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		),
	}
	return func(w http.ResponseWriter, r *http.Request) error {
		target := r.PathValue("target")
		if target == "" {
			return routing.ErrMissingParam{ID: target, View: "Single"}
		}
		id := r.PathValue("id")
		if id == "" {
			return routing.ErrMissingParam{ID: id, View: fmt.Sprintf("Single %s", target)}
		}
		comp, ok := singleMap[target]
		if !ok {
			return routing.ErrNotFound{URL: r.URL}
		}
		templ.Handler(layouts.Page(comp(target, id))).ServeHTTP(w, r)
		return nil
	}, nil
}
