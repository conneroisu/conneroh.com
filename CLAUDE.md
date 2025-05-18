# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a personal portfolio website built with Go, using the following tech stack:
- Backend: Go
- Templates: templ
- CSS: Tailwind CSS
- Frontend Interactivity: Alpine.js
- Dynamic Content: HTMX
- Build/Development: Nix, air (for live reloading)

The site organizes content into three main types:
- Posts - Blog articles
- Projects - Portfolio projects
- Tags - Skills and categories

## Repository Structure

The project follows this structure:
- `cmd/conneroh/`: Main web application
  - `layouts/`: Layout templates
  - `views/`: View templates
  - `components/`: Component templates
  - `_static/`: Static assets
- `internal/`: Private application code
  - `data/`: Data access layer
    - `assets/`: Assets Data Structure definitions
    - `docs/`: Markdown content (posts, projects, tags)
  - `routing/`: HTTP routing

## Development Commands

### Using Nix Development Shell (Recommended)

```bash
# Enter development shell 
nix develop

# Generate all code and assets
generate-all

# Run the application with live reloading
run
```

### Common Development Tasks

```bash
# Generate CSS files
generate-css

# Generate database from markdown files
generate-db

# Generate JS files
generate-js

# Format code
format

# Run tests
tests

# Lint code
lint

# Reset the database
reset-db
```

### Content Management

Content is stored as Markdown files in `internal/data/docs/` directories:
- `internal/data/docs/posts/`: Blog posts
- `internal/data/docs/projects/`: Project descriptions
- `internal/data/docs/tags/`: Tag descriptions

After adding or modifying content, update the database with:
```bash
generate-db
```

Alternatively, you can run the script directly:
```bash
./generate-db.sh
```

## Build and Deploy

### Local Build

```bash
# Build the application (in nix development shell)
go build -o ./conneroh.com .
```

### Deployment

The site is deployed using Fly.io:

```bash
# Deploy to production
nix run .#deployPackage
```

## Content Workflow

When working with content, remember:
1. Content is written in Markdown files in the `internal/data/docs/` directory
2. After changes, run `generate-db` to update the database
3. Run `generate-css` if template changes affect the CSS
4. Use `run` command to see changes with hot-reloading