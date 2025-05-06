# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Environment

The project uses Nix for a reproducible development environment. This provides a consistent set of dependencies and tools. You can access common development commands through the Nix shell.

## Key Commands

```bash
# Enter development shell (most commands are available after this)
nix develop

# Core Build and Development Commands
generate-all       # Generate all files in parallel (CSS, DB, JS)
generate-css       # Update the generated HTML and CSS files
generate-db        # Update the generated Go files from the MD docs
generate-js        # Generate JS files 
generate-docs      # Update the generated documentation files
run                # Run the application with air for hot reloading
tests              # Run all Go tests
format             # Format code files
lint               # Run Nix/Go linting steps

# Database Management
reset-db           # Reset the database
update             # Initialize/update the database with content from Markdown files

# Deployment
nix run .#deployPackage  # Deploy to Fly.io
```

## Project Architecture

This is a personal portfolio website built with Go, using a server-side rendering approach with several modern web technologies:

- **Backend**: Go with `templ` for HTML templates
- **Frontend**: TailwindCSS for styling, Alpine.js for client-side interactivity, HTMX for dynamic content loading
- **Database**: SQLite with Bun ORM for data persistence
- **Content**: Markdown files for blog posts, projects, and tag descriptions

### Core Components

1. **Main Server (`cmd/conneroh/`)**: 
   - Handles HTTP routing and request processing
   - Renders templates with data from the database
   - Implements route handlers for different content types

2. **Content Management (`cmd/update/`)**: 
   - Processes Markdown files to populate the database
   - Handles relationships between content types
   - Regenerates content when changes are detected

3. **Data Models (`internal/assets/`)**: 
   - Defines the core data structures: Posts, Projects, Tags
   - Manages relationships between these entities
   - Handles database initialization and schema

4. **Templating System**:
   - Uses `templ` for type-safe HTML templates
   - Layouts define the overall page structure
   - Views implement specific page types
   - Components are reusable UI elements

### Content Organization

The site organizes content into three main types:

1. **Posts**: Blog articles with markdown content
2. **Projects**: Portfolio project descriptions
3. **Tags**: Skills and technologies that can be linked to posts and projects

Content is stored as Markdown files in `internal/data/docs/` with the following structure:
- `posts/` - Blog posts
- `projects/` - Project descriptions
- `tags/` - Tag descriptions

Each Markdown file includes YAML frontmatter with metadata and relationships to other content types.

### Workflow for Content Updates

1. Edit Markdown files in `internal/data/docs/`
2. Run `generate-db` or `update` to process changes
3. The database is updated with the new content
4. Run `generate-css` if template changes were made
5. Start the application with `run` to see changes

## Development Workflow

1. **Setup**: Enter the development environment with `nix develop`
2. **Code Generation**: Run `generate-all` to ensure all generated files are up to date
3. **Run the App**: Use `run` to start the application with hot reloading
4. **Testing**: Run `tests` to execute all tests
5. **Formatting**: Use `format` to format Go, JS, CSS, and other files
6. **Linting**: Run `lint` to check for code issues

## Important Implementation Details

1. **Template Rendering**: The application uses `templ` for type-safe HTML templates with Go
2. **Dynamic UI**: HTMX is used for navigation and content loading without full page refreshes
3. **Client Interactivity**: Alpine.js provides lightweight client-side interactivity
4. **Responsive Design**: TailwindCSS is used for responsive styling
5. **Database**: SQLite with Bun ORM for efficient data access and management
6. **Content Processing**: Markdown files are processed with frontmatter to extract metadata

## Content Management

Content is managed through Markdown files with YAML frontmatter:

```markdown
---
title: Example Title
slug: example-slug
description: Short description
created_at: 2025-03-27T05:48:53.000-06:00
updated_at: 2025-03-27T14:13:10.000-06:00
banner_path: /dist/img/example-banner.jpg # Optional
tags:
  - go
  - web-development
projects:
  - related-project-slug # Optional related projects
---

# Content

The actual content in Markdown format...
```

After editing content files, run `update` to refresh the database.