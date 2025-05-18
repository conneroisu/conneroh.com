package cache

import (
	"context"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/rotisserie/eris"
)

// dbUpdatePostRelationships updates relationships for a post (to be called from the DB worker).
func (p *Processor) dbUpdatePostRelationships(ctx context.Context, post *assets.Post) error {
	for _, tagSlug := range post.TagSlugs { // Create tag relationships
		tag, err := p.dbFindTagBySlug(ctx, tagSlug, post.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.PostToTag{
				PostID: post.ID,
				TagID:  tag.ID,
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

	for _, postSlug := range post.PostSlugs { // Create post relationships
		relatedPost, err := p.dbFindPostBySlug(ctx, postSlug, post.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.PostToPost{
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

	for _, projectSlug := range post.ProjectSlugs { // Create project relationships
		project, err := p.dbFindProjectBySlug(ctx, projectSlug, post.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.PostToProject{
				PostID:    post.ID,
				ProjectID: project.ID,
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

// dbUpdateProjectRelationships updates relationships for a project (to be called from the DB worker).
func (p *Processor) dbUpdateProjectRelationships(ctx context.Context, project *assets.Project) error {
	for _, projectSlug := range project.ProjectSlugs { // Create project relationships
		relatedProject, err := p.dbFindProjectBySlug(ctx, projectSlug, project.Slug)
		if err != nil {
			return err
		}
		_, err = p.db.NewInsert().
			Model(&assets.ProjectToProject{
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

	for _, tagSlug := range project.TagSlugs { // Create tag relationships
		tag, err := p.dbFindTagBySlug(ctx, tagSlug, project.Slug)
		if err != nil {
			return err
		}
		_, err = p.db.NewInsert().
			Model(&assets.ProjectToTag{
				ProjectID: project.ID,
				TagID:     tag.ID,
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

	for _, postSlug := range project.PostSlugs { // Create post relationships
		post, err := p.dbFindPostBySlug(ctx, postSlug, project.Slug)
		if err != nil {
			return err
		}
		_, err = p.db.NewInsert().
			Model(&assets.PostToProject{
				ProjectID: project.ID,
				PostID:    post.ID,
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

	return nil
}

// dbUpdateTagRelationships updates relationships for a tag (to be called from the DB worker).
func (p *Processor) dbUpdateTagRelationships(ctx context.Context, tag *assets.Tag) error {
	for _, tagSlug := range tag.TagSlugs { // Create tag relationships
		relatedTag, err := p.dbFindTagBySlug(ctx, tagSlug, tag.Slug)
		if err != nil {
			return err
		}
		_, err = p.db.NewInsert().
			Model(&assets.TagToTag{
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

	for _, postSlug := range tag.PostSlugs { // Create post relationships
		post, err := p.dbFindPostBySlug(ctx, postSlug, tag.Slug)
		if err != nil {
			return err
		}

		_, err = p.db.NewInsert().
			Model(&assets.PostToTag{
				TagID:  tag.ID,
				PostID: post.ID,
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

	for _, projectSlug := range tag.ProjectSlugs { // Create project relationships
		project, err := p.dbFindProjectBySlug(ctx, projectSlug, tag.Slug)
		if err != nil {
			return err
		}
		_, err = p.db.NewInsert().
			Model(&assets.ProjectToTag{
				TagID:     tag.ID,
				ProjectID: project.ID,
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
