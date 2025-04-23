package assets

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/rotisserie/eris"
	"github.com/uptrace/bun"
	"gopkg.in/yaml.v3"
)

// DBName returns the name/file of the database.
func DBName() string {
	return "file:master.db?_pragma=busy_timeout=5000&_pragma=journal_mode=WAL&_pragma=mmap_size=30000000000&_pragma=page_size=32768&_pragma=synchronous=FULL"
}

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
		Title        string     `yaml:"title"`
		Slug         string     `yaml:"slug"`
		Description  string     `yaml:"description"`
		Content      string     `yaml:"-"`
		BannerPath   string     `yaml:"banner_path"`
		RawContent   []byte     `yaml:"-"`
		Icon         string     `yaml:"icon"`
		CreatedAt    CustomTime `yaml:"created_at"`
		UpdatedAt    CustomTime `yaml:"updated_at"`
		TagSlugs     []string   `yaml:"tags"`
		PostSlugs    []string   `yaml:"posts"`
		ProjectSlugs []string   `yaml:"projects"`
		Hash         string     `bun:"hash" yaml:"-"`
		X            float64    `json:"x"`
		Y            float64    `json:"y"`
		Z            float64    `json:"z"`
		Posts        []*Post    `yaml:"-"`
		Tags         []*Tag     `yaml:"-"`
		Projects     []*Project `yaml:"-"`
	}
)

type (
	// Post is a post with all its projects and tags.
	Post struct {
		bun.BaseModel `bun:"posts"`

		ID int64 `bun:"id,pk,autoincrement" yaml:"-"`

		Hash string  `bun:"hashed,unique" yaml:"-"`
		X    float64 `json:"x"`
		Y    float64 `json:"y"`
		Z    float64 `json:"z"`

		Title       string     `yaml:"title"`
		Slug        string     `yaml:"slug"`
		Description string     `yaml:"description"`
		Content     string     `yaml:"-"`
		BannerPath  string     `yaml:"banner_path"`
		RawContent  []byte     `yaml:"-"`
		Icon        string     `yaml:"icon"`
		CreatedAt   CustomTime `yaml:"created_at"`
		UpdatedAt   CustomTime `yaml:"updated_at"`

		TagSlugs     []string `bun:"tag_slugs"`
		PostSlugs    []string `bun:"post_slugs"`
		ProjectSlugs []string `bun:"project_slugs"`

		// M2M relationships
		Tags     []*Tag     `yaml:"-" bun:"m2m:post_to_tags,join:Post=Tag"`
		Posts    []*Post    `yaml:"-" bun:"m2m:post_to_posts,join:SourcePost=TargetPost"`
		Projects []*Project `yaml:"-" bun:"m2m:post_to_projects,join:Post=Project"`
	}

	// Project is a project with all its posts and tags.
	Project struct {
		bun.BaseModel `bun:"projects"`

		ID int64 `bun:"id,pk,autoincrement" yaml:"-"`

		Hash string  `bun:"hashed,unique" yaml:"-"`
		X    float64 `json:"x"`
		Y    float64 `json:"y"`
		Z    float64 `json:"z"`

		Title       string     `yaml:"title"`
		Slug        string     `yaml:"slug"`
		Description string     `yaml:"description"`
		Content     string     `yaml:"-"`
		BannerPath  string     `yaml:"banner_path"`
		RawContent  []byte     `yaml:"-"`
		Icon        string     `yaml:"icon"`
		CreatedAt   CustomTime `yaml:"created_at"`
		UpdatedAt   CustomTime `yaml:"updated_at"`

		TagSlugs     []string `bun:"tag_slugs"`
		PostSlugs    []string `bun:"post_slugs"`
		ProjectSlugs []string `bun:"project_slugs"`

		// M2M relationships
		Tags     []*Tag     `yaml:"-" bun:"m2m:project_to_tags,join:Project=Tag"`
		Posts    []*Post    `yaml:"-" bun:"m2m:post_to_projects,join:Project=Post"`
		Projects []*Project `yaml:"-" bun:"m2m:project_to_projects,join:SourceProject=TargetProject"`
	}

	// Tag is a tag with all its posts and projects.
	Tag struct {
		bun.BaseModel `bun:"tags"`

		ID int64 `bun:"id,pk,autoincrement"`

		Hash string  `bun:"hashed,unique"`
		X    float64 `json:"x"`
		Y    float64 `json:"y"`
		Z    float64 `json:"z"`

		Title       string
		Slug        string
		Description string
		Content     string
		BannerPath  string
		RawContent  []byte
		Icon        string
		CreatedAt   CustomTime
		UpdatedAt   CustomTime

		TagSlugs     []string `bun:"tag_slugs"`
		PostSlugs    []string `bun:"post_slugs"`
		ProjectSlugs []string `bun:"project_slugs"`

		// M2M relationships
		Tags     []*Tag     `yaml:"-" bun:"m2m:tag_to_tags,join:SourceTag=TargetTag"`
		Posts    []*Post    `yaml:"-" bun:"m2m:post_to_tags,join:Tag=Post"`
		Projects []*Project `yaml:"-" bun:"m2m:project_to_tags,join:Tag=Project"`
	}

	// Asset is a media asset.
	Asset struct {
		bun.BaseModel `bun:"assets"`

		ID   int64  `bun:"id,pk,autoincrement"`
		Path string `yaml:"-" bun:"path,unique"`
		Hash string `bun:"hashed,unique"`

		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	}
)

type (
	// PostToTag represents a many-to-many relationship between posts and tags
	PostToTag struct {
		bun.BaseModel `bun:"post_to_tags"`

		PostID int64 `bun:",pk"`
		Post   *Post `bun:"rel:belongs-to,join:post_id=id"`
		TagID  int64 `bun:",pk"`
		Tag    *Tag  `bun:"rel:belongs-to,join:tag_id=id"`
	}

	// PostToPost represents a many-to-many relationship between posts and other posts
	PostToPost struct {
		bun.BaseModel `bun:"post_to_posts"`

		SourcePostID int64 `bun:",pk"`
		SourcePost   *Post `bun:"rel:belongs-to,join:source_post_id=id"`
		TargetPostID int64 `bun:",pk"`
		TargetPost   *Post `bun:"rel:belongs-to,join:target_post_id=id"`
	}

	// PostToProject represents a many-to-many relationship between posts and projects
	PostToProject struct {
		bun.BaseModel `bun:"post_to_projects"`

		PostID    int64    `bun:",pk"`
		Post      *Post    `bun:"rel:belongs-to,join:post_id=id"`
		ProjectID int64    `bun:",pk"`
		Project   *Project `bun:"rel:belongs-to,join:project_id=id"`
	}

	// ProjectToTag represents a many-to-many relationship between projects and tags
	ProjectToTag struct {
		bun.BaseModel `bun:"project_to_tags"`

		ProjectID int64    `bun:",pk"`
		Project   *Project `bun:"rel:belongs-to,join:project_id=id"`
		TagID     int64    `bun:",pk"`
		Tag       *Tag     `bun:"rel:belongs-to,join:tag_id=id"`
	}

	// ProjectToProject represents a many-to-many relationship between projects and other projects
	ProjectToProject struct {
		bun.BaseModel `bun:"project_to_projects"`

		SourceProjectID int64    `bun:",pk"`
		SourceProject   *Project `bun:"rel:belongs-to,join:source_project_id=id"`
		TargetProjectID int64    `bun:",pk"`
		TargetProject   *Project `bun:"rel:belongs-to,join:target_project_id=id"`
	}

	// TagToTag represents a many-to-many relationship between tags and other tags
	TagToTag struct {
		bun.BaseModel `bun:"tag_to_tags"`

		SourceTagID int64 `bun:",pk"`
		SourceTag   *Tag  `bun:"rel:belongs-to,join:source_tag_id=id"`
		TargetTagID int64 `bun:",pk"`
		TargetTag   *Tag  `bun:"rel:belongs-to,join:target_tag_id=id"`
	}

	// TagToPost represents a many-to-many relationship between tags and posts
	TagToPost struct {
		bun.BaseModel `bun:"tag_to_posts"`

		TagID  int64 `bun:",pk"`
		Tag    *Tag  `bun:"rel:belongs-to,join:tag_id=id"`
		PostID int64 `bun:",pk"`
		Post   *Post `bun:"rel:belongs-to,join:post_id=id"`
	}

	// TagToProject represents a many-to-many relationship between tags and projects
	TagToProject struct {
		bun.BaseModel `bun:"tag_to_projects"`

		TagID     int64    `bun:",pk"`
		Tag       *Tag     `bun:"rel:belongs-to,join:tag_id=id"`
		ProjectID int64    `bun:",pk"`
		Project   *Project `bun:"rel:belongs-to,join:project_id=id"`
	}
)

// Defaults sets the default values for the document if they are missing.
func Defaults(doc *Doc) error {
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
	emb *Doc,
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

	if len(emb.RawContent) == 0 {
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

// Hash calculates the hash of a file's content.
func Hash(content []byte) string {
	sum := md5.Sum(content)
	return hex.EncodeToString(sum[:])
}
