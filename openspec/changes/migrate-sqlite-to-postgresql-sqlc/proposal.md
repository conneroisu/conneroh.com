# Migrate from SQLite to PostgreSQL with sqlc

## Why

The codebase currently maintains a hybrid database architecture with SQLite (via Bun ORM) for the application and PostgreSQL (via sqlc) for infrastructure. This creates unnecessary complexity, duplicate model definitions, and maintenance overhead. The PostgreSQL infrastructure is already in place with a complete schema, but sqlc is minimally utilized (only 1 query implemented). Consolidating to a single database system will simplify the architecture, improve type safety through sqlc's code generation, and leverage PostgreSQL's advanced features (PostGIS, JSONB, better concurrency).

## What Changes

- **Remove SQLite database**: Delete `master.db` and all SQLite-specific code
- **Remove Bun ORM**: Delete Bun ORM dependencies, model definitions in `internal/assets/static.go`, and initialization code in `cmd/conneroh/root.go`
- **Implement complete sqlc queries**: Create full CRUD operations for all entities (posts, projects, tags, employments, companies) in `internal/db/queries/`
- **Update connection management**: Replace Bun DB initialization with pgx connection pool for PostgreSQL
- **Keep Directus schema**: Maintain all Directus CMS tables in the PostgreSQL schema for future use
- **Update documentation**: Remove references to SQLite, Bun ORM, and the deleted `cmd/update` command throughout `project.md` and other docs
- **Fix Nix configuration**: Update `flake.nix` to remove broken `generate-db` references to the deleted `cmd/update` command
- **Document data migration**: Provide manual migration guidance (no automated tooling)

**BREAKING CHANGES**:
- Database backend changes from SQLite to PostgreSQL
- All database access patterns change from Bun ORM to sqlc
- Application requires PostgreSQL connection configuration (environment variables)
- Development environment requires PostgreSQL instance

## Impact

### Affected Specs
- **database-access** (NEW) - Defines PostgreSQL connection management, sqlc usage patterns, and query organization
- **content-pipeline** (NEW) - Documents how markdown content flows into the database (replaces missing `cmd/update` functionality)
- **development-environment** (MODIFIED) - Requires PostgreSQL setup in addition to Infisical

### Affected Code
- `cmd/conneroh/root.go` - Database initialization and connection management
- `cmd/conneroh/routes.go` - Handler dependencies (DB instance passing)
- `internal/assets/` - Remove Bun models, update to use sqlc-generated types
- `internal/routing/handlers.go` - Update to use sqlc query methods
- `internal/db/queries/*.sql` - Add comprehensive query definitions (currently only `companies.sql` exists)
- `flake.nix` - Fix broken `generate-db` command references
- `go.mod` - Remove Bun ORM, modernc.org/sqlite dependencies
- `openspec/project.md` - Update tech stack and architecture sections
- `sqlc.yaml` - Already correctly configured for PostgreSQL

### Migration Considerations
- Existing SQLite data must be manually migrated to PostgreSQL
- Content can be regenerated from markdown source files
- PostgreSQL instance must be configured before running the application
- Nix development environment will include PostgreSQL setup guidance
