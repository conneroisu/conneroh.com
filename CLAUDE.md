# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a personal portfolio website for Conner Ohnesorge built as a sophisticated Go web application. It demonstrates modern web development patterns with server-side rendering using templ templates, HTMX for seamless navigation, Alpine.js for client-side interactivity, and a markdown-driven content management system that generates a SQLite database.

### Core Architecture
- **Backend**: Go 1.24 with type-safe templ templates
- **Database**: SQLite with Bun ORM for relationships and full-text search
- **Frontend**: Progressive enhancement using HTMX + Alpine.js + Tailwind CSS
- **Content**: Markdown files with YAML frontmatter automatically converted to database
- **Deployment**: Containerized deployment on Fly.io with PR preview environments
- **Development**: Nix flake for reproducible development environment with comprehensive tooling

## Key Commands

### Development Environment (Nix)
All commands should be run using the Nix development environment with `nix develop -c <command>`:

```bash
# Available commands:
nix develop -c clean              # Clean project
nix develop -c format             # Format code files
nix develop -c generate-all       # Generate all files in parallel
nix develop -c generate-css       # Update generated HTML and CSS files
nix develop -c generate-db        # Update generated Go files from markdown docs
nix develop -c generate-docs      # Update generated documentation files
nix develop -c generate-js        # Generate JS files
nix develop -c lint               # Run Nix/Go linting steps
nix develop -c run                # Run application with air for hot reloading
nix develop -c tests              # Run Vitest, Go, and Playwright tests
nix develop -c update             # Update database with content changes
nix develop -c reset-db           # Reset the database
```

### Standard Development Workflow
```bash
# 1. Setup and initialize
nix develop -c generate-all       # Generate all required files
nix develop -c update             # Initialize/update database from markdown
nix develop -c run                # Start server with hot reloading

# 2. After content changes (editing markdown files in internal/data/)
nix develop -c update             # Regenerate database from markdown files

# 3. Run tests
nix develop -c tests              # Go unit tests, Playwright tests, and Vitest
```

### Build and Run (without Nix)
```bash
# Generate templates and run
go generate ./...
go run main.go

# Update database manually
doppler run -- go run ./cmd/update


# Run tests manually
bun install && bun test
go test ./...
```

## Architecture

### Content Management System
- **Markdown-driven content** in `internal/data/` with frontmatter metadata
- **Automatic database generation** from markdown files using `cmd/update`
- **Three content types**: Posts (blog), Projects (portfolio), Employments (work history)
- **Tag system** for categorization with bidirectional relationships
- **SQLite database** with Bun ORM for type-safe queries and relationships

### Template and Rendering System
- **templ templates** for type-safe HTML generation with Go integration
- **Component-based architecture** in `cmd/conneroh/components/`
- **Layout inheritance** with base layouts in `cmd/conneroh/layouts/`
- **Code generation** converts .templ files to Go code via `templ generate`

### Frontend Architecture
- **HTMX** for seamless navigation and partial page updates
- **Alpine.js** for client-side interactivity and state management
- **Tailwind CSS** with custom class generation
- **Progressive enhancement** approach - works without JS

### Project Structure
```
cmd/conneroh/          # Main web application
├── layouts/           # Page layout templates (templ)
├── views/             # Page view templates (templ)
├── components/        # Reusable component templates (templ)
├── classes/           # CSS class generation
└── _static/           # Static assets and dist files

internal/
├── assets/            # Asset management, data types, validation
├── data/              # Content files and media assets
│   ├── posts/         # Blog post markdown files
│   ├── projects/      # Project markdown files  
│   ├── employments/   # Employment history markdown
│   ├── tags/          # Tag definition markdown
│   └── assets/        # Media files (images, videos, SVGs)
├── routing/           # HTTP routing, handlers, pagination
└── logger/            # Logging utilities

tests/
├── browser/           # Browser integration tests (Playwright)
└── unit/              # Unit tests
```

### Code Generation Workflow
1. **templ generate** converts .templ files to Go code
2. **cmd/update** parses markdown files and populates SQLite database
3. **cmd/update-css** generates optimized CSS classes
4. **go generate** runs all generation steps

### Database Schema
- Content entities (Post, Project, Employment, Tag) with full-text search
- Relationship tables for many-to-many associations
- Automatic schema creation and migration via Bun ORM
- Database file: `master.db`

## Development Notes

### Content Updates
When editing markdown files in `internal/data/`, always run `nix develop -c update` to regenerate the database. The application reads from the database, not directly from markdown files.

### Template Development
- Edit .templ files in `cmd/conneroh/` subdirectories
- Run `nix develop -c generate-css` to update generated Go code and CSS
- Hot reloading with Air automatically detects changes when using `nix develop -c run`

### Testing Strategy
- Browser tests cover HTMX navigation, Alpine.js interactions, and responsive design
- Unit tests focus on utility functions and data processing
- CI runs tests automatically before deployment to Fly.io

### Asset Management
- Images optimized to WebP format in `internal/data/assets/`
- SVG icons stored in `internal/data/assets/svg/`
- Static assets served from `cmd/conneroh/_static/dist/`

## Technical Implementation Details

### Performance & Caching Strategy
- **Component-level caching**: Full page components cached in memory (home, project lists)
- **Slug-mapped caching**: Individual content items cached by slug for instant retrieval
- **Global content caches**: `allPosts`, `allProjects`, `allTags`, `allEmployments` loaded once at startup
- **Concurrent search**: Pool-based goroutines (max 10) for search operations with relevance scoring
- **Hash-based optimization**: Change detection throughout build pipeline to minimize unnecessary work

### Database & Content Management
- **Schema**: Post, Project, Employment, Tag entities with many-to-many relationships via junction tables
- **Full-text search**: Built-in SQLite FTS capabilities across title, description, content, and tags
- **Content pipeline**: Markdown → Hash detection → Database generation → Asset upload to Tigris S3
- **Concurrent processing**: 20 workers for file operations during content updates
- **Custom time handling**: Supports both date-only and RFC3339 formats in frontmatter

### Security & Error Handling
- **Input validation**: MIME type validation for assets (images, PDFs only)
- **Error management**: Eris error wrapping with context throughout the stack
- **Graceful shutdown**: 10-second timeout with proper signal handling
- **Production security**: Doppler secrets management, CORS headers configured

### HTMX & Frontend Patterns
- **Target-based updates**: Use `#bodiody` for main content swaps
- **URL state management**: Always include `hx-push-url` for navigation
- **Preloading strategy**: Hover/touch events for performance optimization
- **Progressive enhancement**: Semantic HTML base with HTMX/Alpine.js layered on top

### templ Template System
- **Type safety**: End-to-end type safety from database to rendered HTML
- **Component architecture**: Reusable components in `components/` directory
- **Layout inheritance**: Base layouts with header/footer separation
- **Twerge integration**: CSS class management utility for dynamic styling

## Development Guidelines

### Code Style & Patterns
- **Follow existing naming conventions**: CamelCase for exported functions, lowercase for internal
- **Error handling**: Always wrap errors with Eris for better debugging context
- **Concurrent processing**: Use `github.com/sourcegraph/conc` for goroutine pools
- **Database queries**: Leverage Bun ORM's type-safe query builder over raw SQL

### Content Management Workflow
- **Always regenerate database** after editing markdown files using `nix develop -c update`
- **Test content changes** locally before committing (database must be updated for changes to appear)
- **Asset optimization**: New images should be converted to WebP format
- **Frontmatter validation**: Ensure all required fields (title, slug, description, dates) are present

### Testing Requirements
- **Browser tests**: Cover HTMX navigation, Alpine.js interactions, responsive design
- **Accessibility**: Test ARIA attributes, skip links, keyboard navigation
- **Performance**: Verify caching behavior and search responsiveness
- **CI integration**: All tests must pass before deployment

### Deployment Considerations
- **Environment variables**: Use Doppler for production secrets, avoid hardcoding
- **Database bundling**: SQLite file is included in container image (no external DB needed)
- **Auto-scaling**: Fly.io configured for 0-2 machines, 1GB RAM each
- **PR previews**: Automatic preview app creation/destruction for pull requests

## Common Pitfalls & Solutions

### Content Not Updating
- **Problem**: Changed markdown files but content doesn't appear on site
- **Solution**: Run `nix develop -c update` to regenerate database from markdown files

### Template Compilation Errors
- **Problem**: templ templates failing to compile
- **Solution**: Check for syntax errors, ensure proper Go integration, run `nix develop -c generate-css`

### Cache Issues
- **Problem**: Changes not reflecting immediately
- **Solution**: Component caching is aggressive; restart server or clear cache maps in handlers

### HTMX Navigation Problems
- **Problem**: Links not working with HTMX
- **Solution**: Ensure `hx-target="#bodiody"`, `hx-swap="outerHTML"`, and `hx-push-url` attributes are present

### Build Failures
- **Problem**: Nix build failing
- **Solution**: Check hash files in `internal/cache/`, may need to clear and regenerate

## File Location Reference

### Key Files to Know
- `main.go` - Application entry point
- `cmd/conneroh/root.go` - HTTP server configuration and graceful shutdown
- `cmd/conneroh/handlers.go` - Request handlers with caching logic
- `cmd/conneroh/routes.go` - Route definitions and static file serving
- `cmd/update/main.go` - Content pipeline for markdown → database conversion
- `internal/assets/static.go` - Database schema and type definitions
- `internal/routing/handlers.go` - Search and filtering logic with concurrent processing
- `flake.nix` - Nix development environment and build scripts
- `master.db` - SQLite database file (generated, do not edit directly)

### Content Directories
- `internal/data/posts/` - Blog post markdown files
- `internal/data/projects/` - Project portfolio markdown files
- `internal/data/employments/` - Work history markdown files
- `internal/data/tags/` - Tag definition markdown files
- `internal/data/assets/` - Media files (images, videos) for content

### Generated Files (Do Not Edit)
- `cmd/conneroh/**/*_templ.go` - Generated from .templ files
- `cmd/conneroh/classes/classes.html` - Generated CSS classes
- `cmd/conneroh/_static/dist/` - Built assets
- `internal/cache/*.hash` - Build cache files
