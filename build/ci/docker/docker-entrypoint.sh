#!/bin/bash

set -e

# Statically checking Go source for errors and warnings.
gometalinter.v2 --disable-all -E vet -E gofmt -E misspell -E ineffassign -E goimports -E deadcode -E gocyclo --vendor ./...

# Running tests according to current version.
if [[ "$VERSION" == "snapshot" ]]; then
    go test -race --cover --covermode=atomic ./...
else
    echo "" > build/ci/out/coverage.txt
    for d in $(go list ./... | grep -v vendor); do
        go test -race -coverprofile=profile.out -covermode=atomic $d
        if [ -f profile.out ]; then
            cat profile.out >>  build/ci/out/coverage.txt
            rm profile.out
        fi
    done
fi