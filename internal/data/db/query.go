package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/conneroisu/conneroh.com/internal/copygen"
	"github.com/rotisserie/eris"
	"github.com/uptrace/bun"
)

// GetPosts returns all posts.
// [bun]  09:04:40.056   CREATE TABLE            531µs  CREATE TABLE IF NOT EXISTS "assets" ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "path" VARCHAR, "hashed" VARCHAR, "x" DOUBLE PRECISION, "y" DOUBLE PRECISION, "z" DOUBLE PRECISION, UNIQUE ("path"), UNIQUE ("hashed"))
// [bun]  09:04:40.057   CREATE TABLE             31µs  CREATE TABLE IF NOT EXISTS "posts" ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "path" VARCHAR, "hashed" VARCHAR, "x" DOUBLE PRECISION, "y" DOUBLE PRECISION, "z" DOUBLE PRECISION, "title" VARCHAR, "slug" VARCHAR, "description" VARCHAR, "content" VARCHAR, "banner_path" VARCHAR, "raw_content" VARCHAR, "icon" VARCHAR, "created_at" VARCHAR, "updated_at" VARCHAR, "tag_slugs" VARCHAR, "post_slugs" VARCHAR, "project_slugs" VARCHAR, UNIQUE ("path"))
// [bun]  09:04:40.057   CREATE TABLE             13µs  CREATE TABLE IF NOT EXISTS "projects" ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "path" VARCHAR, "hashed" VARCHAR, "x" DOUBLE PRECISION, "y" DOUBLE PRECISION, "z" DOUBLE PRECISION,"title" VARCHAR, "slug" VARCHAR, "description" VARCHAR, "content" VARCHAR, "banner_path" VARCHAR, "raw_content" VARCHAR, "icon" VARCHAR, "created_at" VARCHAR, "updated_at" VARCHAR, "tag_slugs" VARCHAR, "post_slugs" VARCHAR, "project_slugs" VARCHAR, UNIQUE ("path"))
// [bun]  09:04:40.057   CREATE TABLE             14µs  CREATE TABLE IF NOT EXISTS "tags" ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "path" VARCHAR, "hashed" VARCHAR, "x" DOUBLE PRECISION, "y" DOUBLE PRECISION, "z" DOUBLE PRECISION, "title" VARCHAR, "slug" VARCHAR, "description" VARCHAR, "content" VARCHAR, "banner_path" VARCHAR, "raw_content" VARCHAR, "icon" VARCHAR, "created_at" VARCHAR, "updated_at" VARCHAR, "tag_slugs" VARCHAR, "post_slugs" VARCHAR, "project_slugs" VARCHAR, UNIQUE ("path"))
func GetPosts(ctx context.Context, db *bun.DB) ([]*assets.Post, error) {
	posts := []*assets.Post{}
	// TODO: Should join on the slugs that each post has.
	err := db.NewSelect().
		Model(&posts).
		Relation("Tags").
		Relation("Posts").
		Relation("Projects").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// SavePost saves a post and its relationships to the database
func SavePost(ctx context.Context, db *bun.DB, post *assets.Post) error {
	// Save the post first
	_, err := db.NewInsert().
		Model(post).
		On("CONFLICT (path) DO UPDATE").
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
		relatedTag, err := FindTagBySlug(ctx, db, tagSlug, post.Slug, 0)
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
		relatedPost, err := FindPostBySlug(ctx, db, postSlug, post.Slug, 0)
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
		relatedProject, err := FindProjectBySlug(ctx, db, projectSlug, post.Slug, 0)
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

// SaveProject saves a project and its relationships to the database
func SaveProject(ctx context.Context, db *bun.DB, project *assets.Project) error {
	slog.Info("saving project", slog.String("hashed", project.Hash))
	// Save the project first
	_, err := db.NewInsert().
		Model(project).
		On("CONFLICT (path) DO UPDATE").
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
		relatedTag, err := FindTagBySlug(ctx, db, tagSlug, project.Slug, 0)
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
		relatedPost, err := FindPostBySlug(ctx, db, postSlug, project.Slug, 0)
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
		relatedProject, err := FindProjectBySlug(ctx, db, projectSlug, project.Slug, 0)
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

// SaveTag saves a tag and its relationships to the database
func SaveTag(ctx context.Context, db *bun.DB, tag *assets.Tag) error {
	// Save the tag first
	_, err := db.NewInsert().
		Model(tag).
		On("CONFLICT (path) DO UPDATE").
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
		relatedTag, err := FindTagBySlug(ctx, db, tagSlug, tag.Slug, 0)
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
		relatedPost, err := FindPostBySlug(ctx, db, postSlug, tag.Slug, 0)
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
		relatedProject, err := FindProjectBySlug(ctx, db, projectSlug, tag.Slug, 0)
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

// FindTagBySlug finds a tag by its slug.
func FindTagBySlug(ctx context.Context, db *bun.DB, slug, origin string, tries int) (*assets.Tag, error) {
	var relatedTag assets.Tag
	err := db.NewSelect().
		Model(&relatedTag).
		Where("slug = ?", slug).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		if tries < 3 {
			tries++
			time.Sleep(time.Second)
			return FindTagBySlug(ctx, db, slug, origin, tries+1)
		}
		return nil, eris.Wrapf(err, "failed to find tag by slug %s while processing %s", slug, origin)
	}
	if err != nil {
		return nil, eris.Wrap(err, "failed to find tag by slug")
	}
	return &relatedTag, nil
}

// FindPostBySlug finds a post by its slug.
func FindPostBySlug(ctx context.Context, db *bun.DB, slug, origin string, tries int) (*assets.Post, error) {
	var relatedPost assets.Post
	err := db.NewSelect().
		Model(&relatedPost).
		Where("slug = ?", slug).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		if tries < 3 {
			tries++
			time.Sleep(time.Second)
			return FindPostBySlug(ctx, db, slug, origin, tries+1)
		}
		return nil, eris.Wrapf(err, "failed to find post by slug %s while processing %s", slug, origin)
	}
	if err != nil {
		return nil, eris.Wrap(err, "failed to find post by slug")
	}
	return &relatedPost, nil
}

// FindProjectBySlug finds a project by its slug.
func FindProjectBySlug(ctx context.Context, db *bun.DB, slug, origin string, tries int) (*assets.Project, error) {
	var relatedProject assets.Project
	err := db.NewSelect().
		Model(&relatedProject).
		Where("slug = ?", slug).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		if tries < 3 {
			tries++
			time.Sleep(time.Second)
			return FindProjectBySlug(ctx, db, slug, origin, tries+1)
		}
		return nil, eris.Wrapf(err, "failed to find project by slug %s while processing %s", slug, origin)
	}
	if err != nil {
		return nil, err
	}
	return &relatedProject, nil
}

// Save saves a document to the database.
func Save(ctx context.Context, db *bun.DB, path string, doc *assets.Doc) error {
	slog.Debug("saving document", "path", path)
	defer slog.Debug("saved document", "path", path)
	switch {
	case strings.HasPrefix(path, assets.PostsLoc) && strings.HasSuffix(path, ".md"):
		var post assets.Post
		copygen.ToPost(&post, doc)
		return SavePost(ctx, db, &post)
	case strings.HasPrefix(path, assets.ProjectsLoc) && strings.HasSuffix(path, ".md"):
		var project assets.Project
		copygen.ToProject(&project, doc)
		return SaveProject(ctx, db, &project)
	case strings.HasPrefix(path, assets.TagsLoc) && strings.HasSuffix(path, ".md"):
		var tag assets.Tag
		copygen.ToTag(&tag, doc)
		return SaveTag(ctx, db, &tag)
	}
	return nil
}
