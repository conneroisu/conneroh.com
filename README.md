# Conner Ohnesorge's Portfolio
[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)
<img class="badge" tag="github.com/conneroisu/conneroh.com" src="https://goreportcard.com/badge/github.com/conneroisu/conneroh.com">

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
│   │   ├── assets/      # Assets Data Structure definitions
│   │   ├── docs/        # Markdown content (posts, projects, tags)
│   │   │   ├── posts/   # Blog posts
│   │   │   ├── projects/# Project descriptions
│   │   │   └── tags/    # Tag descriptions
│   │   └── gen/         # Generated data structures
│   └── routing/         # HTTP routing
└── [various config files]
```

The site organizes content into three main types:

- **Posts** - Blog articles
- **Projects** - Portfolio projects
- **Tags** - Skills and categories

Content relationships are maintained through associations between these entities.

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
nix-generate-all

# Initialize the database
update

# Run the application with live reloading
run
```

#### Development Shell Usage

```bash
<!-- BEGIN_MARKER -->
clean - Clean Project
dx - Edit flake.nix
format - Format code files
generate-all - Generate all files in parallel
generate-css - Update the generated html and css files.
generate-db - Update the generated go files from the md docs.
generate-docs - Update the generated documentation files.
generate-js - Generate JS files
generate-reload - Code Generation Steps for specific directory changes.
gx - Edit go.mod
interpolate - Interpolate templates; Usage: interpolate input_file start_marker end_marker replacement_text
lint - Run Nix/Go Linting Steps.
reset-db - Reset the database
run - Run the application with air for hot reloading
test - Run Vitest tests
test-ci - Run Vitest tests for CI
test-ui - Run Vitest with UI
tests - Run all go tests
<!-- END_MARKER -->
```

## Content Management

Content is managed through Markdown files located in `internal/data/`. The format for content files is:

### Blog Post Example (`internal/data/posts/example-post.md`):

```markdown
---
title: Example Post Title
slug: example-post-slug
description: Short description of the post
created_at: 2025-03-27T05:48:53.000-06:00
updated_at: 2025-03-27T14:13:10.000-06:00
banner_path: /dist/img/example-banner.jpg # Optional
tags:
  - go
  - web-development
projects:
  - related-project-slug # Optional related projects
---

# Markdown Content

The actual content of the post in Markdown format...
```

### Updating Content

To update the database with new or modified content:

```bash
# With Nix
update

# Without Nix
go run ./cmd/update
```

## Testing

The project uses Vitest with Playwright for comprehensive testing including unit tests and browser-based integration tests.

### Running Tests

```bash
# Enter development shell (if using Nix)
nix develop

# Install dependencies
bun install

# Run all tests
test
# or
bun test

# Run tests with UI
test-ui
# or
bun test:ui

# Run tests once (CI mode)
test-ci
# or
bun test:run

# Coverage report
bun test:coverage

# Run comprehensive test suite (includes app startup)
nix run .#runTests
```

### Test Structure

- `tests/browser/` - Browser integration tests using Playwright
- `tests/unit/` - Unit tests for utility functions
- `tests/setup.ts` - Test setup and utilities

Tests cover:
- Homepage functionality and navigation
- Project and blog post pages
- HTMX dynamic content loading
- Alpine.js interactivity
- Responsive design across viewports
- Accessibility compliance

### CI Integration

Tests run automatically before deployment:
- On PR creation/updates before preview deployment
- On main branch pushes before production deployment

## Deployment

The site is deployed using Fly.io. The deployment process is automated in the GitHub workflow located in `.github/workflows/fly-deploy.yml`.

To deploy manually:

```bash
nix run .#deployPackage
```

## Technical Implementation Details

### Template Rendering with templ

The application uses templ for type-safe HTML templates:

```go
// Example of a templ component (simplified)
templ Post(post assets.Post) {
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
  hx-get="/morph/projects"
  hx-target="#bodiody"
  hx-swap="outerHTML"
  hx-push-url="/projects"
>
  Projects
</a>
```

## Contributing

This project is personal, but suggestions and bug reports are most welcome. Please open an issue or submit a pull request.

## Author

Conner Ohnesorge
