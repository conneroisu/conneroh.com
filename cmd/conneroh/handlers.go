package conneroh

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

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
		routing.PostPluralPath,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
		"",
		1,
		postPages,
	)
	postPages = (len(gen.AllPosts) + routing.MaxListLargeItems - 1) / routing.MaxListLargeItems
	projects  = views.List(
		routing.ProjectPluralPath,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
		"",
		1,
		projectPages,
	)
	projectPages = (len(gen.AllProjects) + routing.MaxListLargeItems - 1) / routing.MaxListLargeItems
	tags         = views.List(
		routing.TagsPluralPath,
		&gen.AllPosts,
		&gen.AllProjects,
		&gen.AllTags,
		"",
		1,
		tagPages,
	)
	tagPages = (len(gen.AllTags) + routing.MaxListLargeItems - 1) / routing.MaxListLargeItems
)

func filter[T any](
	items []T,
	query string,
	page int,
	pageSize int,
	titleGetter func(T) string,
) ([]T, int) {
	p := pool.New().WithMaxGoroutines(maxSearchRoutines)

	// Use a mutex to safely collect results
	var mu sync.Mutex
	filtered := make([]T, 0)

	for i := range items {
		item := items[i]
		p.Go(func() {
			title := titleGetter(item)
			if strings.Contains(strings.ToLower(title), strings.ToLower(query)) {
				mu.Lock()
				filtered = append(filtered, item)
				mu.Unlock()
			}
		})
	}

	p.Wait()

	// Sort results to ensure consistent pagination
	sort.Slice(filtered, func(i, j int) bool {
		return titleGetter(filtered[i]) < titleGetter(filtered[j])
	})

	// Paginate the filtered results
	return paginate(filtered, page, pageSize)
}

func paginate[T any](
	items []T,
	page int,
	pageSize int,
) ([]T, int) {
	if len(items) == 0 || pageSize <= 0 {
		return []T{}, 0
	}

	// Calculate total number of pages (use exact division with ceiling)
	totalPages := (len(items) + pageSize - 1) / pageSize

	page = max(1, page)
	page = min(page, totalPages)

	// Calculate start and end indices for the current page
	startIndex := (page - 1) * pageSize
	endIndex := min(startIndex+pageSize, len(items))

	// Return the paginated subset and the total page count
	return items[startIndex:endIndex], totalPages
}

func listHandler(
	target routing.PluralPath,
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
		case routing.PostPluralPath:
			filtered, totalPages := filter(gen.AllPosts, query, page, routing.MaxListLargeItems, func(post *gen.Post) string {
				return post.Title
			})
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
		case routing.ProjectPluralPath:
			filtered, totalPages := filter(gen.AllProjects, query, page, routing.MaxListLargeItems, func(project *gen.Project) string {
				return project.Title
			})
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
		case routing.TagsPluralPath:
			filtered, totalPages := filter(gen.AllTags, query, page, routing.MaxListSmallItems, func(tag *gen.Tag) string {
				return tag.Title
			})
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

func globalSearchHandler(
	posts []*gen.Post,
	projects []*gen.Project,
	tags []*gen.Tag,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
