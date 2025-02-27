package conneroh

import (
	"context"
	"net/http"

	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/master"
	"github.com/conneroisu/conneroh.com/internal/routing"
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
	"GET /tag/{id}":                   Tag,
	"GET /hateoas/morph/{view}":       Morph,
	"GET /hateoas/morphs/{view}/{id}": Morphs,
}

// AddRoutes adds all routes to the router.
func AddRoutes(
	ctx context.Context,
	h *http.ServeMux,
	db *data.Database[master.Queries],
) error {
	fullPosts, err := db.Queries.FullPostsList(ctx)
	if err != nil {
		return err
	}
	fullTags, err := db.Queries.FullTagsList(ctx)
	if err != nil {
		return err
	}
	fullProjects, err := db.Queries.FullProjectsList(ctx)
	if err != nil {
		return err
	}
	projectSlugMap, err := db.Queries.FullProjectsSlugMapGet(ctx)
	if err != nil {
		return err
	}
	tagSlugMap, err := db.Queries.FullTagsSlugMapGet(ctx)
	if err != nil {
		return err
	}
	postSlugMap, err := db.Queries.FullPostsSlugMapGet(ctx)
	if err != nil {
		return err
	}
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
