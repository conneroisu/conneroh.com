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
	"GET /dist/":       Dist,
	"GET /favicon.ico": Favicon,
	"GET /{$}":         Home,
	// "GET /project/{id}":                  Project,
	// "GET /post/{id}":                     Post,
	// "GET /tag/{id...}":                   Tag,
	"GET /{target}/{id...}":              Single,
	"GET /{targets}":                     List,
	"GET /hateoas/morph/{view...}":       ListMorph,
	"GET /hateoas/morphs/{view}/{id...}": SingleMorph,
}

// AddRoutes adds all routes to the router.
func AddRoutes(
	ctx context.Context,
	h *http.ServeMux,
	db *data.Database[master.Queries],
) error {
	var (
		posts    []master.Post
		projects []master.Project
		tags     []master.Tag

		fullPosts      *[]master.FullPost
		fullProjects   *[]master.FullProject
		fullTags       *[]master.FullTag
		projectSlugMap *map[string]master.FullProject
		tagSlugMap     *map[string]master.FullTag
		postSlugMap    *map[string]master.FullPost
	)
	eg := errgroup.Group{}

	// get shared data
	eg.Go(func() (err error) {
		posts, err = db.Queries.PostsList(ctx)
		return err
	})
	eg.Go(func() (err error) {
		projects, err = db.Queries.ProjectsList(ctx)
		return err
	})
	eg.Go(func() (err error) {
		tags, err = db.Queries.TagsListAlphabetical(ctx)
		return err
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	eg.Go(func() (err error) {
		fullPosts, err = db.Queries.FullPostsList(ctx, posts)
		return err
	})
	eg.Go(func() (err error) {
		fullTags, err = db.Queries.FullTagsList(ctx, tags)
		return err
	})
	eg.Go(func() (err error) {
		fullProjects, err = db.Queries.FullProjectsList(ctx, projects)
		return err
	})
	eg.Go(func() (err error) {
		projectSlugMap, err = db.Queries.FullProjectsSlugMapGet(ctx, projects)
		return err
	})
	eg.Go(func() (err error) {
		tagSlugMap, err = db.Queries.FullTagsSlugMapGet(ctx, tags)
		return err
	})
	eg.Go(func() (err error) {
		postSlugMap, err = db.Queries.FullPostsSlugMapGet(ctx, posts)
		return err
	})

	if err := eg.Wait(); err != nil {
		return err
	}

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
