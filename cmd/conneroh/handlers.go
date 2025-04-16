package conneroh

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/components"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	"github.com/conneroisu/conneroh.com/internal/hx"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/gorilla/schema"
	"github.com/sourcegraph/conc/pool"
)

// ContactForm is the struct schema for the contact form.
type ContactForm struct {
	Name    string `schema:"name,required"`
	Email   string `schema:"email,required"`
	Subject string `schema:"subject,required"`
	Message string `schema:"message,required"`
}

var encoder = schema.NewEncoder()

const (
	maxSearchRoutines = 10
)

func handleContactForm(w http.ResponseWriter, r *http.Request) {
	var form ContactForm
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = encoder.Encode(form, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: Send email
	templ.Handler(components.ThankYou()).ServeHTTP(w, r)
}

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
			filtered, totalPages := filter(gen.AllPosts, query, page, routing.MaxListLargeItems, func(post *assets.Post) string {
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
				templ.Handler(views.ListResults(
					target,
					&filtered,
					nil,
					nil,
					page,
					totalPages,
				)).ServeHTTP(w, r)
			}
		case routing.ProjectPluralPath:
			filtered, totalPages := filter(gen.AllProjects, query, page, routing.MaxListLargeItems, func(project *assets.Project) string {
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
				templ.Handler(views.ListResults(
					target,
					nil,
					&filtered,
					nil,
					page,
					totalPages,
				)).ServeHTTP(w, r)
			}
		case routing.TagsPluralPath:
			filtered, totalPages := filter(gen.AllTags, query, page, routing.MaxListSmallItems, func(tag *assets.Tag) string {
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
				templ.Handler(views.ListResults(
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

func searchHandler(
	posts []*assets.Post,
	projects []*assets.Project,
	tags []*assets.Tag,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
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
	return routing.Paginate(filtered, page, pageSize)
}
