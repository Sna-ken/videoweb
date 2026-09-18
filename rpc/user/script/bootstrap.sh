#!/usr/bin/env bash

CURDIR=$(cd "$(dirname "$0")" && pwd)
RUNTIME_ROOT=${1:-$CURDIR}

export KITEX_RUNTIME_ROOT="$RUNTIME_ROOT"
export KITEX_LOG_DIR="$RUNTIME_ROOT/log"

mkdir -p "$KITEX_LOG_DIR/app" "$KITEX_LOG_DIR/rpc"
exec "$CURDIR/bin/user"