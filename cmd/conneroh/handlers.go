package conneroh

import (
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/a-h/templ"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/rotisserie/eris"
	"github.com/sourcegraph/conc/pool"
	"github.com/uptrace/bun"
)

const (
	maxSearchRoutines = 10
)

func listHandler(
	target routing.PluralPath,
	db *bun.DB,
) routing.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		header := r.Header.Get(routing.HdrRequest)
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
					nil,
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
					nil,
					page,
					totalPages,
				)).ServeHTTP(w, r)
			}
		case routing.EmploymentPluralPath:
			if len(allEmployments) == 0 {
				err = db.NewSelect().Model(&allEmployments).
					Order("updated_at").
					Relation("Tags").
					Relation("Posts").
					Relation("Projects").
					Relation("Employments").
					Scan(r.Context())
				if err != nil {
					return eris.Wrap(
						err,
						"failed to scan employments for employment list",
					)
				}
			}
			filtered, totalPages := filter(
				allEmployments,
				query,
				page,
				routing.MaxListLargeItems,
				func(employment *assets.Employment) string {
					return employment.Title
				},
			)
			if header == "" {
				templ.Handler(layouts.Page(views.List(
					target,
					nil,
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

// HandleHome handles the home page. aka /{$}.
func HandleHome(db *bun.DB) func(w http.ResponseWriter, r *http.Request) error {
	var homePage *templ.Component

	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		if homePage != nil {
			routing.MorphableHandler(
				layouts.Page,
				*homePage,
			).ServeHTTP(w, r)

			return nil
		}
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
					"failed to scan posts for home page",
				)
			}
		}
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
					"failed to scan projects for home page",
				)
			}
		}
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
					"failed to scan tags for home page",
				)
			}
		}
		if len(allEmployments) == 0 {
			err = db.NewSelect().Model(&allEmployments).
				Order("updated_at").
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Relation("Employments").
				Scan(r.Context())
			if err != nil {
				return eris.Wrap(
					err,
					"failed to scan employments for home page",
				)
			}
		}
		home := views.Home(
			&allPosts,
			&allProjects,
			&allTags,
			&allEmployments,
		)
		homePage = &home
		routing.MorphableHandler(
			layouts.Page,
			home,
		).ServeHTTP(w, r)

		return nil
	}
}

// HandleProjects handles the projects page. aka /projects.
func HandleProjects(db *bun.DB) routing.APIFunc {
	// Handler Component Cache
	var projectList *templ.Component

	return func(w http.ResponseWriter, r *http.Request) error {
		if projectList != nil {
			routing.MorphableHandler(
				layouts.Page,
				*projectList,
			).ServeHTTP(w, r)

			return nil
		}
		if len(allProjects) == 0 {
			err := db.NewSelect().Model(&allProjects).
				Order("created_at").
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Scan(r.Context())
			if err != nil {
				return eris.Wrap(
					err,
					"failed to scan projects for projects page",
				)
			}
		}
		projects := views.List(
			routing.ProjectPluralPath,
			nil,
			&allProjects,
			nil,
			nil,
			"",
			1,
			(len(allProjects)+routing.MaxListLargeItems-1)/routing.MaxListLargeItems,
		)
		projectList = &projects
		routing.MorphableHandler(
			layouts.Page,
			projects,
		).ServeHTTP(w, r)

		return nil
	}
}

// HandlePost handles the post page. aka /post/{slug...}.
func HandlePost(db *bun.DB) routing.APIFunc {
	// Handler Component Slug-Mapped Cache
	postMap := map[string]templ.Component{}

	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			p    assets.Post
			comp templ.Component
			ok   bool
			slug = routing.Slug(r)
		)
		if comp, ok = postMap[slug]; ok {
			routing.MorphableHandler(
				layouts.Page,
				comp,
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
		comp = views.Post(&p)
		postMap[slug] = comp
		routing.MorphableHandler(
			layouts.Page,
			comp,
		).ServeHTTP(w, r)

		return nil
	}
}

// HandleProject handles the project page. aka /project/{slug...}.
func HandleProject(db *bun.DB) routing.APIFunc {
	// Handler Component Slug-Mapped Cache
	projectMap := map[string]templ.Component{}

	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			p    assets.Project
			c    templ.Component
			ok   bool
			slug = routing.Slug(r)
		)
		if c, ok = projectMap[slug]; ok {
			routing.MorphableHandler(
				layouts.Page,
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
		projectMap[slug] = c
		routing.MorphableHandler(
			layouts.Page,
			c,
		).ServeHTTP(w, r)

		return nil
	}
}

// HandleTags handles the tags page. aka /tags.
func HandleTags(db *bun.DB) routing.APIFunc {
	// Handler Component Cache
	var tagList *templ.Component

	return func(w http.ResponseWriter, r *http.Request) error {
		if tagList != nil {
			routing.MorphableHandler(
				layouts.Page,
				*tagList,
			).ServeHTTP(w, r)

			return nil
		}
		err := db.NewSelect().Model(&allTags).
			Order("created_at").
			Relation("Tags").
			Relation("Posts").
			Relation("Projects").
			Scan(r.Context())
		if err != nil {
			return eris.Wrap(
				err,
				"failed to scan tags for tags page",
			)
		}
		tags := views.List(
			routing.TagsPluralPath,
			nil,
			nil,
			&allTags,
			nil,
			"",
			1,
			(len(allTags)+routing.MaxListSmallItems-1)/routing.MaxListSmallItems,
		)
		tagList = &tags
		routing.MorphableHandler(
			layouts.Page,
			tags,
		).ServeHTTP(w, r)

		return nil
	}
}

// HandleTag handles the tag page. aka /tag/{slug...}.
func HandleTag(db *bun.DB) routing.APIFunc {
	// Handler Component Slug-Mapped Cache
	tagMap := map[string]templ.Component{}

	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			tag  assets.Tag
			comp templ.Component
			ok   bool
			slug = routing.Slug(r)
		)
		if comp, ok = tagMap[slug]; ok {
			routing.MorphableHandler(
				layouts.Page,
				comp,
			).ServeHTTP(w, r)

			return nil
		}
		err := db.NewSelect().Model(&tag).
			Where("slug = ?", routing.Slug(r)).
			Relation("Tags").
			Relation("Posts").
			Relation("Projects").
			Limit(1).Scan(r.Context())
		if err != nil {
			return eris.Wrap(
				err,
				"failed to scan tag",
			)
		}
		comp = views.Tag(&tag)
		tagMap[slug] = comp
		routing.MorphableHandler(
			layouts.Page,
			comp,
		).ServeHTTP(w, r)

		return nil
	}
}

// HandlePosts handles the posts page. aka /posts.
func HandlePosts(db *bun.DB) routing.APIFunc {
	// Handler Component Cache
	var postList *templ.Component

	return func(w http.ResponseWriter, r *http.Request) error {
		if postList != nil {
			routing.MorphableHandler(
				layouts.Page,
				*postList,
			).ServeHTTP(w, r)

			return nil
		}
		if len(allPosts) == 0 {
			err := db.NewSelect().Model(&allPosts).
				Order("created_at").
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Scan(r.Context())
			if err != nil {
				return eris.Wrap(
					err,
					"failed to scan posts for posts page",
				)
			}
		}
		posts := views.List(
			routing.PostPluralPath,
			&allPosts,
			nil,
			nil,
			nil,
			"",
			1,
			(len(allPosts)+routing.MaxListLargeItems-1)/routing.MaxListLargeItems,
		)
		postList = &posts
		routing.MorphableHandler(
			layouts.Page,
			posts,
		).ServeHTTP(w, r)

		return nil
	}
}

// HandleEmployments handles the employments page. aka /employments.
func HandleEmployments(db *bun.DB) routing.APIFunc {
	// Handler Component Cache
	var employmentList *templ.Component

	return func(w http.ResponseWriter, r *http.Request) error {
		if employmentList != nil {
			routing.MorphableHandler(
				layouts.Page,
				*employmentList,
			).ServeHTTP(w, r)

			return nil
		}
		if len(allEmployments) == 0 {
			err := db.NewSelect().Model(&allEmployments).
				Order("created_at").
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Relation("Employments").
				Scan(r.Context())
			if err != nil {
				return eris.Wrap(
					err,
					"failed to scan employments for employments page",
				)
			}
		}
		employments := views.List(
			routing.EmploymentPluralPath,
			nil,
			nil,
			nil,
			&allEmployments,
			"",
			1,
			(len(allEmployments)+routing.MaxListLargeItems-1)/routing.MaxListLargeItems,
		)
		employmentList = &employments
		routing.MorphableHandler(
			layouts.Page,
			employments,
		).ServeHTTP(w, r)

		return nil
	}
}

// HandleEmployment handles the employment page. aka /employment/{slug...}.
func HandleEmployment(db *bun.DB) routing.APIFunc {
	// Handler Component Slug-Mapped Cache
	employmentMap := map[string]templ.Component{}

	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			emp  assets.Employment
			comp templ.Component
			ok   bool
			slug = routing.Slug(r)
		)
		if comp, ok = employmentMap[slug]; ok {
			routing.MorphableHandler(
				layouts.Page,
				comp,
			).ServeHTTP(w, r)

			return nil
		}
		err := db.NewSelect().Model(&emp).
			Where("slug = ?", slug).
			Relation("Tags").
			Relation("Posts").
			Relation("Projects").
			Relation("Employments").
			Limit(1).Scan(r.Context())
		if err != nil {
			return eris.Wrap(
				err,
				"failed to scan employment",
			)
		}
		comp = views.Employment(&emp)
		employmentMap[slug] = comp
		routing.MorphableHandler(
			layouts.Page,
			comp,
		).ServeHTTP(w, r)

		return nil
	}
}

// filter returns a paginated slice of items matching the search query, ranked by relevance across multiple fields.
// The function scores and filters items concurrently, prioritizing matches in the title, description, content, tags, and icon fields depending on the item type.
// Results are sorted by descending relevance before pagination. If the query is empty, all items are returned paginated.
func filter[T any](
	items []T,
	query string,
	page int,
	pageSize int,
	titleGetter func(T) string,
) ([]T, int) {
	// If query is empty, return all items
	if query == "" {
		sort.Slice(items, func(i, j int) bool {
			return strings.ToLower(titleGetter(items[i])) <
				strings.ToLower(titleGetter(items[j]))
		})

		return routing.Paginate(items, page, pageSize)
	}

	p := pool.New().WithMaxGoroutines(maxSearchRoutines)
	query = strings.ToLower(query)

	// Use a mutex to safely collect results with relevance score
	type searchResult struct {
		item  T
		score int // Higher is better match
	}

	var mu sync.Mutex
	results := make([]searchResult, 0)

	for i := range items {
		item := items[i]
		p.Go(func() {
			var score int
			var matchFound bool

			// Check title match (highest priority)
			title := titleGetter(item)
			if title != "" {
				titleLower := strings.ToLower(title)
				if strings.Contains(titleLower, query) {
					score += 100
					matchFound = true
					// Exact title match gets bonus points
					if titleLower == query {
						score += 50
					}
				}
			}

			// Check in other fields based on type
			switch v := any(item).(type) {
			case *assets.Post:
				// Check description
				if strings.Contains(strings.ToLower(v.Description), query) {
					score += 50
					matchFound = true
				}
				// Check content
				if strings.Contains(strings.ToLower(v.Content), query) {
					score += 30
					matchFound = true
				}
				// Check tags
				for _, tag := range v.Tags {
					if tag != nil && strings.Contains(strings.ToLower(tag.Title), query) {
						score += 25
						matchFound = true
					}
				}
			case *assets.Project:
				// Check description
				if strings.Contains(strings.ToLower(v.Description), query) {
					score += 50
					matchFound = true
				}
				// Check content
				if strings.Contains(strings.ToLower(v.Content), query) {
					score += 30
					matchFound = true
				}
				// Check tags
				for _, tag := range v.Tags {
					if tag != nil && strings.Contains(strings.ToLower(tag.Title), query) {
						score += 25
						matchFound = true
					}
				}
			case *assets.Tag:
				// Check description
				if strings.Contains(strings.ToLower(v.Description), query) {
					score += 50
					matchFound = true
				}
				// Check content
				if strings.Contains(strings.ToLower(v.Content), query) {
					score += 30
					matchFound = true
				}
				// Check icon name if it matches query
				if v.Icon != "" && strings.Contains(strings.ToLower(v.Icon), query) {
					score += 10
					matchFound = true
				}
			case *assets.Employment:
				// Check description
				if strings.Contains(strings.ToLower(v.Description), query) {
					score += 50
					matchFound = true
				}
				// Check content
				if strings.Contains(strings.ToLower(v.Content), query) {
					score += 30
					matchFound = true
				}
				// Check tags
				for _, tag := range v.Tags {
					if tag != nil && strings.Contains(strings.ToLower(tag.Title), query) {
						score += 25
						matchFound = true
					}
				}
			}

			// If any match found, add to results
			if matchFound {
				mu.Lock()
				results = append(results, searchResult{item, score})
				mu.Unlock()
			}
		})
	}

	p.Wait()

	// Sort results by relevance score (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	// Extract just the items from the sorted results
	filtered := make([]T, len(results))
	for i, result := range results {
		filtered[i] = result.item
	}

	// Paginate the filtered results
	return routing.Paginate(filtered, page, pageSize)
}
