package assets

import (
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"gopkg.in/yaml.v3"
)

// DBName returns the name/file of the database.
func DBName() string {
	return "file:master.db"
}

// CustomTime allows us to customize the YAML time parsing.
type CustomTime struct{ time.Time }

// UnmarshalYAML implements the yaml.Unmarshaler interface.
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
		Title           string        `yaml:"title"`
		Path            string        `yaml:"-"`
		Slug            string        `yaml:"slug"`
		Description     string        `yaml:"description"`
		Content         string        `yaml:"-"`
		BannerPath      string        `yaml:"banner_path"`
		Icon            string        `yaml:"icon"`
		CreatedAt       CustomTime    `yaml:"created_at"`
		UpdatedAt       CustomTime    `yaml:"updated_at"`
		EndDate         *CustomTime   `yaml:"end_date"`
		TagSlugs        []string      `yaml:"tags"`
		PostSlugs       []string      `yaml:"posts"`
		ProjectSlugs    []string      `yaml:"projects"`
		EmploymentSlugs []string      `yaml:"employments"`
		Hash            string        `yaml:"-"`
		Posts           []*Post       `yaml:"-"`
		Tags            []*Tag        `yaml:"-"`
		Projects        []*Project    `yaml:"-"`
		Employments     []*Employment `yaml:"-"`
	}
	// Cache represents an asset cache.
	Cache struct {
		bun.BaseModel `bun:"caches"`

		ID   int64  `bun:"id,pk,autoincrement"`
		Path string `bun:"path,unique"`
		Hash string `bun:"hashed,unique"`
	}
	// Post is a post with all its projects and tags.
	Post struct {
		bun.BaseModel `bun:"posts"`

		ID int64 `bun:"id,pk,autoincrement" `

		Title       string     `bun:"title"`
		Slug        string     `bun:"slug,unique"`
		Description string     `bun:"description"`
		Content     string     `bun:"content"`
		BannerPath  string     `bun:"banner_path"`
		CreatedAt   CustomTime `bun:"created_at"`

		TagSlugs        []string
		PostSlugs       []string
		ProjectSlugs    []string
		EmploymentSlugs []string

		// M2M relationships
		Tags        []*Tag        `bun:"m2m:post_to_tags,join:Post=Tag"`
		Posts       []*Post       `bun:"m2m:post_to_posts,join:SourcePost=TargetPost"`
		Projects    []*Project    `bun:"m2m:post_to_projects,join:Post=Project"`
		Employments []*Employment `bun:"m2m:employment_to_posts,join:Post=Employment"`
	}

	// Project is a project with all its posts and tags.
	Project struct {
		bun.BaseModel `bun:"projects"`

		ID int64 `bun:"id,pk,autoincrement" yaml:"-"`

		Title       string     `bun:"title"`
		Slug        string     `bun:"slug,unique"`
		Description string     `bun:"description"`
		Content     string     `bun:"content"`
		BannerPath  string     `bun:"banner_path"`
		CreatedAt   CustomTime `bun:"created_at"`

		TagSlugs        []string `bun:"tag_slugs"`
		PostSlugs       []string `bun:"post_slugs"`
		ProjectSlugs    []string `bun:"project_slugs"`
		EmploymentSlugs []string `bun:"employment_slugs"`

		// M2M relationships
		Tags        []*Tag        `bun:"m2m:project_to_tags,join:Project=Tag"`
		Posts       []*Post       `bun:"m2m:post_to_projects,join:Project=Post"`
		Projects    []*Project    `bun:"m2m:project_to_projects,join:SourceProject=TargetProject"`
		Employments []*Employment `bun:"m2m:employment_to_projects,join:Project=Employment"`
	}

	// Tag is a tag with all its posts and projects.
	Tag struct {
		bun.BaseModel `bun:"tags"`

		ID int64 `bun:"id,pk,autoincrement"`

		Title       string     `bun:"title"`
		Slug        string     `bun:"slug,unique"`
		Description string     `bun:"description"`
		Content     string     `bun:"content"`
		BannerPath  string     `bun:"banner_path"`
		Icon        string     `bun:"icon"`
		CreatedAt   CustomTime `bun:"created_at"`

		TagSlugs        []string `bun:"tag_slugs"`
		PostSlugs       []string `bun:"post_slugs"`
		ProjectSlugs    []string `bun:"project_slugs"`
		EmploymentSlugs []string `bun:"employment_slugs"`

		// M2M relationships
		Tags        []*Tag        `bun:"m2m:tag_to_tags,join:SourceTag=TargetTag"`
		Posts       []*Post       `bun:"m2m:post_to_tags,join:Tag=Post"`
		Projects    []*Project    `bun:"m2m:project_to_tags,join:Tag=Project"`
		Employments []*Employment `bun:"m2m:employment_to_tags,join:Tag=Employment"`
	}

	// Employment is an employment with all its posts, projects, and tags.
	Employment struct {
		bun.BaseModel `bun:"employments"`

		ID int64 `bun:"id,pk,autoincrement"`

		Title       string      `bun:"title"`
		Slug        string      `bun:"slug,unique"`
		Description string      `bun:"description"`
		Content     string      `bun:"content"`
		BannerPath  string      `bun:"banner_path"`
		CreatedAt   CustomTime  `bun:"created_at"`
		EndDate     *CustomTime `bun:"end_date"`

		TagSlugs        []string `bun:"tag_slugs"`
		PostSlugs       []string `bun:"post_slugs"`
		ProjectSlugs    []string `bun:"project_slugs"`
		EmploymentSlugs []string `bun:"employment_slugs"`

		// M2M relationships
		Tags        []*Tag        `bun:"m2m:employment_to_tags,join:Employment=Tag"`
		Posts       []*Post       `bun:"m2m:employment_to_posts,join:Employment=Post"`
		Projects    []*Project    `bun:"m2m:employment_to_projects,join:Employment=Project"`
		Employments []*Employment `bun:"m2m:employment_to_employments,join:SourceEmployment=TargetEmployment"`
	}

	// PostToTag represents a many-to-many relationship between posts and tags.
	PostToTag struct {
		bun.BaseModel `bun:"post_to_tags"`

		PostID int64 `bun:"post_id,pk"`
		Post   *Post `bun:"rel:belongs-to,join:post_id=id"`
		TagID  int64 `bun:"tag_id,pk"`
		Tag    *Tag  `bun:"rel:belongs-to,join:tag_id=id"`
	}

	// PostToProject represents a many-to-many relationship between posts and projects.
	PostToProject struct {
		bun.BaseModel `bun:"post_to_projects"`

		PostID    int64    `bun:"post_id,pk"`
		Post      *Post    `bun:"rel:belongs-to,join:post_id=id"`
		ProjectID int64    `bun:"project_id,pk"`
		Project   *Project `bun:"rel:belongs-to,join:project_id=id"`
	}

	// ProjectToTag represents a many-to-many relationship between projects and tags.
	ProjectToTag struct {
		bun.BaseModel `bun:"project_to_tags"`

		ProjectID int64    `bun:"project_id,pk"`
		Project   *Project `bun:"rel:belongs-to,join:project_id=id"`
		TagID     int64    `bun:"tag_id,pk"`
		Tag       *Tag     `bun:"rel:belongs-to,join:tag_id=id"`
	}

	// PostToPost represents a many-to-many relationship between posts and other posts.
	PostToPost struct {
		bun.BaseModel `bun:"post_to_posts"`

		SourcePostID int64 `bun:"source_post_id,pk"`
		SourcePost   *Post `bun:"rel:belongs-to,join:source_post_id=id"`
		TargetPostID int64 `bun:"target_post_id,pk"`
		TargetPost   *Post `bun:"rel:belongs-to,join:target_post_id=id"`
	}

	// ProjectToProject represents a many-to-many relationship between projects and other projects.
	ProjectToProject struct {
		bun.BaseModel `bun:"project_to_projects"`

		SourceProjectID int64    `bun:"source_project_id,pk"`
		SourceProject   *Project `bun:"rel:belongs-to,join:source_project_id=id"`
		TargetProjectID int64    `bun:"target_project_id,pk"`
		TargetProject   *Project `bun:"rel:belongs-to,join:target_project_id=id"`
	}

	// TagToTag represents a many-to-many relationship between tags and other tags.
	TagToTag struct {
		bun.BaseModel `bun:"tag_to_tags"`

		SourceTagID int64 `bun:"source_tag_id,pk"`
		SourceTag   *Tag  `bun:"rel:belongs-to,join:source_tag_id=id"`
		TargetTagID int64 `bun:"target_tag_id,pk"`
		TargetTag   *Tag  `bun:"rel:belongs-to,join:target_tag_id=id"`
	}

	// EmploymentToTag represents a many-to-many relationship between employments and tags.
	EmploymentToTag struct {
		bun.BaseModel `bun:"employment_to_tags"`

		EmploymentID int64       `bun:"employment_id,pk"`
		Employment   *Employment `bun:"rel:belongs-to,join:employment_id=id"`
		TagID        int64       `bun:"tag_id,pk"`
		Tag          *Tag        `bun:"rel:belongs-to,join:tag_id=id"`
	}

	// EmploymentToPost represents a many-to-many relationship between employments and posts.
	EmploymentToPost struct {
		bun.BaseModel `bun:"employment_to_posts"`

		EmploymentID int64       `bun:"employment_id,pk"`
		Employment   *Employment `bun:"rel:belongs-to,join:employment_id=id"`
		PostID       int64       `bun:"post_id,pk"`
		Post         *Post       `bun:"rel:belongs-to,join:post_id=id"`
	}

	// EmploymentToProject represents a many-to-many relationship between employments and projects.
	EmploymentToProject struct {
		bun.BaseModel `bun:"employment_to_projects"`

		EmploymentID int64       `bun:"employment_id,pk"`
		Employment   *Employment `bun:"rel:belongs-to,join:employment_id=id"`
		ProjectID    int64       `bun:"project_id,pk"`
		Project      *Project    `bun:"rel:belongs-to,join:project_id=id"`
	}

	// EmploymentToEmployment represents a many-to-many relationship between employments and other employments.
	EmploymentToEmployment struct {
		bun.BaseModel `bun:"employment_to_employments"`

		SourceEmploymentID int64       `bun:"source_employment_id,pk"`
		SourceEmployment   *Employment `bun:"rel:belongs-to,join:source_employment_id=id"`
		TargetEmploymentID int64       `bun:"target_employment_id,pk"`
		TargetEmployment   *Employment `bun:"rel:belongs-to,join:target_employment_id=id"`
	}
)

// GetTitle returns the title of the embedding.
func (emb *Doc) GetTitle() string {
	return emb.Title
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

func (emb *Post) String() string {
	return fmt.Sprintf("Post: %s %s %s %d", emb.Title, emb.Slug, emb.Description, emb.ID)
}

func (emb *Project) String() string {
	return fmt.Sprintf("Project: %s %s %s %d", emb.Title, emb.Slug, emb.Description, emb.ID)
}

func (emb *Tag) String() string {
	return fmt.Sprintf("Tag: %s %s %s %d", emb.Title, emb.Slug, emb.Description, emb.ID)
}

func (emb *Employment) String() string {
	return fmt.Sprintf("Employment: %s %s %s %d", emb.Title, emb.Slug, emb.Description, emb.ID)
}
