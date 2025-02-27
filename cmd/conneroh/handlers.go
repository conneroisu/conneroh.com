package conneroh

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

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
	return func(w http.ResponseWriter, r *http.Request) error {
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

// // PluralHandler handles a singular resource (e.g. project, post, tag, home)
// func PluralHandler[
// 	T master.FullPost | master.FullProject | master.FullTag,
// ](
// 	ctx context.Context,
// 	db *data.Database[master.Queries],
// 	fullPosts *[]master.FullPost,
// 	fullProjects *[]master.FullProject,
// 	fullTags *[]master.FullTag,
// 	fullPostSlugMap *map[string]master.FullPost,
// 	fullProjectSlugMap *map[string]master.FullProject,
// 	fullTagSlugMap *map[string]master.FullTag,
// ) (routing.APIFn, error) {
// 	return func(w http.ResponseWriter, r *http.Request) error {
// 		return nil
// 	}, nil
// }
//
// // SingularHandler handles a singular resource (e.g. project, post, tag)
// func SingularHandler[
// 	T master.FullPost | master.FullProject | master.FullTag,
// ](
// 	ctx context.Context,
// 	db *data.Database[master.Queries],
// 	singular T,
// 	fullPosts *[]master.FullPost,
// 	fullProjects *[]master.FullProject,
// 	fullTags *[]master.FullTag,
// 	fullPostSlugMap *map[string]master.FullPost,
// 	fullProjectSlugMap *map[string]master.FullProject,
// 	fullTagSlugMap *map[string]master.FullTag,
// ) (routing.APIFn, error) {
// 	return func(w http.ResponseWriter, r *http.Request) error {
// 		return nil
// 	}, nil
// }

// Posts is the posts handler.
func Posts(
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
		// Get tag filters from the URL query
		tagsParam := r.URL.Query().Get("tags")

		// Parse +tag and -tag filters
		includeTags, excludeTags := parseTagFilters(tagsParam)

		// Apply filtering
		filteredPosts := *fullPosts
		if len(includeTags) > 0 || len(excludeTags) > 0 {
			filteredPosts = filterPostsByTags(*fullPosts, includeTags, excludeTags)
		}

		// Set up context for the tag parameter
		ctx := context.WithValue(r.Context(), tagsParamContextKey, tagsParam)
		ctx = context.WithValue(ctx, currentURLContextKey, r.URL.String())

		// Render the posts template with filtered posts
		component := views.Page(views.Posts(&filteredPosts, fullProjects, fullTags, fullPostSlugMap, fullProjectSlugMap, fullTagSlugMap))
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r.WithContext(ctx))
		return nil
	}, nil
}

// Post is the post handler.
func Post(
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
		id := r.PathValue("id")
		if id == "" {
			return routing.ErrMissingParam{ID: id, View: "Post"}
		}
		post, ok := (*fullPostSlugMap)[id]
		if !ok {
			return routing.ErrNotFound{URL: r.URL}
		}
		templ.Handler(
			views.Page(views.Post(
				&post,
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

// parseTagFilters extracts include and exclude tags from the tag parameter
func parseTagFilters(tagsParam string) (includeTags, excludeTags []string) {
	if tagsParam == "" {
		return nil, nil
	}

	tags := strings.Fields(tagsParam)
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if strings.HasPrefix(tag, "+") {
			if tagName := strings.TrimPrefix(tag, "+"); tagName != "" {
				includeTags = append(includeTags, strings.ToLower(tagName))
			}
		} else if strings.HasPrefix(tag, "-") {
			if tagName := strings.TrimPrefix(tag, "-"); tagName != "" {
				excludeTags = append(excludeTags, strings.ToLower(tagName))
			}
		} else if tag != "" {
			// If no prefix, assume include
			includeTags = append(includeTags, strings.ToLower(tag))
		}
	}
	return includeTags, excludeTags
}

// filterPostsByTags filters posts based on include and exclude tag lists
func filterPostsByTags(posts []master.FullPost, includeTags, excludeTags []string) []master.FullPost {
	if len(includeTags) == 0 && len(excludeTags) == 0 {
		return posts // No filtering needed
	}

	var filtered []master.FullPost
	for _, post := range posts {
		// Convert post tags to lowercase for case-insensitive matching
		postTags := make([]string, 0, len(post.Tags))
		for _, tag := range post.Tags {
			postTags = append(postTags, strings.ToLower(tag.Name))
		}

		// Check if post should be excluded
		excluded := false
		for _, excludeTag := range excludeTags {
			if slices.Contains(postTags, excludeTag) {
				excluded = true
			}
			if excluded {
				break
			}
		}
		if excluded {
			continue
		}

		// Check if post should be included
		if len(includeTags) > 0 {
			allTagsFound := true
			for _, includeTag := range includeTags {
				tagFound := slices.Contains(postTags, includeTag)
				if !tagFound {
					allTagsFound = false
					break
				}
			}
			if !allTagsFound {
				continue
			}
		}

		filtered = append(filtered, post)
	}

	return filtered
}

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

// Tags is the tags handler.
func Tags(
	_ context.Context,
	_ *data.Database[master.Queries],
	fullPosts *[]master.FullPost,
	fullProjects *[]master.FullProject,
	fullTags *[]master.FullTag,
	fullPostSlugMap *map[string]master.FullPost,
	fullProjectSlugMap *map[string]master.FullProject,
	fullTagSlugMap *map[string]master.FullTag,
) (routing.APIFn, error) {
	handler := templ.Handler(
		views.Page(views.Tags(
			fullPosts,
			fullProjects,
			fullTags,
			fullPostSlugMap,
			fullProjectSlugMap,
			fullTagSlugMap,
		)),
	)
	return func(w http.ResponseWriter, r *http.Request) error {
		handler.ServeHTTP(w, r)
		return nil
	}, nil
}

// Tag is the tag handler.
func Tag(
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
		id := r.PathValue("id")
		if id == "" {
			return routing.ErrMissingParam{ID: id, View: "Tag"}
		}
		tag, ok := (*fullTagSlugMap)[id]
		if !ok {
			return routing.ErrNotFound{URL: r.URL}
		}
		templ.Handler(
			views.Page(views.Tag(
				&tag,
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
