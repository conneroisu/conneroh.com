# Implementation Tasks

## 1. Modify flake.nix
- [ ] 1.1 Add inline `bunNix` derivation using `pkgs.runCommand` in the `packages` section before the `default` package
- [ ] 1.2 Ensure the `bunNix` derivation references `bun2nix.packages.${system}.default` as a build input
- [ ] 1.3 Verify the command generates output to `$out` using `--lock-file ${./bun.lock} --output-file $out`
- [ ] 1.4 Confirm the generated `bunNix` is passed to the Bun package configuration as `inherit bunNix;`

## 2. Clean up package.json
- [ ] 2.1 Remove the `postinstall` script entry from `package.json`
- [ ] 2.2 Verify no other scripts depend on `bun.nix` being pre-generated

## 3. Update repository files
- [ ] 3.1 Delete the committed `bun.nix` file from the repository
- [ ] 3.2 Add `bun.nix` to `.gitignore` to prevent accidental future commits
- [ ] 3.3 Stage all changes for commit

## 4. Validation
- [ ] 4.1 Run `nix flake check` to verify the flake builds correctly
- [ ] 4.2 Build the default package with `nix build .#default` to confirm it works
- [ ] 4.3 Verify that `bun.nix` is generated on-the-fly during the build (check build output)
- [ ] 4.4 Run existing tests with `nix develop -c tests` to ensure functionality is unchanged
- [ ] 4.5 Confirm the built package still functions correctly

## 5. Documentation
- [ ] 5.1 Update this change proposal with any implementation notes or learnings
- [ ] 5.2 Verify all task checkboxes are marked complete before archiving
