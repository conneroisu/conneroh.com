// Package main updates the database with new vault content.
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/conneroisu/conneroh.com/cmd/conneroh"
	"github.com/conneroisu/conneroh.com/internal/data"
	"github.com/conneroisu/conneroh.com/internal/data/docs"
	"github.com/conneroisu/conneroh.com/internal/data/master"
	mathjax "github.com/litao91/goldmark-mathjax"
	ollama "github.com/prathyushnallamothu/ollamago"
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
	"gonum.org/v1/gonum/mat"
)

const (
	embeddingSize = 768
)

var client = ollama.NewClient(
	ollama.WithTimeout(time.Minute * 5),
)

// Project a single embedding to 3D using a projection matrix
func projectTo3D(embedding []float64, projectionMatrix *mat.Dense) (x, y, z float64) {
	// Create a vector from the embedding
	embedVec := mat.NewVecDense(len(embedding), embedding)

	// Project to 3D
	result := mat.NewVecDense(3, nil)
	result.MulVec(projectionMatrix.T(), embedVec)

	// Extract x, y, z coordinates
	x = result.AtVec(0)
	y = result.AtVec(1)
	z = result.AtVec(2)

	return
}

// Generate a random projection matrix for demonstration
// In practice, this would be calculated using PCA or another technique
func generateProjectionMatrix(inputDim, outputDim int) *mat.Dense {
	data := make([]float64, inputDim*outputDim)
	for i := range data {
		data[i] = rand.Float64()*2 - 1 // Random values between -1 and 1
	}
	return mat.NewDense(inputDim, outputDim, data)
}
func embeddingUpsert(
	ctx context.Context,
	db *data.Database[master.Queries],
	input string,
	id int64,
) int64 {
	resp, err := client.Embeddings(ctx, ollama.EmbeddingsRequest{
		Model:  "nomic-embed-text",
		Prompt: input,
	})
	if err != nil {
		panic(err)
	}
	jsVal, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	val := string(jsVal)
	proj := generateProjectionMatrix(embeddingSize, 3)
	x, y, z := projectTo3D(resp.Embedding, proj)
	if id == 0 {
		id, err = db.Queries.EmbeddingsCreate(ctx, master.EmbeddingsCreateParams{
			Embedding: val,
			X:         strconv.FormatFloat(x, 'f', -1, 64),
			Y:         strconv.FormatFloat(y, 'f', -1, 64),
			Z:         strconv.FormatFloat(z, 'f', -1, 64),
		})
		if err != nil {
			panic(err)
		}
		return id
	}
	err = db.Queries.EmbeddingsUpdate(ctx, master.EmbeddingsUpdateParams{
		Embedding: val,
		X:         strconv.FormatFloat(x, 'f', -1, 64),
		Y:         strconv.FormatFloat(y, 'f', -1, 64),
		Z:         strconv.FormatFloat(z, 'f', -1, 64),
	})
	if err != nil {
		panic(err)
	}
	return id
}

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

// UpsertTag inserts or updates the tag document in the database.
func (md *Markdown) UpsertTag(
	ctx context.Context,
	db *data.Database[master.Queries],
) (err error) {
	slog.Info("upserting tag", "slug", md.Slug, "icon", md.Icon)
	defer slog.Info("upserted tag", "slug", md.Slug, "icon", md.Icon)
	var tag master.Tag
	tag, err = db.Queries.TagGetBySlug(ctx, md.Slug)
	if err == nil {
		var embedID = tag.EmbeddingID
		if tag.RawContent != md.RawContent {
			embedID = embeddingUpsert(ctx, db, md.RawContent, tag.EmbeddingID)
		}
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
				EmbeddingID: embedID,
			},
		)
	}

	if errors.Is(err, sql.ErrNoRows) {
		slog.Debug("tag does not exist - creating", "slug", md.Slug)
		return db.Queries.TagCreate(ctx, master.TagCreateParams{
			Slug:        md.Slug,
			Title:       md.Title,
			Description: md.Description,
			RawContent:  md.RawContent,
			Content:     md.RenderContent,
			Icon:        md.Icon,
			EmbeddingID: embeddingUpsert(ctx, db, md.RawContent, 0),
		})
	}
	return
}

// UpsertPost inserts or updates the post document in the database.
func (md *Markdown) UpsertPost(
	ctx context.Context,
	db *data.Database[master.Queries],
) (post master.Post, err error) {
	slog.Info("upserting post", "slug", md.Slug)
	defer slog.Info("upserted post", "slug", md.Slug)

	post, err = db.Queries.PostGetBySlug(ctx, md.Slug)
	if err == nil {
		var embedID = post.EmbeddingID
		if post.RawContent != md.RawContent {
			embedID = embeddingUpsert(ctx, db, md.RawContent, post.EmbeddingID)
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
				EmbeddingID: embedID,
			},
		)
	}
	// err is not nil (new post)
	if errors.Is(err, sql.ErrNoRows) {
		return db.Queries.PostCreate(ctx, master.PostCreateParams{
			Slug:        md.Slug,
			Title:       md.Title,
			BannerUrl:   md.BannerURL,
			Description: md.Description,
			RawContent:  md.RawContent,
			Content:     md.RenderContent,
			EmbeddingID: embeddingUpsert(ctx, db, md.RawContent, 0),
		})
	}
	return
}

// UpsertProject inserts or updates a project in the database.
func (md *Markdown) UpsertProject(
	ctx context.Context,
	db *data.Database[master.Queries],
) (project master.Project, err error) {
	if md == nil {
		return master.Project{}, errors.New("markdown is nil for a project")
	}
	slog.Info("upserting project", "slug", md.Slug)
	defer slog.Info("upserted project", "slug", md.Slug)
	if md.Title == "" {
		return master.Project{}, errors.New("title is empty")
	}
	project, err = db.Queries.ProjectGetBySlug(ctx, md.Slug)
	if err == nil {
		var embedID = project.EmbeddingID
		if project.RawContent != md.RawContent {
			embedID = embeddingUpsert(ctx, db, md.RawContent, project.EmbeddingID)
		}
		return db.Queries.ProjectUpdate(
			ctx,
			master.ProjectUpdateParams{
				ID:          project.ID,
				Slug:        md.Slug,
				Title:       md.Title,
				BannerUrl:   md.BannerURL,
				Description: md.Description,
				Content:     md.RenderContent,
				EmbeddingID: embedID,
			},
		)
	}
	// err is not nil (new project)
	if errors.Is(err, sql.ErrNoRows) {
		return db.Queries.ProjectCreate(ctx, master.ProjectCreateParams{
			Slug:        md.Slug,
			Title:       md.Title,
			BannerUrl:   md.BannerURL,
			Description: md.Description,
			Content:     md.RenderContent,
			EmbeddingID: embeddingUpsert(ctx, db, md.RawContent, 0),
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

	parsedTags, err := pathParse("tags", docs.Tags)
	if err != nil {
		return err
	}
	for _, parsed := range parsedTags {
		err = parsed.UpsertTag(ctx, db)
		if err != nil {
			return fmt.Errorf("failed to upsert tag %s: %w", parsed.Title, err)
		}
	}

	parsedPosts, err := pathParse("posts", docs.Posts)
	if err != nil {
		return err
	}
	var post master.Post
	for _, parsed := range parsedPosts {
		post, err = parsed.UpsertPost(ctx, db)
		if err != nil {
			return fmt.Errorf("failed to upsert post %s: %w", parsed.Title, err)
		}
		err = db.Queries.UpsertPostTags(ctx, parsed.Tags, post.ID)
		if err != nil {
			return fmt.Errorf("failed to upsert post tags: %v", err)
		}
	}

	parsedProjects, err := pathParse("projects", docs.Projects)
	if err != nil {
		return fmt.Errorf("failed to parse projects: %v", err)
	}

	for _, parsed := range parsedProjects {
		project, err := parsed.UpsertProject(ctx, db)
		if err != nil {
			return fmt.Errorf("failed to upsert project %s: %w", parsed.Title, err)
		}
		err = db.Queries.UpsertProjectTags(ctx, parsed.Tags, project.ID)
		if err != nil {
			return fmt.Errorf("failed to upsert project tags: %v", err)
		}
	}
	return nil
}
