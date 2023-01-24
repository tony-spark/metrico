#!/bin/bash

build_version=$(git describe --exact-match --tags 2> /dev/null || git rev-parse --abbrev-ref HEAD)
build_date=$(TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')
build_commit=$(git rev-parse --short HEAD)

echo "$build_version $build_date $build_commit"
go build -ldflags "-X main.buildVersion=$build_version -X main.buildDate=$build_date -X main.buildCommit=$build_commit" "$@"
