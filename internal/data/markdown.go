package data

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"

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

	// Posts are related posts of the document. (Never used on tags)
	Posts []string `yaml:"posts,omitempty"`
	// Projects are related projects of the document. (Never used on tags)
	Projects []string `yaml:"projects,omitempty"`

	// RawContent is the content of the document.
	RawContent []byte `yaml:"-"`
	// RenderContent is the content of the document to be rendered as HTML.
	RenderContent []byte `yaml:"-"`
	// Emebbedding is the content of the document to be embedded in the page.
	EmbeddingContent [][]float64 `yaml:"-"`
}

// Parse parses the markdown file at the given path.
func Parse(path string) (*Markdown, error) {
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
	fm.RenderContent = buf.Bytes()
	if fm.Description == "" {
		// It is a tag page.
		fm.Description = string(fm.RenderContent)
	}
	fm.RawContent = b

	return &fm, nil
}

// UpsertPost inserts or updates the post document in the database.
func (md *Markdown) UpsertPost(
	ctx context.Context,
	db *Database[master.Queries],
	client *api.Client,
) error {
	// check if the post already exists
	post, err := db.Queries.PostGetBySlug(ctx, md.Slug)
	if err == nil {
		// same content
		if post.RawContent == string(md.RawContent) {
			return db.Queries.PostUpdate(ctx, master.PostUpdateParams{
				ID:          post.ID,
				Slug:        md.Slug,
				Title:       md.Title,
				BannerUrl:   md.BannerURL,
				Description: md.Description,
				RawContent:  string(md.RawContent),
				Content:     string(md.RenderContent),
				EmbeddingID: post.EmbeddingID,
			})
		}
		id, err := UpsertEmbedding(ctx, db, client, string(md.RawContent))
		if err != nil {
			return err
		}
		// same embedding
		if id == post.EmbeddingID {
			return db.Queries.PostUpdate(ctx, master.PostUpdateParams{
				ID:          post.ID,
				Slug:        md.Slug,
				Title:       md.Title,
				BannerUrl:   md.BannerURL,
				Description: md.Description,
				RawContent:  string(md.RawContent),
				Content:     string(md.RenderContent),
				EmbeddingID: id,
			})
		}
	}
	// err is not nil
	if errors.Is(err, sql.ErrNoRows) {
		return md.insertPost(ctx, db)
	}
	return err
}

// UpsertEmbedding inserts or updates the embedding document in the database.
func UpsertEmbedding(
	ctx context.Context,
	db *Database[master.Queries],
	client *api.Client,
	content string,
) (int64, error) {
	embed, err := client.Embeddings(ctx, &api.EmbeddingRequest{
		Model:  "granite-embedding:278m",
		Prompt: content,
	})
	if err != nil {
		return int64(0), err
	}
	jsVal, err := json.Marshal(embed)
	if err != nil {
		return 0, err
	}
	id, err := db.Queries.EmbeddingsCreate(ctx, string(jsVal))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (md *Markdown) insertPost(
	ctx context.Context,
	db *Database[master.Queries],
) error {
	return nil
}
