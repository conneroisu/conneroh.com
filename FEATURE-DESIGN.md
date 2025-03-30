# Design of New Search Feature

## Actual Search

- plus and minus tags/projects/posts
- search by tags/projects/posts
- raw text search using embeddings

### Search by tags
- `-tag:programming-language/go` will not show posts with `programming-language/go` tag
- `tag:programming-language/go` will show posts with `programming-language/go` tag

### Search by projects
- `-project:github` will not show posts with `github` project
- `project:github` will show posts with `github` project  

### Search by posts
- `-post:hello-world` will not show posts with `hello-world` post
- `post:hello-world` will show posts with `hello-world` post

HTMX endpoint for list pages: `/search/{tag|project|post}`
HTMX endpoint for global search: `/search/all`

## Auto-complete
- autocomplete tags/projects/posts

### Autocomplete tags
- `-tag:` should show all tags in an autocomplete list
- `tag:` should show all tags in an autocomplete list

### Autocomplete projects
- `-project:` should show all projects in an autocomplete list
- `project:` should show all projects in an autocomplete list

### Autocomplete posts
- `-post:` should show all posts in an autocomplete list
- `post:` should show all posts in an autocomplete list

HTMX endpoint for list pages: `/autocomplete/{tag|project|post}`
HTMX endpoint for global autocomplete: `/autocomplete/all`
