package conneroh

import (
	"context"
	"log/slog"
	"net/http"

	static "github.com/conneroisu/conneroh.com/cmd/conneroh/_static"
	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"github.com/uptrace/bun"
)

var (
	// Instance Caches
	allPosts    = []*assets.Post{}
	allProjects = []*assets.Project{}
	allTags     = []*assets.Tag{}
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
		routing.Make(HandleHome(db)))
	h.HandleFunc(
		"POST /contact",
		routing.Make(handleContactForm()))
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
		routing.Make(HandlePosts(db)))
	h.Handle(
		"GET /search/posts",
		routing.Make(listHandler(routing.PostPluralPath, db)))

	h.HandleFunc(
		"GET /post/{slug...}",
		routing.Make(HandlePost(db)))

	h.HandleFunc(
		"GET /projects",
		routing.Make(HandleProjects(db)))
	h.Handle(
		"GET /search/projects",
		routing.Make(listHandler(routing.ProjectPluralPath, db)))

	h.HandleFunc(
		"GET /project/{slug...}",
		routing.Make(HandleProject(db)))

	h.HandleFunc(
		"GET /tags",
		routing.Make(HandleTags(db)))

	h.Handle(
		"GET /search/tags",
		routing.Make(listHandler(routing.TagsPluralPath, db)))

	h.HandleFunc(
		"GET /tag/{slug...}",
		routing.Make(HandleTag(db)))

	return nil
}
