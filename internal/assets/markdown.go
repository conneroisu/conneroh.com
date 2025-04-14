package assets

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/conneroisu/conneroh.com/internal/data/gen"
	mathjax "github.com/litao91/goldmark-mathjax"
	enclave "github.com/quail-ink/goldmark-enclave"
	"github.com/quail-ink/goldmark-enclave/core"
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

// NewRenderer creates a new markdown parser.
func NewRenderer(
	ctx context.Context,
	pCtx parser.Context,
	fs afero.Fs,
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
					Resolver: newResolver(ctx, fs),
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
	Markdown goldmark.Markdown
	pCtx     parser.Context
}

// Convert converts the provided markdown to HTML.
func (m *DefaultRenderer) Convert(
	path string,
	data []byte,
) (gen.Embedded, error) {
	var (
		buf      = bytes.NewBufferString("")
		metadata *frontmatter.Data
		emb      gen.Embedded
		err      error
	)
	err = m.Markdown.Convert(data, buf, parser.WithContext(m.pCtx))
	if err != nil {
		return emb, eris.Wrapf(
			err,
			"failed to convert %s's markdown to HTML",
			path,
		)
	}

	// Get frontmatter
	metadata = frontmatter.Get(m.pCtx)
	if metadata == nil {
		return emb, eris.Errorf(
			"frontmatter is nil for %s",
			path,
		)
	}

	// Decode frontmatter
	err = metadata.Decode(&emb)
	if err != nil {
		return emb, eris.Wrapf(
			err,
			"failed to decode frontmatter of %s",
			path,
		)
	}

	// Set slug and content
	emb.Slug = slugify(path)
	emb.Content = buf.String()
	emb.RawContent = string(data)

	// Get frontmatter
	metadata = frontmatter.Get(m.pCtx)
	if metadata == nil {
		return emb, eris.Errorf(
			"frontmatter is nil for %s",
			path,
		)
	}

	// Decode frontmatter
	err = metadata.Decode(&emb)
	if err != nil {
		return emb, eris.Wrapf(
			err,
			"failed to decode frontmatter of %s",
			path,
		)
	}
	return emb, nil
}

// resolver is a wikilink.Resolver that resolves pages and media referenced by
// wikilinks to their destinations.
type resolver struct {
	Ctx context.Context
	fs  afero.Fs
}

// newResolver creates a new wikilink resolver.
func newResolver(
	ctx context.Context,
	fs afero.Fs,
) *resolver {
	return &resolver{
		fs:  fs,
		Ctx: ctx,
	}
}

// ResolveWikilink returns the address of the page that the provided
// wikilink points to. The destination will be URL-escaped before
// being placed into a link.
func (r *resolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	var (
		err       error
		targetStr = string(n.Target)
	)
	_, err = afero.ReadFile(r.fs, fmt.Sprintf("/assets/%s", targetStr))
	if err != nil {
		return nil, err
	}
	return fmt.Appendf(nil,
		"https://conneroh.fly.storage.tigris.dev/%s",
		targetStr,
	), nil
}

func slugify(s string) string {
	var path string
	var ok bool
	path, ok = strings.CutPrefix(s, assetsLoc)
	if ok {
		return path
	}
	return strings.TrimSuffix(pathify(s), filepath.Ext(s))
}
func pathify(s string) string {
	var path string
	var ok bool
	path, ok = strings.CutPrefix(s, postsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, projectsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, tagsLoc)
	if ok {
		return path
	}
	path, ok = strings.CutPrefix(s, assetsLoc)
	if ok {
		return path
	}
	return s // Return original path instead of panicking
}
