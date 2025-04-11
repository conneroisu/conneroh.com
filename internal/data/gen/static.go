package gen

import (
	"fmt"
	"strings"
	"time"

	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v3"
)

const (
	// EmbedLength is the length of the full embedding.
	EmbedLength = 768
)

var (
	_ Embeddable = (*Post)(nil)
	_ Embeddable = (*Project)(nil)
	_ Embeddable = (*Tag)(nil)
	_ Embeddable = (*Employment)(nil)
)

// CustomTime allows us to customize the YAML time parsing.
type CustomTime struct{ time.Time }

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (ct *CustomTime) UnmarshalYAML(value *yaml.Node) error {
	// Try parsing as date-only format first
	t, err := time.Parse("2006-01-02", value.Value)
	if err == nil {
		ct.Time = t
		return nil
	}

	// If that fails, try the standard RFC3339 format
	t, err = time.Parse(time.RFC3339, value.Value)
	if err != nil {
		return fmt.Errorf("cannot parse %q as date: %v", value.Value, err)
	}

	ct.Time = t
	return nil
}

type (
	// Embeddable is an interface for embedding content.
	Embeddable interface {
		GetEmb() *Embedded
		PagePath() string
	}
	// Post is a post with all its projects and tags.
	Post struct {
		Embedded
	}
	// Project is a project with all its posts and tags.
	Project struct {
		Embedded
	}
	// Tag is a tag with all its posts and projects.
	Tag struct {
		Embedded
	}
	// Employment is an employment of a tag.
	Employment struct {
		Embedded
	}

	// Embedded is a base struct for all embeddedable structs.
	Embedded struct {
		Title           string `yaml:"title"`
		Slug            string `yaml:"slug"`
		Description     string `yaml:"description"`
		Content         string
		BannerPath      string `yaml:"banner_path"`
		RawContent      string
		Icon            string    `yaml:"icon"`
		CreatedAt       time.Time `yaml:"created_at"`
		UpdatedAt       time.Time `yaml:"updated_at"`
		X               float64
		Y               float64
		Z               float64
		Vec             [EmbedLength]float64
		TagSlugs        []string      `yaml:"tags"`
		PostSlugs       []string      `yaml:"posts"`
		ProjectSlugs    []string      `yaml:"projects"`
		EmploymentSlugs []string      `yaml:"employments"`
		Posts           []*Post       `yaml:"-" structgen:"PostSlugs"`
		Tags            []*Tag        `yaml:"-" structgen:"TagSlugs"`
		Projects        []*Project    `yaml:"-" structgen:"ProjectSlugs"`
		Employments     []*Employment `yaml:"-" structgen:"EmploymentSlugs"`
	}
)

// GetEmb returns the embedding struct itself.
func (emb *Embedded) GetEmb() *Embedded {
	return emb
}

// PagePath returns the path to the post page.
func (emb *Post) PagePath() string {
	return "/post/" + emb.Slug
}

// PagePath returns the path to the project page.
func (emb *Project) PagePath() string {
	return "/project/" + emb.Slug
}

// PagePath returns the path to the tag page.
func (emb *Tag) PagePath() string {
	return "/tag/" + emb.Slug
}

// PagePath returns the path to the employment page.
func (emb *Employment) PagePath() string {
	return "/employment/" + emb.Slug
}

// New creates a new instance of the given type.
func New[
	T Post | Project | Tag,
](emb *Embedded) *T {
	return &T{Embedded: *emb}
}

var (
	// ErrValueMissing is returned when a value is missing.
	ErrValueMissing = eris.Errorf("missing value")

	// ErrValueInvalid is returned when the slug is invalid.
	ErrValueInvalid = eris.Errorf("invalid value")
)

// Defaults sets the default values for the embedding if they are missing.
func Defaults(emb *Embedded) error {
	if emb == nil {
		return eris.Wrap(
			ErrValueMissing,
			"whole embedding is nil",
		)
	}
	// Set default icon if not provided
	if emb.Icon == "" {
		emb.Icon = "tag"
	}

	return nil
}

// Validate validate the given embedding.
func Validate(
	emb *Embedded,
) error {
	if emb.Title == "" {
		return eris.Wrapf(
			ErrValueMissing,
			"%s is missing title",
			emb.RawContent,
		)
	}

	if emb.Slug == "" {
		return eris.Wrapf(
			ErrValueMissing,
			"%s is missing slug",
			emb.Title,
		)
	}

	if emb.Description == "" {
		return eris.Wrapf(
			ErrValueMissing,
			"%s is missing description",
			emb.Slug,
		)
	}

	if emb.Content == "" {
		return eris.Wrapf(
			ErrValueMissing,
			"%s is missing content",
			emb.Slug,
		)
	}

	if emb.RawContent == "" {
		return eris.Wrapf(
			ErrValueMissing,
			"%s is missing raw content",
			emb.Slug,
		)
	}

	if strings.Contains(emb.Slug, " ") {
		return eris.Wrapf(
			ErrValueInvalid,
			"slug %s contains spaces",
			emb.Slug,
		)
	}
	return nil
}

// GetTitle returns the title of the embedding.
func (emb *Embedded) GetTitle() string {
	return emb.Title
}
