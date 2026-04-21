<p align="center">
  <img src="docs/imgs/logo-dark.png" alt="Ollanta logo" width="420">
</p>

<p align="center">
  <img src="https://raw.githubusercontent.com/scovl/Ollanta/refs/heads/main/docs/imgs/o01.png" alt="Ollanta screenshot">
</p>

Ollanta is a multi-language static analysis platform written in Go. It analyzes source code, reports quality issues, computes code metrics, and evaluates configurable quality gates so teams can enforce coding standards locally, in CI, and in centralized review workflows.

Inspired by [OpenStaticAnalyzer](https://github.com/sed-inf-u-szeged/OpenStaticAnalyzer), [Semgrep](https://semgrep.dev/), [SonarQube](https://www.sonarqube.org/), Ollanta is organized as a modular platform where each concern lives in its own Go module.

---

## Supported Languages

| Language   | Engine            | Rules |
|------------|-------------------|-------|
| Go         | Native (`go/ast`) | 8 rules: large functions, naked returns, naming, cognitive complexity, nesting depth, magic numbers, too many params, TODO comments |
| JavaScript | Tree-sitter       | 4 rules: large functions, console.log, strict equality, too many params |
| Python     | Tree-sitter       | 5 rules: large functions, broad except, mutable default args, comparison to None, too many params |

---

## Architecture

Ollanta follows a **hexagonal (ports & adapters)** layout — inner modules have no external dependencies; adapters plug in at the edges. See [docs/architecture.md](docs/architecture.md) for the full module layout.

For contributor workflow, validation commands, and repository conventions, see [CONTRIBUTIONS.md](CONTRIBUTIONS.md) and [docs/contributing.md](docs/contributing.md).

---

## Quick Start

### Local scanner UI

```sh
ollanta \
  -project-dir . \
  -project-key my-project \
  -format all \
  -serve
```

This runs a local scan and opens the embedded web UI at `http://localhost:7777`.

### Docker scanner UI

```sh
docker compose up serve
```

This builds the scanner image, scans the mounted project directory, and serves the embedded UI on port `7777`.

### Centralized server stack

```sh
docker compose --profile server up -d
```

This starts PostgreSQL, ZincSearch, the `ollantaweb` API on port `8080`, plus three background roles: `ollantaworker` for scan compute, `ollantaindexer` for durable search projection, and `ollantawebhookworker` for durable webhook delivery.

If you want scanner push to wait for the final server-side result instead of stopping at `202 Accepted`, use `-server-wait` in the CLI. In Docker Compose, the `push` service can be toggled with `OLLANTA_SERVER_WAIT=true`.

---

## Prerequisites

- **Go 1.21+** with CGO enabled
- **GCC** (required by the Tree-sitter runtime)
  - Linux/macOS: `gcc` from your system package manager
  - Windows: [MSYS2](https://www.msys2.org/) → `pacman -S mingw-w64-x86_64-gcc`
- **Docker** (optional) — for container-based scanning or running the server stack

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

### CLI scan

```sh
ollanta \
  -project-dir /path/to/myproject \
  -project-key  my-project \
  -format       all
```

### Interactive web UI

```sh
ollanta \
  -project-dir /path/to/myproject \
  -project-key  my-project \
  -format       all \
  -serve -port 7777
```

Opens a local web UI at `http://localhost:7777` with the scan results.

### Fix with AI in the local UI

The local scanner UI now includes a `Fix with AI` tab inside each issue detail. The flow is:

1. Open an issue in the local UI on port `7777`
2. Choose a configured AI agent
3. Generate a fix preview for the affected snippet
4. Review the diff
5. Apply the change directly to your local file

Agent configuration is local to the scanner process. The simplest setup is OpenAI-compatible:

```sh
export OPENAI_API_KEY=your_api_key
export OLLANTA_AI_OPENAI_MODEL=gpt-4.1-mini
```

For multiple agents, configure `OLLANTA_AI_AGENTS` with a JSON array:

```json
[
  {
    "id": "openai-fast",
    "label": "OpenAI Fast",
    "provider": "openai",
    "model": "gpt-4.1-mini",
    "base_url": "https://api.openai.com/v1",
    "api_key_env": "OPENAI_API_KEY"
  },
  {
    "id": "openai-strong",
    "label": "OpenAI Strong",
    "provider": "openai",
    "model": "gpt-4.1",
    "base_url": "https://api.openai.com/v1",
    "api_key_env": "OPENAI_API_KEY"
  }
]
```

For local development without an external provider, enable the built-in mock agent:

```sh
export OLLANTA_AI_ENABLE_MOCK=1
```

The local UI only applies a fix after explicit confirmation. If the target file changed after preview generation, Ollanta rejects the apply action and asks for a new preview.

If you run the scanner through Docker Compose, export the AI-related environment variables before recreating `serve` so they are available to the scanner process inside the container.

```sh
export OLLANTA_AI_ENABLE_MOCK=1
docker compose up -d --build --force-recreate serve
```

### Push results to a centralized server

```sh
ollanta \
  -project-dir /path/to/myproject \
  -project-key  my-project \
  -format       all \
  -server       http://localhost:8080
```

Posts the report to the ollantaweb API. Exits with code 1 if the quality gate fails.

### CLI flags

| Flag            | Default      | Description |
|-----------------|--------------|-------------|
| `-project-dir`  | `.`          | Root directory to scan |
| `-project-key`  | *(dir name)* | Identifier used in reports |
| `-sources`      | `./...`      | Comma-separated source patterns |
| `-exclusions`   | *(none)*     | Comma-separated glob patterns to exclude |
| `-format`       | `all`        | Output: `summary`, `json`, `sarif`, `all` |
| `-debug`        | `false`      | Enable verbose debug output |
| `-serve`        | `false`      | Open interactive web UI after scan |
| `-port`         | `7777`       | Port for `-serve` |
| `-bind`         | `127.0.0.1`  | Bind address for `-serve` (use `0.0.0.0` in Docker) |
| `-server`       | *(none)*     | URL of ollantaweb server to push results to |
| `-branch`       | *(auto)*     | Explicit branch name for the analysis scope |
| `-commit-sha`   | *(auto)*     | Explicit commit SHA stored with the scan |
| `-pull-request-key` | *(auto)* | Pull request identifier for pull-request analysis mode |
| `-pull-request-branch` | *(auto)* | Source branch associated with the pull request |
| `-pull-request-base` | *(auto)* | Target branch associated with the pull request |

### Branch and pull request analysis

Ollanta now tracks scan state per analysis scope instead of flattening everything into one project timeline.

- Branch mode uses `-branch` when provided, otherwise it auto-detects the Git branch for `-project-dir`.
- Pull request mode is enabled when `-pull-request-key`, `-pull-request-branch`, and `-pull-request-base` are present, whether from flags or CI metadata.
- Supported CI pull-request environments are `OLLANTA_PULL_REQUEST_*`, GitHub Actions, GitLab CI, and Azure Pipelines.
- Detached HEAD executions fail fast unless `-branch` is provided explicitly.
- On the server, blank historic branch values remain visible through the resolved default branch for backward compatibility.
- Each successful ingest replaces the latest code snapshot for that branch or pull request scope. The web UI exposes that data through the `Branches`, `Pull Requests`, `Code`, and `Project Information` views.
- Snapshot storage is bounded to `128 KB` per file and `4 MB` total per scope. Truncated or omitted files are reported in the snapshot metadata.

### Output formats

| Format    | Description |
|-----------|-------------|
| `summary` | Human-readable table to stdout |
| `json`    | `.ollanta/report.json` |
| `sarif`   | `.ollanta/report.sarif` (GitHub Code Scanning compatible) |
| `all`     | Both `json` and `sarif` |

---

## Docker

### Scan with Docker

```sh
# Scan current directory and open UI at http://localhost:7777
docker compose up serve

# Scan a specific project
PROJECT_DIR=/path/to/myapp PROJECT_KEY=myapp docker compose up serve

# One-shot scan (no UI, just write report files)
docker compose run --rm scan-only
```

If you changed the scanner frontend under `ollantascanner/server/static`, rebuild the frontend bundle first and then recreate `serve`:

```sh
cd ollantascanner/server/static
npm run build

cd ../../..
docker compose up -d --build --force-recreate serve
```

### Centralized server stack

Start PostgreSQL, ZincSearch, the ollantaweb API server, and the three background roles for compute, index projection, and webhook delivery:

```sh
docker compose --profile server up -d
```

Then push scan results from any machine:

```sh
ollanta -project-dir . -project-key my-project -server http://your-server:8080
```

Or via Docker:

```sh
OLLANTA_SERVER=http://your-server:8080 docker compose --profile push run --build --rm push
```

### Environment variables

| Variable              | Default                  | Description |
|-----------------------|--------------------------|-------------|
| `PROJECT_DIR`         | `.`                      | Host directory to scan |
| `PROJECT_KEY`         | `project`                | Project identifier |
| `PORT`                | `7777`                   | Scanner UI port |
| `PG_PASSWORD`         | `ollanta_dev`            | PostgreSQL password |
| `ZINC_USER`           | `admin`                  | ZincSearch admin user |
| `ZINC_PASSWORD`       | `ollanta_dev`            | ZincSearch admin password |
| `OLLANTA_SEARCH_BACKEND` | `zincsearch`          | Search backend (`zincsearch` or `postgres`) |
| `OLLANTA_SERVER`      | `http://ollantaweb:8080` | API server URL (for push mode) |
| `OLLANTA_TOKEN`       | `ollanta-dev-scanner-token` | Scanner token used by `push` |
| `OLLANTA_SCANNER_TOKEN` | `ollanta-dev-scanner-token` | Shared secret accepted by `ollantaweb` |
| `OLLANTA_AI_ENABLE_MOCK` | *(empty)*             | Enables the built-in mock AI agent in the local scanner UI |
| `OLLANTA_AI_AGENTS`   | *(empty)*                | JSON array describing configured local AI agents |
| `OLLANTA_AI_OPENAI_MODEL` | *(empty)*            | Default OpenAI-compatible model for simple setups |
| `OLLANTA_AI_OPENAI_BASE_URL` | `https://api.openai.com/v1` | Base URL for OpenAI-compatible APIs |
| `OLLANTA_AI_OPENAI_LABEL` | `OpenAI`             | Display label shown in the local UI |
| `OPENAI_API_KEY`      | *(empty)*                | API key used by OpenAI-compatible agents |

---

## Server API (ollantaweb)

Full REST API reference at [docs/api.md](docs/api.md). All `/api/v1` routes require a `Bearer` token or API token (`olt_…`) in the `Authorization` header.

## Rollout guidance for existing projects

When enabling branch and pull request analysis on a project that already has historic scans:

1. Set `main_branch` on the project if the default branch is not `main` or `master`.
2. Run a fresh analysis on the default branch so the server establishes the current baseline and code snapshot for that scope.
3. Update CI pull-request jobs to pass explicit flags or expose supported CI environment variables.
4. For detached HEAD runners, always provide `-branch` so the scanner can resolve branch scope deterministically.
5. Historic scans with blank branch values continue to appear on the resolved default branch, so migration is backward-compatible.

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

17 built-in rules across Go, JavaScript, and Python. See [docs/rules.md](docs/rules.md) for the full reference and instructions on adding new rules.

---

## Quality Gates

Quality gates evaluate numeric metrics against configurable thresholds after every scan. The scanner CLI exits with code 1 when the gate fails. See [docs/quality-gates.md](docs/quality-gates.md).

## Kubernetes

Ollanta is designed for cloud-native operation: stateless app pods, externalized state, pluggable search backend, and independent scaling per component. Full deployment guide with manifests at [docs/kubernetes.md](docs/kubernetes.md).

## Authentication

Three mechanisms: local (JWT), OAuth (GitHub, GitLab, Google), and API tokens (`olt_…`). See [docs/authentication.md](docs/authentication.md).

## Webhooks

Projects can register outbound webhooks that fire on scan events, with HMAC-SHA256 signature verification and automatic retry. See [docs/webhooks.md](docs/webhooks.md).

## Issue Tracking

Each scan is compared against a previous baseline to classify issues as **new**, **unchanged**, **closed**, or **reopened**. See [docs/issue-tracking.md](docs/issue-tracking.md).

---

## License

Apache-2.0 — see [LICENSE](LICENSE).
