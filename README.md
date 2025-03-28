# Conner Ohnesorge's Portfolio

A personal portfolio website built with Go and modern web technologies, showcasing projects, blog posts, and skills.

## Overview

This portfolio website is built using Go with templ for generating HTML templates, Alpine.js for client-side interactivity, and Tailwind CSS for styling.

It follows a modern server-side rendering approach with HTMX for dynamic content loading without full page refreshes.

## Features

- **Server-side rendered** pages with templ templates
- **Content organization** by projects, posts, and tags
- **Dynamic content loading** with HTMX for a smooth user experience
- **Responsive design** with Tailwind CSS
- **Client-side interactions** powered by Alpine.js
- **Full-text search** capabilities
- **Tag-based navigation** to easily find related content
- **Markdown content** for easy authoring and maintenance

## Tech Stack

- **Backend**: Go
- **Templates**: templ
- **CSS**: Tailwind CSS
- **Frontend Interactivity**: Alpine.js
- **Dynamic Content**: HTMX
- **Build/Development**: Nix, air (for live reloading)

## Project Structure

```
├── cmd/                 # Entry points for executables
│   ├── conneroh/        # Main web application
│   │   ├── layouts/     # Layout templates
│   │   ├── views/       # View templates
│   │   ├── components/  # Component templates
│   │   └── _static/     # Static assets
│   └── update/          # Content update utility
├── internal/            # Private application code
│   ├── data/            # Data access layer
│   │   ├── docs/        # Markdown content
│   │   │   ├── posts/   # Blog posts
│   │   │   ├── projects/# Project descriptions
│   │   │   └── tags/    # Tag descriptions
│   │   └── master/      # Database schema and queries
│   └── routing/         # HTTP routing
├── shells/              # Development environment shells
└── [various config files]
```

## Database Schema

The application uses a SQLite database with the following main tables:

- `posts` - Blog articles
- `projects` - Portfolio projects
- `tags` - Skills and categories
- `post_tags` - Relationship between posts and tags
- `project_tags` - Relationship between projects and tags
- `embeddings` - Vector embeddings for search functionality. Nearly all content is indexed with embeddings.

## Setup and Development

### Prerequisites

- Go 1.24 or later
- Nix (optional, for reproducible development environment)

### Development with Nix (Recommended)

Nix provides a consistent development environment with all dependencies locked:

```bash
# Clone the repository
git clone https://github.com/conneroisu/conneroh.com.git
cd conneroh.com

# Enter development shell
nix develop

# Generate code and assets
generate-all

# Initialize the database
update

# Run the application with live reloading
run
```

### Development without Nix

```bash
# Clone the repository
git clone https://github.com/conneroisu/conneroh.com.git
cd conneroh.com

# Install dependencies
go mod download

# Install required tools
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/cosmtrek/air@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
npm install -g tailwindcss
npm install -g bun

# Generate code
templ generate
go generate ./...

# Build CSS and JS
tailwindcss -i input.css -o cmd/conneroh/_static/dist/style.css
bun build index.js --minify --outdir cmd/conneroh/_static/dist/

# Initialize the database
go run ./cmd/update

# Run the application with hot reload
air
```

### Content Management

Content is managed through Markdown files located in `internal/data/docs/`. The format for content files is:

#### Blog Post Example (`internal/data/docs/posts/example-post.md`):

```markdown
---
title: Example Post Title
slug: example-post-slug
description: Short description of the post
created_at: 1640995200  # Unix timestamp
banner_url: https://example.com/image.jpg  # Optional
tags:
  - go
  - web-development
projects:
  - related-project-slug  # Optional related projects
---

# Markdown Content

The actual content of the post in Markdown format...
```

#### Updating Content

To update the database with new or modified content:

```bash
# With Nix
update

# Without Nix
go run ./cmd/update
```


### Template Rendering with templ

The application uses templ for type-safe HTML templates:

```go
// Example of a templ component (simplified)
templ Post(post gen.Post) {
    <article>
        <h1>{post.Title}</h1>
        <div class="content">
            @templ.Raw(post.Content)
        </div>
    </article>
}
```

### Dynamic UI with HTMX and Alpine.js

The site uses HTMX for navigation and content loading:

```html
<!-- Example of HTMX usage for navigation -->
<a 
  hx-get="/hateoas/morph/projects" 
  hx-target="#bodiody" 
  hx-swap="outerHTML" 
  hx-push-url="/projects"
>
  Projects
</a>
```

Alpine.js is used for client-side interactivity:

```html
<!-- Example of Alpine.js usage for tabs -->
<div x-data="{ activeTab: 'posts' }">
  <button @click="activeTab = 'posts'" :class="{ 'active': activeTab === 'posts' }">
    Posts
  </button>
  <div x-show="activeTab === 'posts'">
    Posts content here...
  </div>
</div>
```

## Deployment

The application can be deployed as a standalone binary or as a Docker container.

### Building a Docker Container

```bash
# With Nix
nix build .#packages.x86_64-linux.C-conneroh
# Load into Docker
docker load < result
# Run the container
docker run -p 8080:8080 conneroh.com:latest
```

### Running the Docker Container

```bash
docker run -p 8080:8080 conneroh.com:latest
```

### Direct Binary

```bash
# Build the application
go build -o conneroh ./cmd/conneroh

# Run the built application
./conneroh
```

## Contributing

This project is personal, but suggestions and bug reports are welcome. Please open an issue or submit a pull request.


## Resources

- [ preload ](https://htmx.org/extensions/preload/)

## License

This project is personal and not licensed for public use without permission. Use at your own risk.

## Author

Conner Ohnesorge
