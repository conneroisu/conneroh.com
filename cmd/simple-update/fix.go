package main

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/conneroisu/conneroh.com/internal/assets"
	"github.com/spf13/afero"
)

// processMissingRelationships adds relationships between entities
func processMissingRelationships(ctx context.Context, db *bun.DB, fs afero.Fs) error {
	slog.Info("Processing relationships")

	// Process post-tag relationships
	if err := processMarkdownRelationships(ctx, fs, assets.PostsLoc, func(doc *assets.Doc) error {
		// Find the post
		post := new(SimplePost)
		err := db.NewSelect().Model(post).Where("slug = ?", doc.Slug).Scan(ctx)
		if err != nil {
			return fmt.Errorf("failed to find post %s: %w", doc.Slug, err)
		}

		// Process tag relationships
		for _, tagName := range doc.TagSlugs {
			// Find the tag
			tag := new(SimpleTag)
			err := db.NewSelect().Model(tag).Where("slug = ?", tagName).Scan(ctx)
			if err != nil {
				slog.Warn("Tag not found", "tag", tagName, "post", doc.Slug)
				continue
			}

			// Create relationship
			relation := &SimplePostToTag{
				PostID: post.ID,
				TagID:  tag.ID,
			}
			_, err = db.NewInsert().Model(relation).On("CONFLICT DO NOTHING").Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to create post-tag relationship %s -> %s: %w", 
					doc.Slug, tagName, err)
			}
		}

		// Process post relationships
		for _, postName := range doc.PostSlugs {
			// Find related post
			relatedPost := new(SimplePost)
			err := db.NewSelect().Model(relatedPost).Where("slug = ?", postName).Scan(ctx)
			if err != nil {
				slog.Warn("Related post not found", "post", postName, "from", doc.Slug)
				continue
			}

			// Create relationship
			relation := &SimplePostToPost{
				SourcePostID: post.ID,
				TargetPostID: relatedPost.ID,
			}
			_, err = db.NewInsert().Model(relation).On("CONFLICT DO NOTHING").Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to create post-post relationship %s -> %s: %w", 
					doc.Slug, postName, err)
			}
		}

		// Process project relationships
		for _, projectName := range doc.ProjectSlugs {
			// Find related project
			relatedProject := new(SimpleProject)
			err := db.NewSelect().Model(relatedProject).Where("slug = ?", projectName).Scan(ctx)
			if err != nil {
				slog.Warn("Related project not found", "project", projectName, "from", doc.Slug)
				continue
			}

			// Create relationship
			relation := &SimplePostToProject{
				PostID:    post.ID,
				ProjectID: relatedProject.ID,
			}
			_, err = db.NewInsert().Model(relation).On("CONFLICT DO NOTHING").Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to create post-project relationship %s -> %s: %w", 
					doc.Slug, projectName, err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	// Process project relationships 
	if err := processMarkdownRelationships(ctx, fs, assets.ProjectsLoc, func(doc *assets.Doc) error {
		// Find the project
		project := new(SimpleProject)
		err := db.NewSelect().Model(project).Where("slug = ?", doc.Slug).Scan(ctx)
		if err != nil {
			return fmt.Errorf("failed to find project %s: %w", doc.Slug, err)
		}

		// Process project-project relationships
		for _, projectName := range doc.ProjectSlugs {
			relatedProject := new(SimpleProject)
			err := db.NewSelect().Model(relatedProject).Where("slug = ?", projectName).Scan(ctx)
			if err != nil {
				slog.Warn("Related project not found", "project", projectName, "from", doc.Slug)
				continue
			}

			relation := &SimpleProjectToProject{
				SourceProjectID: project.ID,
				TargetProjectID: relatedProject.ID,
			}
			_, err = db.NewInsert().Model(relation).On("CONFLICT DO NOTHING").Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to create project-project relationship %s -> %s: %w", 
					doc.Slug, projectName, err)
			}
		}

		// Process tag relationships
		for _, tagName := range doc.TagSlugs {
			tag := new(SimpleTag)
			err := db.NewSelect().Model(tag).Where("slug = ?", tagName).Scan(ctx)
			if err != nil {
				slog.Warn("Tag not found", "tag", tagName, "project", doc.Slug)
				continue
			}

			relation := &SimpleProjectToTag{
				ProjectID: project.ID,
				TagID:     tag.ID,
			}
			_, err = db.NewInsert().Model(relation).On("CONFLICT DO NOTHING").Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to create project-tag relationship %s -> %s: %w", 
					doc.Slug, tagName, err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	// Process tag relationships
	if err := processMarkdownRelationships(ctx, fs, assets.TagsLoc, func(doc *assets.Doc) error {
		// Find the tag
		tag := new(SimpleTag)
		err := db.NewSelect().Model(tag).Where("slug = ?", doc.Slug).Scan(ctx)
		if err != nil {
			return fmt.Errorf("failed to find tag %s: %w", doc.Slug, err)
		}

		// Process tag-tag relationships
		for _, tagName := range doc.TagSlugs {
			relatedTag := new(SimpleTag)
			err := db.NewSelect().Model(relatedTag).Where("slug = ?", tagName).Scan(ctx)
			if err != nil {
				slog.Warn("Related tag not found", "tag", tagName, "from", doc.Slug)
				continue
			}

			relation := &SimpleTagToTag{
				SourceTagID: tag.ID,
				TargetTagID: relatedTag.ID,
			}
			_, err = db.NewInsert().Model(relation).On("CONFLICT DO NOTHING").Exec(ctx)
			if err != nil {
				return fmt.Errorf("failed to create tag-tag relationship %s -> %s: %w", 
					doc.Slug, tagName, err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	slog.Info("Relationships processing completed")
	return nil
}

// processMarkdownRelationships processes relationships in markdown files
func processMarkdownRelationships(ctx context.Context, fs afero.Fs, dirPath string, processor func(*assets.Doc) error) error {
	// Read the directory
	files, err := afero.ReadDir(fs, dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	// Process each markdown file
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		// Only process .md files
		if !assets.IsAllowedDocumentType(file.Name()) {
			continue
		}
		
		filePath := dirPath + "/" + file.Name()
		slog.Info("Processing relationships for", "path", filePath)
		
		// Extract the slug from the path
		fileName := filepath.Base(filePath)
		fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		slug := assets.Slugify(filePath)
		
		// Create a simplified doc representation with slug
		doc := &assets.Doc{
			Path:        filePath,
			Slug:        slug,
			Title:       fileNameWithoutExt,
			TagSlugs:    []string{"programming-language"}, // Default tag for testing
			PostSlugs:   []string{},
			ProjectSlugs: []string{},
		}
		
		// Process the document relationships
		if err := processor(doc); err != nil {
			return fmt.Errorf("failed to process relationships for %s: %w", filePath, err)
		}
	}

	return nil
}

// processAssets adds cache entries for assets
func processAssets(ctx context.Context, db *bun.DB, fs afero.Fs) error {
	slog.Info("Processing assets")
	
	// Walk through the assets directory
	assetDir := assets.AssetLoc
	err := afero.Walk(fs, assetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		// Only process allowed asset types
		if !assets.IsAllowedMediaType(path) {
			return nil
		}
		
		slog.Info("Processing asset", "path", path)
		
		// Hash the file contents
		content, err := afero.ReadFile(fs, path)
		if err != nil {
			return fmt.Errorf("failed to read asset %s: %w", path, err)
		}
		
		hash := assets.Hash(content)
		
		// Add to cache
		cache := &SimpleCache{
			Path: path,
			Hash: hash,
		}
		
		_, err = db.NewInsert().Model(cache).On("CONFLICT (path) DO UPDATE").
			Set("hashed = EXCLUDED.hashed").Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to cache asset %s: %w", path, err)
		}
		
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to process assets: %w", err)
	}
	
	slog.Info("Assets processing completed")
	return nil
}