package conneroh

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/master"
	"github.com/conneroisu/conneroh.com/internal/routing"
	"golang.org/x/sync/errgroup"
)

// RouteMap is a map of all routes.
var RouteMap = routing.APIMap{
	"GET /dist/":                      Dist,
	"GET /favicon.ico":                Favicon,
	"GET /{$}":                        Home,
	"GET /projects":                   Projects,
	"GET /posts":                      Posts,
	"GET /tags":                       Tags,
	"GET /project/{id}":               Project,
	"GET /post/{id}":                  Post,
	"GET /tag/{id...}":                Tag,
	"GET /hateoas/morph/{view...}":    Morph,
	"GET /hateoas/morphs/{view}/{id}": Morphs,
}

// AddRoutes adds all routes to the router.
func AddRoutes(
	ctx context.Context,
	h *http.ServeMux,
	db *data.Database[master.Queries],
) error {
	slog.Info("getting full data")
	var (
		fullPosts      *[]master.FullPost
		fullProjects   *[]master.FullProject
		fullTags       *[]master.FullTag
		projectSlugMap *map[string]master.FullProject
		tagSlugMap     *map[string]master.FullTag
		postSlugMap    *map[string]master.FullPost
	)
	eg := errgroup.Group{}
	eg.Go(func() (err error) {
		fullPosts, err = db.Queries.FullPostsList(ctx)
		return err
	})
	eg.Go(func() (err error) {
		fullTags, err = db.Queries.FullTagsList(ctx)
		return err
	})
	eg.Go(func() (err error) {
		fullProjects, err = db.Queries.FullProjectsList(ctx)
		return err
	})
	eg.Go(func() (err error) {
		projectSlugMap, err = db.Queries.FullProjectsSlugMapGet(ctx)
		return err
	})
	eg.Go(func() (err error) {
		tagSlugMap, err = db.Queries.FullTagsSlugMapGet(ctx)
		return err
	})
	eg.Go(func() (err error) {
		postSlugMap, err = db.Queries.FullPostsSlugMapGet(ctx)
		return err
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	slog.Info("got full data")
	slog.Info("adding routes")
	defer slog.Info("added routes")
	return RouteMap.AddRoutes(
		ctx,
		h,
		db,
		fullPosts,
		fullProjects,
		fullTags,
		postSlugMap,
		projectSlugMap,
		tagSlugMap,
	)
}
