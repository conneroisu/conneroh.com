# database-access Specification Delta

## ADDED Requirements

### Requirement: PostgreSQL Connection Management (REQ-DB-001)

The system SHALL establish and manage connections to PostgreSQL using the pgx/v5 connection pool with configuration loaded from environment variables.

**Priority**: Critical

#### Scenario: Successful connection with environment variables

- **GIVEN** the following environment variables are set: `DB_CLIENT=postgresql`, `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_DATABASE`
- **WHEN** the application initializes the database connection
- **THEN** a pgx connection pool is created with the provided credentials
- **AND** the connection is validated with a ping operation
- **AND** the connection pool is available for query execution
- **AND** connection lifecycle is managed (proper cleanup on shutdown)

#### Scenario: Connection failure with missing credentials

- **GIVEN** required database environment variables are not set
- **WHEN** the application attempts to initialize the database connection
- **THEN** initialization fails with a clear error message indicating missing configuration
- **AND** the error message lists the required environment variables
- **AND** the application exits gracefully without starting the web server

#### Scenario: Connection failure to unavailable database

- **GIVEN** database environment variables are set with invalid host/port
- **WHEN** the application attempts to connect to PostgreSQL
- **THEN** the connection attempt fails after timeout
- **AND** a clear error message indicates connection failure
- **AND** the application exits gracefully

### Requirement: sqlc Code Generation (REQ-DB-002)

The system SHALL use sqlc to generate type-safe Go code from SQL queries, ensuring compile-time validation and type safety for all database operations.

**Priority**: Critical

#### Scenario: Generate code from SQL queries

- **GIVEN** SQL query files exist in `internal/db/queries/*.sql`
- **WHEN** `sqlc generate` is executed
- **THEN** Go code is generated in `internal/db/` package
- **AND** each query produces a method on the `*Queries` struct
- **AND** query parameters are type-safe Go function arguments
- **AND** query results are type-safe Go structs
- **AND** nullable columns use pgx/v5 nullable types (pgtype.Text, pgtype.UUID, etc.)

#### Scenario: Invalid SQL fails generation

- **GIVEN** a SQL query file contains syntax errors or references non-existent tables
- **WHEN** `sqlc generate` is executed
- **THEN** sqlc reports validation errors
- **AND** generation fails without producing invalid code
- **AND** error messages reference the specific file and line with the issue

### Requirement: Complete CRUD Operations for Core Entities (REQ-DB-003)

The system SHALL provide complete Create, Read, Update, and Delete operations for all core content entities: posts, projects, tags, employments, and companies.

**Priority**: High

#### Scenario: Query all posts with pagination

- **GIVEN** multiple posts exist in the database
- **WHEN** the application queries posts with limit and offset parameters
- **THEN** sqlc-generated `ListPosts` method returns paginated results
- **AND** results include all post fields (id, title, slug, description, content, metadata, timestamps)
- **AND** results are ordered by created_at descending
- **AND** related tags are accessible via separate tag junction table queries

#### Scenario: Create new post with full metadata

- **GIVEN** a post with title, slug, description, markdown content, and metadata
- **WHEN** the application calls sqlc-generated `CreatePost` method
- **THEN** a new row is inserted into the posts table
- **AND** the generated UUID is returned
- **AND** timestamps are set to current time
- **AND** JSONB metadata is stored correctly

#### Scenario: Get single post by slug

- **GIVEN** a post exists with a specific slug
- **WHEN** the application queries by slug using `GetPostBySlug`
- **THEN** the matching post row is returned
- **AND** all fields are populated with correct types
- **AND** if no post matches, pgx.ErrNoRows is returned

#### Scenario: Update post content

- **GIVEN** an existing post identified by ID
- **WHEN** the application calls `UpdatePost` with modified fields
- **THEN** the post row is updated with new values
- **AND** the updated_at timestamp is refreshed
- **AND** the row count indicates successful update

#### Scenario: Delete post by ID

- **GIVEN** an existing post
- **WHEN** the application calls `DeletePost` with the post ID
- **THEN** the post row is removed from the database
- **AND** related junction table entries (posts_tags) are cascade deleted
- **AND** the deletion is confirmed via row count

#### Scenario: Query projects with associated companies

- **GIVEN** projects with employment relationships to companies
- **WHEN** the application queries projects
- **THEN** projects are returned with all fields
- **AND** related companies are queryable via employment foreign keys
- **AND** PostGIS location data is accessible through go-postgis types

#### Scenario: List tags with usage counts

- **GIVEN** tags associated with multiple posts and projects
- **WHEN** the application queries tags
- **THEN** tags are returned with names and slugs
- **AND** usage counts can be calculated via junction table queries
- **AND** tags are orderable by usage frequency

### Requirement: Transaction Support (REQ-DB-004)

The system SHALL support database transactions for operations that require atomicity across multiple queries.

**Priority**: High

#### Scenario: Create post with tags atomically

- **GIVEN** a new post and a list of tag IDs
- **WHEN** the application creates a post within a transaction
- **THEN** the post is inserted
- **AND** junction table entries (posts_tags) are inserted for each tag
- **AND** if any operation fails, all changes are rolled back
- **AND** if all succeed, the transaction is committed

#### Scenario: Transaction rollback on error

- **GIVEN** a transaction with multiple database operations
- **WHEN** one operation fails (e.g., unique constraint violation)
- **THEN** the transaction is rolled back
- **AND** no partial data is committed
- **AND** the database state remains unchanged
- **AND** the error is propagated to the caller

### Requirement: Query Organization and Naming (REQ-DB-005)

SQL queries SHALL be organized in separate files by entity domain, following a consistent naming convention that clearly indicates the operation and target.

**Priority**: Medium

#### Scenario: Query files organized by entity

- **GIVEN** the codebase contains database queries
- **WHEN** a developer navigates `internal/db/queries/`
- **THEN** each entity has a dedicated SQL file (posts.sql, projects.sql, tags.sql, employments.sql, companies.sql)
- **AND** junction table queries are in relationship-specific files (posts_tags.sql, employments_companies.sql)
- **AND** query names follow the pattern: `{Action}{Entity}[By{Field}]` (e.g., GetPostBySlug, ListPosts, CreatePost)

#### Scenario: sqlc generates predictable method names

- **GIVEN** a SQL query with name `-- name: ListPostsByTag :many`
- **WHEN** sqlc generates Go code
- **THEN** a method `ListPostsByTag(ctx, tagID)` is created
- **AND** the method signature matches the query parameters and return type
- **AND** naming is consistent across all entities

### Requirement: PostGIS Geometry Support (REQ-DB-006)

The system SHALL support PostGIS geometry types for location data, using the go-postgis library for type mapping.

**Priority**: Medium

#### Scenario: Query employment with location

- **GIVEN** an employment record with a PostGIS geometry location
- **WHEN** the application queries the employment
- **THEN** the location field is mapped to `postgis.Point` type
- **AND** latitude and longitude are accessible from the Point struct
- **AND** SRID 4326 (WGS84) coordinate system is maintained

#### Scenario: Insert employment with geometry

- **GIVEN** an employment with latitude/longitude coordinates
- **WHEN** the application inserts the record with a PostGIS point
- **THEN** the geometry is stored as `geometry(Geometry,4326)` in PostgreSQL
- **AND** the data is queryable with PostGIS spatial functions
- **AND** the go-postgis library handles serialization correctly

### Requirement: Directus CMS Schema Preservation (REQ-DB-007)

The PostgreSQL schema SHALL include all Directus CMS tables, maintained for future integration or reference.

**Priority**: Low

#### Scenario: Directus tables remain intact

- **GIVEN** the PostgreSQL schema contains Directus tables (directus_users, directus_files, etc.)
- **WHEN** the application connects to the database
- **THEN** Directus tables are present but not actively used by application code
- **AND** migrations do not drop or modify Directus tables
- **AND** schema documentation acknowledges Directus presence

#### Scenario: Application queries ignore Directus tables

- **GIVEN** sqlc query files in the codebase
- **WHEN** queries are executed
- **THEN** only application tables (posts, projects, tags, employments, companies) are accessed
- **AND** no queries reference Directus tables
- **AND** Directus tables do not interfere with application functionality
