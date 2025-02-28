# Routing

Package routing provides implementations for routing.

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# routing

```go
import "github.com/conneroisu/conneroh.com/internal/routing"
```

Package routing provides implementations for routing.

## Index

- [type APIFn](<#APIFn>)
- [type APIHandler](<#APIHandler>)
- [type APIMap](<#APIMap>)
  - [func \(m APIMap\) AddRoutes\(ctx context.Context, mux \*http.ServeMux, db \*data.Database\[master.Queries\], fullPosts \*\[\]master.FullPost, fullProjects \*\[\]master.FullProject, fullTags \*\[\]master.FullTag, postsSlugMap \*map\[string\]master.FullPost, projectsSlugMap \*map\[string\]master.FullProject, tagsSlugMap \*map\[string\]master.FullTag\) error](<#APIMap.AddRoutes>)
- [type ErrMissingParam](<#ErrMissingParam>)
  - [func \(e ErrMissingParam\) Error\(\) string](<#ErrMissingParam.Error>)
- [type ErrNotFound](<#ErrNotFound>)
  - [func \(e ErrNotFound\) Error\(\) string](<#ErrNotFound.Error>)


<a name="APIFn"></a>
## type [APIFn](<https://github.com/conneroisu/conneroh/blob/main/internal/routing/main.go#L13>)

APIFn is a function that handles an API request.

```go
type APIFn func(http.ResponseWriter, *http.Request) error
```

<a name="APIHandler"></a>
## type [APIHandler](<https://github.com/conneroisu/conneroh/blob/main/internal/routing/main.go#L16-L25>)

APIHandler is a function that returns an APIFn.

```go
type APIHandler func(
    ctx context.Context,
    db *data.Database[master.Queries],
    fullPosts *[]master.FullPost,
    fullProjects *[]master.FullProject,
    fullTags *[]master.FullTag,
    postsSlugMap *map[string]master.FullPost,
    projectsSlugMap *map[string]master.FullProject,
    tagsSlugMap *map[string]master.FullTag,
) (APIFn, error)
```

<a name="APIMap"></a>
## type [APIMap](<https://github.com/conneroisu/conneroh/blob/main/internal/routing/main.go#L28>)

APIMap is a map of API functions.

```go
type APIMap map[string]APIHandler
```

<a name="APIMap.AddRoutes"></a>
### func \(APIMap\) [AddRoutes](<https://github.com/conneroisu/conneroh/blob/main/internal/routing/main.go#L31-L41>)

```go
func (m APIMap) AddRoutes(ctx context.Context, mux *http.ServeMux, db *data.Database[master.Queries], fullPosts *[]master.FullPost, fullProjects *[]master.FullProject, fullTags *[]master.FullTag, postsSlugMap *map[string]master.FullPost, projectsSlugMap *map[string]master.FullProject, tagsSlugMap *map[string]master.FullTag) error
```

AddRoutes adds all routes to the router.

<a name="ErrMissingParam"></a>
## type [ErrMissingParam](<https://github.com/conneroisu/conneroh/blob/main/internal/routing/errors.go#L23-L26>)

ErrMissingParam is an error that is returned when a required parameter is missing.

```go
type ErrMissingParam struct {
    ID   string
    View string
}
```

<a name="ErrMissingParam.Error"></a>
### func \(ErrMissingParam\) [Error](<https://github.com/conneroisu/conneroh/blob/main/internal/routing/errors.go#L29>)

```go
func (e ErrMissingParam) Error() string
```

Error implements the error interface on ErrMissingParam.

<a name="ErrNotFound"></a>
## type [ErrNotFound](<https://github.com/conneroisu/conneroh/blob/main/internal/routing/errors.go#L9-L11>)

ErrNotFound is an error that is returned when a resource is not found.

```go
type ErrNotFound struct {
    URL *url.URL
}
```

<a name="ErrNotFound.Error"></a>
### func \(ErrNotFound\) [Error](<https://github.com/conneroisu/conneroh/blob/main/internal/routing/errors.go#L14>)

```go
func (e ErrNotFound) Error() string
```

Error implements the error interface on ErrNotFound.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->
