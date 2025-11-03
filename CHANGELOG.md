# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog and this project adheres to Semantic Versioning.

## [Unreleased]

## [0.2.1] - 2025-11-03
### Added
- Normalization stats helpers: `AxisStats`, `NormStats`, `ComputeNormStats2D/3D/4D`.
- Builders that reuse stats: `Build2DWithStats`, `Build3DWithStats`, `Build4DWithStats`.
- pkg.go.dev examples: `ExampleBuild2DWithStats`, `ExampleBuild4DWithStats`.
- Tests for stats parity, min==max safety, and dynamic update with reused stats.
- Docs: API reference section “KDTree Normalization Stats (reuse across updates)”; updated multi-dimensional docs with WithStats snippet.

### Changed
- Bumped version to `0.2.1`.

### Previously added in Unreleased
- README badges (pkg.go.dev, CI, Go Report Card, govulncheck) and KDTree performance/concurrency notes.
- Examples directory with runnable programs: 1D ping, 2D ping+hop, 3D ping+hop+geo, 4D ping+hop+geo+score.
- CI workflow (Go 1.22/1.23): tidy check, build, vet, test -race, build examples, govulncheck, golangci-lint.
- Lint configuration (.golangci.yml) with a pragmatic ruleset.
- Contributor docs: CONTRIBUTING.md, CODE_OF_CONDUCT.md, SECURITY.md.
- pkg.go.dev example functions for KDTree usage and helpers.
- Fuzz tests and benchmarks for KDTree (Nearest/KNearest/Radius and metrics).

## [0.2.0] - 2025-10-??
### Added
- KDTree public API with generic payloads and helper builders (Build2D/3D/4D).
- Docs pages for DHT examples and multi-dimensional KDTree usage.

[Unreleased]: https://github.com/Snider/Poindexter/compare/v0.2.1...HEAD
[0.2.1]: https://github.com/Snider/Poindexter/releases/tag/v0.2.1
[0.2.0]: https://github.com/Snider/Poindexter/releases/tag/v0.2.0
