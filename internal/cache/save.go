package cache

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/rotisserie/eris"
)

// dbSavePost saves a post to the database (to be called from the DB worker).
func (p *Processor) dbSavePost(ctx context.Context, post *assets.Post) error {
	// Save post with a timeout context to prevent hanging
	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Save post
	_, err := p.db.NewInsert().
		Model(post).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("created_at = EXCLUDED.created_at").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(insertCtx)
	slog.Info(
		"saved post",
		slog.String("slug", post.Slug),
		slog.String("banner_path", post.BannerPath),
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return eris.New("database operation timed out")
		}

		return eris.Wrapf(err, "failed to save post: %s", post.Slug)
	}

	// Update the entity cache
	p.entityCache.SetPost(post)

	return nil
}

// dbSaveProject saves a project to the database (to be called from the DB worker).
func (p *Processor) dbSaveProject(ctx context.Context, project *assets.Project) error {
	// Save project with a timeout context to prevent hanging
	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Save project
	_, err := p.db.NewInsert().
		Model(project).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("created_at = EXCLUDED.created_at").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(insertCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return eris.New("database operation timed out")
		}

		return eris.Wrapf(err, "failed to save project: %s", project.Slug)
	}

	// Update the entity cache
	p.entityCache.SetProject(project)

	return nil
}

// dbSaveTag saves a tag to the database (to be called from the DB worker).
func (p *Processor) dbSaveTag(ctx context.Context, tag *assets.Tag) error {
	// Save tag with a timeout context to prevent hanging
	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Save tag
	_, err := p.db.NewInsert().
		Model(tag).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("created_at = EXCLUDED.created_at").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(insertCtx)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return eris.New("database operation timed out")
		}

		return eris.Wrapf(err, "failed to save tag: %s", tag.Slug)
	}

	// Update the entity cache
	p.entityCache.SetTag(tag)

	return nil
}
