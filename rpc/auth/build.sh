#!/usr/bin/env bash
set -e

PROJECT_ROOT=$(cd "$(dirname "$0")/../.." && pwd)
OUTPUT_DIR="$PROJECT_ROOT/rpc/auth/output"

mkdir -p "$OUTPUT_DIR/bin"
cp "$PROJECT_ROOT/rpc/auth/script/"* "$OUTPUT_DIR" 2>/dev/null || true
chmod +x "$OUTPUT_DIR/bootstrap.sh"
go build -o "$OUTPUT_DIR/bin/auth" "$PROJECT_ROOT/cmd/auth"