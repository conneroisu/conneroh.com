<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# assets

```go
import "github.com/conneroisu/conneroh.com/internal/assets"
```

Package assets contains the main data.

## Index

- [Constants](<#constants>)
- [Variables](<#variables>)
- [func BucketPath\(path string\) string](<#BucketPath>)
- [func ComputeHash\(content \[\]byte\) string](<#ComputeHash>)
- [func DBName\(\) string](<#DBName>)
- [func Defaults\(doc \*Doc\) error](<#Defaults>)
- [func GetContentType\(path string\) string](<#GetContentType>)
- [func InitDB\(ctx context.Context, db \*bun.DB\) error](<#InitDB>)
- [func IsAllowedAsset\(path string\) bool](<#IsAllowedAsset>)
- [func IsAllowedDocumentType\(path string\) bool](<#IsAllowedDocumentType>)
- [func IsAllowedMediaType\(path string\) bool](<#IsAllowedMediaType>)
- [func NewMD\(fs afero.Fs\) goldmark.Markdown](<#NewMD>)
- [func Pathify\(s string\) string](<#Pathify>)
- [func RegisterModels\(db \*bun.DB\)](<#RegisterModels>)
- [func Slugify\(s string\) string](<#Slugify>)
- [func UploadToS3\(ctx context.Context, client Tigris, bucket string, path string, data \[\]byte\) error](<#UploadToS3>)
- [func Validate\(path string, emb \*Doc\) error](<#Validate>)
- [type Cache](<#Cache>)
- [type CustomTime](<#CustomTime>)
  - [func \(ct \*CustomTime\) UnmarshalYAML\(value \*yaml.Node\) error](<#CustomTime.UnmarshalYAML>)
- [type DefaultTigrisClient](<#DefaultTigrisClient>)
  - [func NewTigris\(getenv func\(string\) string\) \(\*DefaultTigrisClient, error\)](<#NewTigris>)
- [type DirMatchItem](<#DirMatchItem>)
  - [func HashDirMatch\(ctx context.Context, fs afero.Fs, path string, db \*bun.DB\) \(\[\]DirMatchItem, error\)](<#HashDirMatch>)
  - [func MatchItem\(fs afero.Fs, path string\) \(DirMatchItem, error\)](<#MatchItem>)
- [type Doc](<#Doc>)
  - [func ParseMarkdown\(md goldmark.Markdown, item DirMatchItem\) \(\*Doc, error\)](<#ParseMarkdown>)
  - [func \(emb \*Doc\) GetTitle\(\) string](<#Doc.GetTitle>)
- [type Employment](<#Employment>)
  - [func \(emb \*Employment\) PagePath\(\) string](<#Employment.PagePath>)
  - [func \(emb \*Employment\) String\(\) string](<#Employment.String>)
- [type EmploymentToEmployment](<#EmploymentToEmployment>)
- [type EmploymentToPost](<#EmploymentToPost>)
- [type EmploymentToProject](<#EmploymentToProject>)
- [type EmploymentToTag](<#EmploymentToTag>)
- [type Post](<#Post>)
  - [func \(emb \*Post\) PagePath\(\) string](<#Post.PagePath>)
  - [func \(emb \*Post\) String\(\) string](<#Post.String>)
- [type PostToPost](<#PostToPost>)
- [type PostToProject](<#PostToProject>)
- [type PostToTag](<#PostToTag>)
- [type Project](<#Project>)
  - [func \(emb \*Project\) PagePath\(\) string](<#Project.PagePath>)
  - [func \(emb \*Project\) String\(\) string](<#Project.String>)
- [type ProjectToProject](<#ProjectToProject>)
- [type ProjectToTag](<#ProjectToTag>)
- [type RelationshipFn](<#RelationshipFn>)
  - [func UpsertEmployment\(ctx context.Context, db \*bun.DB, employment \*Employment\) \(RelationshipFn, error\)](<#UpsertEmployment>)
  - [func UpsertEmploymentRelationships\(db \*bun.DB, employment \*Employment\) RelationshipFn](<#UpsertEmploymentRelationships>)
  - [func UpsertPost\(ctx context.Context, db \*bun.DB, post \*Post\) \(RelationshipFn, error\)](<#UpsertPost>)
  - [func UpsertPostRelationships\(db \*bun.DB, post \*Post\) RelationshipFn](<#UpsertPostRelationships>)
  - [func UpsertProject\(ctx context.Context, db \*bun.DB, project \*Project\) \(RelationshipFn, error\)](<#UpsertProject>)
  - [func UpsertProjectRelationships\(db \*bun.DB, project \*Project\) RelationshipFn](<#UpsertProjectRelationships>)
  - [func UpsertTag\(ctx context.Context, db \*bun.DB, tag \*Tag\) \(RelationshipFn, error\)](<#UpsertTag>)
  - [func UpsertTagRelationships\(db \*bun.DB, tag \*Tag\) RelationshipFn](<#UpsertTagRelationships>)
- [type Tag](<#Tag>)
  - [func \(emb \*Tag\) PagePath\(\) string](<#Tag.PagePath>)
  - [func \(emb \*Tag\) String\(\) string](<#Tag.String>)
- [type TagToTag](<#TagToTag>)
- [type Tigris](<#Tigris>)


## Constants

<a name="VaultLoc"></a>

```go
const (
    // VaultLoc is the location of the vault.
    // This is the location of the documents and assets.
    VaultLoc = "internal/data/"
    // AssetsLoc is the location of the assets relative to the vault.
    AssetsLoc = "assets/"
    // PostsLoc is the location of the posts relative to the vault.
    PostsLoc = "posts/"
    // TagsLoc is the location of the tags relative to the vault.
    TagsLoc = "tags/"
    // ProjectsLoc is the location of the projects relative to the vault.
    ProjectsLoc = "projects/"
    // EmploymentsLoc is the location of the employments relative to the vault.
    EmploymentsLoc = "employments/"
)
```

## Variables

<a name="EmpPost"></a>

```go
var (
    // EmpPost is a pointer to a Post.
    EmpPost = new(Post)
    // EmpTag is a pointer to a Tag.
    EmpTag = new(Tag)
    // EmpProject is a pointer to a Project.
    EmpProject = new(Project)
    // EmpEmployment is a pointer to an Employment.
    EmpEmployment = new(Employment)
    // EmpCache is a pointer to a Cache.
    EmpCache = new(Cache)
    // EmpPostToTag is a pointer to a PostToTag.
    EmpPostToTag = new(PostToTag)
    // EmpPostToPost is a pointer to a PostToPost.
    EmpPostToPost = new(PostToPost)
    // EmpPostToProject is a pointer to a PostToProject.
    EmpPostToProject = new(PostToProject)
    // EmpProjectToTag is a pointer to a ProjectToTag.
    EmpProjectToTag = new(ProjectToTag)
    // EmpProjectToProject is a pointer to a ProjectToProject.
    EmpProjectToProject = new(ProjectToProject)
    // EmpTagToTag is a pointer to a TagToTag.
    EmpTagToTag = new(TagToTag)
    // EmpEmploymentToTag is a pointer to an EmploymentToTag.
    EmpEmploymentToTag = new(EmploymentToTag)
    // EmpEmploymentToPost is a pointer to an EmploymentToPost.
    EmpEmploymentToPost = new(EmploymentToPost)
    // EmpEmploymentToProject is a pointer to an EmploymentToProject.
    EmpEmploymentToProject = new(EmploymentToProject)
    // EmpEmploymentToEmployment is a pointer to an EmploymentToEmployment.
    EmpEmploymentToEmployment = new(EmploymentToEmployment)
)
```

<a name="ErrValueMissing"></a>

```go
var (
    // ErrValueMissing is returned when a value is missing.
    ErrValueMissing = eris.Errorf("missing value")

    // ErrValueInvalid is returned when the slug is invalid.
    ErrValueInvalid = eris.Errorf("invalid value")

    // ErrMissingCreds is returned when the credentials are missing.
    ErrMissingCreds = errors.New("missing credentials (ENV)")

    // ErrInvalidCreds is returned when the credentials are invalid.
    ErrInvalidCreds = errors.New("invalid credentials (ENV)")
)
```

<a name="AllowedAssetTypes"></a>

```go
var (
    // AllowedAssetTypes is a list of allowed asset types.
    AllowedAssetTypes = []string{
        "image/png",
        "image/jpeg",
        "image/jpg",
        "image/gif",
        "image/webp",
        "image/avif",
        "image/tiff",
        "image/svg+xml",
        "application/pdf",
    }
    // AllowedDocumentTypes is a list of allowed document types.
    AllowedDocumentTypes = []string{
        "text/markdown",
        "text/markdown; charset=utf-8",
    }
)
```

<a name="BucketPath"></a>
## func [BucketPath](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/paths.go#L106>)

```go
func BucketPath(path string) string
```

BucketPath returns the path to the bucket for a given file path.

<a name="ComputeHash"></a>
## func [ComputeHash](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/hash.go#L21>)

```go
func ComputeHash(content []byte) string
```

ComputeHash generates an MD5 hash of the given content. Note: MD5 is not cryptographically secure; only use for content fingerprinting.

<a name="DBName"></a>
## func [DBName](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L12>)

```go
func DBName() string
```

DBName returns the name/file of the database.

<a name="Defaults"></a>
## func [Defaults](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/validation.go#L10>)

```go
func Defaults(doc *Doc) error
```

Defaults sets the default values for the document if they are missing.

<a name="GetContentType"></a>
## func [GetContentType](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/paths.go#L96>)

```go
func GetContentType(path string) string
```

GetContentType returns the content type for a file extension.

<a name="InitDB"></a>
## func [InitDB](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/emp.go#L61-L64>)

```go
func InitDB(ctx context.Context, db *bun.DB) error
```

InitDB initializes the database.

<a name="IsAllowedAsset"></a>
## func [IsAllowedAsset](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/types.go#L38>)

```go
func IsAllowedAsset(path string) bool
```

IsAllowedAsset returns true if the provided path is an allowed asset.

<a name="IsAllowedDocumentType"></a>
## func [IsAllowedDocumentType](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/types.go#L43>)

```go
func IsAllowedDocumentType(path string) bool
```

IsAllowedDocumentType returns true if the provided path is an allowed document type.

<a name="IsAllowedMediaType"></a>
## func [IsAllowedMediaType](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/types.go#L30>)

```go
func IsAllowedMediaType(path string) bool
```

IsAllowedMediaType returns true if the provided path is an allowed asset type.

<a name="NewMD"></a>
## func [NewMD](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/markdown.go#L28-L30>)

```go
func NewMD(fs afero.Fs) goldmark.Markdown
```

NewMD creates a new markdown parser.

<a name="Pathify"></a>
## func [Pathify](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/paths.go#L26>)

```go
func Pathify(s string) string
```

Pathify returns the slugified path of a document or media asset.

<a name="RegisterModels"></a>
## func [RegisterModels](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/emp.go#L79>)

```go
func RegisterModels(db *bun.DB)
```

RegisterModels registers all the M2M relationship models with Bun.

<a name="Slugify"></a>
## func [Slugify](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/paths.go#L38>)

```go
func Slugify(s string) string
```

Slugify returns the path to the document page or media asset page.

<a name="UploadToS3"></a>
## func [UploadToS3](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/s3.go#L65-L71>)

```go
func UploadToS3(ctx context.Context, client Tigris, bucket string, path string, data []byte) error
```

UploadToS3 uploads a file to S3.

<a name="Validate"></a>
## func [Validate](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/validation.go#L20-L23>)

```go
func Validate(path string, emb *Doc) error
```

Validate validate the given embedding.

<a name="Cache"></a>
## type [Cache](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L64-L70>)

Cache represents an asset cache.

```go
type Cache struct {
    bun.BaseModel `bun:"caches"`

    ID   int64  `bun:"id,pk,autoincrement"`
    Path string `bun:"path,unique"`
    Hash string `bun:"hashed,unique"`
}
```

<a name="CustomTime"></a>
## type [CustomTime](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L17>)

CustomTime allows us to customize the YAML time parsing.

```go
type CustomTime struct{ time.Time }
```

<a name="CustomTime.UnmarshalYAML"></a>
### func \(\*CustomTime\) [UnmarshalYAML](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L20>)

```go
func (ct *CustomTime) UnmarshalYAML(value *yaml.Node) error
```

UnmarshalYAML implements the yaml.Unmarshaler interface.

<a name="DefaultTigrisClient"></a>
## type [DefaultTigrisClient](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/s3.go#L31>)

DefaultTigrisClient is a wrapper for the AWS S3 client.

```go
type DefaultTigrisClient struct{ *s3.Client }
```

<a name="NewTigris"></a>
### func [NewTigris](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/s3.go#L34>)

```go
func NewTigris(getenv func(string) string) (*DefaultTigrisClient, error)
```

NewTigris creates a new Tigris client.

<a name="DirMatchItem"></a>
## type [DirMatchItem](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/hash.go#L43-L47>)

DirMatchItem contains the path and content of a file.

```go
type DirMatchItem struct {
    Path    string
    Content string
    Hash    string
}
```

<a name="HashDirMatch"></a>
### func [HashDirMatch](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/hash.go#L53-L58>)

```go
func HashDirMatch(ctx context.Context, fs afero.Fs, path string, db *bun.DB) ([]DirMatchItem, error)
```

HashDirMatch takes an fs, path, and a db.

It returns a slice of paths if the hash of the directory does not match the hash in the database.

<a name="MatchItem"></a>
### func [MatchItem](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/hash.go#L29>)

```go
func MatchItem(fs afero.Fs, path string) (DirMatchItem, error)
```

MatchItem takes a path and returns a DirMatchItem.

<a name="Doc"></a>
## type [Doc](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L42-L62>)

Doc is a base struct for all embeddedable structs.

```go
type Doc struct {
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
```

<a name="ParseMarkdown"></a>
### func [ParseMarkdown](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/markdown.go#L120>)

```go
func ParseMarkdown(md goldmark.Markdown, item DirMatchItem) (*Doc, error)
```

ParseMarkdown parses a markdown document.

<a name="Doc.GetTitle"></a>
### func \(\*Doc\) [GetTitle](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L275>)

```go
func (emb *Doc) GetTitle() string
```

GetTitle returns the title of the embedding.

<a name="Employment"></a>
## type [Employment](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L148-L171>)

Employment is an employment with all its posts, projects, and tags.

```go
type Employment struct {
    bun.BaseModel `bun:"employments"`

    ID  int64 `bun:"id,pk,autoincrement"`

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
```

<a name="Employment.PagePath"></a>
### func \(\*Employment\) [PagePath](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L295>)

```go
func (emb *Employment) PagePath() string
```

PagePath returns the path to the employment page.

<a name="Employment.String"></a>
### func \(\*Employment\) [String](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L311>)

```go
func (emb *Employment) String() string
```



<a name="EmploymentToEmployment"></a>
## type [EmploymentToEmployment](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L264-L271>)

EmploymentToEmployment represents a many\-to\-many relationship between employments and other employments.

```go
type EmploymentToEmployment struct {
    bun.BaseModel `bun:"employment_to_employments"`

    SourceEmploymentID int64       `bun:"source_employment_id,pk"`
    SourceEmployment   *Employment `bun:"rel:belongs-to,join:source_employment_id=id"`
    TargetEmploymentID int64       `bun:"target_employment_id,pk"`
    TargetEmployment   *Employment `bun:"rel:belongs-to,join:target_employment_id=id"`
}
```

<a name="EmploymentToPost"></a>
## type [EmploymentToPost](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L244-L251>)

EmploymentToPost represents a many\-to\-many relationship between employments and posts.

```go
type EmploymentToPost struct {
    bun.BaseModel `bun:"employment_to_posts"`

    EmploymentID int64       `bun:"employment_id,pk"`
    Employment   *Employment `bun:"rel:belongs-to,join:employment_id=id"`
    PostID       int64       `bun:"post_id,pk"`
    Post         *Post       `bun:"rel:belongs-to,join:post_id=id"`
}
```

<a name="EmploymentToProject"></a>
## type [EmploymentToProject](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L254-L261>)

EmploymentToProject represents a many\-to\-many relationship between employments and projects.

```go
type EmploymentToProject struct {
    bun.BaseModel `bun:"employment_to_projects"`

    EmploymentID int64       `bun:"employment_id,pk"`
    Employment   *Employment `bun:"rel:belongs-to,join:employment_id=id"`
    ProjectID    int64       `bun:"project_id,pk"`
    Project      *Project    `bun:"rel:belongs-to,join:project_id=id"`
}
```

<a name="EmploymentToTag"></a>
## type [EmploymentToTag](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L234-L241>)

EmploymentToTag represents a many\-to\-many relationship between employments and tags.

```go
type EmploymentToTag struct {
    bun.BaseModel `bun:"employment_to_tags"`

    EmploymentID int64       `bun:"employment_id,pk"`
    Employment   *Employment `bun:"rel:belongs-to,join:employment_id=id"`
    TagID        int64       `bun:"tag_id,pk"`
    Tag          *Tag        `bun:"rel:belongs-to,join:tag_id=id"`
}
```

<a name="Post"></a>
## type [Post](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L72-L94>)

Post is a post with all its projects and tags.

```go
type Post struct {
    bun.BaseModel `bun:"posts"`

    ID  int64 `bun:"id,pk,autoincrement" `

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
```

<a name="Post.PagePath"></a>
### func \(\*Post\) [PagePath](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L280>)

```go
func (emb *Post) PagePath() string
```

PagePath returns the path to the post page.

<a name="Post.String"></a>
### func \(\*Post\) [String](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L299>)

```go
func (emb *Post) String() string
```



<a name="PostToPost"></a>
## type [PostToPost](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L204-L211>)

PostToPost represents a many\-to\-many relationship between posts and other posts.

```go
type PostToPost struct {
    bun.BaseModel `bun:"post_to_posts"`

    SourcePostID int64 `bun:"source_post_id,pk"`
    SourcePost   *Post `bun:"rel:belongs-to,join:source_post_id=id"`
    TargetPostID int64 `bun:"target_post_id,pk"`
    TargetPost   *Post `bun:"rel:belongs-to,join:target_post_id=id"`
}
```

<a name="PostToProject"></a>
## type [PostToProject](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L184-L191>)

PostToProject represents a many\-to\-many relationship between posts and projects.

```go
type PostToProject struct {
    bun.BaseModel `bun:"post_to_projects"`

    PostID    int64    `bun:"post_id,pk"`
    Post      *Post    `bun:"rel:belongs-to,join:post_id=id"`
    ProjectID int64    `bun:"project_id,pk"`
    Project   *Project `bun:"rel:belongs-to,join:project_id=id"`
}
```

<a name="PostToTag"></a>
## type [PostToTag](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L174-L181>)

PostToTag represents a many\-to\-many relationship between posts and tags.

```go
type PostToTag struct {
    bun.BaseModel `bun:"post_to_tags"`

    PostID int64 `bun:"post_id,pk"`
    Post   *Post `bun:"rel:belongs-to,join:post_id=id"`
    TagID  int64 `bun:"tag_id,pk"`
    Tag    *Tag  `bun:"rel:belongs-to,join:tag_id=id"`
}
```

<a name="Project"></a>
## type [Project](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L97-L119>)

Project is a project with all its posts and tags.

```go
type Project struct {
    bun.BaseModel `bun:"projects"`

    ID  int64 `bun:"id,pk,autoincrement" yaml:"-"`

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
```

<a name="Project.PagePath"></a>
### func \(\*Project\) [PagePath](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L285>)

```go
func (emb *Project) PagePath() string
```

PagePath returns the path to the project page.

<a name="Project.String"></a>
### func \(\*Project\) [String](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L303>)

```go
func (emb *Project) String() string
```



<a name="ProjectToProject"></a>
## type [ProjectToProject](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L214-L221>)

ProjectToProject represents a many\-to\-many relationship between projects and other projects.

```go
type ProjectToProject struct {
    bun.BaseModel `bun:"project_to_projects"`

    SourceProjectID int64    `bun:"source_project_id,pk"`
    SourceProject   *Project `bun:"rel:belongs-to,join:source_project_id=id"`
    TargetProjectID int64    `bun:"target_project_id,pk"`
    TargetProject   *Project `bun:"rel:belongs-to,join:target_project_id=id"`
}
```

<a name="ProjectToTag"></a>
## type [ProjectToTag](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L194-L201>)

ProjectToTag represents a many\-to\-many relationship between projects and tags.

```go
type ProjectToTag struct {
    bun.BaseModel `bun:"project_to_tags"`

    ProjectID int64    `bun:"project_id,pk"`
    Project   *Project `bun:"rel:belongs-to,join:project_id=id"`
    TagID     int64    `bun:"tag_id,pk"`
    Tag       *Tag     `bun:"rel:belongs-to,join:tag_id=id"`
}
```

<a name="RelationshipFn"></a>
## type [RelationshipFn](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/upsert.go#L15>)

RelationshipFn is a function that updates relationships.

```go
type RelationshipFn func(context.Context) error
```

<a name="UpsertEmployment"></a>
### func [UpsertEmployment](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/upsert.go#L554-L558>)

```go
func UpsertEmployment(ctx context.Context, db *bun.DB, employment *Employment) (RelationshipFn, error)
```

UpsertEmployment saves an employment to the database \(to be called from the DB worker\).

<a name="UpsertEmploymentRelationships"></a>
### func [UpsertEmploymentRelationships](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/upsert.go#L591-L594>)

```go
func UpsertEmploymentRelationships(db *bun.DB, employment *Employment) RelationshipFn
```

UpsertEmploymentRelationships updates relationships for an employment \(to be called from the DB worker\).

<a name="UpsertPost"></a>
### func [UpsertPost](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/upsert.go#L18-L22>)

```go
func UpsertPost(ctx context.Context, db *bun.DB, post *Post) (RelationshipFn, error)
```

UpsertPost saves a post to the database \(to be called from the DB worker\).

<a name="UpsertPostRelationships"></a>
### func [UpsertPostRelationships](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/upsert.go#L119-L122>)

```go
func UpsertPostRelationships(db *bun.DB, post *Post) RelationshipFn
```

UpsertPostRelationships updates relationships for a post \(to be called from the DB worker\).

<a name="UpsertProject"></a>
### func [UpsertProject](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/upsert.go#L55-L59>)

```go
func UpsertProject(ctx context.Context, db *bun.DB, project *Project) (RelationshipFn, error)
```

UpsertProject saves a project to the database \(to be called from the DB worker\).

<a name="UpsertProjectRelationships"></a>
### func [UpsertProjectRelationships](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/upsert.go#L264-L267>)

```go
func UpsertProjectRelationships(db *bun.DB, project *Project) RelationshipFn
```

UpsertProjectRelationships updates relationships for a project.

<a name="UpsertTag"></a>
### func [UpsertTag](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/upsert.go#L87-L91>)

```go
func UpsertTag(ctx context.Context, db *bun.DB, tag *Tag) (RelationshipFn, error)
```

UpsertTag saves a tag to the database \(to be called from the DB worker\).

<a name="UpsertTagRelationships"></a>
### func [UpsertTagRelationships](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/upsert.go#L409-L412>)

```go
func UpsertTagRelationships(db *bun.DB, tag *Tag) RelationshipFn
```

UpsertTagRelationships updates relationships for a tag .

<a name="Tag"></a>
## type [Tag](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L122-L145>)

Tag is a tag with all its posts and projects.

```go
type Tag struct {
    bun.BaseModel `bun:"tags"`

    ID  int64 `bun:"id,pk,autoincrement"`

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
```

<a name="Tag.PagePath"></a>
### func \(\*Tag\) [PagePath](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L290>)

```go
func (emb *Tag) PagePath() string
```

PagePath returns the path to the tag page.

<a name="Tag.String"></a>
### func \(\*Tag\) [String](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L307>)

```go
func (emb *Tag) String() string
```



<a name="TagToTag"></a>
## type [TagToTag](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/static.go#L224-L231>)

TagToTag represents a many\-to\-many relationship between tags and other tags.

```go
type TagToTag struct {
    bun.BaseModel `bun:"tag_to_tags"`

    SourceTagID int64 `bun:"source_tag_id,pk"`
    SourceTag   *Tag  `bun:"rel:belongs-to,join:source_tag_id=id"`
    TargetTagID int64 `bun:"target_tag_id,pk"`
    TargetTag   *Tag  `bun:"rel:belongs-to,join:target_tag_id=id"`
}
```

<a name="Tigris"></a>
## type [Tigris](<https://github.com/conneroisu/conneroh.com/blob/main/internal/assets/s3.go#L22-L28>)

Tigris is an minimal interface for AWS clients.

```go
type Tigris interface {
    PutObject(
        ctx context.Context,
        params *s3.PutObjectInput,
        optFns ...func(*s3.Options),
    ) (*s3.PutObjectOutput, error)
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->