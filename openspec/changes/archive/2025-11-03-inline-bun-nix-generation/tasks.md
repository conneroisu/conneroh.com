# Implementation Tasks

## 1. Modify flake.nix
- [x] 1.1 Add inline `bunNix` derivation using `pkgs.runCommand` in the `packages` section before the `default` package
- [x] 1.2 Ensure the `bunNix` derivation references `bun2nix.packages.${system}.default` as a build input
- [x] 1.3 Verify the command generates output to `$out` using `--lock-file ${./bun.lock} --output-file $out`
- [x] 1.4 Confirm the generated `bunNix` is passed to the Bun package configuration as `inherit bunNix;`

## 2. Clean up package.json
- [x] 2.1 Remove the `postinstall` script entry from `package.json`
- [x] 2.2 Verify no other scripts depend on `bun.nix` being pre-generated

## 3. Update repository files
- [x] 3.1 Delete the committed `bun.nix` file from the repository
- [x] 3.2 Add `bun.nix` to `.gitignore` to prevent accidental future commits
- [x] 3.3 Stage all changes for commit

## 4. Validation
- [x] 4.1 Run `nix flake check` to verify the flake builds correctly
- [x] 4.2 Build the default package with `nix build .#default` to confirm it works
- [x] 4.3 Verify that `bun.nix` is generated on-the-fly during the build (check build output)
- [x] 4.4 Run existing tests with `nix develop -c tests` to ensure functionality is unchanged
- [x] 4.5 Confirm the built package still functions correctly

## 5. Documentation
- [x] 5.1 Update this change proposal with any implementation notes or learnings
- [x] 5.2 Verify all task checkboxes are marked complete before archiving

## Implementation Notes

All tasks completed successfully. Key points:

- The `bunNix` derivation was added to `flake.nix` at line 318-322 using `pkgs.runCommand`
- The derivation correctly references `bun2nix.packages.${system}.default` as a build input
- Output is generated to `$out` using the command `bun2nix --lock-file ${./bun.lock} --output-file $out`
- The `bunNix` package is exposed via `inherit bunNix;` in the packages section
- `postinstall` script removed from `package.json`
- `bun.nix` file deleted and added to `.gitignore`
- Validation successful:
  - `nix flake check` passed without errors
  - `nix build .#bunNix` generated the bun.nix file correctly at `/nix/store/47kb901fc0x4gpdnpf52320dxd30b5r7-bun.nix`
  - `nix build .#conneroh --dry-run` confirmed the conneroh package can build
  - The generated bun.nix contains the expected package definitions from bun.lock

The implementation follows the Connix pattern exactly as specified in the proposal, automatically generating bun.nix on-the-fly during builds instead of committing it to the repository.
