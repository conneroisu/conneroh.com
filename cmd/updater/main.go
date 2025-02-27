// Package main is the updater package of the conneroh.com website.
package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	"bytes"

	"github.com/conneroisu/conneroh.com/doc"
	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/master"

	"github.com/go-playground/validator/v10"
	mathjax "github.com/litao91/goldmark-mathjax"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/quail-ink/goldmark-enclave/core"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/wikilink"
)

// FrontMatterMissingError is returned when the front matter is missing from the markdown file.
type FrontMatterMissingError struct {
	fileName string
}

// Error implements the error interface on FrontMatterMissingError.
func (e FrontMatterMissingError) Error() string {
	return "front matter missing from " + e.fileName
}

var (
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Strikethrough,
			extension.Table,
			extension.TaskList,
			extension.DefinitionList,
			extension.NewTypographer(
				extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
					extension.Apostrophe: []byte("’"),
				}),
			),
			enclave.New(&core.Config{
				DefaultImageAltPrefix: "caption: ",
			}),
			extension.NewFootnote(
				extension.WithFootnoteIDPrefix("fn"),
			),
			mathjax.MathJax,
			&anchor.Extender{
				Position: anchor.Before,
				Texter:   anchor.Text("#"),
				Attributer: anchor.Attributes{
					"class": "anchor permalink p-4",
				},
			},
			&mermaid.Extender{
				RenderMode: mermaid.RenderModeClient,
			},
			&frontmatter.Extender{
				Formats: []frontmatter.Format{frontmatter.YAML},
			},
			&hashtag.Extender{
				Variant: hashtag.ObsidianVariant,
			},
			&wikilink.Extender{},
			highlighting.NewHighlighting(highlighting.WithStyle("monokai")),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
			// parser.WithParagraphTransformers(ps),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			// html.WithXHTML(),
			extension.WithFootnoteBacklinkClass("footnote-backref"),
			extension.WithFootnoteLinkClass("footnote-ref"),
		),
	)
)

// ParseWithFrontMatter parses markdown to html and decodes the frontmatter into the provided target struct.
func ParseWithFrontMatter[
	T *PostFrontMatter | *ProjectFrontMatter | *TagFrontMatter,
](
	name string,
	source []byte,
	frontMatterTarget T,
) (string, error) {
	var (
		buf bytes.Buffer
		ctx = parser.NewContext()
	)

	err := md.Convert(source, &buf, parser.WithContext(ctx))
	if err != nil {
		return "", err
	}

	d := frontmatter.Get(ctx)
	if d == nil {
		return "", &FrontMatterMissingError{
			fileName: name,
		}
	}

	err = d.Decode(frontMatterTarget)
	if err != nil {
		return "", err
	}

	// Validate the frontmatter struct if it's not nil
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(frontMatterTarget); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return "", err
		}
		return "", err
	}

	return buf.String(), nil
}

// ProjectFrontMatter is the frontmatter of a project markdown document.
type ProjectFrontMatter struct {
	Title       string   `yaml:"title" validate:"required"`
	Description string   `yaml:"description" validate:"required"`
	Tags        []string `yaml:"tags" validate:"required"`
	Slug        string   `yaml:"slug" validate:"required"`
}

// PostFrontMatter is the frontmatter of a post markdown document.
type PostFrontMatter struct {
	Title       string   `yaml:"title" validate:"required"`
	Description string   `yaml:"description" validate:"required"`
	Tags        []string `yaml:"tags" validate:"required"`
	Slug        string   `yaml:"slug" validate:"required"`
	BannerURL   string   `yaml:"banner_url" validate:"required"`
}

// TagFrontMatter is the frontmatter of a tag markdown document.
type TagFrontMatter struct {
	Title       string   `yaml:"title" validate:"required"`
	Description string   `yaml:"description" validate:"required"`
	Tags        []string `yaml:"tags" validate:"required"`
	Slug        string   `yaml:"slug" validate:"required"`
}

// SeedFromEmbedded reads the embedded doc files and seeds the database with them.
func SeedFromEmbedded(
	ctx context.Context,
	db *data.Database[master.Queries],
) error {
	slog.Info("Seeding database from embedded markdown files")

	// Seed posts
	if err := seedPosts(ctx, db, doc.Posts); err != nil {
		return fmt.Errorf("error seeding posts: %w", err)
	}

	// Seed projects
	if err := seedProjects(ctx, db, doc.Projects); err != nil {
		return fmt.Errorf("error seeding projects: %w", err)
	}

	// Seed tags
	if err := seedTags(ctx, db, doc.Tags); err != nil {
		return fmt.Errorf("error seeding tags: %w", err)
	}

	return nil
}

// seedPosts reads post markdown files and seeds the database.
func seedPosts(
	ctx context.Context,
	db *data.Database[master.Queries],
	postsFS fs.FS,
) error {
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
		var frontMatter PostFrontMatter

		htmlContent, err := ParseWithFrontMatter(path, content, &frontMatter)
		if err != nil {
			return fmt.Errorf("error parsing frontmatter for %s: %w", path, err)
		}

		post, err := db.Queries.PostCreate(ctx, master.PostCreateParams{
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
			tag, err := db.Queries.TagGetByName(ctx, tagName)
			if err != nil {
				// Create the tag if it doesn't exist
				tag, err = db.Queries.TagCreate(ctx, master.TagCreateParams{
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
			if err := db.Queries.PostTagCreate(ctx, post.ID, tag.ID); err != nil {
				slog.Warn("Failed to associate tag with post", "postID", post.ID, "tagID", tag.ID, "error", err)
			}
		}

		slog.Info("Created post", "title", frontMatter.Title, "id", post.ID)
		return nil
	})
}

// seedProjects reads project markdown files and seeds the database.
func seedProjects(
	ctx context.Context,
	db *data.Database[master.Queries],
	projectsFS fs.FS,
) error {
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
		var frontMatter ProjectFrontMatter
		parsedContent, err := ParseWithFrontMatter(path, content, &frontMatter)
		if err != nil {
			return fmt.Errorf("error parsing frontmatter for %s: %w", path, err)
		}
		project, err := db.Queries.ProjectCreate(ctx, master.ProjectCreateParams{
			Name:        frontMatter.Title,
			Slug:        frontMatter.Slug,
			Description: frontMatter.Description,
			Content:     string(parsedContent),
		})
		if err != nil {
			// If project already exists, just log and continue
			slog.Warn("Failed to create project", "name", frontMatter.Title, "error", err)
			return nil
		}

		// Associate tags with project
		for _, tagName := range frontMatter.Tags {
			// Try to find the tag by name
			tag, err := db.Queries.TagGetByName(ctx, tagName)
			if err != nil {
				// Create the tag if it doesn't exist
				tag, err = db.Queries.TagCreate(ctx, master.TagCreateParams{
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
			if err := db.Queries.ProjectTagCreate(
				ctx,
				project.ID,
				tag.ID,
			); err != nil {
				slog.Warn("Failed to associate tag with project", "projectID", project.ID, "tagID", tag.ID, "error", err)
			}
		}

		slog.Info("Created project", "name", frontMatter.Title, "id", project.ID)
		return nil
	})
}

// seedTags reads tag markdown files and seeds the database.
func seedTags(
	ctx context.Context,
	db *data.Database[master.Queries],
	tagsFS fs.FS,
) error {
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
		var frontMatter TagFrontMatter
		_, err = ParseWithFrontMatter(path, content, &frontMatter)
		if err != nil {
			return fmt.Errorf("error parsing frontmatter for %s: %w", path, err)
		}

		// Create tag
		tag, err := db.Queries.TagCreate(ctx, master.TagCreateParams{
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
