package cache

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/rotisserie/eris"
)

// dbFindTagBySlug finds a tag by its slug (to be called from the DB worker).
func (p *Processor) dbFindTagBySlug(ctx context.Context, slug, origin string) (*assets.Tag, error) {
	// First check in-memory cache
	if tag, ok := p.entityCache.GetTag(slug); ok {
		return tag, nil
	}

	// Use timeout to prevent hanging
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var tag assets.Tag

	err := p.db.NewSelect().
		Model(&tag).
		Where("slug = ?", slug).
		Scan(queryCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, eris.New("database operation timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, eris.Errorf(
				"tag not found: %s (referenced from %s)",
				slug,
				origin,
			)
		}

		return nil, eris.Wrapf(err, "failed to find tag: %s", slug)
	}

	// Update cache
	p.entityCache.SetTag(&tag)

	return &tag, nil
}

// dbFindPostBySlug finds a post by its slug (to be called from the DB worker).
func (p *Processor) dbFindPostBySlug(ctx context.Context, slug, origin string) (*assets.Post, error) {
	// First check in-memory cache
	if post, ok := p.entityCache.GetPost(slug); ok {
		return post, nil
	}

	// Use timeout to prevent hanging
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var post assets.Post

	err := p.db.NewSelect().
		Model(&post).
		Where("slug = ?", slug).
		Scan(queryCtx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, eris.New("database operation timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, eris.Errorf(
				"post not found: %s (referenced from %s)",
				slug,
				origin,
			)
		}

		return nil, eris.Wrapf(err, "failed to find post: %s", slug)
	}

	// Update cache
	p.entityCache.SetPost(&post)

	return &post, nil
}

// dbFindProjectBySlug finds a project by its slug (to be called from the DB worker).
func (p *Processor) dbFindProjectBySlug(ctx context.Context, slug, origin string) (*assets.Project, error) {
	// First check in-memory cache
	if project, ok := p.entityCache.GetProject(slug); ok {
		return project, nil
	}

	// Use timeout to prevent hanging
	queryCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var project assets.Project

	err := p.db.NewSelect().
		Model(&project).
		Where("slug = ?", slug).
		Scan(queryCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, eris.New("database operation timed out")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, eris.Errorf(
				"project not found: %s (referenced from %s)",
				slug,
				origin,
			)
		}

		return nil, eris.Wrapf(err, "failed to find project: %s", slug)
	}

	// Update cache
	p.entityCache.SetProject(&project)

	return &project, nil
}
