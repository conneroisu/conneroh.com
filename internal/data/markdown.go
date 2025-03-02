package data

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"strings"

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

// ParseAll parses all markdown files in the database.
func ParseAll(ctx context.Context, db *Database[master.Queries]) error {
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
			slog.Info("parsing tag", "path", path)
			parsed, err := Parse(path, true)
			if err != nil {
				return err
			}
			return parsed.UpsertTag(ctx, db, client)
		},
	)
	if err != nil {
		return err
	}

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
			slog.Info("parsing post", "path", path)
			parsed, err := Parse(path, false)
			if err != nil {
				return err
			}
			return parsed.UpsertPost(ctx, db, client)
		},
	)
	if err != nil {
		return err
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
			slog.Info("parsing project", "path", path)
			parsed, err := Parse(path, false)
			if err != nil {
				return err
			}
			return parsed.UpsertProject(ctx, db, client)
		},
	)
	if err != nil {
		return err
	}
	return nil
}

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
func Parse(path string, tag bool) (*Markdown, error) {
	b, err := os.ReadFile(path)
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
	fm.RenderContent = buf.String()
	if fm.Description == "" {
		if !tag {
			return nil, fmt.Errorf("description is empty for %s", path)
		}
		// It is a tag page.
		fm.Description = fm.RawContent
	}
	fm.RawContent = string(b)

	return &fm, nil
}

func (md *Markdown) assertTags(
	ctx context.Context,
	db *Database[master.Queries],
) error {
	for _, tag := range md.Tags {
		_, err := db.Queries.TagGetBySlug(ctx, tag)
		if err == nil || errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf(
				"failed to find referenced tag %s: %w",
				tag,
				err,
			)
		}
	}
	return nil
}

// UpsertEmbedding inserts or updates the embedding document in the database.
func UpsertEmbedding(
	ctx context.Context,
	db *Database[master.Queries],
	client *api.Client,
	content string,
	existingID int64,
) (int64, error) {
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
	db *Database[master.Queries],
	client *api.Client,
) error {
	var id int64
	tag, err := db.Queries.TagGetBySlug(ctx, md.Slug)
	if err == nil {
		if tag.Description == md.RawContent {
			return db.Queries.TagUpdate(
				ctx,
				master.TagUpdateParams{
					ID:          tag.ID,
					Slug:        md.Slug,
					Title:       md.Title,
					Description: md.Description,
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
				EmbeddingID: id,
			},
		)
	}
	// err is not nil (new tag)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = UpsertEmbedding(ctx, db, client, md.RawContent, 0)
		if err != nil {
			return err
		}
		return db.Queries.TagCreate(ctx, master.TagCreateParams{
			Slug:        md.Slug,
			Title:       md.Title,
			Description: md.Description,
			EmbeddingID: id,
		})
	}
	return err
}

// UpsertPost inserts or updates the post document in the database.
func (md *Markdown) UpsertPost(
	ctx context.Context,
	db *Database[master.Queries],
	client *api.Client,
) error {
	var id int64
	err := md.assertTags(ctx, db)
	if err != nil {
		return err
	}
	// check if the post already exists
	post, err := db.Queries.PostGetBySlug(ctx, md.Slug)
	if err == nil {
		// same content
		if post.RawContent == string(md.RawContent) {
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
			return err
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
			return err
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
	return err
}

// UpsertProject inserts or updates a project in the database.
func (md *Markdown) UpsertProject(
	ctx context.Context,
	db *Database[master.Queries],
	client *api.Client,
) error {
	var id int64
	err := md.assertTags(ctx, db)
	if err != nil {
		return err
	}
	project, err := db.Queries.ProjectGetBySlug(ctx, md.Slug)
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
			return err
		}

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
			return err
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
	return nil
}
