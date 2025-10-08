# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Environment Setup

This project uses Nix for reproducible development environments. Enter the development shell with:

```bash
nix develop
```

## Common Development Commands

### Code Generation & Building
- `generate-all` - Generate all artifacts (CSS, DB, JS) in parallel
- `generate-css` - Update HTML and CSS files using templ and Tailwind
- `generate-db` - Update generated Go files from markdown content (requires doppler)
- `generate-js` - Build/minify JavaScript files
- `generate-reload` - Run code generation for specific directory changes

### Development Server
- `run` - Start the application with hot reloading using air

### Testing
- `tests` - Run all Go tests (generates templates first)

### Code Quality
- `lint` - Run comprehensive linting (Go, Nix, Rust)
  - Includes: golangci-lint, statix, deadnix, nix flake check
- `format` - Format code files using treefmt (Go, Nix, Rust)

### Database
- `update` - Initialize/update the database with content changes
- `reset-db` - Reset the database

### Utilities
- `clean` - Clean project (git clean -fdx)
- `dx` - Edit flake.nix
- `gx` - Edit go.mod

## Project Architecture

This is a Go-based personal portfolio website using modern web technologies:

### Core Technologies
- **Backend**: Go 1.24
- **Templates**: templ for type-safe HTML templates
- **CSS**: Tailwind CSS v4
- **Frontend Interactivity**: Alpine.js + HTMX
- **Database**: SQLite with Bun ORM
- **Build System**: Nix flakes

### Directory Structure
- `cmd/conneroh/` - Main web application entry point
  - `views/` - Generated templ templates
  - `components/` - Reusable template components
  - `layouts/` - Base layout templates
  - `_static/` - Static assets (CSS, JS, images)
- `cmd/update/` - Content management utility
- `internal/` - Private application code
  - `data/` - Data access layer and content management
    - `assets/` - Asset data structures
    - `docs/` - Markdown content (posts, projects, tags)
    - `gen/` - Generated data structures
  - `routing/` - HTTP routing and handlers
  - `logger/` - Logging configuration
  - `copygen/` - Code generation utilities

### Content Management

Content is stored as Markdown files in `internal/data/docs/`:
- `posts/` - Blog articles
- `projects/` - Portfolio projects
- `tags/` - Skill/category descriptions
- `employments/` - Work history

Each content file includes YAML frontmatter with metadata (title, slug, dates, tags, projects).

### Database Schema

The application uses SQLite with the following main entities:
- Posts (blog articles)
- Projects (portfolio items)
- Tags (skills/categories)
- Employments (work history)

Relationships are maintained through association tables.

### Template System

Uses templ for type-safe templates with component-based architecture:
- Layouts provide base structure
- Components are reusable UI elements
- Views render specific pages (home, posts, projects, tags)

### Deployment

Deployed to Fly.io with automated deployments via GitHub Actions.
- Production app: `conneroh-com.fly.dev`
- PR previews available via `pr-preview deploy <pr-number>`

## Development Workflow

1. Enter nix development environment
2. Run `generate-all` to generate initial artifacts
3. Use `run` to start development server with hot reloading
4. Edit markdown content in `internal/data/docs/`
5. Run `update` to regenerate database from content
6. Use `tests` to verify changes
7. Run `lint` before committing

## Testing Strategy

- Go unit tests for business logic
- Integration tests for HTTP handlers
- Browser tests using Playwright for UI interactions
- Tests are located in `tests/` directory
- Run with `tests` command (includes template generation)

## Key Dependencies

- `github.com/a-h/templ` - Type-safe templates
- `github.com/uptrace/bun` - SQL ORM
- `modernc.org/sqlite` - SQLite driver
- `github.com/yuin/goldmark` - Markdown processing
- `github.com/gorilla/schema` - Form decoding
- `github.com/aws/aws-sdk-go-v2` - S3 integration
