# Inline Bun.nix Generation Proposal

## Why

Currently, the project maintains a separate `bun.nix` file (33KB) that is generated via a `postinstall` script in `package.json` using the `bun2nix` CLI tool. This approach has several drawbacks:

1. **Build-time dependency**: The `bun.nix` file must be pre-generated and committed to the repository
2. **Manual regeneration**: Developers must remember to run `bun install` or manually invoke `bun2nix` whenever `bun.lock` changes
3. **Repository bloat**: A large generated file (33KB) is committed to version control
4. **Inconsistency risk**: The `bun.nix` file can drift out of sync with `bun.lock` if not regenerated

Following the pattern demonstrated in the Connix project, we can inline the `bun.nix` generation directly within `flake.nix` using `pkgs.runCommand`, eliminating the need for the separate file and the `postinstall` script.

## What Changes

- **Remove** `bun.nix` file from repository
- **Remove** `postinstall` script from `package.json`
- **Add** inline `bun.nix` generation within `flake.nix` using `pkgs.runCommand` in the packages section
- **Update** `.gitignore` to exclude `bun.nix` if accidentally generated locally

The inline generation pattern will follow the Connix approach:

```nix
bunNix = pkgs.runCommand "bun.nix" {
  buildInputs = [bun2nix.packages.${system}.default];
} ''
  bun2nix --lock-file ${./bun.lock} --output-file $out
'';
```

This `bunNix` derivation is then passed to the Bun package build configuration.

## Impact

**Affected specs:**
- `nix-build-system` (MODIFIED)

**Affected code:**
- `flake.nix` - Add inline `bunNix` derivation generation
- `package.json` - Remove `postinstall` script
- `bun.nix` - Delete file (no longer needed)
- `.gitignore` - Add `bun.nix` to prevent accidental commits

**Breaking changes:** None

**Benefits:**
- Automatic synchronization between `bun.lock` and Bun package dependencies
- Cleaner repository (one less generated file)
- More declarative build process
- Follows established patterns from similar projects (Connix)
- Reduces cognitive load for contributors

**Risks:**
- Minimal: The change is purely structural and doesn't affect runtime behavior
- If `bun2nix` tool changes its CLI interface, we'll need to update the inline command (same risk as before)
