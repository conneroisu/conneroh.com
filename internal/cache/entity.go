// Package cache provides a simplified caching system for assets
package cache

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/rotisserie/eris"
)

// dbFindEntityBySlug is a generic function for finding entities by slug
// This function is used to replace the duplicated code in dbFindTagBySlug, dbFindPostBySlug, and dbFindProjectBySlug.
func (p *Processor) dbFindEntityBySlug(
	ctx context.Context,
	slug, origin string,
	entityType string,
	getFromCache func(string) (interface{}, bool),
	setToCache func(interface{}),
	model interface{},
) (interface{}, error) {
	// First check in-memory cache
	if entity, ok := getFromCache(slug); ok {
		return entity, nil
	}

	// Use timeout to prevent hanging
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := p.db.NewSelect().
		Model(model).
		Where("slug = ?", slug).
		Scan(queryCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, eris.New("database operation timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, eris.Errorf(
				"%s not found: %s (referenced from %s)",
				entityType,
				slug,
				origin,
			)
		}

		return nil, eris.Wrapf(err, "failed to find %s: %s", entityType, slug)
	}

	// Update cache
	setToCache(model)

	return model, nil
}

// updateRelationships is a generic function to update entity relationships.
func (p *Processor) updateRelationships(
	ctx context.Context,
	entity interface{},
	relationType string,
) error {
	switch relationType {
	case "post":
		post, ok := entity.(*assets.Post)
		if !ok {
			return eris.New("invalid entity type for post relationship update")
		}

		return p.dbUpdatePostRelationships(ctx, post)
	case "project":
		project, ok := entity.(*assets.Project)
		if !ok {
			return eris.New("invalid entity type for project relationship update")
		}

		return p.dbUpdateProjectRelationships(ctx, project)
	case "tag":
		tag, ok := entity.(*assets.Tag)
		if !ok {
			return eris.New("invalid entity type for tag relationship update")
		}

		return p.dbUpdateTagRelationships(ctx, tag)
	default:
		return eris.Errorf("unknown relationship type: %s", relationType)
	}
}
