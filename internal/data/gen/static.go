package gen

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// EmbedLength is the length of the embedding.
	EmbedLength = 768
)

// CustomTime allows us to customize the YAML time parsing
type CustomTime struct {
	time.Time
}

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

// New creates a new instance of the given type.
func New[
	T Post | Project | Tag,
](emb *Embedded) *T {
	return &T{Embedded: *emb}
}
