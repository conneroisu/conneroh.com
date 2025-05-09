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

- [Constants](<#constants>)
- [func ComputeAllURLs\(base string\) \[\]string](<#ComputeAllURLs>)
- [func GeneratePagination\(currentPage, totalPages, maxDisplay int\) \[\]string](<#GeneratePagination>)
- [func GetPostURL\(base string, post \*assets.Post\) string](<#GetPostURL>)
- [func GetProjectURL\(base string, project \*assets.Project\) string](<#GetProjectURL>)
- [func GetTagURL\(base string, tag \*assets.Tag\) string](<#GetTagURL>)
- [func MorphableHandler\(full templ.Component, morph templ.Component\) http.HandlerFunc](<#MorphableHandler>)
- [func Paginate\[T any\]\(items \[\]T, page int, pageSize int\) \(\[\]T, int\)](<#Paginate>)
- [type PluralPath](<#PluralPath>)


## Constants

<a name="MaxListLargeItems"></a>

```go
const (
    // MaxListLargeItems is the maximum number of items in a list view.
    MaxListLargeItems = 9
    // MaxListSmallItems is the maximum number of items in a list view.
    MaxListSmallItems = 27

    // MaxMobilePageNumber is the maximum number of pages to display before ... is shown.
    MaxMobilePageNumber = 5
    // MaxDesktopPageNumber is the maximum number of pages to display before ... is shown.
    MaxDesktopPageNumber = 10

    // Ellipsis represents pagination gaps.
    Ellipsis = "..."
)
```

<a name="ComputeAllURLs"></a>
## func [ComputeAllURLs](<https://github.com/conneroisu/conneroh.com/blob/main/internal/routing/urls.go#L10>)

```go
func ComputeAllURLs(base string) []string
```

ComputeAllURLs computes all URLs for all posts, projects, and tags given a base URL.

<a name="GeneratePagination"></a>
## func [GeneratePagination](<https://github.com/conneroisu/conneroh.com/blob/main/internal/routing/pagination.go#L47>)

```go
func GeneratePagination(currentPage, totalPages, maxDisplay int) []string
```

GeneratePagination generates a pagination list of page numbers.

<a name="GetPostURL"></a>
## func [GetPostURL](<https://github.com/conneroisu/conneroh.com/blob/main/internal/routing/targets.go#L23>)

```go
func GetPostURL(base string, post *assets.Post) string
```

GetPostURL returns the URL for a post.

<a name="GetProjectURL"></a>
## func [GetProjectURL](<https://github.com/conneroisu/conneroh.com/blob/main/internal/routing/targets.go#L28>)

```go
func GetProjectURL(base string, project *assets.Project) string
```

GetProjectURL returns the URL for a project.

<a name="GetTagURL"></a>
## func [GetTagURL](<https://github.com/conneroisu/conneroh.com/blob/main/internal/routing/targets.go#L33>)

```go
func GetTagURL(base string, tag *assets.Tag) string
```

GetTagURL returns the URL for a tag.

<a name="MorphableHandler"></a>
## func [MorphableHandler](<https://github.com/conneroisu/conneroh.com/blob/main/internal/routing/handlers.go#L12-L15>)

```go
func MorphableHandler(full templ.Component, morph templ.Component) http.HandlerFunc
```

MorphableHandler returns a handler that checks for the presence of the hx\-trigger header and serves either the full or morphed view.

<a name="Paginate"></a>
## func [Paginate](<https://github.com/conneroisu/conneroh.com/blob/main/internal/routing/pagination.go#L23-L27>)

```go
func Paginate[T any](items []T, page int, pageSize int) ([]T, int)
```

Paginate paginates a list of items.

<a name="PluralPath"></a>
## type [PluralPath](<https://github.com/conneroisu/conneroh.com/blob/main/internal/routing/targets.go#L11>)

PluralPath is the target of a plural view. string

```go
type PluralPath = string
```

<a name="PostPluralPath"></a>

```go
const (
    // PostPluralPath is the target of a plural post view.
    PostPluralPath PluralPath = "posts"
    // ProjectPluralPath is the target of a plural project view.
    ProjectPluralPath PluralPath = "projects"
    // TagsPluralPath is the target of a plural tag view.
    TagsPluralPath PluralPath = "tags"
)
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->
