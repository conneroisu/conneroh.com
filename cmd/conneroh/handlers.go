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
	"github.com/rotisserie/eris"
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
) routing.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		header := r.Header.Get(hx.HdrRequest)
		query := r.URL.Query().Get("search")
		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return eris.Wrap(
				err,
				"failed to parse page",
			)
		}
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
					return eris.Wrap(
						err,
						"failed to scan posts for post list",
					)
				}
			}
			filtered, totalPages := filter(
				allPosts,
				query,
				page,
				routing.MaxListLargeItems,
				func(post *assets.Post) string {
					return post.Title
				},
			)
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
					return eris.Wrap(
						err,
						"failed to scan projects for project list",
					)
				}
			}
			filtered, totalPages := filter(
				allProjects,
				query,
				page,
				routing.MaxListLargeItems,
				func(project *assets.Project) string {
					return project.Title
				},
			)
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
					return eris.Wrap(
						err,
						"failed to scan tags for tag list",
					)
				}
			}
			filtered, totalPages := filter(
				allTags,
				query,
				page,
				routing.MaxListSmallItems,
				func(tag *assets.Tag) string {
					return tag.Title
				},
			)
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
		return nil
	}
}

// HandleHome handles the home page. aka /{$}
func HandleHome(db *bun.DB) func(w http.ResponseWriter, r *http.Request) error {
	var homePage *templ.Component
	return func(w http.ResponseWriter, r *http.Request) error {
		if homePage != nil {
			routing.MorphableHandler(
				layouts.Page(*homePage),
				*homePage,
			).ServeHTTP(w, r)
			return nil
		}
		err := db.NewSelect().Model(&allPosts).
			Order("updated_at").
			Relation("Tags").
			Relation("Posts").
			Relation("Projects").
			Scan(r.Context())
		if err != nil {
			return eris.Wrap(
				err,
				"failed to scan posts for home page",
			)
		}
		err = db.NewSelect().Model(&allProjects).
			Order("updated_at").
			Relation("Tags").
			Relation("Posts").
			Relation("Projects").
			Scan(r.Context())
		if err != nil {
			return eris.Wrap(
				err,
				"failed to scan projects for home page",
			)
		}
		err = db.NewSelect().Model(&allTags).
			Order("updated_at").
			Relation("Tags").
			Relation("Posts").
			Relation("Projects").
			Scan(r.Context())
		if err != nil {
			return eris.Wrap(
				err,
				"failed to scan tags for home page",
			)
		}
		var home = views.Home(
			&allPosts,
			&allProjects,
			&allTags,
		)
		routing.MorphableHandler(
			layouts.Page(home),
			home,
		).ServeHTTP(w, r)
		homePage = &home
		return nil
	}
}

// HandleProjects handles the projects page. aka /projects
func HandleProjects(db *bun.DB) routing.APIFunc {
	// Handler Component Cache
	var projectList *templ.Component
	return func(w http.ResponseWriter, r *http.Request) error {
		if projectList != nil {
			routing.MorphableHandler(
				layouts.Page(*projectList),
				*projectList,
			).ServeHTTP(w, r)
			return nil
		}
		if len(allProjects) == 0 {
			_, err := db.NewSelect().Model(assets.EmpProject).
				Order("created_at").
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Exec(r.Context(), &allProjects)
			if err != nil {
				return eris.Wrap(
					err,
					"failed to scan projects for projects page",
				)
			}
		}
		var projects = views.List(
			routing.ProjectPluralPath,
			nil,
			&allProjects,
			nil,
			"",
			1,
			(len(allProjects)+routing.MaxListLargeItems-1)/routing.MaxListLargeItems,
		)
		routing.MorphableHandler(layouts.Page(projects), projects).ServeHTTP(w, r)
		projectList = &projects
		return nil
	}
}

// HandlePost handles the post page. aka /post/{slug...}
func HandlePost(db *bun.DB) routing.APIFunc {
	// Handler Component Slug-Mapped Cache
	var postMap = map[string]templ.Component{}
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			p    assets.Post
			c    templ.Component
			ok   bool
			slug = routing.Slug(r)
		)
		if c, ok = postMap[slug]; ok {
			routing.MorphableHandler(
				layouts.Page(c),
				c,
			).ServeHTTP(w, r)
			return nil
		}
		err := db.NewSelect().Model(&p).
			Where("slug = ?", slug).
			Relation("Tags").
			Relation("Posts").
			Relation("Projects").
			Limit(1).Scan(r.Context())
		if err != nil {
			slog.Error("failed to scan post", "err", err)
			return eris.Wrap(
				err,
				"failed to scan post",
			)
		}
		c = views.Post(&p)
		routing.MorphableHandler(layouts.Page(c), c).
			ServeHTTP(w, r)
		postMap[slug] = c
		return nil
	}
}

// HandleProject handles the project page. aka /project/{slug...}
func HandleProject(db *bun.DB) routing.APIFunc {
	// Handler Component Slug-Mapped Cache
	var projectMap = map[string]templ.Component{}
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			p    assets.Project
			c    templ.Component
			ok   bool
			slug = routing.Slug(r)
		)
		if c, ok = projectMap[slug]; ok {
			routing.MorphableHandler(
				layouts.Page(c),
				c,
			).ServeHTTP(w, r)
			return nil
		}
		err := db.NewSelect().Model(&p).
			Where("slug = ?", slug).
			Relation("Tags").
			Relation("Posts").
			Relation("Projects").
			Limit(1).Scan(r.Context())
		if err != nil {
			return eris.Wrap(
				err,
				"failed to scan project",
			)
		}
		c = views.Project(&p)
		routing.MorphableHandler(
			layouts.Page(c),
			c,
		).ServeHTTP(w, r)
		projectMap[slug] = c
		return nil
	}
}

// HandleTags handles the tags page. aka /tags
func HandleTags(db *bun.DB) routing.APIFunc {
	// Handler Component Cache
	var tagList *templ.Component
	return func(w http.ResponseWriter, r *http.Request) error {
		if tagList != nil {
			routing.MorphableHandler(
				layouts.Page(*tagList),
				*tagList,
			).ServeHTTP(w, r)
		}
		_, err := db.NewSelect().Model(assets.EmpTag).
			Order("created_at").
			Relation("Tags").
			Relation("Posts").
			Relation("Projects").
			Exec(r.Context(), &allTags)
		if err != nil {
			slog.Error("failed to scan tags", "err", err)
			return eris.Wrap(
				err,
				"failed to scan tags for tags page",
			)
		}
		var tags = views.List(
			routing.TagsPluralPath,
			nil,
			nil,
			&allTags,
			"",
			1,
			(len(allTags)+routing.MaxListSmallItems-1)/routing.MaxListSmallItems,
		)
		routing.MorphableHandler(layouts.Page(tags), tags).ServeHTTP(w, r)
		tagList = &tags
		return nil
	}
}

// HandleTag handles the tag page. aka /tag/{slug...}
func HandleTag(db *bun.DB) routing.APIFunc {
	// Handler Component Slug-Mapped Cache
	var tagMap = map[string]templ.Component{}
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			t    assets.Tag
			c    templ.Component
			ok   bool
			slug = routing.Slug(r)
		)
		if c, ok = tagMap[slug]; ok {
			routing.MorphableHandler(
				layouts.Page(c),
				c,
			).ServeHTTP(w, r)
			return nil
		}
		_, err := db.NewSelect().Model(assets.EmpTag).
			Where("slug = ?", routing.Slug(r)).
			Relation("Tags").
			Relation("Posts").
			Relation("Projects").
			Limit(1).Exec(r.Context(), &t)
		if err != nil {
			return eris.Wrap(
				err,
				"failed to scan tag",
			)
		}
		c = views.Tag(&t)
		routing.MorphableHandler(
			layouts.Page(c),
			c,
		).ServeHTTP(w, r)
		tagMap[slug] = c
		return nil
	}
}

// HandlePosts handles the posts page. aka /posts
func HandlePosts(db *bun.DB) routing.APIFunc {
	// Handler Component Cache
	var postList *templ.Component
	return func(w http.ResponseWriter, r *http.Request) error {
		if postList != nil {
			routing.MorphableHandler(
				layouts.Page(*postList),
				*postList,
			).ServeHTTP(w, r)
			return nil
		}
		_, err := db.NewSelect().Model(assets.EmpPost).
			Order("created_at").
			Relation("Tags").
			Relation("Posts").
			Relation("Projects").
			Exec(r.Context(), &allPosts)
		if err != nil {
			return eris.Wrap(
				err,
				"failed to scan posts for posts page",
			)
		}
		var posts = views.List(
			routing.PostPluralPath,
			&allPosts,
			nil,
			nil,
			"",
			1,
			(len(allPosts)+routing.MaxListLargeItems-1)/routing.MaxListLargeItems,
		)
		routing.MorphableHandler(layouts.Page(posts), posts).ServeHTTP(w, r)
		postList = &posts
		return nil
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
