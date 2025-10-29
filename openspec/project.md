# Project Context

## Purpose
Personal portfolio website showcasing projects, blog posts, and technical skills. Built with modern web technologies using server-side rendering for optimal performance and SEO. The site demonstrates clean architecture, type-safe templates, and progressive enhancement through HTMX and Alpine.js.

## Tech Stack

### Backend
- **Go 1.24+** - Primary backend language
- **templ** - Type-safe HTML templating
- **Bun ORM** - Database abstraction layer
- **SQLite** - Embedded database (modernc.org/sqlite)

### Frontend
- **Alpine.js** - Client-side reactivity and interactions
- **HTMX 2.0** - Dynamic content loading without full page refreshes
- **Tailwind CSS v4** - Utility-first CSS framework
- **Preline** - UI component library

### Build & Development
- **Nix** - Reproducible development environment
- **Air** - Live reloading during development
- **Bun** - JavaScript package management and bundling
- **golangci-lint** - Go code linting

### Testing
- **Playwright** - Browser integration testing
- **happy-dom** - Lightweight DOM implementation for testing

### Infrastructure
- **Fly.io** - Application deployment
- **AWS S3** - Asset storage
- **GitHub Actions** - CI/CD automation

### Content Processing
- **Goldmark** - Markdown parser with extensions:
  - goldmark-frontmatter - YAML frontmatter parsing
  - goldmark-highlighting - Syntax highlighting
  - goldmark-mathjax - Math rendering
  - goldmark-mermaid - Diagram rendering
  - goldmark-wikilink - Wiki-style links
  - goldmark-hashtag - Hashtag support
  - goldmark-obsidian-callout - Callout blocks

## Project Conventions

### Code Style
- **templ files**: 2-space indentation (see .editorconfig)
- **Go code**: Follow golangci-lint configuration with enabled linters:
  - asasalint, exhaustive, godox, nlreturn, bidichk, bodyclose
  - dupl, gocritic, godot, errname, copyloopvar, gosec, goconst
  - intrange, perfsprint, usestdlibvars, staticcheck
  - misspell, revive, unconvert, wastedassign, whitespace, govet
- **Generated code**: Excluded from linting (vendor/, third_party/, builtin/, examples/)
- **TypeScript**: Follows standard TypeScript 5.7+ conventions

### Architecture Patterns
- **Server-Side Rendering (SSR)**: Primary rendering approach using templ
- **Progressive Enhancement**: Core functionality works without JS, enhanced with Alpine.js/HTMX
- **Content-First Architecture**: Markdown files as source of truth for posts/projects
- **Type Safety**: templ provides compile-time type checking for templates
- **Code Generation**: Database schemas and assets auto-generated from markdown content
- **Clean Architecture**: Separation of concerns via cmd/ (executables) and internal/ (private code)

### Project Structure
```
cmd/                 # Entry points
├── conneroh/        # Main web application
│   ├── layouts/     # Layout templates
│   ├── views/       # View templates
│   ├── components/  # Component templates
│   └── _static/     # Static assets
├── update/          # Content update utility
└── update-css/      # CSS update utility

internal/            # Private application code
├── assets/          # Asset data structures
├── cache/           # Caching layer
├── copygen/         # Code generation utilities
├── data/            # Data access layer
│   ├── assets/      # Generated asset definitions
│   ├── docs/        # Markdown content
│   │   ├── posts/   # Blog posts
│   │   ├── projects/# Portfolio projects
│   │   └── tags/    # Skill tags
│   └── gen/         # Generated data structures
├── logger/          # Logging utilities
└── routing/         # HTTP routing
```

### Content Management
- **Markdown Files**: Source of truth for all content
- **Frontmatter Format**: YAML with required fields (title, slug, description, created_at, updated_at)
- **Content Types**: Posts, Projects, Tags
- **Relationships**: Tags and projects linked via frontmatter references
- **Code Generation**: Run `update` command to regenerate database from markdown files

### Testing Strategy
- **Browser Tests**: Playwright for integration testing
  - Homepage functionality and navigation
  - Project and blog post pages
  - HTMX dynamic content loading
  - Alpine.js interactivity
  - Responsive design across viewports
  - Accessibility compliance
- **CI Integration**: Tests run on PR creation and before deployments
- **Test Commands**:
  - `nix develop -c tests` - Run all tests
  - Tests located in `tests/browser/` and `tests/unit/`

### Development Workflow
- **Nix Shell**: Enter with `nix develop`
- **Available Commands**:
  - `run` - Run app with hot reloading
  - `generate-all` - Generate all files in parallel
  - `generate-db` - Update generated go files from markdown docs
  - `generate-css` - Update generated HTML and CSS files
  - `generate-js` - Generate JS files
  - `update` / `reset-db` - Database management
  - `tests` - Run all go tests
  - `lint` - Run Nix/Go linting
  - `format` - Format code files
  - `clean` - Clean project

### Git Workflow
- **Main Branch**: `main` (production deployment target)
- **CI/CD**: GitHub Actions workflow (`.github/workflows/fly-deploy.yml`)
- **Deployment**: Automated on main branch pushes
- **Manual Deploy**: `nix run .#deployPackage`

## Domain Context

### Content Organization
- **Posts**: Blog articles with markdown content, metadata, and tag associations
- **Projects**: Portfolio projects with descriptions, technologies, and related posts
- **Tags**: Skills and categories that connect posts and projects

### Navigation Pattern
- **HTMX Morphing**: Smooth page transitions using `hx-get`, `hx-target="#bodiody"`, `hx-swap="outerHTML"`
- **URL Management**: `hx-push-url` for proper browser history
- **Preloading**: htmx-ext-preload for faster navigation

### Content Metadata
- **Timestamps**: ISO 8601 format with timezone (e.g., 2025-03-27T05:48:53.000-06:00)
- **Slugs**: URL-friendly identifiers for posts/projects/tags
- **Banner Images**: Optional banner_path in frontmatter

## Important Constraints

### Technical Constraints
- **Nix-First**: All development must work within Nix development environment
- **Type Safety**: All templates must use templ (no raw HTML generation)
- **Generated Code**: Never manually edit generated files (files in gen/ directories)
- **Content Format**: All content must be valid markdown with proper YAML frontmatter
- **Database**: SQLite only - no external database dependencies

### Build Constraints
- **Reproducible Builds**: All dependencies must be Nix-compatible
- **Code Generation**: Template/CSS/DB changes require regeneration steps
- **Static Assets**: Assets must be in cmd/conneroh/_static/

### Performance Constraints
- **Core Web Vitals**: Optimize for LCP, FID, CLS
- **Progressive Enhancement**: Site must be functional without JavaScript
- **SSR-First**: No client-side routing or hydration delays

## External Dependencies

### Required Services
- **Fly.io**: Production deployment platform
- **AWS S3**: Asset storage (configured via AWS SDK v2)
- **GitHub**: Source control and CI/CD (GitHub Actions)

### Optional Integrations
- **MathJax**: Mathematical equation rendering
- **Mermaid**: Diagram generation from markdown
- **Chroma**: Syntax highlighting for code blocks

### MCP Servers (Claude Code)
- **context7**: Library documentation lookup
- **laravel-boost**: Laravel ecosystem support (note: this project is Go-based, not Laravel)
