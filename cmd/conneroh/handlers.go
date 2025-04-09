package conneroh

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/hx"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/sourcegraph/conc/pool"
)

const (
	maxSearchRoutines = 10
)

var (
	home = views.Home(
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
	)
	posts = views.List(
		routing.PluralTargetPost,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
		"",
		1,
		postPages,
	)
	postPages = len(gen.AllPosts) / routing.MaxListLargeItems
	projects  = views.List(
		routing.PluralTargetProject,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
		"",
		1,
		postPages,
	)
	projectPages = len(gen.AllProjects) / routing.MaxListLargeItems
	tags         = views.List(
		routing.PluralTargetTag,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
		"",
		1,
		tagPages,
	)
	tagPages = len(gen.AllTags) / routing.MaxListLargeItems

	listMap = map[routing.PluralTarget]int{
		routing.PluralTargetPost:    routing.MaxListLargeItems,
		routing.PluralTargetProject: routing.MaxListLargeItems,
		routing.PluralTargetTag:     routing.MaxListSmallItems,
	}
)

func filterPosts(
	posts []*gen.Post,
	query string,
	page int,
) ([]*gen.Post, int) {
	p := pool.New().WithMaxGoroutines(maxSearchRoutines)
	filtered := make([]*gen.Post, 0)
	for _, post := range posts {
		post := post
		p.Go(func() {
			if strings.Contains(post.Title, query) {
				filtered = append(filtered, post)
			}
		})
	}
	p.Wait()
	// Sort results to ensure consistent pagination
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Title < filtered[j].Title
	})

	// Paginate the filtered results
	return paginate(filtered, page, routing.PluralTargetProject)
}

func filterProjects(
	projects []*gen.Project,
	query string,
	page int,
) ([]*gen.Project, int) {
	p := pool.New().WithMaxGoroutines(maxSearchRoutines)
	filtered := make([]*gen.Project, 0)
	for _, project := range projects {
		project := project
		p.Go(func() {
			if strings.Contains(project.Title, query) {
				filtered = append(filtered, project)
			}
		})
	}
	p.Wait()
	// Sort results to ensure consistent pagination
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Title < filtered[j].Title
	})

	// Paginate the filtered results
	return paginate(filtered, page, routing.PluralTargetProject)
}

func filterTags(
	tags []*gen.Tag,
	query string,
	page int,
) ([]*gen.Tag, int) {
	p := pool.New().WithMaxGoroutines(maxSearchRoutines)

	filtered := make([]*gen.Tag, 0)

	for _, tag := range tags {
		tag := tag
		p.Go(func() {
			if strings.Contains(strings.ToLower(tag.Title), strings.ToLower(query)) {
				filtered = append(filtered, tag)
			}
		})
	}
	p.Wait()

	// Sort results to ensure consistent pagination
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Title < filtered[j].Title
	})

	// Paginate the filtered results
	return paginate(filtered, page, routing.PluralTargetTag)
}

func paginate[T any](
	items []T,
	page int,
	target routing.PluralTarget,
) ([]T, int) {
	if len(items) == 0 {
		return []T{}, 0
	}

	// Get the page size based on the target type
	pageSize := listMap[target]
	if pageSize <= 0 {
		panic(fmt.Sprintf("invalid page size for target %s", target))
	}

	// Calculate total number of pages
	totalPages := (len(items) + pageSize - 1) / pageSize

	// Validate page number
	if page < 1 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	// Calculate start and end indices for the current page
	startIndex := (page - 1) * pageSize
	endIndex := min(startIndex+pageSize, len(items))

	// Return the paginated subset and the total page count
	if startIndex < len(items) {
		return items[startIndex:endIndex], totalPages
	}

	return []T{}, totalPages
}

func listHandler(
	target routing.PluralTarget,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(hx.HdrRequest)
		query := r.URL.Query().Get("search")
		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return
		}
		switch target {
		case routing.PluralTargetPost:
			filtered, totalPages := filterPosts(gen.AllPosts, query, page)
			if header == "" {
				templ.Handler(layouts.Page(views.List(
					target,
					&filtered,
					nil,
					nil,
					query,
					page,
					totalPages,
				))).ServeHTTP(w, r)
			} else {
				templ.Handler(views.Results(
					target,
					&filtered,
					nil,
					nil,
					page,
					totalPages,
				)).ServeHTTP(w, r)
			}
		case routing.PluralTargetProject:
			filtered, totalPages := filterProjects(gen.AllProjects, query, page)
			if header == "" {
				templ.Handler(layouts.Page(views.List(
					target,
					nil,
					&filtered,
					nil,
					query,
					page,
					totalPages,
				))).ServeHTTP(w, r)
			} else {
				templ.Handler(views.Results(
					target,
					nil,
					&filtered,
					nil,
					page,
					totalPages,
				)).ServeHTTP(w, r)
			}
		case routing.PluralTargetTag:
			filtered, totalPages := filterTags(gen.AllTags, query, page)
			if header == "" {
				templ.Handler(layouts.Page(views.List(
					target,
					nil,
					nil,
					&filtered,
					query,
					page,
					totalPages,
				))).ServeHTTP(w, r)
			} else {
				templ.Handler(views.Results(
					target,
					nil,
					nil,
					&filtered,
					page,
					totalPages,
				)).ServeHTTP(w, r)
			}
		}
	}
}
