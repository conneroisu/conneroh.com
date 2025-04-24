package conneroh

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	static "github.com/conneroisu/conneroh.com/cmd/conneroh/_static"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/layouts"
	"github.com/conneroisu/conneroh.com/cmd/conneroh/views"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/uptrace/bun"
)

var (
	homePage    *templ.Component
	postList    *templ.Component
	projectList *templ.Component
	tagList     *templ.Component
	allPosts    = []*assets.Post{}
	postMap     = map[string]templ.Component{}
	allProjects = []*assets.Project{}
	projectMap  = map[string]templ.Component{}
	allTags     = []*assets.Tag{}
	tagMap      = map[string]templ.Component{}
)

// AddRoutes adds all routes to the router.
func AddRoutes(
	_ context.Context,
	h *http.ServeMux,
	db *bun.DB,
) error {
	slog.Info("adding routes")
	defer slog.Info("added routes")

	h.HandleFunc(
		"GET /{$}",
		func(w http.ResponseWriter, r *http.Request) {
			if homePage != nil {
				routing.MorphableHandler(
					layouts.Page(*homePage),
					*homePage,
				).ServeHTTP(w, r)
				return
			}
			err := db.NewSelect().Model(&allPosts).
				Order("updated_at").
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Scan(r.Context())
			if err != nil {
				slog.Error("failed to scan posts", "err", err, "posts", postList)
				return
			}
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
			err = db.NewSelect().Model(&allTags).
				Order("updated_at").
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Scan(r.Context())
			if err != nil {
				slog.Error("failed to scan tags", "err", err)
				return
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
		})
	h.HandleFunc(
		"POST /contact",
		handleContactForm)
	h.Handle(
		"GET /dist/",
		http.FileServer(http.FS(static.Dist)))
	h.HandleFunc(
		"GET /favicon.ico",
		func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write(static.Favicon)
		})
	h.HandleFunc(
		"GET /search/all",
		searchHandler(db))

	h.HandleFunc(
		"GET /posts",
		func(w http.ResponseWriter, r *http.Request) {
			if postList != nil {
				routing.MorphableHandler(layouts.Page(*postList), *postList).ServeHTTP(w, r)
				return
			}
			_, err := db.NewSelect().Model(assets.EmpPost).
				Order("created_at").
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Exec(r.Context(), &allPosts)
			if err != nil {
				slog.Error("failed to exec select posts", "err", err, "posts", postList)
				return
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
		})
	h.Handle(
		"GET /search/posts",
		listHandler(routing.PostPluralPath, db))
	h.HandleFunc(
		"GET /post/{slug...}",
		func(w http.ResponseWriter, r *http.Request) {
			var (
				p  assets.Post
				c  templ.Component
				ok bool
			)
			if c, ok = postMap[p.Slug]; ok {
				routing.MorphableHandler(
					layouts.Page(c),
					c,
				).ServeHTTP(w, r)
				return
			}
			_, err := db.NewSelect().Model(assets.EmpPost).
				Where("slug = ?", routing.Slug(r)).
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Limit(1).Exec(r.Context(), &p)
			if err != nil {
				slog.Error("failed to scan post", "err", err)
				return
			}
			c = views.Post(&p)
			routing.MorphableHandler(layouts.Page(c), c).
				ServeHTTP(w, r)
			postMap[p.Slug] = c
		})

	h.HandleFunc(
		"GET /projects",
		func(w http.ResponseWriter, r *http.Request) {
			if projectList != nil {
				routing.MorphableHandler(
					layouts.Page(*projectList),
					*projectList,
				).ServeHTTP(w, r)
				return
			}
			_, err := db.NewSelect().Model(assets.EmpProject).
				Order("created_at").
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Exec(r.Context(), &allProjects)
			if err != nil {
				slog.Error("failed to scan projects", "err", err, "projects", projectList)
				return
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
		})
	h.Handle(
		"GET /search/projects",
		listHandler(routing.ProjectPluralPath, db))

	h.HandleFunc(
		"GET /project/{slug...}",
		func(w http.ResponseWriter, r *http.Request) {
			var (
				p  assets.Project
				c  templ.Component
				ok bool
			)
			if c, ok = projectMap[p.Slug]; ok {
				routing.MorphableHandler(
					layouts.Page(c),
					c,
				).ServeHTTP(w, r)
				return
			}
			_, err := db.NewSelect().Model(assets.EmpProject).
				Where("slug = ?", routing.Slug(r)).
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Limit(1).Exec(r.Context(), &p)
			if err != nil {
				slog.Error("failed to scan project", "err", err)
				return
			}
			c = views.Project(&p)
			routing.MorphableHandler(
				layouts.Page(c),
				c,
			).ServeHTTP(w, r)
			projectMap[p.Slug] = c
		})

	h.HandleFunc(
		"GET /tags",
		func(w http.ResponseWriter, r *http.Request) {
			if tagList != nil {
				routing.MorphableHandler(layouts.Page(*tagList), *tagList).ServeHTTP(w, r)
				return
			}
			_, err := db.NewSelect().Model(assets.EmpTag).
				Order("created_at").
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Exec(r.Context(), &allTags)
			if err != nil {
				slog.Error("failed to scan tags", "err", err)
				return
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
		})

	h.Handle(
		"GET /search/tags",
		listHandler(routing.TagsPluralPath, db))

	h.HandleFunc(
		"GET /tag/{slug...}",
		func(w http.ResponseWriter, r *http.Request) {
			var (
				t  assets.Tag
				c  templ.Component
				ok bool
			)
			if c, ok = tagMap[t.Slug]; ok {
				routing.MorphableHandler(
					layouts.Page(c),
					c,
				).ServeHTTP(w, r)
				return
			}
			_, err := db.NewSelect().Model(assets.EmpTag).
				Where("slug = ?", routing.Slug(r)).
				Relation("Tags").
				Relation("Posts").
				Relation("Projects").
				Limit(1).Exec(r.Context(), &t)
			if err != nil {
				slog.Error("failed to scan tag", "err", err)
				return
			}
			c = views.Tag(&t)
			routing.MorphableHandler(
				layouts.Page(c),
				c,
			).ServeHTTP(w, r)
			tagMap[t.Slug] = c
		})

	return nil
}
