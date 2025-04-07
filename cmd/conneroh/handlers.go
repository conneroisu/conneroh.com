package conneroh

import (
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/routing"
)

const hName = "HX-Trigger-Name"

func filterPosts(
	posts []*gen.Post,
	query string,
) []*gen.Post {
	filtered := make([]*gen.Post, 0)
	for _, post := range posts {
		if strings.Contains(post.Title, query) {
			filtered = append(filtered, post)
		}
	}
	return filtered
}

func filterProjects(
	projects []*gen.Project,
	query string,
) []*gen.Project {
	return projects
}

func filterTags(
	tags []*gen.Tag,
	query string,
) []*gen.Tag {
	return tags
}

func searchHandler(
	target routing.PluralTarget,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("search")
		header := r.Header.Get(hName)
		switch target {
		case routing.PluralTargetPost:
			filtered := filterPosts(gen.AllPosts, query)
			if header == "" {
				templ.Handler(layouts.Page(views.List(target, &filtered, nil, nil, query))).ServeHTTP(w, r)
			} else {
				templ.Handler(views.Results(target, &gen.AllPosts, nil, nil)).ServeHTTP(w, r)
			}
		case routing.PluralTargetProject:
			filtered := filterProjects(gen.AllProjects, query)
			if header == "" {
				templ.Handler(layouts.Page(views.List(target, nil, &filtered, nil, query))).ServeHTTP(w, r)
			} else {
				templ.Handler(views.Results(target, nil, &gen.AllProjects, nil)).ServeHTTP(w, r)
			}
		case routing.PluralTargetTag:
			filtered := filterTags(gen.AllTags, query)
			if header == "" {
				templ.Handler(layouts.Page(views.List(target, nil, nil, &filtered, query))).ServeHTTP(w, r)
			} else {
				templ.Handler(views.Results(target, nil, nil, &gen.AllTags)).ServeHTTP(w, r)
			}
		}
	}
}
