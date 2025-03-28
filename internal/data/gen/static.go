package gen

import (
	"fmt"
	"time"

	"slices"

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
		Title           string               `yaml:"title"`
		Slug            string               `yaml:"slug"`
		Description     string               `yaml:"description"`
		Content         string               `yaml:"-"`
		BannerPath      string               `yaml:"banner_path"`
		RawContent      string               `yaml:"-"`
		Icon            string               `yaml:"icon"`
		CreatedAt       time.Time            `yaml:"created_at"`
		UpdatedAt       time.Time            `yaml:"updated_at"`
		X               float64              `yaml:"-"`
		Y               float64              `yaml:"-"`
		Z               float64              `yaml:"-"`
		Vec             [EmbedLength]float64 `yaml:"-"`
		TagSlugs        []string             `yaml:"tags"`
		PostSlugs       []string             `yaml:"posts"`
		ProjectSlugs    []string             `yaml:"projects"`
		EmploymentSlugs []string             `yaml:"employments"`
		Posts           []*Post              `yaml:"-" structgen:"PostSlugs"`
		Tags            []*Tag               `yaml:"-" structgen:"TagSlugs"`
		Projects        []*Project           `yaml:"-" structgen:"ProjectSlugs"`
		Employments     []*Employment        `yaml:"-" structgen:"EmploymentSlugs"`
	}
)

// New creates a new instance of the given type.
func New[
	T Post | Project | Tag,
](emb Embedded) *T {
	return &T{
		Embedded: emb,
	}
}

func init() {
	for _, t := range AllTags {
		for _, a := range AllPosts {
			// for each post that has this tag add it the tag's posts
			if !slices.Contains(a.TagSlugs, t.Slug) {
				a.TagSlugs = append(a.TagSlugs, t.Slug)
			}
		}
		for _, p := range AllProjects {
			// for each project that has this tag add it the tag's projects
			if !slices.Contains(p.TagSlugs, t.Slug) {
				p.TagSlugs = append(p.TagSlugs, t.Slug)
			}
		}
	}
}
