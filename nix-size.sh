
#!/usr/bin/env bash
set -euo pipefail

# Optional first arg: system (e.g. x86_64-linux, aarch64-linux)
SYSTEM="${1:-x86_64-linux}"

# Make a temp dir for all our no-link builds
TMPDIR=$(mktemp -d)
cleanup() { rm -rf "$TMPDIR"; }
trap cleanup EXIT

echo "▶ Checking flake output sizes for system: $SYSTEM"
echo

# Grab all package keys under .packages.$SYSTEM
ATTRS=$(nix flake show --json . \
  | jq -r ".packages.\"$SYSTEM\" | keys[]")

for attr in $ATTRS; do
  echo "── $attr ──"

  # Build it into $TMPDIR/$attr (no symlink)
  nix build --no-link ".#packages.$SYSTEM.$attr" -o "$TMPDIR/$attr"

  # Query NAR (archive) size and closure size in bytes
  nix path-info --json -shS "$TMPDIR/$attr" \
    | jq -r '.[] | "archive: \(.narSize) bytes, closure: \(.closureSize) bytes"'

  echo
done

echo "▶ Checked flake output sizes for system: $SYSTEM"
