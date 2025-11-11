# development-environment Specification Delta

## ADDED Requirements

### Requirement: PostgreSQL Database Configuration (REQ-DEV-002)

The Nix development shell SHALL provide guidance for configuring PostgreSQL database connection credentials required by the application.

**Priority**: High

**Rationale**: The application now depends on PostgreSQL instead of SQLite. Developers need clear instructions for database setup.

#### Scenario: PostgreSQL environment variables documented

- **GIVEN** a developer enters the Nix development shell with `nix develop`
- **WHEN** they need to configure database access
- **THEN** documentation or shell messages explain the required environment variables:
  - DB_CLIENT (set to "postgresql")
  - DB_USER (PostgreSQL username)
  - DB_PASSWORD (PostgreSQL password)
  - DB_HOST (PostgreSQL host, e.g., localhost or remote host)
  - DB_PORT (PostgreSQL port, typically 5432)
  - DB_DATABASE (database name)
- **AND** guidance is provided on sourcing these from Infisical (if integrated) or creating a `.env` file

#### Scenario: Missing PostgreSQL configuration causes clear error

- **GIVEN** required database environment variables are not set
- **WHEN** a developer runs the application in the Nix shell
- **THEN** the application fails to start with a clear error message
- **AND** the error lists the missing environment variables
- **AND** guidance is provided on how to configure them

#### Scenario: Local PostgreSQL setup is documented

- **GIVEN** a developer wants to run PostgreSQL locally for development
- **WHEN** they consult project documentation (README, CLAUDE.md, or shell startup message)
- **THEN** instructions are provided for:
  - Installing PostgreSQL (via Nix or system package manager)
  - Creating a development database
  - Running schema migrations (`internal/db/migrations/001_init.sql`)
  - Populating content via `sync-content` command
- **AND** instructions mention remote PostgreSQL options (Neon, Fly Postgres, etc.) as alternatives
