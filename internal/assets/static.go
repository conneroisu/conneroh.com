package assets

import (
	"fmt"
	"strings"
	"time"

	"github.com/rotisserie/eris"
	"github.com/uptrace/bun"
	"gopkg.in/yaml.v3"
)

const (
	// EmbedLength is the length of the full embedding.
	EmbedLength = 768
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
	// Doc is a base struct for all embeddedable structs.
	Doc struct {
		bun.BaseModel `bun:"docs"`

		ID    int64  `yaml:"-" bun:"id,pk,autoincrement"`
		CType string `yaml:"-" bun:"ctype,notnull"`

		Title        string     `yaml:"title"`
		Slug         string     `yaml:"slug"`
		Description  string     `yaml:"description"`
		Content      string     `yaml:"-"`
		BannerPath   string     `yaml:"banner_path"`
		RawContent   string     `yaml:"-"`
		Icon         string     `yaml:"icon"`
		CreatedAt    CustomTime `yaml:"created_at"`
		UpdatedAt    CustomTime `yaml:"updated_at"`
		TagSlugs     []string   `yaml:"tags"`
		PostSlugs    []string   `yaml:"posts"`
		ProjectSlugs []string   `yaml:"projects"`
		Posts        []*Post    `yaml:"-"`
		Tags         []*Tag     `yaml:"-"`
		Projects     []*Project `yaml:"-"`
	}
	// Emb is a struct for embedding.
	Emb struct {
		bun.BaseModel `bun:"embeddings"`

		Hash string `bun:"hash" yaml:"-"`
		// Hash string  `bun:"hash,pk,unique" yaml:"-"`
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	}
	// Raw is a full raw doc.
	Raw struct {
		ID int64 `yaml:"-"`
		Doc
		Emb
		Posts    []*Post    `yaml:"-"`
		Tags     []*Tag     `yaml:"-"`
		Projects []*Project `yaml:"-"`
	}
)

type (
	// Post is a post with all its projects and tags.
	Post struct {
		bun.BaseModel `bun:"posts"`

		ID int64 `bun:"id,pk,autoincrement" yaml:"-"`
		Doc

		EmbeddingID int64 `yaml:"-"`
		Embbedding  Emb   `yaml:"-"`

		TagSlugs     []string   `yaml:"tags"`
		PostSlugs    []string   `yaml:"posts"`
		ProjectSlugs []string   `yaml:"projects"`
		Posts        []*Post    `yaml:"-" bun:"rel:has-many,join:post_slugs=slug"`
		Tags         []*Tag     `yaml:"-" bun:"rel:has-many,join:tag_slugs=slug"`
		Projects     []*Project `yaml:"-" bun:"rel:has-many,join:project_slugs=slug"`
	}
)

type (
	// Project is a project with all its posts and tags.
	Project struct {
		bun.BaseModel `bun:"projects"`

		ID int64 `bun:"id,pk,autoincrement" yaml:"-"`
		Doc

		EmbbeddingID int64 `yaml:"-"`
		Embedding    Emb   `yaml:"-"`

		TagSlugs     []string   `yaml:"tags"`
		PostSlugs    []string   `yaml:"posts"`
		ProjectSlugs []string   `yaml:"projects"`
		Posts        []*Post    `yaml:"-" bun:"rel:has-many,join:post_slugs=slug"`
		Tags         []*Tag     `yaml:"-" bun:"rel:has-many,join:tag_slugs=slug"`
		Projects     []*Project `yaml:"-" bun:"rel:has-many,join:project_slugs=slug"`
	}
)

type (
	// Tag is a tag with all its posts and projects.
	Tag struct {
		bun.BaseModel `bun:"tags"`

		ID int64 `bun:"id,pk,autoincrement" yaml:"-"`
		Doc
		Embedding Emb `yaml:"-"`

		TagSlugs     []string   `yaml:"tags"`
		PostSlugs    []string   `yaml:"posts"`
		ProjectSlugs []string   `yaml:"projects"`
		Posts        []*Post    `yaml:"-" bun:"rel:has-many,join:post_slugs=slug"`
		Tags         []*Tag     `yaml:"-" bun:"rel:has-many,join:tag_slugs=slug"`
		Projects     []*Project `yaml:"-" bun:"rel:has-many,join:project_slugs=slug"`
	}
)

// Defaults sets the default values for the embedding if they are missing.
func Defaults(doc *Raw) error {
	// Set default icon if not provided
	if doc.Icon == "" {
		doc.Icon = "tag"
	}

	return nil
}

// GetTitle returns the title of the embedding.
func (emb *Doc) GetTitle() string {
	return emb.Title
}

// Validate validate the given embedding.
func Validate(
	path string,
	emb *Raw,
) error {
	errs := []error{}
	if emb.Slug == "" {
		errs = append(errs, eris.Wrapf(
			ErrValueMissing,
			"%s is missing slug",
			path,
		))
	}

	if emb.Description == "" {
		errs = append(errs, eris.Wrapf(
			ErrValueMissing,
			"%s is missing description",
			path,
		))
	}

	if emb.Content == "" {
		errs = append(errs, eris.Wrapf(
			ErrValueMissing,
			"%s is missing content",
			path,
		))
	}

	if emb.RawContent == "" {
		errs = append(errs, eris.Wrapf(
			ErrValueMissing,
			"%s is missing raw content",
			path,
		))
	}

	if strings.Contains(emb.Slug, " ") {
		errs = append(errs, eris.Wrapf(
			ErrValueInvalid,
			"slug %s contains spaces",
			path,
		))
	}

	var err error
	for _, er := range errs {
		if er != nil {
			err = eris.Wrapf(err, "failed validating %s", path)
		}
	}
	if err != nil {
		return err
	}
	return nil
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

func (emb *Post) String() string {
	return fmt.Sprintf("%s %s %s %d", emb.Title, emb.Slug, emb.Description, emb.ID)
}

func (emb *Project) String() string {
	return fmt.Sprintf("%s %s %s %d", emb.Title, emb.Slug, emb.Description, emb.ID)
}

func (emb *Tag) String() string {
	return fmt.Sprintf("%s %s %s %d", emb.Title, emb.Slug, emb.Description, emb.ID)
}
