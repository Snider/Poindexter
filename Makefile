# Maintainer Makefile for Poindexter
# Usage: `make <target>`
# Many targets are CI-parity helpers for local use.

# Tools (override with env if needed)
GO           ?= go
GOLANGCI_LINT?= golangci-lint
GORELEASER   ?= goreleaser
MKDOCS       ?= mkdocs

# Params
FUZZTIME ?= 10s
BENCHOUT ?= bench.txt
COVEROUT ?= coverage.out
COVERHTML?= coverage.html

.PHONY: help all
all: help
help: ## List available targets
	@awk 'BEGIN {FS = ":.*##"}; /^[a-zA-Z0-9_.-]+:.*##/ {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: tidy
tidy: ## Run `go mod tidy`
	$(GO) mod tidy

.PHONY: tidy-check
tidy-check: ## Run tidy and ensure go.mod/go.sum unchanged
	$(GO) mod tidy
	@git diff --exit-code -- go.mod go.sum

.PHONY: fmt
fmt: ## Format code with go fmt
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: build
build: ## Build all packages
	$(GO) build ./...

# WebAssembly build outputs
DIST_DIR ?= dist
WASM_OUT ?= $(DIST_DIR)/poindexter.wasm
WASM_EXEC ?= $(shell $(GO) env GOROOT)/lib/wasm/wasm_exec.js

.PHONY: wasm-build
wasm-build: ## Build WebAssembly module to $(WASM_OUT)
	@mkdir -p $(DIST_DIR)
	GOOS=js GOARCH=wasm $(GO) build -o $(WASM_OUT) ./wasm
	@set -e; \
	if [ -n "$$WASM_EXEC" ] && [ -f "$$WASM_EXEC" ]; then \
	  cp "$$WASM_EXEC" $(DIST_DIR)/wasm_exec.js; \
	else \
	  CAND1="$$($(GO) env GOROOT)/lib/wasm/wasm_exec.js"; \
	  CAND2="$$($(GO) env GOROOT)/libexec/lib/wasm/wasm_exec.js"; \
	  if [ -f "$$CAND1" ]; then cp "$$CAND1" $(DIST_DIR)/wasm_exec.js; \
	  elif [ -f "$$CAND2" ]; then cp "$$CAND2" $(DIST_DIR)/wasm_exec.js; \
	  else echo "Warning: could not locate wasm_exec.js under GOROOT or WASM_EXEC; please copy it manually"; fi; \
	fi
	@echo "WASM built: $(WASM_OUT)"

.PHONY: npm-pack
npm-pack: wasm-build ## Prepare npm package folder with dist artifacts
	@mkdir -p npm/poindexter-wasm
	@rm -rf npm/poindexter-wasm/dist
	@cp -R $(DIST_DIR) npm/poindexter-wasm/dist
	@cp LICENSE npm/poindexter-wasm/LICENSE
	@cp README.md npm/poindexter-wasm/PROJECT_README.md
	@echo "npm package prepared in npm/poindexter-wasm"

.PHONY: examples
examples: ## Build all example programs under examples/
	@if [ -d examples ]; then $(GO) build ./examples/...; else echo "No examples/ directory"; fi

.PHONY: test
test: ## Run unit tests
	$(GO) test ./...

.PHONY: race
race: ## Run tests with race detector
	$(GO) test -race ./...

.PHONY: cover
cover: ## Run tests with race + coverage and summarize
	$(GO) test -race -coverprofile=$(COVEROUT) -covermode=atomic  ./...
	@$(GO) tool cover -func=$(COVEROUT) | tail -n 1

.PHONY: coverfunc
coverfunc: ## Print per-function coverage from $(COVEROUT)
	@$(GO) tool cover -func=$(COVEROUT)

.PHONY: cover-kdtree
cover-kdtree: ## Print coverage details for kdtree.go only
	@$(GO) tool cover -func=$(COVEROUT) | grep 'kdtree.go' || true

.PHONY: coverhtml
coverhtml: cover ## Generate HTML coverage report at $(COVERHTML)
	@$(GO) tool cover -html=$(COVEROUT) -o $(COVERHTML)
	@echo "Wrote $(COVERHTML)"

.PHONY: fuzz
fuzz: ## Run Go fuzz tests for $(FUZZTIME)
	@set -e; \
	PKGS="$$($(GO) list ./...)"; \
	for pkg in $$PKGS; do \
	  FUZZES="$$($(GO) test -list '^Fuzz' $$pkg | grep '^Fuzz' || true)"; \
	  if [ -z "$$FUZZES" ]; then \
	    echo "==> Skipping $$pkg (no fuzz targets)"; \
	    continue; \
	  fi; \
	  for fz in $$FUZZES; do \
	    echo "==> Fuzzing $$pkg :: $$fz for $(FUZZTIME)"; \
	    $(GO) test -run=NONE -fuzz=^$${fz}$$ -fuzztime=$(FUZZTIME) $$pkg; \
	  done; \
	done

.PHONY: bench
bench: ## Run benchmarks and write $(BENCHOUT)
	$(GO) test -bench . -benchmem -run=^$$ ./... | tee $(BENCHOUT)

.PHONY: lint
lint: ## Run golangci-lint (requires it installed)
	$(GOLANGCI_LINT) run

.PHONY: vuln
vuln: ## Run govulncheck (requires it installed)
	govulncheck ./...

.PHONY: ci
ci: tidy-check build vet cover examples bench lint vuln wasm-build ## CI-parity local run (includes wasm-build)
	@echo "CI-like checks completed"

.PHONY: release
release: ## Run GoReleaser to publish a tagged release (requires tag and permissions)
	$(GORELEASER) release --clean --config .goreleaser.yaml

.PHONY: snapshot
snapshot: ## Run GoReleaser in snapshot mode (no publish)
	$(GORELEASER) release --skip=publish --clean --config .goreleaser.yaml

.PHONY: docs-serve
docs-serve: ## Serve MkDocs locally (requires mkdocs-material)
	$(MKDOCS) serve -a 127.0.0.1:8000

.PHONY: docs-build
docs-build: ## Build MkDocs site into site/
	$(MKDOCS) build

.PHONY: clean
clean: ## Remove generated files and directories
	rm -rf $(DIST_DIR) $(COVEROUT) $(COVERHTML) $(BENCHOUT)
