# Conner Ohnesorge's Portfolio

A personal portfolio website built with Go and modern web technologies, showcasing projects, blog posts, and skills.

## Overview

This portfolio website is built using Go with templ for generating HTML templates, Alpine.js for client-side interactivity, and Tailwind CSS for styling.

It follows a modern server-side rendering approach with HTMX for dynamic content loading without full page refreshes.

The site organizes content into projects, posts, and tags, all backed by a SQLite database with content managed through Markdown files in the repository.

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
- **Database**: SQLite (via turso)
- **CSS**: Tailwind CSS
- **Frontend Interactivity**: Alpine.js
- **Dynamic Content**: HTMX
- **Icons**: Nerd Fonts
- **Build/Development**: Nix, air (for live reloading)

## Project Structure

```
├── cmd/                 # Entry points for executables
│   ├── conneroh/        # Main web application
│   │   ├── layouts/     # Layout templates
│   │   ├── views/       # View templates
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
- Turso account (for database)

### Environment Variables

The application requires the following environment variables:

- `TURSO_URI` - Turso database URI
- `TURSO_TOKEN` - Turso authentication token

### Development with Nix

If you have Nix installed with flakes enabled:

```bash
# Enter development shell
nix develop

# Run the application with live reloading
run
```

### Development without Nix

```bash
# Install dependencies
go mod download

# Install tools
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/cosmtrek/air@latest
npm install -g tailwindcss

# Generate code
templ generate
go generate ./...

# Build CSS
tailwindcss -i input.css -o cmd/conneroh/_static/dist/style.css

# Run the application with hot reload
air
```

### Content Management

Content is managed through Markdown files located in `internal/data/docs/`. To update the database with new content:

```bash
# With Nix
update

# Without Nix
go run ./cmd/update
```

## Deployment

The application can be deployed as a standalone binary or as a Docker container. A Dockerfile is provided for containerized deployment.

```bash
# Build the application
go build -o conneroh ./cmd/conneroh

# Run the built application
./conneroh
```

## License

This project is personal and not licensed for public use without permission. Use at your own risk.

## Author

Conner Ohnesorge
