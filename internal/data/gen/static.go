package gen

import "time"

//go:generate gomarkdoc -o README.md -e .

const (
	// EmbedLength is the length of the embedding.
	EmbedLength = 768
)

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

	// Embedded is a base struct for all embedded structs.
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
