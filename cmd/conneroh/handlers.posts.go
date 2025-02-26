package conneroh

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/master"
	"github.com/conneroisu/conneroh.com/internal/routing"
)

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
		component := views.Posts(&filteredPosts, fullProjects, fullTags, fullPostSlugMap, fullProjectSlugMap, fullTagSlugMap)
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r.WithContext(ctx))
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
