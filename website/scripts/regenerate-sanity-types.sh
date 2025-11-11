#!/usr/bin/env bash
set -euo pipefail

# Script to regenerate Sanity TypeScript types
# Usage: ./scripts/regenerate-sanity-types.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
SANITY_DIR="$PROJECT_ROOT/sanity"
SRC_DIR="$PROJECT_ROOT/src/lib"

echo "🔄 Regenerating Sanity TypeScript types..."

# Navigate to sanity directory
cd "$SANITY_DIR"

# Generate types
echo "📝 Running Sanity TypeGen..."
bun sanity typegen generate

# Copy to src/lib
echo "📦 Copying types to src/lib..."
cp "$SANITY_DIR/sanity.types.ts" "$SRC_DIR/sanity-types.ts"

echo "✅ Done! Types have been regenerated."
echo ""
echo "Generated file: $SRC_DIR/sanity-types.ts"
echo ""
echo "Next steps:"
echo "1. Review the generated types"
echo "2. Restart your TypeScript server (in your IDE)"
echo "3. Run 'bun run build' to verify type safety"
