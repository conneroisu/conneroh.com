# Design of New Search Feature

## Actual Search

- plus and minus tags/projects/posts
- search by tags/projects/posts
- raw text search using embeddings
- enable recommendations toggle?

### Search by tags

- `-tag:programming-language/go` will not show posts with `programming-language/go` tag
- `tag:programming-language/go` will show posts with `programming-language/go` tag

### Search by projects

- `-project:conneroh.com` will not show posts related to the `conneroh.com` project
- `project:conneroh.com` will show posts with `github` project

### Search by posts

- `-post:hello-world` will not show results related to the `hello-world` post
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

# Enhanced Search Feature Design

## Architecture Overview

The search system will operate in two modes:

1. **List filtering** - Filtering items on list pages (/posts, /projects, /tags)
2. **Global search** - Searching across all content types from any page

Both modes will support advanced query syntax with prefix modifiers and boolean operations.

## Core Search Functionality

### Query Parser

We'll implement a robust query parser that supports:

```go
type SearchQuery struct {
    RawQuery       string            // Original query string
    TextTerms      []string          // Plain text search terms
    TagInclusions  []string          // Tags to include
    TagExclusions  []string          // Tags to exclude
    ProjectInclusions []string       // Projects to include
    ProjectExclusions []string       // Projects to exclude
    PostInclusions []string          // Posts to include
    PostExclusions []string          // Posts to exclude
}
```

The parser will handle:

- Plain text terms for content search
- Prefix modifiers (`tag:`, `-tag:`, `project:`, etc.)
- Quoted phrases (`"exact phrase"`)
- Boolean operations through syntax

### Search Backend Options

#### Option 1: In-memory search (simpler)

- Use Go's built-in string operations for filtering
- Pre-compute term indexes for faster text matching
- Good for moderate content volumes

#### Option 2: Embeddings-based search (more powerful)

- Generate vector embeddings for all content
- Use cosine similarity for semantic search
- Integrate with external embedding API (e.g., OpenAI)
- Better for natural language queries

## User Interface Components

### List Page Search UI

For each list page, the search component will:

1. Allow typing search queries
2. Show autocomplete suggestions
3. Filter items in real-time
4. Display active filters visually
5. Support clearing individual filters

UI mockup for list page search:

```
┌───────────────────────────────────┐
│ 🔍 Search posts...                │
└───────────────────────────────────┘

  Active filters:
  [tag:golang ⨯] [project:website ⨯] [-tag:legacy ⨯]

  Results: 12 posts matching your search
```

### Global Search UI

The global search will:

1. Be accessible from any page
2. Show categorized results (posts, projects, tags)
3. Support keyboard navigation
4. Provide rich result previews
5. Link directly to content

## Implementation Details

### API Endpoints

#### Search Endpoints

- `POST /api/search/posts` - Search posts with query/filters
- `POST /api/search/projects` - Search projects with query/filters
- `POST /api/search/tags` - Search tags with query/filters
- `POST /api/search/all` - Global search across all content

#### Autocomplete Endpoints

- `GET /api/autocomplete/tags?q=go` - Autocomplete tag names
- `GET /api/autocomplete/projects?q=web` - Autocomplete project names
- `GET /api/autocomplete/posts?q=how` - Autocomplete post titles
- `GET /api/autocomplete/all?q=go` - Autocomplete across all types

### Frontend Integration with HTMX

The search UI will use HTMX for seamless updates:

```html
<input
  type="text"
  hx-post="/api/search/posts"
  hx-trigger="keyup changed delay:500ms, search"
  hx-target="#search-results"
  hx-indicator="#search-indicator"
  placeholder="Search posts..."
/>

<div id="search-indicator" class="htmx-indicator">
  <div class="spinner"></div>
</div>

<div id="search-results">
  <!-- Results will be inserted here -->
</div>
```

### Autocomplete Dropdown

The autocomplete will show different types of suggestions:

```html
<div class="search-suggestions">
  <div class="suggestion-group">
    <h4>Tags</h4>
    <div
      class="suggestion-item"
      hx-post="/api/search/posts"
      hx-vals='{"query": "tag:golang"}'
    >
      <span class="tag-icon">🏷️</span> golang
    </div>
  </div>
  <div class="suggestion-group">
    <h4>Projects</h4>
    <div
      class="suggestion-item"
      hx-post="/api/search/posts"
      hx-vals='{"query": "project:website"}'
    >
      <span class="project-icon">📁</span> website
    </div>
  </div>
  <!-- More suggestion groups -->
</div>
```

## Backend Implementation

### Query Parsing

The query parser will break down complex queries:

```go
func ParseSearchQuery(query string) SearchQuery {
    result := SearchQuery{
        RawQuery: query,
    }

    // Split by spaces, respecting quotes
    terms := splitQueryTerms(query)

    for _, term := range terms {
        if strings.HasPrefix(term, "tag:") {
            result.TagInclusions = append(result.TagInclusions, strings.TrimPrefix(term, "tag:"))
        } else if strings.HasPrefix(term, "-tag:") {
            result.TagExclusions = append(result.TagExclusions, strings.TrimPrefix(term, "-tag:"))
        } else if strings.HasPrefix(term, "project:") {
            result.ProjectInclusions = append(result.ProjectInclusions, strings.TrimPrefix(term, "project:"))
        } else if strings.HasPrefix(term, "-project:") {
            result.ProjectExclusions = append(result.ProjectExclusions, strings.TrimPrefix(term, "-project:"))
        } else if strings.HasPrefix(term, "post:") {
            result.PostInclusions = append(result.PostInclusions, strings.TrimPrefix(term, "post:"))
        } else if strings.HasPrefix(term, "-post:") {
            result.PostExclusions = append(result.PostExclusions, strings.TrimPrefix(term, "-post:"))
        } else {
            result.TextTerms = append(result.TextTerms, term)
        }
    }

    return result
}
```

### Search Handler

The main search handler will process queries and return matching items:

```go
func SearchPostsHandler(w http.ResponseWriter, r *http.Request) {
    query := r.FormValue("query")
    parsedQuery := ParseSearchQuery(query)

    // Apply filters to all posts
    filteredPosts := FilterPosts(gen.AllPosts, parsedQuery)

    // Render the search results template
    views.PostSearchResults(filteredPosts).Render(r.Context(), w)
}
```

## Advanced Features

### Filtering Implementation

The filtering logic will handle complex combinations:

```go
func FilterPosts(posts []*gen.Post, query SearchQuery) []*gen.Post {
    var filtered []*gen.Post

    for _, post := range posts {
        // Check tag inclusions
        if len(query.TagInclusions) > 0 {
            hasRequiredTag := false
            for _, requiredTag := range query.TagInclusions {
                if PostHasTagWithSlug(post, requiredTag) {
                    hasRequiredTag = true
                    break
                }
            }
            if !hasRequiredTag {
                continue // Skip this post
            }
        }

        // Check tag exclusions
        shouldExclude := false
        for _, excludedTag := range query.TagExclusions {
            if PostHasTagWithSlug(post, excludedTag) {
                shouldExclude = true
                break
            }
        }
        if shouldExclude {
            continue // Skip this post
        }

        // Similar checks for projects and posts
        // ...

        // Text search on title and content
        if len(query.TextTerms) > 0 {
            matchesText := false
            for _, term := range query.TextTerms {
                if strings.Contains(strings.ToLower(post.Title), strings.ToLower(term)) ||
                   strings.Contains(strings.ToLower(post.Content), strings.ToLower(term)) {
                    matchesText = true
                    break
                }
            }
            if !matchesText {
                continue // Skip this post
            }
        }

        // If we've made it here, the post matches all criteria
        filtered = append(filtered, post)
    }

    return filtered
}
```

### Vector Search Implementation (Optional)

For semantic search capabilities:

```go
type VectorSearchIndex struct {
    PostVectors    map[string][]float32  // Post slug -> embedding vector
    ProjectVectors map[string][]float32  // Project slug -> embedding vector
    TagVectors     map[string][]float32  // Tag slug -> embedding vector
}

func (idx *VectorSearchIndex) Search(query string, limit int) ([]SearchResult, error) {
    // 1. Get embedding for the query
    queryVector, err := getEmbedding(query)
    if err != nil {
        return nil, err
    }

    // 2. Calculate cosine similarity with all content
    var results []SearchResult

    // Posts
    for slug, vector := range idx.PostVectors {
        similarity := cosineSimilarity(queryVector, vector)
        if similarity > 0.7 { // Threshold
            results = append(results, SearchResult{
                Type:       "post",
                Slug:       slug,
                Score:      similarity,
            })
        }
    }

    // Similar for projects and tags
    // ...

    // 3. Sort by score
    sort.Slice(results, func(i, j int) bool {
        return results[i].Score > results[j].Score
    })

    // 4. Return top results
    if len(results) > limit {
        results = results[:limit]
    }

    return results, nil
}
```

## Templates for Search Results

### Post Search Results Template

```templ
templ PostSearchResults(posts []*gen.Post, query SearchQuery) {
    <div class="search-results">
        <div class="search-metadata">
            <span class="result-count">{ len(posts) } posts found</span>
            if len(query.TagInclusions) > 0 || len(query.TagExclusions) > 0 ||
               len(query.ProjectInclusions) > 0 || len(query.ProjectExclusions) > 0 ||
               len(query.TextTerms) > 0 {
                <div class="active-filters">
                    for _, tag := range query.TagInclusions {
                        <span class="filter-tag">
                            tag:{ tag }
                            <button class="remove-filter" hx-post="/api/search/posts"
                                    hx-vals='{"query": "{ removeTermFromQuery(query.RawQuery, "tag:" + tag) }"}'>×</button>
                        </span>
                    }
                    // Similar for other filter types
                </div>
            }
        </div>

        <div class="posts-grid">
            if len(posts) > 0 {
                for _, post := range posts {
                    @views.listPostItem(post)
                }
            } else {
                <div class="no-results">
                    <h3>No posts found matching your search</h3>
                    <p>Try adjusting your search terms or removing some filters</p>
                </div>
            }
        </div>
    </div>
}
```

## Integration With Existing Code

### Adding Search Routes

Add these routes to your existing `routes.go`:

```go
// Add search routes
h.Handle("POST /search/posts", http.HandlerFunc(SearchPostsHandler))
h.Handle("POST /search/projects", http.HandlerFunc(SearchProjectsHandler))
h.Handle("POST /search/tags", http.HandlerFunc(SearchTagsHandler))
h.Handle("POST /search/all", http.HandlerFunc(SearchAllHandler))

// Add autocomplete routes
h.Handle("GET /autocomplete/tags", http.HandlerFunc(AutocompleteTagsHandler))
h.Handle("GET /autocomplete/projects", http.HandlerFunc(AutocompleteProjectsHandler))
h.Handle("GET /autocomplete/posts", http.HandlerFunc(AutocompletePostsHandler))
h.Handle("GET /autocomplete/all", http.HandlerFunc(AutocompleteAllHandler))
```

## User Experience Improvements

### Keyboard Navigation

Enhance the search UI with keyboard navigation:

```js
document.addEventListener("alpine:init", () => {
  Alpine.data("searchUI", () => ({
    selectedIndex: -1,
    suggestions: [],

    init() {
      // Handle keyboard navigation
      this.$watch("suggestions", () => {
        this.selectedIndex = -1;
      });
    },

    onKeyDown(e) {
      if (e.key === "ArrowDown") {
        e.preventDefault();
        this.selectedIndex = Math.min(
          this.selectedIndex + 1,
          this.suggestions.length - 1,
        );
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        this.selectedIndex = Math.max(this.selectedIndex - 1, -1);
      } else if (e.key === "Enter" && this.selectedIndex >= 0) {
        e.preventDefault();
        this.selectSuggestion(this.suggestions[this.selectedIndex]);
      }
    },

    selectSuggestion(suggestion) {
      // Add the suggestion to the query
      // ...
    },
  }));
});
```

### Search History

Store recent searches in local storage:

```js
function saveSearchToHistory(query) {
  const searches = JSON.parse(localStorage.getItem("recentSearches") || "[]");

  // Don't add duplicates
  if (!searches.includes(query)) {
    searches.unshift(query);

    // Keep only the 10 most recent searches
    if (searches.length > 10) {
      searches.pop();
    }

    localStorage.setItem("recentSearches", JSON.stringify(searches));
  }
}
```

## Progressive Enhancement

The search will work without JavaScript, but will be enhanced where possible.

1. **Basic functionality**: Form submission for non-JS browsers
2. **HTMX enhancement**: Real-time results with HTMX and forced keyboard interaction
3. **Alpine enhancement**: Keyboard navigation and advanced UI features
