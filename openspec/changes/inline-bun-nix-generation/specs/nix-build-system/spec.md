# Nix Build System - Spec Delta

## MODIFIED Requirements

### Requirement: Bun Package Dependencies Resolution

The Nix build system SHALL generate the `bun.nix` dependency mapping inline during the flake evaluation phase rather than relying on a pre-committed file.

**Rationale**: Inline generation ensures automatic synchronization with `bun.lock` and reduces repository maintenance overhead.

#### Scenario: Successful inline generation during build

- **GIVEN** a `bun.lock` file exists in the repository
- **WHEN** the Nix flake is evaluated (e.g., via `nix build`, `nix develop`, or `nix flake check`)
- **THEN** the system automatically generates `bun.nix` using `bun2nix --lock-file ${./bun.lock} --output-file $out`
- **AND** the generated derivation is available for use by the Bun package builder
- **AND** no pre-committed `bun.nix` file is required in the repository

#### Scenario: Build fails when bun.lock is missing

- **GIVEN** the `bun.lock` file does not exist
- **WHEN** the Nix flake is evaluated
- **THEN** the build fails with a clear error indicating the missing lockfile
- **AND** provides guidance to run `bun install` to generate the lockfile

#### Scenario: Stale bun.nix is ignored

- **GIVEN** a developer accidentally generates a local `bun.nix` file
- **WHEN** Git operations are performed
- **THEN** the file is ignored via `.gitignore`
- **AND** it does not get committed to the repository

### Requirement: Reproducible Bun Builds

The Nix build system SHALL maintain reproducible builds where the same `bun.lock` input always produces the same build output, regardless of whether `bun.nix` is pre-generated or generated inline.

**Rationale**: Build reproducibility is a core Nix principle that must be maintained through this refactoring.

#### Scenario: Consistent builds across environments

- **GIVEN** the same `bun.lock` file contents
- **WHEN** builds are performed on different machines or at different times
- **THEN** the generated `bun.nix` content is identical
- **AND** the final package hash is identical
- **AND** all Bun dependencies resolve to the same versions

#### Scenario: Lockfile changes trigger rebuild

- **GIVEN** the `bun.lock` file is modified (e.g., dependency update)
- **WHEN** a subsequent Nix build is performed
- **THEN** the `bun.nix` derivation is regenerated
- **AND** the new dependency set is used
- **AND** the package is rebuilt with updated dependencies

## REMOVED Requirements

### Requirement: Pre-committed bun.nix file

**Removed**: The requirement to maintain a committed `bun.nix` file in the repository

**Reason**: Inline generation eliminates the need for version controlling generated files, reducing repository size and maintenance burden.

**Migration**:
- Delete `bun.nix` from repository
- Add to `.gitignore`
- No code changes required for consumers (the file was only used internally by Nix)

### Requirement: Manual bun.nix regeneration

**Removed**: The requirement for developers to manually run `bun2nix` via `postinstall` scripts or CLI commands

**Reason**: Inline generation happens automatically during Nix evaluation, eliminating manual steps.

**Migration**:
- Remove `postinstall` script from `package.json`
- No action required from developers
- Build process remains: `nix build`, `nix develop`, etc.
