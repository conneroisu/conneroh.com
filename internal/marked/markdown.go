package marked

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conneroisu/conneroh.com/internal/credited"
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

var exts = []goldmark.Option{
	goldmark.WithExtensions(
		extension.GFM, extension.Footnote,
		extension.Strikethrough, extension.Table,
		extension.TaskList, extension.DefinitionList,
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
}

// NewMDParser creates a new markdown parser.
func NewMDParser(
	client *credited.AWSClient,
) goldmark.Markdown {
	exts = append(exts, goldmark.WithExtensions(
		&wikilink.Extender{
			Resolver: newResolver(client),
		},
	))
	return goldmark.New(exts...)
}

// resolver is a wikilink.Resolver that resolves pages and media referenced by
// wikilinks to their destinations.
type resolver struct{ client *credited.AWSClient }

// newResolver creates a new wikilink resolver.
func newResolver(client *credited.AWSClient) *resolver { return &resolver{client: client} }

// ResolveWikilink returns the address of the page that the provided
// wikilink points to. The destination will be URL-escaped before
// being placed into a link.
func (c *resolver) ResolveWikilink(n *wikilink.Node) (destination []byte, err error) {
	targetStr := string(n.Target)
	_, err = c.client.GetObject(context.Background(), &s3.GetObjectInput{
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
