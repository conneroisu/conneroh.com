package assets

import (
	"bytes"
	"strings"

	callout "github.com/VojtaStruhar/goldmark-obsidian-callout"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	enclave "github.com/quailyquaily/goldmark-enclave"
	"github.com/quailyquaily/goldmark-enclave/core"
	"github.com/rotisserie/eris"
	"github.com/spf13/afero"
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

// NewMD creates a new markdown parser.
func NewMD(
	fs afero.Fs,
) goldmark.Markdown {
	return goldmark.New(goldmark.WithParserOptions(
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
				Resolver: newResolver(fs),
			},
			callout.ObsidianCallout,
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
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					// chromahtml.WithClasses(true),
					chromahtml.WithLinkableLineNumbers(true, "l"),
					chromahtml.WithLineNumbers(true),
					chromahtml.LineNumbersInTable(true),
				),
			),
		))
}

// resolver is a wikilink.Resolver that resolves pages and media referenced by
// wikilinks to their destinations.
type resolver struct{ fs afero.Fs }

// newResolver creates a new wikilink resolver.
func newResolver(fs afero.Fs) *resolver {
	return &resolver{fs: fs}
}

// ResolveWikilink returns the address of the page that the provided
// wikilink points to. The destination will be URL-escaped before
// being placed into a link.
func (r *resolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	var (
		err       error
		targetStr = string(n.Target)
	)
	// first check if it a url
	if strings.HasPrefix(targetStr, "http") {
		return []byte(targetStr), nil
	}
	_, err = afero.ReadFile(r.fs, "/assets/"+targetStr)
	if err != nil {
		return nil, err
	}

	return []byte(BucketPath(targetStr)), nil
}

// ParseMarkdown parses a markdown document.
func ParseMarkdown(md goldmark.Markdown, item DirMatchItem) (*Doc, error) {
	// Parse document
	doc := &Doc{
		Path:    item.Path,
		Content: item.Content,
		Slug:    Slugify(item.Path),
	}

	var buf bytes.Buffer
	// Set default values
	err := Defaults(doc)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to set defaults for document: %s", item.Path)
	}

	// Parse markdown
	err = md.Convert([]byte(doc.Content), &buf)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to parse markdown: %s", doc.Path)
	}
	// Create a new parser context
	pCtx := parser.NewContext()

	// Convert markdown
	err = md.Convert([]byte(item.Content), &buf, parser.WithContext(pCtx))
	if err != nil {
		return nil, eris.Wrapf(err, "failed to convert markdown: %s", doc.Path)
	}

	// Extract frontmatter
	metadata := frontmatter.Get(pCtx)
	if metadata == nil {
		return nil, eris.Errorf("frontmatter is nil for %s", doc.Path)
	}

	if err := metadata.Decode(doc); err != nil {
		return nil, eris.Wrapf(err, "failed to decode frontmatter: %s", doc.Path)
	}

	doc.Content = buf.String()

	return doc, nil
}
