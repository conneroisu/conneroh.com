package markdown

import (
	"fmt"
	"strings"

	mathjax "github.com/litao91/goldmark-mathjax"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/quail-ink/goldmark-enclave/core"
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
	_, err = afero.ReadFile(r.fs, fmt.Sprintf("/assets/%s", targetStr))
	if err != nil {
		return nil, err
	}

	return fmt.Appendf(nil,
		"https://conneroh.fly.storage.tigris.dev/%s",
		targetStr,
	), nil
}
