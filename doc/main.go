// Package main contains documentation and updater.
package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	"bytes"

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

// Posts contains all posts.
//
//go:embed posts/*.md
var Posts embed.FS

// Projects contains all projects.
//
//go:embed projects/*.md
var Projects embed.FS

// Tags contains all tags.
//
//go:embed tags/*.md
var Tags embed.FS

// Frontmatter is the frontmatter of a tag markdown document.
type Frontmatter struct {
	Title       string   `yaml:"title" validate:"required"`
	Description string   `yaml:"description" validate:"required"`
	Tags        []string `yaml:"tags" validate:"required"`
	Slug        string   `yaml:"slug" validate:"required"`
	BannerURL   string   `yaml:"banner_url" validate:"required"`
}

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
					extension.Apostrophe: []byte("'"),
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
func ParseWithFrontMatter(
	name string,
	source []byte,
) (content string, fm Frontmatter, err error) {
	var (
		buf bytes.Buffer
		ctx = parser.NewContext()
	)
	err = md.Convert(source, &buf, parser.WithContext(ctx))
	if err != nil {
		return
	}
	d := frontmatter.Get(ctx)
	if d == nil {
		err = &FrontMatterMissingError{
			fileName: name,
		}
		return
	}
	err = d.Decode(&fm)
	if err != nil {
		return
	}
	// Validate the frontmatter struct if it's not nil
	v := validator.New(validator.WithRequiredStructEnabled())
	if err = v.Struct(fm); err != nil {
		var ok bool
		err, ok = err.(*validator.InvalidValidationError)
		if ok {
			return
		}
		return
	}
	content = buf.String()
	return
}

// SeedFromEmbedded reads the embedded doc files and seeds the database with them.
func SeedFromEmbedded(
	ctx context.Context,
	db *data.Database[master.Queries],
) error {
	slog.Info("Seeding database from embedded markdown files")

	// Seed tags
	if err := seedTags(ctx, db, Tags); err != nil {
		return fmt.Errorf("error seeding tags: %w", err)
	}

	// Seed posts
	if err := seedPosts(ctx, db, Posts); err != nil {
		return fmt.Errorf("error seeding posts: %w", err)
	}

	// Seed projects
	if err := seedProjects(ctx, db, Projects); err != nil {
		return fmt.Errorf("error seeding projects: %w", err)
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
		htmlContent, frontMatter, err := ParseWithFrontMatter(path, content)
		if err != nil {
			return fmt.Errorf("error parsing frontmatter for %s: %w", path, err)
		}

		// Prepare the post struct
		post := master.Post{
			Title:       frontMatter.Title,
			Description: frontMatter.Description,
			Slug:        frontMatter.Slug,
			Content:     string(htmlContent),
			BannerUrl:   frontMatter.BannerURL,
		}

		// Collect all tags that will be associated with this post
		var postTags []master.Tag
		var tag master.Tag
		for _, tagName := range frontMatter.Tags {
			// Try to find the tag by name
			tag, err = db.Queries.TagGetByName(ctx, tagName)
			if err != nil {
				// Create the tag if it doesn't exist
				tag, err = db.Queries.TagCreate(ctx, master.TagCreateParams{
					Name:        tagName,
					Description: "Auto-generated tag from post " + frontMatter.Title,
					Slug:        tagName,
				})
				if err != nil {
					slog.Warn("Failed to create tag", "name", tagName, "error", err)
					continue
				}
			}
			postTags = append(postTags, tag)
		}

		// Create post with all related tags and projects (empty projects for now)
		fullPost, err := db.Queries.FullPostsCreate(ctx, post, []master.Project{}, postTags)
		if err != nil {
			slog.Warn("failed to create post", "title", frontMatter.Title, "error", err)
			return err
		}

		slog.Info("created post", "title", frontMatter.Title, "id", fullPost.Post.ID)
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
		content, err := fs.ReadFile(projectsFS, path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}
		htmlContent, frontMatter, err := ParseWithFrontMatter(path, content)
		if err != nil {
			return fmt.Errorf("error parsing frontmatter for %s: %w", path, err)
		}

		// Prepare the project struct
		project := master.Project{
			Name:        frontMatter.Title,
			Slug:        frontMatter.Slug,
			Description: frontMatter.Description,
			Content:     string(htmlContent),
			BannerUrl:   frontMatter.BannerURL,
		}

		// Collect all tags that will be associated with this project
		var projectTags []master.Tag
		var tag master.Tag
		for _, tagName := range frontMatter.Tags {
			tag, err = db.Queries.TagGetByName(ctx, tagName)
			if err != nil {
				// Create the tag if it doesn't exist
				tag, err = db.Queries.TagCreate(ctx, master.TagCreateParams{
					Name:        tagName,
					Description: "Auto-generated tag from project " + frontMatter.Title,
					Slug:        tagName,
				})
				if err != nil {
					slog.Warn("Failed to create tag", "name", tagName, "error", err)
					continue
				}
			}
			projectTags = append(projectTags, tag)
		}

		// Create project with all related tags and posts (empty posts for now)
		fullProject, err := db.Queries.FullProjectCreate(ctx, project, []master.Post{}, projectTags)
		if err != nil {
			// If project already exists, just log and continue
			slog.Warn("Failed to create project", "name", frontMatter.Title, "error", err)
			return nil
		}

		slog.Info("Created project", "name", frontMatter.Title, "id", fullProject.Project.ID)
		return nil
	})
}

// seedTags reads tag markdown files and seeds the database.
func seedTags(
	ctx context.Context,
	db *data.Database[master.Queries],
	tagsFS fs.FS,
) error {
	return fs.WalkDir(
		tagsFS,
		"tags",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			content, err := fs.ReadFile(tagsFS, path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", path, err)
			}
			htmlContent, frontMatter, err := ParseWithFrontMatter(path, content)
			if err != nil {
				return fmt.Errorf("error parsing frontmatter for %s: %w", path, err)
			}
			// Prepare tag struct
			tag := master.Tag{
				Name:        frontMatter.Title,
				Description: htmlContent,
				Slug:        frontMatter.Slug,
			}
			// Create the tag with empty posts and projects arrays (will be populated later)
			fullTag, err := db.Queries.FullTagCreate(ctx, tag, []master.Project{}, []master.Post{})
			if err != nil {
				// If tag already exists, just log and continue
				slog.Warn("Failed to create tag", "name", frontMatter.Title, "error", err)
				return nil
			}

			slog.Info("Created tag", "name", frontMatter.Title, "id", fullTag.Tag.ID)
			return nil
		})
}
