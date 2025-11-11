# content-pipeline Specification Delta

## ADDED Requirements

### Requirement: Markdown Content as Source of Truth (REQ-CONTENT-001)

Markdown files with YAML frontmatter SHALL be the authoritative source for all content (posts, projects, tags), with the database serving as a generated representation.

**Priority**: Critical

**Rationale**: Content stored in version-controlled markdown files provides history, diffability, and portability. The database is a materialized view optimized for serving content to the web application.

#### Scenario: Frontmatter defines required metadata

- **GIVEN** a markdown file in `internal/data/docs/posts/`
- **WHEN** the file is parsed
- **THEN** YAML frontmatter contains required fields: title, slug, description, created_at, updated_at
- **AND** optional fields include: banner_path, tags (array of tag slugs)
- **AND** timestamps follow ISO 8601 format with timezone (e.g., 2025-03-27T05:48:53.000-06:00)

#### Scenario: Markdown body contains post content

- **GIVEN** a post markdown file
- **WHEN** the file is read
- **THEN** content after frontmatter is valid markdown
- **AND** markdown supports extensions: goldmark-frontmatter, goldmark-highlighting, goldmark-mathjax, goldmark-mermaid, goldmark-wikilink, goldmark-hashtag, goldmark-obsidian-callout
- **AND** content is stored as raw markdown in the database

#### Scenario: Slug uniqueness is enforced

- **GIVEN** multiple markdown files for the same content type
- **WHEN** files are processed
- **THEN** each slug must be unique within its type (posts, projects, tags)
- **AND** duplicate slugs are detected and reported as errors
- **AND** slug format is URL-safe (lowercase, hyphens, no special characters)

### Requirement: Content Processing Pipeline (REQ-CONTENT-002)

The system SHALL provide a mechanism to process markdown files and populate the PostgreSQL database with structured content data.

**Priority**: Critical

#### Scenario: Parse markdown files into database records

- **GIVEN** markdown files in `internal/data/docs/{posts,projects,tags}/`
- **WHEN** the content pipeline runs
- **THEN** each markdown file is parsed
- **AND** frontmatter is extracted and validated
- **AND** markdown content is extracted
- **AND** database records are created or updated via sqlc queries
- **AND** relationships (post-tag associations) are established

#### Scenario: Handle markdown parsing errors

- **GIVEN** a markdown file with invalid YAML frontmatter
- **WHEN** the content pipeline attempts to parse the file
- **THEN** parsing fails with a clear error message
- **AND** the error includes file path and line number
- **AND** the pipeline continues processing other files
- **AND** the failed file is reported in a summary

#### Scenario: Incremental updates vs full regeneration

- **GIVEN** an existing database with content
- **WHEN** the content pipeline runs
- **THEN** the system can either:
  - Truncate all content tables and regenerate from scratch (full rebuild)
  - Or compare existing records and only update changed files (incremental update)
- **AND** the approach is configurable (e.g., CLI flag or environment variable)
- **AND** full rebuilds are idempotent (same source files produce same database state)

### Requirement: Tag and Relationship Management (REQ-CONTENT-003)

The system SHALL automatically create tags referenced in post/project frontmatter and establish many-to-many relationships via junction tables.

**Priority**: High

#### Scenario: Create tags from frontmatter references

- **GIVEN** a post with frontmatter field `tags: [golang, web-development, htmx]`
- **WHEN** the content pipeline processes the post
- **THEN** tags with slugs `golang`, `web-development`, `htmx` are created if they don't exist
- **AND** tag names are derived from slugs (titlecase transformation)
- **AND** junction table entries (posts_tags) are created linking the post to each tag
- **AND** existing tags are reused (no duplicates)

#### Scenario: Update tag associations on content change

- **GIVEN** a post previously associated with tags [golang, htmx]
- **WHEN** the post frontmatter is updated to `tags: [golang, postgresql]`
- **THEN** the posts_tags junction entries are updated
- **AND** the htmx association is removed
- **AND** a new postgresql tag is created and associated
- **AND** the golang association remains unchanged

#### Scenario: Orphaned tags are preserved

- **GIVEN** a tag that is no longer referenced by any posts or projects
- **WHEN** the content pipeline runs
- **THEN** the tag remains in the database
- **AND** orphaned tags are reported (optional: flag to clean them up)
- **AND** no tags are deleted automatically

### Requirement: Content Type Detection and Routing (REQ-CONTENT-004)

The system SHALL automatically detect content type (post, project, tag) based on directory structure and apply appropriate schema validation.

**Priority**: Medium

#### Scenario: Detect post from directory

- **GIVEN** a markdown file in `internal/data/docs/posts/my-article.md`
- **WHEN** the content pipeline processes files
- **THEN** the file is recognized as a post
- **AND** post-specific frontmatter validation is applied
- **AND** the record is inserted into the posts table
- **AND** post-tag relationships are established

#### Scenario: Detect project from directory

- **GIVEN** a markdown file in `internal/data/docs/projects/my-project.md`
- **WHEN** the content pipeline processes files
- **THEN** the file is recognized as a project
- **AND** project-specific frontmatter validation is applied (includes project-specific fields)
- **AND** the record is inserted into the projects table
- **AND** project-tag relationships are established

#### Scenario: Detect tag definition file

- **GIVEN** a markdown file in `internal/data/docs/tags/golang.md`
- **WHEN** the content pipeline processes files
- **THEN** the file is recognized as a tag definition
- **AND** tag metadata (name, slug, description) is extracted from frontmatter
- **AND** the record is inserted or updated in the tags table
- **AND** tag definitions take precedence over auto-generated tags from references

### Requirement: Validation and Error Reporting (REQ-CONTENT-005)

The content pipeline SHALL validate all markdown files against schema requirements and provide detailed error reports for invalid content.

**Priority**: High

#### Scenario: Validate required frontmatter fields

- **GIVEN** a post markdown file
- **WHEN** the content pipeline parses the file
- **THEN** required fields (title, slug, description, created_at, updated_at) are checked
- **AND** missing fields trigger validation errors
- **AND** errors include field name, file path, and expected format
- **AND** the file is skipped if validation fails

#### Scenario: Validate timestamp format

- **GIVEN** frontmatter with created_at field
- **WHEN** the timestamp is parsed
- **THEN** the format must be ISO 8601 with timezone
- **AND** invalid formats (e.g., missing timezone, wrong delimiter) are rejected
- **AND** error message shows expected format example

#### Scenario: Validation summary report

- **GIVEN** the content pipeline processes multiple markdown files
- **WHEN** processing completes
- **THEN** a summary is printed showing:
  - Total files processed
  - Successful records created/updated
  - Failed files with error counts
  - File paths of failed files
- **AND** exit code is non-zero if any files failed

### Requirement: Manual Execution Model (REQ-CONTENT-006)

The content pipeline SHALL be manually triggered by developers, not automatically run during application startup or deployment.

**Priority**: Medium

**Rationale**: Content updates are deliberate and infrequent. Manual execution provides control and visibility into content changes.

#### Scenario: Developer runs content pipeline on demand

- **GIVEN** markdown files have been added or modified
- **WHEN** a developer runs the content pipeline command (e.g., via Nix shell command or direct Go binary)
- **THEN** the pipeline processes all markdown files
- **AND** database is updated with new/modified content
- **AND** developer reviews the output summary
- **AND** application does not need to restart (database is the live source)

#### Scenario: Pipeline integrated with Nix development workflow

- **GIVEN** a Nix development shell environment
- **WHEN** a developer runs a Nix shell command (e.g., `process-content` or `sync-content`)
- **THEN** the content pipeline binary is executed
- **AND** PostgreSQL connection credentials are sourced from environment
- **AND** output is displayed in the terminal
- **AND** the command is documented in the project README and `flake.nix`

#### Scenario: Content pipeline is independent of web application

- **GIVEN** the web application is running
- **WHEN** the content pipeline is executed
- **THEN** the pipeline operates independently (separate process)
- **AND** database changes are immediately reflected in the web app (via queries)
- **AND** no application restart is required
- **AND** the pipeline does not depend on application code (separate entry point)
