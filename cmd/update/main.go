// Package main updates the database with new vault content.
package main

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/conneroisu/conneroh.com/cmd/conneroh"
	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/docs"
	"github.com/conneroisu/conneroh.com/internal/data/master"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/ollama/ollama/api"
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

var (
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, extension.Footnote,
			extension.Strikethrough, extension.Table,
			extension.TaskList, extension.DefinitionList,
			mathjax.MathJax, &wikilink.Extender{},
			extension.NewTypographer(
				extension.WithTypographicSubstitutions(
					extension.TypographicSubstitutions{
						extension.Apostrophe: []byte("'"),
					}),
			),
			enclave.New(&core.Config{DefaultImageAltPrefix: "caption: "}),
			extension.NewFootnote(
				extension.WithFootnoteIDPrefix("fn"),
			),
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
			highlighting.NewHighlighting(highlighting.WithStyle("monokai")),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			extension.WithFootnoteBacklinkClass("footnote-backref"),
			extension.WithFootnoteLinkClass("footnote-ref"),
		),
	)
)

// Markdown is the frontmatter and content of a markdown document.
type Markdown struct {
	// Title is the title/display name of the document.
	Title string `yaml:"title"`
	// Slug is the slug of the document.
	//
	// It is used to generate the URL.
	//
	// It defaults to the name of the file.
	Slug string `yaml:"-"`
	// Description is a description of the document.
	Description string `yaml:"description,omitempty"`
	// Tags are related tags of the document.
	Tags []string `yaml:"tags"`
	// BannerURL is the URL of the banner image.
	BannerURL string `yaml:"banner_url"`
	// Icon is the URL of the icon image.
	Icon string `yaml:"icon,omitempty"`

	// Posts are related posts of the document. (Never used on tags)
	Posts []string `yaml:"posts,omitempty"`
	// Projects are related projects of the document. (Never used on tags)
	Projects []string `yaml:"projects,omitempty"`

	// RawContent is the content of the document.
	RawContent string `yaml:"-"`
	// RenderContent is the content of the document to be rendered as HTML.
	RenderContent string `yaml:"-"`
}

// Parse parses the markdown file at the given path.
func Parse(fsPath string, embedFs embed.FS) (*Markdown, error) {
	b, err := embedFs.ReadFile(fsPath)
	if err != nil {
		return nil, err
	}
	var fm Markdown
	buf := bytes.NewBufferString("")
	ctx := parser.NewContext()
	err = md.Convert(b, buf, parser.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	d := frontmatter.Get(ctx)
	if d == nil {
		return nil, nil
	}
	err = d.Decode(&fm)
	if err != nil {
		return nil, err
	}
	switch embedFs {
	case docs.Posts:
		fsPath = strings.Replace(fsPath, "posts/", "", 1)
	case docs.Tags:
		fsPath = strings.Replace(fsPath, "tags/", "", 1)
	case docs.Projects:
		fsPath = strings.Replace(fsPath, "projects/", "", 1)
	}
	fsPath = strings.TrimSuffix(fsPath, filepath.Ext(fsPath))
	fm.Slug = fsPath
	fm.RenderContent = buf.String()
	if fm.Description == "" {
		return nil, fmt.Errorf("description is empty for %s", fsPath)
	}
	fm.RawContent = string(b)

	if fm.Icon == "" {
		fm.Icon = "nf-fa-tag"
	}

	return &fm, nil
}

// UpsertEmbedding inserts or updates the embedding document in the database.
func UpsertEmbedding(
	ctx context.Context,
	db *data.Database[master.Queries],
	client *api.Client,
	content string,
	existingID int64,
) (int64, error) {
	slog.Info("upserting embedding", "existingID", existingID)
	defer slog.Info("upserted embedding", "existingID", existingID)
	embed, err := client.Embeddings(ctx, &api.EmbeddingRequest{
		Model:  "granite-embedding:278m",
		Prompt: content,
	})
	if err != nil {
		return 0, err
	}
	jsVal, err := json.Marshal(embed)
	if err != nil {
		return 0, err
	}
	// embedding already exists, update it
	if existingID != 0 {
		err = db.Queries.EmeddingUpdate(ctx, string(jsVal), existingID)
		if err != nil {
			return 0, err
		}
		return existingID, nil
	}
	// embedding does not exist, create it
	id, err := db.Queries.EmbeddingsCreate(ctx, string(jsVal))
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpsertTag inserts or updates the tag document in the database.
func (md *Markdown) UpsertTag(
	ctx context.Context,
	db *data.Database[master.Queries],
	client *api.Client,
) (err error) {
	slog.Info("upserting tag", "slug", md.Slug, "icon", md.Icon)
	defer slog.Info("upserted tag", "slug", md.Slug, "icon", md.Icon)
	var id int64
	var tag master.Tag
	tag, err = db.Queries.TagGetBySlug(ctx, md.Slug)
	if err == nil {
		slog.Debug("tag already exists - updating", "slug", md.Slug, "icon", md.Icon)
		if tag.RawContent == md.RawContent {
			return db.Queries.TagUpdate(
				ctx,
				master.TagUpdateParams{
					ID:          tag.ID,
					Slug:        md.Slug,
					Title:       md.Title,
					RawContent:  md.RawContent,
					Content:     md.RenderContent,
					Description: md.Description,
					Icon:        md.Icon,
					EmbeddingID: tag.EmbeddingID,
				},
			)
		}
		id, err = UpsertEmbedding(
			ctx,
			db,
			client,
			md.RawContent,
			tag.EmbeddingID,
		)
		if err != nil {
			return err
		}
		// updated embedding
		return db.Queries.TagUpdate(
			ctx,
			master.TagUpdateParams{
				ID:          tag.ID,
				Slug:        md.Slug,
				Title:       md.Title,
				Description: md.Description,
				RawContent:  md.RawContent,
				Content:     md.RenderContent,
				EmbeddingID: id,
			},
		)
	}
	// err is not nil (new tag)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Debug("tag does not exist - creating", "slug", md.Slug)
		id, err = UpsertEmbedding(ctx, db, client, md.RawContent, 0)
		if err != nil {
			return
		}
		return db.Queries.TagCreate(ctx, master.TagCreateParams{
			Slug:        md.Slug,
			Title:       md.Title,
			Description: md.Description,
			RawContent:  md.RawContent,
			Content:     md.RenderContent,
			EmbeddingID: id,
		})
	}
	return
}

// UpsertPost inserts or updates the post document in the database.
func (md *Markdown) UpsertPost(
	ctx context.Context,
	db *data.Database[master.Queries],
	client *api.Client,
) (post master.Post, err error) {
	slog.Info("upserting post", "slug", md.Slug)
	defer slog.Info("upserted post", "slug", md.Slug)
	var id int64
	// check if the post already exists
	post, err = db.Queries.PostGetBySlug(ctx, md.Slug)
	if err == nil {
		// same content
		if post.RawContent == md.RawContent {
			return db.Queries.PostUpdate(
				ctx,
				master.PostUpdateParams{
					ID:          post.ID,
					Slug:        md.Slug,
					Title:       md.Title,
					BannerUrl:   md.BannerURL,
					Description: md.Description,
					RawContent:  md.RawContent,
					Content:     md.RenderContent,
					EmbeddingID: post.EmbeddingID,
				},
			)
		}
		id, err = UpsertEmbedding(
			ctx,
			db,
			client,
			md.RawContent,
			post.EmbeddingID,
		)
		if err != nil {
			return
		}
		return db.Queries.PostUpdate(
			ctx,
			master.PostUpdateParams{
				ID:          post.ID,
				Slug:        md.Slug,
				Title:       md.Title,
				BannerUrl:   md.BannerURL,
				Description: md.Description,
				RawContent:  md.RawContent,
				Content:     md.RenderContent,
				EmbeddingID: id,
			},
		)
	}
	// err is not nil (new post)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = UpsertEmbedding(ctx, db, client, md.RawContent, 0)
		if err != nil {
			return
		}
		return db.Queries.PostCreate(ctx, master.PostCreateParams{
			Slug:        md.Slug,
			Title:       md.Title,
			BannerUrl:   md.BannerURL,
			Description: md.Description,
			RawContent:  md.RawContent,
			Content:     md.RenderContent,
			EmbeddingID: id,
		})
	}
	return
}

// UpsertProject inserts or updates a project in the database.
func (md *Markdown) UpsertProject(
	ctx context.Context,
	db *data.Database[master.Queries],
	client *api.Client,
) (project master.Project, err error) {
	if md == nil {
		return master.Project{}, errors.New("markdown is nil for a project")
	}
	slog.Info("upserting project", "slug", md.Slug)
	defer slog.Info("upserted project", "slug", md.Slug)
	var id int64
	if md.Title == "" {
		return master.Project{}, errors.New("title is empty")
	}
	project, err = db.Queries.ProjectGetBySlug(ctx, md.Slug)
	if err == nil {
		// project already exists with the same content, update metadata
		if project.Content == md.RawContent {
			return db.Queries.ProjectUpdate(
				ctx,
				master.ProjectUpdateParams{
					ID:          project.ID,
					Slug:        md.Slug,
					Title:       md.Title,
					BannerUrl:   md.BannerURL,
					Description: md.Description,
					Content:     md.RenderContent,
					RawContent:  md.RawContent,
					EmbeddingID: project.EmbeddingID,
				},
			)
		}
		// project already exists with a different content, update embedding
		id, err = UpsertEmbedding(
			ctx,
			db,
			client,
			string(md.RawContent),
			project.EmbeddingID,
		)
		if err != nil {
			return
		}
		return db.Queries.ProjectUpdate(ctx, master.ProjectUpdateParams{
			ID:          project.ID,
			Slug:        md.Slug,
			Title:       md.Title,
			BannerUrl:   md.BannerURL,
			Description: md.Description,
			Content:     md.RenderContent,
			EmbeddingID: id,
		})
	}
	// err is not nil (new project)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = UpsertEmbedding(
			ctx,
			db,
			client,
			md.RawContent,
			0,
		)
		if err != nil {
			return
		}
		return db.Queries.ProjectCreate(ctx, master.ProjectCreateParams{
			Slug:        md.Slug,
			Title:       md.Title,
			BannerUrl:   md.BannerURL,
			Description: md.Description,
			Content:     md.RenderContent,
			EmbeddingID: id,
		})
	}
	slog.Error("project update failed", "err", err)
	return
}

func main() {
	ctx := context.Background()
	err := Run(ctx, os.Getenv)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Run parses all markdown files in the database.
func Run(
	ctx context.Context,
	getenv func(string) string,
) error {
	db, err := conneroh.NewDb(getenv)
	if err != nil {
		return err
	}
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return err
	}

	err = fs.WalkDir(
		docs.Tags,
		"tags",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			parsed, err := Parse(path, docs.Tags)
			if err != nil {
				slog.Error("failed to parse tag", "path", path, "err", err)
				return err
			}
			return parsed.UpsertTag(ctx, db, client)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update tags: %v", err)
	}

	var post master.Post
	err = fs.WalkDir(
		docs.Posts,
		"posts",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			parsed, err := Parse(path, docs.Posts)
			if err != nil {
				return err
			}
			post, err = parsed.UpsertPost(ctx, db, client)
			if err != nil {
				return err
			}
			return db.Queries.UpsertPostTags(ctx, parsed.Tags, post.ID)
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update posts: %v", err)
	}

	err = fs.WalkDir(
		docs.Projects,
		"projects",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			parsed, err := Parse(path, docs.Projects)
			if err != nil {
				return err
			}
			project, err := parsed.UpsertProject(ctx, db, client)
			if err != nil {
				return err
			}
			err = db.Queries.UpsertProjectTags(ctx, parsed.Tags, project.ID)
			if err != nil {
				return fmt.Errorf("failed to upsert project tags: %v", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update projects: %v", err)
	}
	return nil
}
