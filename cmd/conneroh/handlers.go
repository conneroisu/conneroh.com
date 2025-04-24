package conneroh

import (
	"log/slog"
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
	"github.com/conneroisu/conneroh.com/internal/hx"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/gorilla/schema"
	"github.com/sourcegraph/conc/pool"
	"github.com/uptrace/bun"
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

func listHandler(
	target routing.PluralPath,
	db *bun.DB,
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
		slog.Info("searching", "target", target)
		switch target {
		case routing.PostPluralPath:
			if len(allPosts) == 0 {
				err = db.NewSelect().Model(&allPosts).
					Order("updated_at").
					Relation("Tags").
					Relation("Posts").
					Relation("Projects").
					Scan(r.Context())
				if err != nil {
					slog.Error("failed to scan posts", "err", err, "posts", postList)
					return
				}
			}
			filtered, totalPages := filter(allPosts, query, page, routing.MaxListLargeItems, func(post *assets.Post) string {
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
			if len(allProjects) == 0 {
				err = db.NewSelect().Model(&allProjects).
					Order("updated_at").
					Relation("Tags").
					Relation("Posts").
					Relation("Projects").
					Scan(r.Context())
				if err != nil {
					slog.Error("failed to scan projects", "err", err, "projects", projectList)
					return
				}
			}
			filtered, totalPages := filter(allProjects, query, page, routing.MaxListLargeItems, func(project *assets.Project) string {
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
			if len(allTags) == 0 {
				err = db.NewSelect().Model(&allTags).
					Order("updated_at").
					Relation("Tags").
					Relation("Posts").
					Relation("Projects").
					Scan(r.Context())
				if err != nil {
					slog.Error("failed to scan tags", "err", err, "tags", tagList)
					return
				}
			}
			filtered, totalPages := filter(allTags, query, page, routing.MaxListSmallItems, func(tag *assets.Tag) string {
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
	db *bun.DB,
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
