# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog and this project adheres to Semantic Versioning.

## [Unreleased]
### Added
- Dual-backend benchmarks (Linear vs Gonum) with deterministic datasets (uniform/clustered) in 2D/4D for N=1k/10k; artifacts uploaded in CI as `bench-linear.txt` and `bench-gonum.txt`.
- Documentation: Performance guide updated to cover backend selection, how to run both backends, CI artifact links, and guidance on when each backend is preferred.
- Documentation: Performance guide now includes a Sample results table sourced from a recent local run.
- Documentation: README gained a “Backend selection” section with default behavior, build tag usage, overrides, and supported metrics notes.
- Documentation: API reference (`docs/api.md`) now documents `KDBackend`, `WithBackend`, default selection, and supported metrics for the optimized backend.
- Examples: Added `examples/wasm-browser/` minimal browser demo (ESM + HTML) for the WASM build.
- pkg.go.dev Examples: `ExampleNewKDTreeFromDim_Insert`, `ExampleKDTree_TiesBehavior`, `ExampleKDTree_Radius_none`.
- Lint: enable `errcheck` in `.golangci.yml` with test-file exclusion to reduce noise.
- CI: enable module cache in `actions/setup-go` to speed up workflows.

### Fixed
- go vet failures in examples due to misnamed `Example*` functions; renamed to avoid referencing non-existent methods and identifiers.
- Stabilized `ExampleKDTree_Nearest` to avoid a tie case; adjusted query and expected output.
- Relaxed floating-point equality in `TestWeightedCosineDistance_Basics` to use an epsilon, avoiding spurious failures on some toolchains.

## [0.3.0] - 2025-11-03
### Added
- New distance metrics: `CosineDistance` and `WeightedCosineDistance` (1 - cosine similarity), with robust zero-vector handling and bounds.
- N-D normalization helpers: `ComputeNormStatsND`, `BuildND`, `BuildNDWithStats` for arbitrary dimensions, with validation errors (`ErrInvalidFeatures`, `ErrInvalidWeights`, `ErrInvalidInvert`, `ErrStatsDimMismatch`).
- Tests: unit tests for cosine/weighted-cosine metrics; parity tests between `Build4D` and `BuildND`; error-path tests; extended fuzz to include cosine metrics.
- pkg.go.dev examples: `ExampleBuildND`, `ExampleBuildNDWithStats`, `ExampleCosineDistance`.

### Changed
- Version bumped to `0.3.0`.
- README: list Cosine among supported metrics.

## [0.2.1] - 2025-11-03
### Added
- Normalization stats helpers: `AxisStats`, `NormStats`, `ComputeNormStats2D/3D/4D`.
- Builders that reuse stats: `Build2DWithStats`, `Build3DWithStats`, `Build4DWithStats`.
- CI: coverage integration (`-coverprofile`), Codecov upload and badge.
- CI: benchmark runs publish artifacts per Go version.
- Docs: Performance page (`docs/perf.md`) and MkDocs nav entry.
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

[Unreleased]: https://github.com/Snider/Poindexter/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/Snider/Poindexter/releases/tag/v0.3.0
[0.2.1]: https://github.com/Snider/Poindexter/releases/tag/v0.2.1
[0.2.0]: https://github.com/Snider/Poindexter/releases/tag/v0.2.0
