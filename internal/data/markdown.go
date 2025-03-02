package data

import (
	"bytes"

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
	// Title is the title of the document.
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

	// Content is the content of the document.
	Content []byte `yaml:"-"`
}

func parse(b []byte) (*Markdown, error) {
	var fm Markdown
	buf := bytes.NewBufferString("")
	ctx := parser.NewContext()
	err := md.Convert(b, buf, parser.WithContext(ctx))
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

	if fm.Description == "" {
		// It is a tag page.
		fm.Description = buf.String()
	}
	return &fm, nil
}
