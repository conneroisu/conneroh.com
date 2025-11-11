# Implementation Tasks

## 1. Database Schema and sqlc Setup

- [ ] 1.1 Review and validate existing PostgreSQL schema in `internal/db/migrations/001_init.sql`
- [ ] 1.2 Ensure schema includes all necessary tables, indexes, and constraints for application needs
- [ ] 1.3 Verify PostGIS extension is enabled and geometry columns are properly configured
- [ ] 1.4 Validate `sqlc.yaml` configuration for PostgreSQL and pgx/v5 driver settings

## 2. Implement sqlc Query Definitions

- [ ] 2.1 Create `internal/db/queries/posts.sql` with comprehensive post queries:
  - [ ] ListPosts (paginated, ordered by created_at DESC)
  - [ ] GetPostByID
  - [ ] GetPostBySlug
  - [ ] CreatePost
  - [ ] UpdatePost
  - [ ] DeletePost
- [ ] 2.2 Create `internal/db/queries/projects.sql` with project queries:
  - [ ] ListProjects (paginated, ordered by display order or created_at)
  - [ ] GetProjectByID
  - [ ] GetProjectBySlug
  - [ ] CreateProject
  - [ ] UpdateProject
  - [ ] DeleteProject
- [ ] 2.3 Create `internal/db/queries/tags.sql` with tag queries:
  - [ ] ListTags (ordered by name or usage count)
  - [ ] GetTagByID
  - [ ] GetTagBySlug
  - [ ] CreateTag
  - [ ] UpdateTag
  - [ ] DeleteTag
- [ ] 2.4 Create `internal/db/queries/employments.sql` with employment queries (including PostGIS location):
  - [ ] ListEmployments (ordered by start date)
  - [ ] GetEmploymentByID
  - [ ] CreateEmployment
  - [ ] UpdateEmployment
  - [ ] DeleteEmployment
- [ ] 2.5 Create `internal/db/queries/companies.sql` (already has GetCompany, add remaining):
  - [ ] ListCompanies
  - [ ] CreateCompany
  - [ ] UpdateCompany
  - [ ] DeleteCompany
- [ ] 2.6 Create `internal/db/queries/posts_tags.sql` for many-to-many relationships:
  - [ ] LinkPostToTag
  - [ ] UnlinkPostFromTag
  - [ ] GetTagsForPost
  - [ ] GetPostsForTag
  - [ ] UnlinkAllTagsFromPost (for tag replacement during updates)
- [ ] 2.7 Create `internal/db/queries/employments_companies.sql` for employment-company relationships:
  - [ ] LinkEmploymentToCompany
  - [ ] UnlinkEmploymentFromCompany
  - [ ] GetCompaniesForEmployment
  - [ ] GetEmploymentsForCompany
- [ ] 2.8 Run `sqlc generate` and verify generated code in `internal/db/*.sql.go`
- [ ] 2.9 Review generated types in `internal/db/models.go` for correctness

## 3. PostgreSQL Connection Management

- [ ] 3.1 Remove Bun ORM initialization from `cmd/conneroh/root.go`
- [ ] 3.2 Implement PostgreSQL connection pool setup using pgx/v5:
  - [ ] Load connection config from environment variables (DB_CLIENT, DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_DATABASE)
  - [ ] Create pgxpool connection with appropriate timeouts and pool settings
  - [ ] Implement connection validation (ping on startup)
  - [ ] Implement graceful shutdown with connection cleanup
- [ ] 3.3 Pass pgx connection pool to route handlers and middleware in `cmd/conneroh/routes.go`
- [ ] 3.4 Create `internal/db` helper to initialize sqlc Queries struct with connection pool

## 4. Update Route Handlers to Use sqlc

- [ ] 4.1 Update `internal/routing/handlers.go` to use sqlc-generated query methods
- [ ] 4.2 Replace Bun model references with sqlc-generated model types
- [ ] 4.3 Update pagination logic in `internal/routing/pagination.go` to work with sqlc queries
- [ ] 4.4 Ensure transaction support for operations that need atomicity (e.g., creating post with tags)
- [ ] 4.5 Update error handling to work with pgx errors (e.g., pgx.ErrNoRows for 404 responses)

## 5. Remove Bun ORM and SQLite Dependencies

- [ ] 5.1 Delete Bun model definitions in `internal/assets/static.go` (Post, Project, Tag, Employment, Company, Cache models)
- [ ] 5.2 Remove Bun ORM import statements throughout codebase
- [ ] 5.3 Remove SQLite driver imports (modernc.org/sqlite)
- [ ] 5.4 Update `go.mod`:
  - [ ] Remove `github.com/uptrace/bun` dependency
  - [ ] Remove `github.com/uptrace/bun/dialect/sqlitedialect` dependency
  - [ ] Remove `modernc.org/sqlite` dependency
  - [ ] Verify `github.com/jackc/pgx/v5` and `github.com/cridenour/go-postgis` are present
- [ ] 5.5 Run `go mod tidy` to clean up unused dependencies
- [ ] 5.6 Delete `master.db` file (add to .gitignore if not already present)

## 6. Create Content Pipeline Tool

- [ ] 6.1 Create new entry point `cmd/sync-content/main.go` for markdown-to-database pipeline
- [ ] 6.2 Implement markdown file discovery and parsing:
  - [ ] Recursively find .md files in `internal/data/docs/{posts,projects,tags}/`
  - [ ] Parse YAML frontmatter using goldmark-frontmatter
  - [ ] Extract markdown body content
- [ ] 6.3 Implement content type detection based on directory structure
- [ ] 6.4 Implement validation logic:
  - [ ] Validate required frontmatter fields per content type
  - [ ] Validate slug uniqueness
  - [ ] Validate timestamp formats (ISO 8601 with timezone)
- [ ] 6.5 Implement database operations using sqlc queries:
  - [ ] Insert or update posts, projects, tags
  - [ ] Handle tag auto-creation from references
  - [ ] Establish junction table relationships (posts_tags)
- [ ] 6.6 Implement transaction support for atomic operations
- [ ] 6.7 Implement error reporting and validation summary
- [ ] 6.8 Add CLI flags:
  - [ ] `--full-rebuild` (truncate and regenerate) vs incremental update
  - [ ] `--dry-run` (validate without database changes)
  - [ ] `--content-dir` (override default content directory)
- [ ] 6.9 Test content pipeline with existing markdown files

## 7. Update Nix Configuration

- [ ] 7.1 Update `flake.nix`:
  - [ ] Remove broken `generate-db` references to deleted `cmd/update`
  - [ ] Add new `sync-content` command that runs `cmd/sync-content`
  - [ ] Update `generate-all` to remove `generate-db` reference or replace with `sync-content`
  - [ ] Ensure PostgreSQL environment variables are sourced (or passed via doppler/infisical)
- [ ] 7.2 Test all Nix shell commands:
  - [ ] `nix develop` (shell startup)
  - [ ] `run` (application with hot reloading)
  - [ ] `sync-content` (new content pipeline command)
  - [ ] `tests` (ensure tests pass with new DB setup)
  - [ ] `lint` and `format`

## 8. Update Documentation

- [ ] 8.1 Update `openspec/project.md`:
  - [ ] Change "SQLite" to "PostgreSQL" in Tech Stack section
  - [ ] Change "Bun ORM" to "sqlc with pgx/v5" in Backend section
  - [ ] Remove references to `cmd/update` command
  - [ ] Add `cmd/sync-content` documentation
  - [ ] Update database-related commands in Development Workflow section
  - [ ] Update "Important Constraints" to reflect PostgreSQL requirement
- [ ] 8.2 Update README (if exists) with PostgreSQL setup instructions
- [ ] 8.3 Document environment variable requirements:
  - [ ] DB_CLIENT (postgresql)
  - [ ] DB_USER
  - [ ] DB_PASSWORD
  - [ ] DB_HOST
  - [ ] DB_PORT
  - [ ] DB_DATABASE
- [ ] 8.4 Document data migration process (manual, from SQLite to PostgreSQL)
- [ ] 8.5 Update `CLAUDE.md` or other project instructions to reference PostgreSQL and sqlc

## 9. Testing and Validation

- [ ] 9.1 Write unit tests for content pipeline:
  - [ ] Test markdown parsing with various frontmatter scenarios
  - [ ] Test validation logic (missing fields, invalid timestamps, duplicate slugs)
  - [ ] Test error reporting
- [ ] 9.2 Write integration tests for database operations:
  - [ ] Test all sqlc CRUD queries
  - [ ] Test transaction rollback on errors
  - [ ] Test tag auto-creation and relationship management
- [ ] 9.3 Update existing tests in `internal/routing/pagination_test.go` to work with sqlc
- [ ] 9.4 Run full test suite and ensure all tests pass
- [ ] 9.5 Test Playwright browser tests (in `tests/browser/`) to ensure UI still works with PostgreSQL backend

## 10. Data Migration (Manual)

- [ ] 10.1 Document manual migration steps for existing SQLite data:
  - [ ] Export existing content from SQLite (or rely on markdown source of truth)
  - [ ] Run `sync-content` to populate PostgreSQL from markdown files
  - [ ] Verify data integrity (compare record counts, spot-check content)
- [ ] 10.2 Test migration on a staging environment before production
- [ ] 10.3 Backup existing SQLite database before deleting

## 11. Deployment Preparation

- [ ] 11.1 Update deployment configuration (Fly.io) to include PostgreSQL connection credentials
- [ ] 11.2 Ensure PostgreSQL database is provisioned and accessible from application
- [ ] 11.3 Run database migrations (`internal/db/migrations/001_init.sql`) on production PostgreSQL
- [ ] 11.4 Test application in staging environment with PostgreSQL
- [ ] 11.5 Verify all routes and functionality work correctly
- [ ] 11.6 Monitor logs for database connection issues or errors

## 12. Final Verification

- [ ] 12.1 Code review: Ensure all Bun ORM references are removed
- [ ] 12.2 Code review: Ensure all SQLite references are removed
- [ ] 12.3 Verify no `master.db` file is present or referenced
- [ ] 12.4 Run linting (`golangci-lint`) and fix any issues
- [ ] 12.5 Run formatting (`go fmt`, `nix fmt`) and ensure code is clean
- [ ] 12.6 Verify all OpenSpec proposal acceptance criteria are met
- [ ] 12.7 Update proposal status and archive after deployment
