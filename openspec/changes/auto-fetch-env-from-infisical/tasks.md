# Implementation Tasks

## 1. Modify flake.nix shellHook
- [ ] 1.1 Add conditional check `[ ! -f .env ]` before attempting to fetch secrets
- [ ] 1.2 Add Infisical export command: `infisical export --env prod > .env 2>/dev/null`
- [ ] 1.3 Add success message when secrets are fetched: `echo "✓ Successfully fetched secrets from Infisical"`
- [ ] 1.4 Add failure warning with helpful instructions if Infisical fetch fails
- [ ] 1.5 Preserve existing `set -a; source .env; set +a` logic after the conditional block
- [ ] 1.6 Ensure proper error handling prevents shell startup failure

## 2. Testing
- [ ] 2.1 Verify shell starts successfully when `.env` already exists (no changes)
- [ ] 2.2 Test shell startup with missing `.env` and valid Infisical authentication
- [ ] 2.3 Test shell startup with missing `.env` and invalid/missing Infisical authentication
- [ ] 2.4 Confirm existing `.env` files are never overwritten
- [ ] 2.5 Verify environment variables are correctly loaded after auto-fetch
- [ ] 2.6 Test that the shell continues to start even when Infisical fetch fails

## 3. Documentation
- [ ] 3.1 Update implementation notes with any learnings or edge cases discovered
- [ ] 3.2 Consider updating project documentation if onboarding instructions change
- [ ] 3.3 Verify all task checkboxes are marked complete before marking change as done

## Implementation Notes

(To be added after implementation)
