# Contributing to Poindexter

Thanks for your interest in contributing! This document describes how to build, test, lint, and propose changes.

## Getting started

- Go 1.22+ (1.23 preferred)
- `git clone https://github.com/Snider/Poindexter`
- `cd Poindexter`

## Build and test

- Tidy deps: `go mod tidy`
- Build: `go build ./...`
- Run tests: `go test ./...`
- Run race tests: `go test -race ./...`
- Run examples: `go run ./examples/...`

## Lint and vet

We use golangci-lint in CI. To run locally:

```
# Install once
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Run
golangci-lint run
```

Also run `go vet ./...` periodically.

## Fuzzing and benchmarks

- Fuzz (manually): `go test -run=NONE -fuzz=Fuzz -fuzztime=10s`
- Benchmarks: `go test -bench=. -benchmem`

## Pull requests

- Create a branch from `main`.
- Ensure `go mod tidy` produces no changes.
- Ensure `go test -race ./...` passes.
- Ensure `golangci-lint run` has no issues.
- Update CHANGELOG.md (Unreleased section) with a brief summary.

## Coding style

- Follow standard Go formatting and idioms.
- Public APIs must have doc comments starting with the identifier name and should be concise.
- Avoid breaking changes in minor versions; use SemVer.

## Release process

Maintainers:
- Update CHANGELOG.md.
- Tag releases `vX.Y.Z`.
- Consider updating docs and README badges if needed.

Thanks for helping improve Poindexter!