# Dependency and Supply Chain Audit Report

## 1. Dependency Analysis

### 1.1. Direct Dependencies

- **Go:** The project uses Go version `1.23`, as specified in the `go.mod` file. There are no direct Go module dependencies.
- **npm:** The WebAssembly component of the project, located in `npm/poindexter-wasm`, has no direct npm dependencies listed in its `package.json` file.

### 1.2. Transitive Dependencies

- **Go:** Since there are no direct dependencies, there are no transitive Go dependencies. This was confirmed by running `go mod why -m all`.
- **npm:** An `npm audit` was performed, and it confirmed that there are no transitive dependencies.

### 1.3. License Compliance

- The project itself is licensed under the MIT license.
- Since there are no external dependencies, there are no third-party licenses to track or comply with.

## 2. Lock Files

- **Go:** The `go.mod` file is present, but since there are no dependencies, a `go.sum` file is not generated.
- **npm:** A `package-lock.json` file has been added to the repository to ensure reproducible builds, although there are currently no dependencies.

## 3. Supply Chain Risks

### 3.1. Package Sources

- **Go:** The project does not use any external Go modules.
- **npm:** The project does not use any external npm packages.

### 3.2. Build Process

- The build process is managed by a `Makefile` and automated with GitHub Actions.
- The CI/CD pipeline, defined in `.github/workflows/ci.yml` and `.github/workflows/release.yml`, is comprehensive and includes:
  - Linting (`golangci-lint`)
  - Vetting (`go vet`)
  - Testing (including race detection)
  - Code coverage analysis
  - Vulnerability scanning (`govulncheck`)
  - WebAssembly build and smoke testing
- Releases are automated using `goreleaser`, which helps ensure a consistent and reproducible build process.

## 4. Vulnerability Analysis

A vulnerability scan was performed using `govulncheck`. The scan identified 13 vulnerabilities in the Go standard library for the version used in this project (`1.23`).

### 4.1. Identified Vulnerabilities

| CVE ID          | Severity | Description                                                                 | Remediation Priority |
|-----------------|----------|-----------------------------------------------------------------------------|----------------------|
| `GO-2026-4340`  | High     | Handshake messages may be processed at the incorrect encryption level       | High                 |
| `GO-2025-4175`  | Medium   | Improper application of excluded DNS name constraints                       | Medium               |
| `GO-2025-4155`  | Medium   | Excessive resource consumption when printing error string for host cert     | Medium               |
| `GO-2025-4013`  | Medium   | Panic when validating certificates with DSA public keys                     | Medium               |
| `GO-2025-4012`  | Medium   | Lack of limit when parsing cookies can cause memory exhaustion              | Medium               |
| `GO-2025-4011`  | Medium   | Parsing DER payload can cause memory exhaustion                             | Medium               |
| `GO-2025-4010`  | Medium   | Insufficient validation of bracketed IPv6 hostnames                         | Medium               |
| `GO-2025-4009`  | Medium   | Quadratic complexity when parsing some invalid inputs                       | Medium               |
| `GO-2025-4008`  | Medium   | ALPN negotiation error contains attacker controlled information             | Medium               |
| `GO-2025-4007`  | Medium   | Quadratic complexity when checking name constraints                         | Medium               |
| `GO-2025-3751`  | Medium   | Sensitive headers not cleared on cross-origin redirect                      | Medium               |
| `GO-2025-3750`  | Low      | Inconsistent handling of O_CREATE\|O_EXCL on Unix and Windows             | Low                  |
| `GO-2025-3749`  | Low      | Usage of ExtKeyUsageAny disables policy validation                          | Low                  |

### 4.2. Remediation

The identified vulnerabilities are all in the Go standard library. The recommended remediation is to update the Go version to the latest stable release, which includes patches for these vulnerabilities. Given that some of these vulnerabilities are rated as "High" severity, this should be a high-priority action.
