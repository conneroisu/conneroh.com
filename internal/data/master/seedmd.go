package master

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	"github.com/conneroisu/conneroh.com/doc"
	"github.com/conneroisu/conneroh.com/internal/md"
)

// SeedFromEmbedded reads the embedded doc files and seeds the database with them.
func (q *Queries) SeedFromEmbedded(ctx context.Context) error {
	slog.Info("Seeding database from embedded markdown files")

	// Seed posts
	if err := q.seedPosts(ctx, doc.Posts); err != nil {
		return fmt.Errorf("error seeding posts: %w", err)
	}

	// Seed projects
	if err := q.seedProjects(ctx, doc.Projects); err != nil {
		return fmt.Errorf("error seeding projects: %w", err)
	}

	// Seed tags
	if err := q.seedTags(ctx, doc.Tags); err != nil {
		return fmt.Errorf("error seeding tags: %w", err)
	}

	return nil
}

// seedPosts reads post markdown files and seeds the database.
func (q *Queries) seedPosts(ctx context.Context, postsFS fs.FS) error {
	return fs.WalkDir(postsFS, "posts", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Read file content
		content, err := fs.ReadFile(postsFS, path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		// Parse frontmatter
		var frontMatter md.PostFrontMatter

		htmlContent, err := md.ParseWithFrontMatter(path, content, &frontMatter)
		if err != nil {
			return fmt.Errorf("error parsing frontmatter for %s: %w", path, err)
		}

		post, err := q.PostCreate(ctx, PostCreateParams{
			Title:       frontMatter.Title,
			Description: frontMatter.Description,
			Slug:        frontMatter.Slug,
			Content:     string(htmlContent),
			BannerUrl:   "/assets/img/posts/default.jpg",
		})
		if err != nil {
			// If post already exists (e.g., unique constraint violation), just log and continue
			slog.Warn("Failed to create post", "title", frontMatter.Title, "error", err)
			return nil
		}

		// Associate tags with post
		for _, tagName := range frontMatter.Tags {
			// Try to find the tag by name
			tag, err := q.TagGetByName(ctx, tagName)
			if err != nil {
				// Create the tag if it doesn't exist
				tag, err = q.TagCreate(ctx, TagCreateParams{
					Name:        tagName,
					Description: "Auto-generated tag from post",
					Slug:        tagName,
				})
				if err != nil {
					slog.Warn("Failed to create tag", "name", tagName, "error", err)
					continue
				}
			}

			// Associate tag with post
			if err := q.PostTagCreate(ctx, post.ID, tag.ID); err != nil {
				slog.Warn("Failed to associate tag with post", "postID", post.ID, "tagID", tag.ID, "error", err)
			}
		}

		slog.Info("Created post", "title", frontMatter.Title, "id", post.ID)
		return nil
	})
}

// seedProjects reads project markdown files and seeds the database.
func (q *Queries) seedProjects(ctx context.Context, projectsFS fs.FS) error {
	return fs.WalkDir(projectsFS, "projects", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Read file content
		content, err := fs.ReadFile(projectsFS, path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		// Parse frontmatter
		var frontMatter md.ProjectFrontMatter
		_, err = md.ParseWithFrontMatter(path, content, &frontMatter)
		if err != nil {
			return fmt.Errorf("error parsing frontmatter for %s: %w", path, err)
		}
		project, err := q.ProjectCreate(ctx, ProjectCreateParams{
			Name:        frontMatter.Title,
			Slug:        frontMatter.Slug,
			Description: frontMatter.Description,
		})
		if err != nil {
			// If project already exists, just log and continue
			slog.Warn("Failed to create project", "name", frontMatter.Title, "error", err)
			return nil
		}

		// Associate tags with project
		for _, tagName := range frontMatter.Tags {
			// Try to find the tag by name
			tag, err := q.TagGetByName(ctx, tagName)
			if err != nil {
				// Create the tag if it doesn't exist
				tag, err = q.TagCreate(ctx, TagCreateParams{
					Name:        tagName,
					Description: "Auto-generated tag from project",
					Slug:        tagName,
				})
				if err != nil {
					slog.Warn("Failed to create tag", "name", tagName, "error", err)
					continue
				}
			}

			// Associate tag with project
			if err := q.ProjectTagCreate(ctx, project.ID, tag.ID); err != nil {
				slog.Warn("Failed to associate tag with project", "projectID", project.ID, "tagID", tag.ID, "error", err)
			}
		}

		slog.Info("Created project", "name", frontMatter.Title, "id", project.ID)
		return nil
	})
}

// seedTags reads tag markdown files and seeds the database.
func (q *Queries) seedTags(ctx context.Context, tagsFS fs.FS) error {
	return fs.WalkDir(tagsFS, "tags", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Read file content
		content, err := fs.ReadFile(tagsFS, path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		// Parse frontmatter
		var frontMatter md.TagFrontMatter
		_, err = md.ParseWithFrontMatter(path, content, &frontMatter)
		if err != nil {
			return fmt.Errorf("error parsing frontmatter for %s: %w", path, err)
		}

		// Create tag
		tag, err := q.TagCreate(ctx, TagCreateParams{
			Name:        frontMatter.Title,
			Description: frontMatter.Description,
			Slug:        frontMatter.Slug,
		})
		if err != nil {
			// If tag already exists, just log and continue
			slog.Warn("Failed to create tag", "name", frontMatter.Title, "error", err)
			return nil
		}

		slog.Info("Created tag", "name", frontMatter.Title, "id", tag.ID)
		return nil
	})
}
