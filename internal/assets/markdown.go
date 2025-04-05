package assets

import (
	"bytes"

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

// MDParser is a markdown parser that supports wikilinks.
type MDParser struct {
	goldmark.Markdown
	pCtx parser.Context
}

// NewMDParser creates a new markdown parser.
func NewMDParser(
	pCtx parser.Context,
	client *credited.AWSClient,
) *MDParser {
	return &MDParser{
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
					Resolver: newResolver(client),
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

// Convert converts the provided markdown to HTML.
func (m *MDParser) Convert(
	asset Asset,
) (emb gen.Embedded, err error) {
	var (
		buf      = bytes.NewBufferString("")
		metadata *frontmatter.Data
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
