package cache

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"sync"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/copygen"
	"github.com/rotisserie/eris"
	"github.com/uptrace/bun"
)

// SavePost saves a post and its relationships to the database
func SavePost(
	ctx context.Context,
	db *bun.DB,
	post *assets.Post,
	msgCh MsgChannel,
	wg *sync.WaitGroup,
) error {
	// Save the post first
	_, err := db.NewInsert().
		Model(post).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("slug = EXCLUDED.slug").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("raw_content = EXCLUDED.raw_content").
		Set("icon = EXCLUDED.icon").
		Set("created_at = EXCLUDED.created_at").
		Set("updated_at = EXCLUDED.updated_at").
		Set("hashed = EXCLUDED.hashed").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(ctx)
	if err != nil {
		return err
	}

	attrFn := func() error {
		// Clear existing relationships
		_, err = db.NewDelete().
			Model(assets.EmpPostToTag).
			Where("post_id = ?", post.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDelete().
			Model(assets.EmpPostToPost).
			Where("source_post_id = ?", post.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDelete().
			Model(assets.EmpPostToProject).
			Where("post_id = ?", post.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Create tag relationships
		for _, tagSlug := range post.TagSlugs {
			slog.Debug("tag rel tag", "tagSlug", tagSlug)
			// Find the tag ID by its slug
			relatedTag, err := FindTagBySlug(ctx, db, tagSlug, post.Slug)
			if err != nil {
				return err
			}

			// Create relationship
			_, err = db.NewInsert().
				Model(&assets.PostToTag{
					PostID: post.ID,
					TagID:  relatedTag.ID,
				}).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		// Create post relationships
		for _, postSlug := range post.PostSlugs {
			slog.Debug("post rel post", "postSlug", postSlug)
			// Find the related post ID by its slug
			relatedPost, err := FindPostBySlug(ctx, db, postSlug, post.Slug)
			if err != nil {
				return err
			}

			// Create relationship
			_, err = db.NewInsert().
				Model(&assets.PostToPost{
					SourcePostID: post.ID,
					TargetPostID: relatedPost.ID,
				}).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		// Create project relationships
		for _, projectSlug := range post.ProjectSlugs {
			slog.Debug("project rel post", "projectSlug", projectSlug)
			// Find the project ID by its slug
			relatedProject, err := FindProjectBySlug(ctx, db, projectSlug, post.Slug)
			if err != nil {
				return err
			}

			// Create relationship
			_, err = db.NewInsert().
				Model(&assets.PostToProject{
					PostID:    post.ID,
					ProjectID: relatedProject.ID,
				}).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}

	go msgCh.CtxSend(ctx, Msg{
		Path:  post.Slug,
		Type:  MsgTypeAction,
		fn:    attrFn,
		Tries: 0,
	}, wg)

	return nil
}

// SaveProject saves a project and its relationships to the database
func SaveProject(
	ctx context.Context,
	db *bun.DB,
	project *assets.Project,
	msgCh MsgChannel,
	wg *sync.WaitGroup,
) error {
	slog.Info("saving project", slog.String("hashed", project.Hash))
	// Save the project first
	_, err := db.NewInsert().
		Model(project).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("slug = EXCLUDED.slug").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("raw_content = EXCLUDED.raw_content").
		Set("icon = EXCLUDED.icon").
		Set("created_at = EXCLUDED.created_at").
		Set("updated_at = EXCLUDED.updated_at").
		Set("hashed = EXCLUDED.hashed").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(ctx)
	if err != nil {
		return err
	}

	attrFn := func() error {
		// Clear existing relationships
		_, err = db.NewDelete().
			Model(assets.EmpProjectToTag).
			Where("project_id = ?", project.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDelete().
			Model(assets.EmpPostToProject).
			Where("project_id = ?", project.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDelete().
			Model(assets.EmpProjectToProject).
			Where("source_project_id = ?", project.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Create tag relationships
		for _, tagSlug := range project.TagSlugs {
			slog.Debug("tag rel project", "tagSlug", tagSlug)
			// Find the tag ID by its slug
			relatedTag, err := FindTagBySlug(ctx, db, tagSlug, project.Slug)
			if err != nil {
				return err
			}

			// Create relationship
			_, err = db.NewInsert().
				Model(&assets.ProjectToTag{
					ProjectID: project.ID,
					TagID:     relatedTag.ID,
				}).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		// Create post relationships
		for _, postSlug := range project.PostSlugs {
			slog.Debug("post rel project", "postSlug", postSlug)
			// Find the related post ID by its slug
			relatedPost, err := FindPostBySlug(ctx, db, postSlug, project.Slug)
			if err != nil {
				return err
			}

			// Create relationship
			_, err = db.NewInsert().
				Model(&assets.PostToProject{
					ProjectID: project.ID,
					PostID:    relatedPost.ID,
				}).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		// Create project relationships
		for _, projectSlug := range project.ProjectSlugs {
			slog.Debug("project rel project", "projectSlug", projectSlug)
			// Find the related project ID by its slug
			relatedProject, err := FindProjectBySlug(ctx, db, projectSlug, project.Slug)
			if err != nil {
				return err
			}
			_, err = db.NewInsert().
				Model(&assets.ProjectToProject{
					SourceProjectID: project.ID,
					TargetProjectID: relatedProject.ID,
				}).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}

	go msgCh.CtxSend(ctx, Msg{
		Type:  MsgTypeAction,
		fn:    attrFn,
		Tries: 0,
		Path:  project.Slug,
	}, wg)

	return nil
}

// SaveTag saves a tag and its relationships to the database
func SaveTag(
	ctx context.Context,
	db *bun.DB,
	tag *assets.Tag,
	msgCh MsgChannel,
	wg *sync.WaitGroup,
) error {
	// Save the tag first
	_, err := db.NewInsert().
		Model(tag).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("slug = EXCLUDED.slug").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("raw_content = EXCLUDED.raw_content").
		Set("icon = EXCLUDED.icon").
		Set("created_at = EXCLUDED.created_at").
		Set("updated_at = EXCLUDED.updated_at").
		Set("hashed = EXCLUDED.hashed").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(ctx)
	if err != nil {
		return err
	}

	attrFn := func() error {
		// Clear existing relationships
		_, err = db.NewDelete().
			Model(assets.EmpPostToTag).
			Where("tag_id = ?", tag.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDelete().
			Model(assets.EmpProjectToTag).
			Where("tag_id = ?", tag.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewDelete().
			Model(assets.EmpTagToTag).
			Where("source_tag_id = ?", tag.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Create tag relationships
		for _, tagSlug := range tag.TagSlugs {
			slog.Debug("tag rel tag", "tagSlug", tagSlug)
			// Find the related tag ID by its slug
			relatedTag, err := FindTagBySlug(ctx, db, tagSlug, tag.Slug)
			if err != nil {
				return err
			}

			// Create relationship
			_, err = db.NewInsert().
				Model(&assets.TagToTag{
					SourceTagID: tag.ID,
					TargetTagID: relatedTag.ID,
				}).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		// Create post relationships
		for _, postSlug := range tag.PostSlugs {
			slog.Debug("post rel tag", "postSlug", postSlug)
			// Find the post ID by its slug
			relatedPost, err := FindPostBySlug(ctx, db, postSlug, tag.Slug)
			if err != nil {
				return err
			}

			// Create relationship
			_, err = db.NewInsert().
				Model(&assets.PostToTag{
					TagID:  tag.ID,
					PostID: relatedPost.ID,
				}).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		// Create project relationships
		for _, projectSlug := range tag.ProjectSlugs {
			slog.Debug("project rel tag", "projectSlug", projectSlug)
			// Find the project ID by its slug
			relatedProject, err := FindProjectBySlug(ctx, db, projectSlug, tag.Slug)
			if err != nil {
				return err
			}

			// Create relationship
			_, err = db.NewInsert().
				Model(&assets.ProjectToTag{
					TagID:     tag.ID,
					ProjectID: relatedProject.ID,
				}).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}

	go msgCh.CtxSend(ctx, Msg{
		Path:  tag.Slug,
		fn:    attrFn,
		Tries: 0,
		Type:  MsgTypeAction,
	}, wg)

	return nil
}

// FindTagBySlug finds a tag by its slug.
func FindTagBySlug(
	ctx context.Context,
	db *bun.DB,
	slug, origin string,
) (*assets.Tag, error) {
	var relatedTag assets.Tag
	err := db.NewSelect().
		Model(&relatedTag).
		Where("slug = ?", slug).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrapf(err, "failed to find tag by slug %s while processing %s", slug, origin)
	}
	if err != nil {
		return nil, eris.Wrap(err, "failed to find tag by slug")
	}
	return &relatedTag, nil
}

// FindPostBySlug finds a post by its slug.
func FindPostBySlug(
	ctx context.Context,
	db *bun.DB,
	slug, origin string,
) (*assets.Post, error) {
	var relatedPost assets.Post
	err := db.NewSelect().
		Model(&relatedPost).
		Where("slug = ?", slug).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrapf(err, "failed to find post by slug %s while processing %s", slug, origin)
	}
	if err != nil {
		return nil, eris.Wrap(err, "failed to find post by slug")
	}
	return &relatedPost, nil
}

// FindProjectBySlug finds a project by its slug.
func FindProjectBySlug(ctx context.Context, db *bun.DB, slug, origin string) (*assets.Project, error) {
	var relatedProject assets.Project
	err := db.NewSelect().
		Model(&relatedProject).
		Where("slug = ?", slug).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrapf(err, "failed to find project by slug %s while processing %s", slug, origin)
	}
	if err != nil {
		return nil, err
	}
	return &relatedProject, nil
}

// SaveAsset saves an asset media to the database.
func SaveAsset(
	ctx context.Context,
	db *bun.DB,
	path string,
	asset *assets.Asset,
	msgCh MsgChannel,
	wg *sync.WaitGroup,
) error {
	slog.Debug("saving asset media", "path", path)
	defer slog.Debug("saved asset media", "path", path)
	_, err := db.NewInsert().
		Model(asset).
		On("CONFLICT (path) DO UPDATE").
		Set("path = EXCLUDED.path").
		Set("hashed = EXCLUDED.hashed").
		Set("x = EXCLUDED.x").
		Set("y = EXCLUDED.y").
		Set("z = EXCLUDED.z").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Save saves a document to the database.
func Save(
	ctx context.Context,
	db *bun.DB,
	path string,
	doc *assets.Doc,
	msgCh MsgChannel,
	wg *sync.WaitGroup,
) error {
	slog.Debug("saving document", "path", path)
	defer slog.Debug("saved document", "path", path)
	switch {
	case strings.HasPrefix(path, assets.PostsLoc) && strings.HasSuffix(path, ".md"):
		var post assets.Post
		copygen.ToPost(&post, doc)
		return SavePost(ctx, db, &post, msgCh, wg)
	case strings.HasPrefix(path, assets.ProjectsLoc) && strings.HasSuffix(path, ".md"):
		var project assets.Project
		copygen.ToProject(&project, doc)
		return SaveProject(ctx, db, &project, msgCh, wg)
	case strings.HasPrefix(path, assets.TagsLoc) && strings.HasSuffix(path, ".md"):
		var tag assets.Tag
		copygen.ToTag(&tag, doc)
		return SaveTag(ctx, db, &tag, msgCh, wg)
	}
	return nil
}
