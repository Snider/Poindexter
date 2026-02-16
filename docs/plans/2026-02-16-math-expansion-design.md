# Poindexter Math Expansion

**Date:** 2026-02-16
**Status:** Approved

## Context

Poindexter serves as the math pillar (alongside Borg=data, Enchantrix=encryption) in the Lethean ecosystem. It currently provides KD-Tree spatial queries, 5 distance metrics, sorting utilities, and normalization helpers.

Analysis of math operations scattered across core/go, core/go-ai, and core/mining revealed common patterns that Poindexter should centralize: descriptive statistics, scaling/interpolation, approximate equality, weighted scoring, and signal generation.

## New Modules

### stats.go — Descriptive statistics
Sum, Mean, Variance, StdDev, MinMax, IsUnderrepresented.
Consumers: ml/coverage.go, lab/handler/chart.go

### scale.go — Normalization and interpolation
Lerp, InverseLerp, Remap, RoundToN, Clamp, MinMaxScale.
Consumers: lab/handler/chart.go, i18n/numbers.go

### epsilon.go — Approximate equality
ApproxEqual, ApproxZero.
Consumers: ml/exact.go

### score.go — Weighted composite scoring
Factor type, WeightedScore, Ratio, Delta, DeltaPercent.
Consumers: ml/heuristic.go, ml/compare.go

### signal.go — Time-series primitives
RampUp, SineWave, Oscillate, Noise (seeded RNG).
Consumers: mining/simulated_miner.go

## Constraints

- Zero external dependencies (WASM-compilable)
- Pure Go, stdlib only (math, math/rand)
- Same package (`poindexter`), flat structure
- Table-driven tests for every function
- No changes to existing files

## Not In Scope

- MLX tensor ops (hardware-accelerated, stays in go-ai)
- DNS tools migration to go-netops (separate PR)
- gonum backend integration (future work)
