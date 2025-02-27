package md

import (
	"bytes"
	"fmt"

	"github.com/go-playground/validator/v10"
	mathjax "github.com/litao91/goldmark-mathjax"
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

type (
	// FrontMatter is the frontmatter of a markdown document.
	FrontMatter struct {
		Title       string   `yaml:"title" validate:"required"`
		Description string   `yaml:"description" validate:"required"`
		Tags        []string `yaml:"tags" validate:"required"`
	}
)

var (
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Strikethrough,
			extension.Table,
			extension.TaskList,
			extension.DefinitionList,
			extension.NewTypographer(
				extension.WithTypographicSubstitutions(extension.TypographicSubstitutions{
					extension.Apostrophe: []byte("’"),
				}),
			),
			extension.NewFootnote(
				extension.WithFootnoteIDPrefix("fn"),
			),
			mathjax.MathJax,
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
			&wikilink.Extender{},
			highlighting.NewHighlighting(highlighting.WithStyle("monokai")),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
			// parser.WithParagraphTransformers(ps),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			// html.WithXHTML(),
			extension.WithFootnoteBacklinkClass("footnote-backref"),
			extension.WithFootnoteLinkClass("footnote-ref"),
		),
	)
)

// Parse parses markdown to html.
func Parse(
	name string,
	source []byte,
) (string, error) {
	var (
		buf     bytes.Buffer
		ctx     = parser.NewContext()
		fM      FrontMatter
		validor = validator.New(validator.WithRequiredStructEnabled())
	)
	err := md.Convert(source, &buf, parser.WithContext(ctx))
	if err != nil {
		return "", err
	}
	d := frontmatter.Get(ctx)
	if d == nil {
		return "", &FrontMatterMissingError{
			fileName: name,
		}
	}
	err = d.Decode(&fM)
	if err != nil {
		return "", err
	}

	// returns nil or ValidationErrors ( []FieldError )
	err = validor.Struct(fM)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return "", err
		}
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}

		// from here you can create your own error messages in whatever language you wish
		return "", err
	}
	return buf.String(), nil
}

// ParseWithFrontMatter parses markdown to html and decodes the frontmatter into the provided target struct.
func ParseWithFrontMatter[
	T *PostFrontMatter | *ProjectFrontMatter | *TagFrontMatter,
](
	name string,
	source []byte,
	frontMatterTarget T,
) (string, error) {
	var (
		buf bytes.Buffer
		ctx = parser.NewContext()
	)

	err := md.Convert(source, &buf, parser.WithContext(ctx))
	if err != nil {
		return "", err
	}

	d := frontmatter.Get(ctx)
	if d == nil {
		return "", &FrontMatterMissingError{
			fileName: name,
		}
	}

	err = d.Decode(frontMatterTarget)
	if err != nil {
		return "", err
	}

	// Validate the frontmatter struct if it's not nil
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(frontMatterTarget); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return "", err
		}
		return "", err
	}

	return buf.String(), nil
}
