package assets

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/rotisserie/eris"
	"github.com/uptrace/bun"
)

// RelationshipFn is a function that updates relationships.
type RelationshipFn func(context.Context) error

// UpsertPost saves a post to the database (to be called from the DB worker).
func UpsertPost(
	ctx context.Context,
	db *bun.DB,
	post *Post,
) (RelationshipFn, error) {
	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := db.NewInsert().
		Model(post).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("created_at = EXCLUDED.created_at").
		Exec(insertCtx)
	slog.Debug(
		"saved post",
		slog.String("slug", post.Slug),
		slog.String("banner_path", post.BannerPath),
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, eris.Wrapf(
				err,
				"database operation timed out for post: %s",
				post.Slug,
			)
		}

		return nil, eris.Wrapf(err, "failed to save post: %s", post.Slug)
	}

	return UpsertPostRelationships(db, post), nil
}

// UpsertProject saves a project to the database (to be called from the DB worker).
func UpsertProject(
	ctx context.Context,
	db *bun.DB,
	project *Project,
) (RelationshipFn, error) {
	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := db.NewInsert().
		Model(project).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("created_at = EXCLUDED.created_at").
		Exec(insertCtx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, eris.Wrapf(
				err,
				"database operation timed out for project: %s",
				project.Slug,
			)
		}

		return nil, eris.Wrapf(err, "failed to save project: %s", project.Slug)
	}

	return UpsertProjectRelationships(db, project), nil
}

// UpsertTag saves a tag to the database (to be called from the DB worker).
func UpsertTag(
	ctx context.Context,
	db *bun.DB,
	tag *Tag,
) (RelationshipFn, error) {
	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := db.NewInsert().
		Model(tag).
		On("CONFLICT (slug) DO UPDATE").
		Set("title = EXCLUDED.title").
		Set("description = EXCLUDED.description").
		Set("content = EXCLUDED.content").
		Set("banner_path = EXCLUDED.banner_path").
		Set("created_at = EXCLUDED.created_at").
		Exec(insertCtx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, eris.Wrapf(
				err,
				"database operation timed out for tag: %s",
				tag.Slug,
			)
		}

		return nil, eris.Wrapf(err, "failed to save tag: %s", tag.Slug)
	}

	return UpsertTagRelationships(db, tag), nil
}

// UpsertPostRelationships updates relationships for a post (to be called from the DB worker).
func UpsertPostRelationships(
	db *bun.DB,
	post *Post,
) RelationshipFn {
	return func(ctx context.Context) error {
		var relatedTag Tag
		for _, tagSlug := range post.TagSlugs {
			err := db.NewSelect().
				Model(&relatedTag).
				Where("slug = ?", tagSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for tag: %s",
						tagSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"tag not found: %s (referenced from %s)",
						tagSlug,
						post.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find tag: %s",
					tagSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&PostToTag{
					PostID: post.ID,
					TagID:  relatedTag.ID,
				}).
				On("CONFLICT (post_id, tag_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create post-tag relationship: %s -> %s",
					post.Slug,
					tagSlug,
				)
			}
		}

		var relatedPost Post
		for _, postSlug := range post.PostSlugs {
			err := db.NewSelect().
				Model(&relatedPost).
				Where("slug = ?", postSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for post: %s",
						postSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"post not found: %s (referenced from %s)",
						postSlug,
						post.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find post: %s",
					postSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&PostToPost{
					SourcePostID: post.ID,
					TargetPostID: relatedPost.ID,
				}).
				On("CONFLICT (source_post_id, target_post_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create post-post relationship: %s -> %s",
					post.Slug,
					postSlug,
				)
			}
		}

		var relatedProject Project
		for _, projectSlug := range post.ProjectSlugs {
			err := db.NewSelect().
				Model(&relatedProject).
				Where("slug = ?", projectSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for project: %s",
						projectSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"project not found: %s (referenced from %s)",
						projectSlug,
						post.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find project: %s",
					projectSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&PostToProject{
					PostID:    post.ID,
					ProjectID: relatedProject.ID,
				}).
				On("CONFLICT (post_id, project_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create post-project relationship: %s -> %s",
					post.Slug,
					projectSlug,
				)
			}
		}

		return nil
	}
}

// UpsertProjectRelationships updates relationships for a project.
func UpsertProjectRelationships(
	db *bun.DB,
	project *Project,
) RelationshipFn {
	return func(ctx context.Context) error {
		var relatedTag Tag
		for _, tagSlug := range project.TagSlugs { // Create tag relationships
			err := db.NewSelect().
				Model(&relatedTag).
				Where("slug = ?", tagSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for tag: %s",
						tagSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"tag not found: %s (referenced from %s)",
						tagSlug,
						project.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find tag: %s",
					tagSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&ProjectToTag{
					ProjectID: project.ID,
					TagID:     relatedTag.ID,
				}).
				On("CONFLICT (project_id, tag_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create project-tag relationship: %s -> %s",
					project.Slug,
					tagSlug,
				)
			}
		}

		var relatedPost Post
		for _, postSlug := range project.PostSlugs { // Create post relationships
			err := db.NewSelect().
				Model(&relatedPost).
				Where("slug = ?", postSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for post: %s",
						postSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"post not found: %s (referenced from %s)",
						postSlug,
						project.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find post: %s",
					postSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&PostToProject{
					ProjectID: project.ID,
					PostID:    relatedPost.ID,
				}).
				On("CONFLICT (project_id, post_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create project-post relationship: %s -> %s",
					project.Slug,
					postSlug,
				)
			}
		}

		var relatedProject Project
		for _, projectSlug := range project.ProjectSlugs { // Create project relationships
			err := db.NewSelect().
				Model(&relatedProject).
				Where("slug = ?", projectSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for project: %s",
						projectSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"project not found: %s (referenced from %s)",
						projectSlug,
						project.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find project: %s",
					projectSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&ProjectToProject{
					SourceProjectID: project.ID,
					TargetProjectID: relatedProject.ID,
				}).
				On("CONFLICT (source_project_id, target_project_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create project-project relationship: %s -> %s",
					project.Slug,
					projectSlug,
				)
			}
		}

		return nil
	}
}

// UpsertTagRelationships updates relationships for a tag .
func UpsertTagRelationships(
	db *bun.DB,
	tag *Tag,
) RelationshipFn {
	return func(ctx context.Context) error {
		var relatedTag Tag
		for _, tagSlug := range tag.TagSlugs { // Create tag relationships
			err := db.NewSelect().
				Model(&relatedTag).
				Where("slug = ?", tagSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for tag: %s",
						tagSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"tag not found: %s (referenced from %s)",
						tagSlug,
						tag.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find tag: %s",
					tagSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&TagToTag{
					SourceTagID: tag.ID,
					TargetTagID: relatedTag.ID,
				}).
				On("CONFLICT (source_tag_id, target_tag_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create tag-tag relationship: %s -> %s",
					tag.Slug,
					tagSlug,
				)
			}
		}

		var relatedPost Post
		for _, postSlug := range tag.PostSlugs { // Create post relationships
			err := db.NewSelect().
				Model(&relatedPost).
				Where("slug = ?", postSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for post: %s",
						postSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"post not found: %s (referenced from %s)",
						postSlug,
						tag.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find post: %s",
					postSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&PostToTag{
					TagID:  tag.ID,
					PostID: relatedPost.ID,
				}).
				On("CONFLICT (tag_id, post_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create tag-post relationship: %s -> %s",
					tag.Slug,
					postSlug,
				)
			}
		}

		var relatedProject Project
		for _, projectSlug := range tag.ProjectSlugs { // Create project relationships
			err := db.NewSelect().
				Model(&relatedProject).
				Where("slug = ?", projectSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for project: %s",
						projectSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"project not found: %s (referenced from %s)",
						projectSlug,
						tag.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find project: %s",
					projectSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&ProjectToTag{
					TagID:     tag.ID,
					ProjectID: relatedProject.ID,
				}).
				On("CONFLICT (tag_id, project_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create tag-project relationship: %s -> %s",
					tag.Slug,
					projectSlug,
				)
			}
		}

		return nil
	}
}

// UpsertEmployment saves an employment to the database (to be called from the DB worker).
func UpsertEmployment(
	ctx context.Context,
	db *bun.DB,
	employment *Employment,
) (RelationshipFn, error) {
	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := db.NewInsert().
		Model(employment).
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
	slog.Debug(
		"saved employment",
		slog.String("slug", employment.Slug),
		slog.String("banner_path", employment.BannerPath),
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, eris.Wrapf(
				err,
				"database operation timed out for employment: %s",
				employment.Slug,
			)
		}

		return nil, eris.Wrapf(err, "failed to save employment: %s", employment.Slug)
	}

	return UpsertEmploymentRelationships(db, employment), nil
}

// UpsertEmploymentRelationships updates relationships for an employment (to be called from the DB worker).
func UpsertEmploymentRelationships(
	db *bun.DB,
	employment *Employment,
) RelationshipFn {
	return func(ctx context.Context) error {
		var relatedTag Tag
		for _, tagSlug := range employment.TagSlugs {
			err := db.NewSelect().
				Model(&relatedTag).
				Where("slug = ?", tagSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for tag: %s",
						tagSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"tag not found: %s (referenced from %s)",
						tagSlug,
						employment.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find tag: %s",
					tagSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&EmploymentToTag{
					EmploymentID: employment.ID,
					TagID:        relatedTag.ID,
				}).
				On("CONFLICT (employment_id, tag_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create employment-tag relationship: %s -> %s",
					employment.Slug,
					tagSlug,
				)
			}
		}

		var relatedPost Post
		for _, postSlug := range employment.PostSlugs {
			err := db.NewSelect().
				Model(&relatedPost).
				Where("slug = ?", postSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for post: %s",
						postSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"post not found: %s (referenced from %s)",
						postSlug,
						employment.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find post: %s",
					postSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&EmploymentToPost{
					EmploymentID: employment.ID,
					PostID:       relatedPost.ID,
				}).
				On("CONFLICT (employment_id, post_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create employment-post relationship: %s -> %s",
					employment.Slug,
					postSlug,
				)
			}
		}

		var relatedProject Project
		for _, projectSlug := range employment.ProjectSlugs {
			err := db.NewSelect().
				Model(&relatedProject).
				Where("slug = ?", projectSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for project: %s",
						projectSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"project not found: %s (referenced from %s)",
						projectSlug,
						employment.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find project: %s",
					projectSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&EmploymentToProject{
					EmploymentID: employment.ID,
					ProjectID:    relatedProject.ID,
				}).
				On("CONFLICT (employment_id, project_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create employment-project relationship: %s -> %s",
					employment.Slug,
					projectSlug,
				)
			}
		}

		var relatedEmployment Employment
		for _, employmentSlug := range employment.EmploymentSlugs {
			err := db.NewSelect().
				Model(&relatedEmployment).
				Where("slug = ?", employmentSlug).
				Scan(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return eris.Wrapf(
						err,
						"database operation timed out for employment: %s",
						employmentSlug,
					)
				}
				if errors.Is(err, sql.ErrNoRows) {
					return eris.Errorf(
						"employment not found: %s (referenced from %s)",
						employmentSlug,
						employment.Slug,
					)
				}

				return eris.Wrapf(
					err,
					"failed to find employment: %s",
					employmentSlug,
				)
			}
			_, err = db.NewInsert().
				Model(&EmploymentToEmployment{
					SourceEmploymentID: employment.ID,
					TargetEmploymentID: relatedEmployment.ID,
				}).
				On("CONFLICT (source_employment_id, target_employment_id) DO NOTHING").
				Exec(ctx)
			if err != nil {
				return eris.Wrapf(
					err,
					"failed to create employment-employment relationship: %s -> %s",
					employment.Slug,
					employmentSlug,
				)
			}
		}

		return nil
	}
}
