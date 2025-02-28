# conneroh

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# conneroh

```go
import "github.com/conneroisu/conneroh.com/cmd/conneroh"
```

Package conneroh provides implementations for conneroh.

## Index

- [Variables](<#variables>)
- [func AddRoutes\(ctx context.Context, h \*http.ServeMux, db \*data.Database\[master.Queries\]\) error](<#AddRoutes>)
- [func Dist\(\_ context.Context, \_ \*data.Database\[master.Queries\], \_ \*\[\]master.FullPost, \_ \*\[\]master.FullProject, \_ \*\[\]master.FullTag, \_ \*map\[string\]master.FullPost, \_ \*map\[string\]master.FullProject, \_ \*map\[string\]master.FullTag\) \(routing.APIFn, error\)](<#Dist>)
- [func Favicon\(\_ context.Context, \_ \*data.Database\[master.Queries\], \_ \*\[\]master.FullPost, \_ \*\[\]master.FullProject, \_ \*\[\]master.FullTag, \_ \*map\[string\]master.FullPost, \_ \*map\[string\]master.FullProject, \_ \*map\[string\]master.FullTag\) \(routing.APIFn, error\)](<#Favicon>)
- [func Home\(ctx context.Context, db \*data.Database\[master.Queries\], fullPosts \*\[\]master.FullPost, fullProjects \*\[\]master.FullProject, fullTags \*\[\]master.FullTag, fullPostsSlugMap \*map\[string\]master.FullPost, fullProjectsSlugMap \*map\[string\]master.FullProject, fullTagsSlugMap \*map\[string\]master.FullTag\) \(routing.APIFn, error\)](<#Home>)
- [func Morph\(ctx context.Context, db \*data.Database\[master.Queries\], fullPosts \*\[\]master.FullPost, fullProjects \*\[\]master.FullProject, fullTags \*\[\]master.FullTag, fullPostSlugMap \*map\[string\]master.FullPost, fullProjectSlugMap \*map\[string\]master.FullProject, fullTagSlugMap \*map\[string\]master.FullTag\) \(routing.APIFn, error\)](<#Morph>)
- [func Morphs\(\_ context.Context, \_ \*data.Database\[master.Queries\], fullPosts \*\[\]master.FullPost, fullProjects \*\[\]master.FullProject, fullTags \*\[\]master.FullTag, fullPostSlugMap \*map\[string\]master.FullPost, fullProjectSlugMap \*map\[string\]master.FullProject, fullTagSlugMap \*map\[string\]master.FullTag\) \(routing.APIFn, error\)](<#Morphs>)
- [func NewServer\(ctx context.Context, db \*data.Database\[master.Queries\]\) http.Handler](<#NewServer>)
- [func Post\(\_ context.Context, \_ \*data.Database\[master.Queries\], fullPosts \*\[\]master.FullPost, fullProjects \*\[\]master.FullProject, fullTags \*\[\]master.FullTag, fullPostSlugMap \*map\[string\]master.FullPost, fullProjectSlugMap \*map\[string\]master.FullProject, fullTagSlugMap \*map\[string\]master.FullTag\) \(routing.APIFn, error\)](<#Post>)
- [func Posts\(\_ context.Context, \_ \*data.Database\[master.Queries\], fullPosts \*\[\]master.FullPost, fullProjects \*\[\]master.FullProject, fullTags \*\[\]master.FullTag, fullPostSlugMap \*map\[string\]master.FullPost, fullProjectSlugMap \*map\[string\]master.FullProject, fullTagSlugMap \*map\[string\]master.FullTag\) \(routing.APIFn, error\)](<#Posts>)
- [func Project\(ctx context.Context, db \*data.Database\[master.Queries\], fullPosts \*\[\]master.FullPost, fullProjects \*\[\]master.FullProject, fullTags \*\[\]master.FullTag, fullPostSlugMap \*map\[string\]master.FullPost, fullProjectSlugMap \*map\[string\]master.FullProject, fullTagSlugMap \*map\[string\]master.FullTag\) \(routing.APIFn, error\)](<#Project>)
- [func Projects\(ctx context.Context, db \*data.Database\[master.Queries\], fullPosts \*\[\]master.FullPost, fullProjects \*\[\]master.FullProject, fullTags \*\[\]master.FullTag, fullPostSlugMap \*map\[string\]master.FullPost, fullProjectSlugMap \*map\[string\]master.FullProject, fullTagSlugMap \*map\[string\]master.FullTag\) \(routing.APIFn, error\)](<#Projects>)
- [func Run\(ctx context.Context, \_ func\(string\) string\) error](<#Run>)
- [func Tag\(\_ context.Context, \_ \*data.Database\[master.Queries\], fullPosts \*\[\]master.FullPost, fullProjects \*\[\]master.FullProject, fullTags \*\[\]master.FullTag, fullPostSlugMap \*map\[string\]master.FullPost, fullProjectSlugMap \*map\[string\]master.FullProject, fullTagSlugMap \*map\[string\]master.FullTag\) \(routing.APIFn, error\)](<#Tag>)
- [func Tags\(\_ context.Context, \_ \*data.Database\[master.Queries\], fullPosts \*\[\]master.FullPost, fullProjects \*\[\]master.FullProject, fullTags \*\[\]master.FullTag, fullPostSlugMap \*map\[string\]master.FullPost, fullProjectSlugMap \*map\[string\]master.FullProject, fullTagSlugMap \*map\[string\]master.FullTag\) \(routing.APIFn, error\)](<#Tags>)


## Variables

<a name="RouteMap"></a>RouteMap is a map of all routes.

```go
var RouteMap = routing.APIMap{
    "GET /dist/":                      Dist,
    "GET /favicon.ico":                Favicon,
    "GET /{$}":                        Home,
    "GET /projects":                   Projects,
    "GET /posts":                      Posts,
    "GET /tags":                       Tags,
    "GET /project/{id}":               Project,
    "GET /post/{id}":                  Post,
    "GET /tag/{id}":                   Tag,
    "GET /hateoas/morph/{view}":       Morph,
    "GET /hateoas/morphs/{view}/{id}": Morphs,
}
```

<a name="AddRoutes"></a>
## func [AddRoutes](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/routes.go#L28-L32>)

```go
func AddRoutes(ctx context.Context, h *http.ServeMux, db *data.Database[master.Queries]) error
```

AddRoutes adds all routes to the router.

<a name="Dist"></a>
## func [Dist](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/handlers.go#L36-L45>)

```go
func Dist(_ context.Context, _ *data.Database[master.Queries], _ *[]master.FullPost, _ *[]master.FullProject, _ *[]master.FullTag, _ *map[string]master.FullPost, _ *map[string]master.FullProject, _ *map[string]master.FullTag) (routing.APIFn, error)
```

Dist is the dist handler for serving/distributing static files.

<a name="Favicon"></a>
## func [Favicon](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/handlers.go#L53-L62>)

```go
func Favicon(_ context.Context, _ *data.Database[master.Queries], _ *[]master.FullPost, _ *[]master.FullProject, _ *[]master.FullTag, _ *map[string]master.FullPost, _ *map[string]master.FullProject, _ *map[string]master.FullTag) (routing.APIFn, error)
```

Favicon is the favicon handler.

<a name="Home"></a>
## func [Home](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/handlers.go#L74-L83>)

```go
func Home(ctx context.Context, db *data.Database[master.Queries], fullPosts *[]master.FullPost, fullProjects *[]master.FullProject, fullTags *[]master.FullTag, fullPostsSlugMap *map[string]master.FullPost, fullProjectsSlugMap *map[string]master.FullProject, fullTagsSlugMap *map[string]master.FullTag) (routing.APIFn, error)
```

Home is the home page handler.

<a name="Morph"></a>
## func [Morph](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/handlers.go#L99-L108>)

```go
func Morph(ctx context.Context, db *data.Database[master.Queries], fullPosts *[]master.FullPost, fullProjects *[]master.FullProject, fullTags *[]master.FullTag, fullPostSlugMap *map[string]master.FullPost, fullProjectSlugMap *map[string]master.FullProject, fullTagSlugMap *map[string]master.FullTag) (routing.APIFn, error)
```

Morph renders a morphed view.

<a name="Morphs"></a>
## func [Morphs](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/handlers.go#L138-L147>)

```go
func Morphs(_ context.Context, _ *data.Database[master.Queries], fullPosts *[]master.FullPost, fullProjects *[]master.FullProject, fullTags *[]master.FullTag, fullPostSlugMap *map[string]master.FullPost, fullProjectSlugMap *map[string]master.FullProject, fullTagSlugMap *map[string]master.FullTag) (routing.APIFn, error)
```

Morphs renders a morphed view.

<a name="NewServer"></a>
## func [NewServer](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/root.go#L31-L34>)

```go
func NewServer(ctx context.Context, db *data.Database[master.Queries]) http.Handler
```

NewServer creates a new web\-ui server

<a name="Post"></a>
## func [Post](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/handlers.go#L289-L298>)

```go
func Post(_ context.Context, _ *data.Database[master.Queries], fullPosts *[]master.FullPost, fullProjects *[]master.FullProject, fullTags *[]master.FullTag, fullPostSlugMap *map[string]master.FullPost, fullProjectSlugMap *map[string]master.FullProject, fullTagSlugMap *map[string]master.FullTag) (routing.APIFn, error)
```

Post is the post handler.

<a name="Posts"></a>
## func [Posts](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/handlers.go#L253-L262>)

```go
func Posts(_ context.Context, _ *data.Database[master.Queries], fullPosts *[]master.FullPost, fullProjects *[]master.FullProject, fullTags *[]master.FullTag, fullPostSlugMap *map[string]master.FullPost, fullProjectSlugMap *map[string]master.FullProject, fullTagSlugMap *map[string]master.FullTag) (routing.APIFn, error)
```

Posts is the posts handler.

<a name="Project"></a>
## func [Project](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/handlers.go#L424-L433>)

```go
func Project(ctx context.Context, db *data.Database[master.Queries], fullPosts *[]master.FullPost, fullProjects *[]master.FullProject, fullTags *[]master.FullTag, fullPostSlugMap *map[string]master.FullPost, fullProjectSlugMap *map[string]master.FullProject, fullTagSlugMap *map[string]master.FullTag) (routing.APIFn, error)
```

Project is the project handler.

<a name="Projects"></a>
## func [Projects](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/handlers.go#L398-L407>)

```go
func Projects(ctx context.Context, db *data.Database[master.Queries], fullPosts *[]master.FullPost, fullProjects *[]master.FullProject, fullTags *[]master.FullTag, fullPostSlugMap *map[string]master.FullPost, fullProjectSlugMap *map[string]master.FullProject, fullTagSlugMap *map[string]master.FullTag) (routing.APIFn, error)
```

Projects is the projects handler.

<a name="Run"></a>
## func [Run](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/root.go#L54-L57>)

```go
func Run(ctx context.Context, _ func(string) string) error
```

Run is the entry point for the application.

<a name="Tag"></a>
## func [Tag](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/handlers.go#L485-L494>)

```go
func Tag(_ context.Context, _ *data.Database[master.Queries], fullPosts *[]master.FullPost, fullProjects *[]master.FullProject, fullTags *[]master.FullTag, fullPostSlugMap *map[string]master.FullPost, fullProjectSlugMap *map[string]master.FullProject, fullTagSlugMap *map[string]master.FullTag) (routing.APIFn, error)
```

Tag is the tag handler.

<a name="Tags"></a>
## func [Tags](<https://github.com/conneroisu/conneroh/blob/main/cmd/conneroh/handlers.go#L458-L467>)

```go
func Tags(_ context.Context, _ *data.Database[master.Queries], fullPosts *[]master.FullPost, fullProjects *[]master.FullProject, fullTags *[]master.FullTag, fullPostSlugMap *map[string]master.FullPost, fullProjectSlugMap *map[string]master.FullProject, fullTagSlugMap *map[string]master.FullTag) (routing.APIFn, error)
```

Tags is the tags handler.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->
