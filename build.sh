#!/bin/sh
# Cross-compiles static binaries for all supported CPUs into ./dist
# Usage: ./build.sh [version]     (needs Go 1.23+, or run inside golang:1.23-alpine)
set -e
VERSION="${1:-dev}"
mkdir -p dist
build() { # name GOARCH [extra env]
  name="$1"; arch="$2"; shift 2
  echo "building lfe-linux-$name"
  env CGO_ENABLED=0 GOOS=linux GOARCH="$arch" "$@" \
    go build -trimpath -ldflags "-s -w -X main.version=$VERSION" -o "dist/lfe-linux-$name" .
}
build amd64  amd64
build arm64  arm64
build armv7  arm    GOARM=7
build armv5  arm    GOARM=5
build mipsle mipsle GOMIPS=softfloat
build mips   mips   GOMIPS=softfloat
ls -l dist
