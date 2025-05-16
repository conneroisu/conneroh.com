#!/usr/bin/env bash
set -euo pipefail

# 1) Load all of the exports that nix-shell generated:
. /env-vars

# 2) (Optional) print a message so you know it ran:
echo "✅ Nix env loaded from /env-vars"

# 3) Keep the container alive by handing off to sleep
exec sleep infinity
