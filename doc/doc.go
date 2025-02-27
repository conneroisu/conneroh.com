package doc

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

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

// Tags contains all tags
//
//go:embed tags/*.md
var Tags embed.FS

// FrontMatterMissingError is returned when the front matter is missing from the markdown file.
type FrontMatterMissingError struct{ fileName string }

func (e FrontMatterMissingError) Error() string { return "front matter missing from " + e.fileName }

var md = goldmark.New(goldmark.WithExtensions(extension.GFM, extension.Footnote, extension.Strikethrough, extension.Table, extension.TaskList, extension.DefinitionList, extension.NewTypographer(extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{extension.Apostrophe: []byte("'")})), enclave.New(&core.Config{DefaultImageAltPrefix: "caption: "}), extension.NewFootnote(extension.WithFootnoteIDPrefix("fn")), mathjax.MathJax, &anchor.Extender{Position: anchor.Before, Texter: anchor.Text("#"), Attributer: anchor.Attributes{"class": "anchor permalink p-4"}}, &mermaid.Extender{RenderMode: mermaid.RenderModeClient}, &frontmatter.Extender{Formats: []frontmatter.Format{frontmatter.YAML}}, &hashtag.Extender{Variant: hashtag.ObsidianVariant}, &wikilink.Extender{}, highlighting.NewHighlighting(highlighting.WithStyle("monokai"))), goldmark.WithParserOptions(parser.WithAutoHeadingID(), parser.WithAttribute()), goldmark.WithRendererOptions(html.WithHardWraps(), extension.WithFootnoteBacklinkClass("footnote-backref"), extension.WithFootnoteLinkClass("footnote-ref")))

// ParseWithFrontMatter parses markdown to html and decodes the frontmatter into the provided target struct.
func ParseWithFrontMatter[T *PostFrontMatter | *ProjectFrontMatter | *TagFrontMatter](name string, source []byte, frontMatterTarget T) (string, error) {
	var buf bytes.Buffer
	ctx := parser.NewContext()
	if err := md.Convert(source, &buf, parser.WithContext(ctx)); err != nil {
		return "", err
	}
	if d := frontmatter.Get(ctx); d == nil {
		return "", &FrontMatterMissingError{fileName: name}
	} else if err := d.Decode(frontMatterTarget); err != nil {
		return "", err
	}
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(frontMatterTarget); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// BaseFrontMatter is the base frontmatter of a markdown document.
type BaseFrontMatter struct {
	Title       string   `yaml:"title" validate:"required"`
	Description string   `yaml:"description" validate:"required"`
	Tags        []string `yaml:"tags" validate:"required"`
	Slug        string   `yaml:"slug" validate:"required"`
}

// ProjectFrontMatter is the frontmatter of a project markdown document.
type ProjectFrontMatter struct{ BaseFrontMatter }

// TagFrontMatter is the frontmatter of a tag markdown document.
type TagFrontMatter struct{ BaseFrontMatter }

// PostFrontMatter is the frontmatter of a post markdown document.
type PostFrontMatter struct {
	BaseFrontMatter
	BannerURL string `yaml:"banner_url" validate:"required"`
}

// SeedFromEmbedded reads the embedded doc files and seeds the database with them.
func SeedFromEmbedded(ctx context.Context, db *data.Database[master.Queries]) error {
	slog.Info("Seeding database from embedded markdown files")
	if err := seedPosts(ctx, db, Posts); err != nil {
		return fmt.Errorf("error seeding posts: %w", err)
	}
	if err := seedProjects(ctx, db, Projects); err != nil {
		return fmt.Errorf("error seeding projects: %w", err)
	}
	if err := seedTags(ctx, db, Tags); err != nil {
		return fmt.Errorf("error seeding tags: %w", err)
	}
	return nil
}
func createTag(ctx context.Context, db *data.Database[master.Queries], name, description string) (master.Tag, error) {
	if tag, err := db.Queries.TagGetByName(ctx, name); err == nil {
		return tag, nil
	}
	return db.Queries.TagCreate(ctx, master.TagCreateParams{Name: name, Description: description, Slug: name})
}
func processMarkdownFile[T *PostFrontMatter | *ProjectFrontMatter | *TagFrontMatter](fsys fs.FS, path string, target T) (string, error) {
	content, err := fs.ReadFile(fsys, path)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", path, err)
	}
	parsedContent, err := ParseWithFrontMatter(path, content, target)
	if err != nil {
		return "", fmt.Errorf("error parsing frontmatter for %s: %w", path, err)
	}
	return parsedContent, nil
}
func seedPosts(ctx context.Context, db *data.Database[master.Queries], postsFS fs.FS) error {
	return fs.WalkDir(postsFS, "posts", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		var frontMatter PostFrontMatter
		htmlContent, err := processMarkdownFile(postsFS, path, &frontMatter)
		if err != nil {
			return err
		}
		post, err := db.Queries.PostCreate(ctx, master.PostCreateParams{Title: frontMatter.Title, Description: frontMatter.Description, Slug: frontMatter.Slug, Content: string(htmlContent), BannerUrl: "/assets/img/posts/default.jpg"})
		if err != nil {
			slog.Warn("Failed to create post", "title", frontMatter.Title, "error", err)
			return nil
		}
		for _, tagName := range frontMatter.Tags {
			tag, err := createTag(ctx, db, tagName, "Auto-generated tag from post "+frontMatter.Title)
			if err != nil {
				slog.Warn("Failed to create tag", "name", tagName, "error", err)
				continue
			}
			if err := db.Queries.PostTagCreate(ctx, post.ID, tag.ID); err != nil {
				slog.Warn("Failed to associate tag with post", "postID", post.ID, "tagID", tag.ID, "error", err)
			}
		}
		slog.Info("Created post", "title", frontMatter.Title, "id", post.ID)
		return nil
	})
}
func seedProjects(ctx context.Context, db *data.Database[master.Queries], projectsFS fs.FS) error {
	return fs.WalkDir(projectsFS, "projects", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		var frontMatter ProjectFrontMatter
		parsedContent, err := processMarkdownFile(projectsFS, path, &frontMatter)
		if err != nil {
			return err
		}
		project, err := db.Queries.ProjectCreate(ctx, master.ProjectCreateParams{Name: frontMatter.Title, Slug: frontMatter.Slug, Description: frontMatter.Description, Content: string(parsedContent)})
		if err != nil {
			slog.Warn("Failed to create project", "name", frontMatter.Title, "error", err)
			return nil
		}
		for _, tagName := range frontMatter.Tags {
			tag, err := createTag(ctx, db, tagName, "Auto-generated tag from project")
			if err != nil {
				slog.Warn("Failed to create tag", "name", tagName, "error", err)
				continue
			}
			if err := db.Queries.ProjectTagCreate(ctx, project.ID, tag.ID); err != nil {
				slog.Warn("Failed to associate tag with project", "projectID", project.ID, "tagID", tag.ID, "error", err)
			}
		}
		slog.Info("Created project", "name", frontMatter.Title, "id", project.ID)
		return nil
	})
}
func seedTags(ctx context.Context, db *data.Database[master.Queries], tagsFS fs.FS) error {
	return fs.WalkDir(tagsFS, "tags", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		var frontMatter TagFrontMatter
		_, err = processMarkdownFile(tagsFS, path, &frontMatter)
		if err != nil {
			return err
		}
		tag, err := db.Queries.TagCreate(ctx, master.TagCreateParams{Name: frontMatter.Title, Description: frontMatter.Description, Slug: frontMatter.Slug})
		if err != nil {
			slog.Warn("Failed to create tag", "name", frontMatter.Title, "error", err)
			return nil
		}
		slog.Info("Created tag", "name", frontMatter.Title, "id", tag.ID)
		return nil
	})
}
