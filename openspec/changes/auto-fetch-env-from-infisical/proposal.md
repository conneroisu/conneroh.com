# Auto-fetch .env from Infisical Proposal

## Why

Currently, the Nix development shell's `shellHook` immediately sources `.env` without checking if it exists. This creates friction for new contributors or when switching between environments:

1. **Silent failures**: If `.env` is missing, `source .env` fails silently and variables aren't loaded
2. **Manual secret fetching**: Developers must remember to manually fetch secrets from Infisical before entering the shell
3. **Inconsistent state**: New checkouts or fresh environments require manual intervention to become functional
4. **Poor DX**: The shell enters successfully but the application fails at runtime due to missing environment variables

By automatically fetching secrets from Infisical when `.env` is missing, we provide a smoother onboarding experience and reduce environment setup friction.

## What Changes

- **Modify** `shellHook` in `flake.nix` to check if `.env` exists before sourcing
- **Add** automatic `infisical export > .env` when `.env` is missing
- **Add** warning message if Infisical fetch fails (but continue shell startup)
- **Ensure** existing `.env` files are never overwritten (manual removal required to re-fetch)

The enhancement will follow this pattern:

```nix
shellHook = ''
  if [ ! -f .env ]; then
    echo "⚠️  .env not found, fetching from Infisical..."
    if infisical export --env prod > .env 2>/dev/null; then
      echo "✓ Successfully fetched secrets from Infisical"
    else
      echo "⚠️  Failed to fetch secrets from Infisical. Please run: infisical export --env prod > .env"
    fi
  fi

  set -a
  source .env
  set +a
'';
```

## Impact

**Affected specs:**
- `nix-development-shell` (MODIFIED - if exists)

**Affected code:**
- `flake.nix` - Modify `shellHook` to add conditional `.env` fetching logic

**Breaking changes:** None

**Benefits:**
- Automatic secret provisioning for new contributors and fresh checkouts
- Graceful handling of missing `.env` with helpful error messages
- No impact on existing workflows (existing `.env` files are preserved)
- Clear feedback when Infisical authentication or fetching fails
- Reduced onboarding friction and documentation burden

**Risks:**
- Minimal: If Infisical is not authenticated or unavailable, the shell still starts with a warning
- Existing `.env` files are never modified, preventing accidental overwrites
- The `infisical` CLI tool is already available in the Nix shell environment (line 184 of `flake.nix`)
