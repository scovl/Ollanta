# Contributing to Ollanta

This document is the fast path for contributing to Ollanta. It focuses on the commands, validation flow, and repository conventions you need when changing scanner, server, frontend, or shared modules.

For architecture details, see [docs/architecture.md](docs/architecture.md). For a docs-friendly version of this workflow, see [docs/contributing.md](docs/contributing.md).

## Prerequisites

- Go 1.21+
- CGO toolchain available locally
- Docker and Docker Compose
- Node.js for the local scanner frontend in `ollantascanner/server/static`
- `golangci-lint` available in PATH

On Windows, the project expects an MSYS2/MinGW toolchain for CGO-backed modules.

## Repository layout

Ollanta is a multi-module Go workspace. Important areas:

- `domain/` and `application/`: hexagonal core
- `adapter/`: HTTP, telemetry, persistence and bridge adapters
- `ollantascanner/`: CLI scanner and local UI server
- `ollantaweb/`: centralized API server
- `ollantarules/`: rule registry and rule metadata
- `ollantaparser/`: tree-sitter boundary and the only true CGo-heavy parser module

## Development workflow

### 1. Make the smallest change that solves the problem

- Prefer updating existing modules over adding new abstractions
- Keep types in their canonical package
- Do not duplicate rule metadata or shared structs
- Do not silently ignore errors

### 2. Rebuild what your change actually touches

General validation:

```sh
make build
make test
make lint
```

Frontend changes in `ollantascanner/server/static`:

```sh
cd ollantascanner/server/static
npm test
npm run build
```

Scanner UI changes served through Docker:

```sh
docker compose up -d --build --force-recreate serve
```

Server-only changes:

```sh
docker compose --profile server build ollantaweb
```

## Linting and testing rules

- Run `make lint` from the repository root
- Do not run `golangci-lint` at the workspace root manually across all modules with a custom glob; the Makefile already runs the module-aware commands
- If you changed the scanner frontend, run both `npm test` and `npm run build`
- If you changed Docker-served scanner assets, recreate `serve` after rebuilding the frontend bundle

## Scanner UI and Fix with AI

The local scanner UI is embedded into the scanner binary via `go:embed`. That means:

1. Changes in `ollantascanner/server/static/src` do not reach the browser until you run `npm run build`
2. Docker users must recreate `serve` after that rebuild
3. Browser cache should not be relied on during development

To enable the mock AI agent locally:

```sh
export OLLANTA_AI_ENABLE_MOCK=1
ollanta -project-dir . -project-key my-project -format all -serve
```

To enable the mock agent with Docker Compose:

```sh
export OLLANTA_AI_ENABLE_MOCK=1
docker compose up -d --build --force-recreate serve
```

## Pull request checklist

- Update docs when the behavior, workflow, or configuration changed
- Run the relevant validation commands before opening the PR
- Call out security implications explicitly when applicable
- Keep the scope focused; avoid unrelated refactors

## Commit guidance

Use conventional commit prefixes:

- `feat:` for new functionality
- `fix:` for bug fixes
- `docs:` for documentation-only changes
- `test:` for tests
- `chore:` for maintenance work

Recommended branch format:

- `username/brief-description`

## Common mistakes to avoid

- Duplicating rule metadata across packages without a clear canonical source
- Importing CGo-heavy packages into layers that must stay CGo-free
- Forgetting to rebuild `ollantascanner/server/static/dist/app.js` after frontend changes
- Assuming a browser reload is enough when the embedded scanner assets were not rebuilt into the binary
- Ignoring `make lint` failures caused by scanner or rule packages