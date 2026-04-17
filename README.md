# Ollanta

Ollanta is a multi-language static analysis platform written in Go. It analyses source code, reports quality issues, computes code metrics, and evaluates configurable quality gates — making it easy to enforce coding standards in any CI/CD pipeline.

Inspired by [OpenStaticAnalyzer](https://github.com/sed-inf-u-szeged/OpenStaticAnalyzer), Ollanta is designed as a modular solution where each concern lives in its own Go module.

---

## Supported Languages

| Language   | Engine          | Rules |
|------------|-----------------|-------|
| Go         | Native (`go/ast`) | `go:no-large-functions`, `go:no-naked-returns`, `go:naming-conventions` |
| JavaScript | Tree-sitter     | `js:no-large-functions` |
| Python     | Tree-sitter     | `py:no-large-functions` |

---

## Architecture

```
ollantacore/      — shared domain types (Component, Issue, Metric)
ollantaparser/    — language parsers (Go AST, Tree-sitter)
ollantarules/     — rule definitions and language sensors
ollantascanner/   — scan orchestration, file discovery, reporting
ollantaengine/    — post-scan analysis (quality gates, issue tracking, metrics aggregation)
```

---

## Prerequisites

- **Go 1.21+** with CGO enabled
- **GCC** (required by the Tree-sitter runtime)
  - Linux/macOS: `gcc` from your system package manager
  - Windows: [MSYS2](https://www.msys2.org/) → `pacman -S mingw-w64-x86_64-gcc`

---

## Build

```sh
# Build all modules
make build

# Run all tests
make test

# Format source code
make fmt

# Run the linter (requires golangci-lint)
make lint
```

On Windows, the Makefile automatically prepends `C:\msys64\mingw64\bin` to PATH and sets `CGO_ENABLED=1`.

---

## Usage

### Run a scan

```sh
go run ./ollantascanner/cmd/ollanta \
  -project-dir /path/to/myproject \
  -project-key  my-project \
  -format       all
```

### Flags

| Flag            | Default     | Description |
|-----------------|-------------|-------------|
| `-project-dir`  | `.`         | Root directory to scan |
| `-project-key`  | *(dir name)*| Identifier used in reports |
| `-sources`      | `./...`     | Comma-separated source patterns |
| `-exclusions`   | *(none)*    | Comma-separated glob patterns to exclude |
| `-format`       | `all`       | Output format: `summary`, `json`, `sarif`, `all` |
| `-debug`        | `false`     | Enable verbose debug output |

### Output formats

| Format    | Description |
|-----------|-------------|
| `summary` | Prints a human-readable table to stdout |
| `json`    | Saves `.ollanta/report.json` |
| `sarif`   | Saves `.ollanta/report.sarif` (compatible with GitHub Code Scanning) |
| `all`     | Both `json` and `sarif` |

Reports are written to a `.ollanta/` directory inside the scanned project root.

---

## Example output (summary)

```
Project : my-project
Files   : 42    Lines : 3 218    NCLOC : 2 104    Comments : 311

ISSUES (7)
  CRITICAL  go:no-naked-returns       handlers/auth.go:87
  MAJOR     go:no-large-functions     handlers/auth.go:12
  MAJOR     go:no-large-functions     services/payment.go:34
  MINOR     go:naming-conventions     models/user_model.go:8
  ...

Quality Gate : ERROR
  ✗  bugs > 0  (actual: 1)
  ✓  coverage ≥ 80
```

---

## Rules reference

### Go

| Rule key                  | Severity | Description |
|---------------------------|----------|-------------|
| `go:no-large-functions`   | Major    | Functions exceeding `max_lines` (default: 40) |
| `go:no-naked-returns`     | Critical | Naked `return` in functions with named return values longer than `min_lines` (default: 5) |
| `go:naming-conventions`   | Minor    | Exported names must use MixedCaps; underscores not allowed per Effective Go |

### JavaScript

| Rule key                | Severity | Description |
|-------------------------|----------|-------------|
| `js:no-large-functions` | Major    | Functions exceeding `max_lines` (default: 40) |

### Python

| Rule key                | Severity | Description |
|-------------------------|----------|-------------|
| `py:no-large-functions` | Major    | Functions exceeding `max_lines` (default: 40) |

---

## Quality Gates

Quality gates evaluate numeric metrics against configurable thresholds after a scan. Each condition uses one of these operators: `gt`, `lt`, `eq`, `gte`, `lte`.

Example gate configuration (Go API):

```go
conditions := []qualitygate.Condition{
    {MetricKey: "bugs",     Operator: qualitygate.OpGreaterThan, ErrorThreshold: 0,  Description: "Zero bugs"},
    {MetricKey: "coverage", Operator: qualitygate.OpLessThan,    ErrorThreshold: 80, Description: "Coverage ≥ 80%"},
}
status := qualitygate.Evaluate(conditions, measures)
if !status.Passed() {
    log.Fatal("Quality gate failed")
}
```

---

## Issue Tracking

The `ollantaengine/tracking` package compares current scan results against a previous baseline to classify each issue as **new**, **unchanged**, **closed**, or **reopened**. Issues are matched by rule key + line hash, with a fallback to file path + line number.

---

## Adding a new rule

1. Create a struct implementing the `ollantarules.Rule` interface in the appropriate language package.
2. Register it in the sensor's `ActiveRules()` method.
3. Add tests in the corresponding `*_test.go` file.

---

## License

MIT
