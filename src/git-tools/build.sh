#!/usr/bin/env bash
set -euo pipefail

# Simple wrapper for build.go
cd "$(dirname "${BASH_SOURCE[0]}")"

# Check if builder binary exists and is newer than source
BUILDER_BIN=".builder"
if [[ ! -f "$BUILDER_BIN" ]] || [[ "build.go" -nt "$BUILDER_BIN" ]]; then
    echo "Building builder..."
    go build -o "$BUILDER_BIN" build.go
fi

# Run the compiled builder
./"$BUILDER_BIN" "$@"
