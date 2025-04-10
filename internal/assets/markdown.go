package assets

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/credited"
	"github.com/conneroisu/conneroh.com/internal/data/gen"
	mathjax "github.com/litao91/goldmark-mathjax"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/quail-ink/goldmark-enclave/core"
	"github.com/rotisserie/eris"
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

type awsClient interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// Renderer is a markdown parser/renderer.
type Renderer interface {
	Convert(Asset) (gen.Embedded, error)
}

// Converter is a markdown parser.
type Converter interface {
	Convert(source []byte, writer io.Writer, opts ...parser.ParseOption) error
}

// NewRenderer creates a new markdown parser.
func NewRenderer(
	ctx context.Context,
	pCtx parser.Context,
	client credited.AWSClient,
) *DefaultRenderer {
	return &DefaultRenderer{
		Markdown: goldmark.New(goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				extension.WithFootnoteBacklinkClass("footnote-backref"),
				extension.WithFootnoteLinkClass("footnote-ref"),
			),
			goldmark.WithExtensions(
				&wikilink.Extender{
					Resolver: newResolver(ctx, client),
				},
				extension.GFM,
				extension.Footnote,
				extension.Strikethrough,
				extension.Table,
				extension.TaskList,
				extension.DefinitionList,
				mathjax.MathJax,
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
			)),
		pCtx: pCtx,
	}
}

// DefaultRenderer is a markdown parser that supports wikilinks.
type DefaultRenderer struct {
	Markdown Converter
	pCtx     parser.Context
}

// Convert converts the provided markdown to HTML.
func (m *DefaultRenderer) Convert(asset Asset) (gen.Embedded, error) {
	var (
		buf      = bytes.NewBufferString("")
		metadata *frontmatter.Data
		emb      gen.Embedded
		err      error
	)
	err = m.Markdown.Convert(asset.Data, buf, parser.WithContext(m.pCtx))
	if err != nil {
		return emb, eris.Wrapf(
			err,
			"failed to convert %s's markdown to HTML",
			asset.Path,
		)
	}

	// Get frontmatter
	metadata = frontmatter.Get(m.pCtx)
	if metadata == nil {
		return emb, eris.Errorf(
			"frontmatter is nil for %s",
			asset.Path,
		)
	}

	// Decode frontmatter
	err = metadata.Decode(&emb)
	if err != nil {
		return emb, eris.Wrapf(
			err,
			"failed to decode frontmatter of %s",
			asset.Path,
		)
	}

	// Set slug and content
	emb.Slug = asset.Slug
	emb.Content = buf.String()
	emb.RawContent = string(asset.Data)

	// Get frontmatter
	metadata = frontmatter.Get(m.pCtx)
	if metadata == nil {
		return emb, eris.Errorf(
			"frontmatter is nil for %s",
			asset.Path,
		)
	}

	// Decode frontmatter
	err = metadata.Decode(&emb)
	if err != nil {
		return emb, eris.Wrapf(
			err,
			"failed to decode frontmatter of %s",
			asset.Path,
		)
	}
	return emb, nil
}

// resolver is a wikilink.Resolver that resolves pages and media referenced by
// wikilinks to their destinations.
type resolver struct {
	Ctx    context.Context
	client awsClient
}

// newResolver creates a new wikilink resolver.
func newResolver(ctx context.Context, client credited.AWSClient) *resolver {
	return &resolver{client: client, Ctx: ctx}
}

// ResolveWikilink returns the address of the page that the provided
// wikilink points to. The destination will be URL-escaped before
// being placed into a link.
func (r *resolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	var (
		err       error
		targetStr = string(n.Target)
	)
	_, err = r.client.GetObject(r.Ctx, &s3.GetObjectInput{
		Bucket: aws.String("conneroh"),
		Key:    aws.String(targetStr),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object(%s): %w", targetStr, err)
	}
	return fmt.Appendf(nil,
		"https://conneroh.fly.storage.tigris.dev/%s",
		targetStr,
	), nil
}
