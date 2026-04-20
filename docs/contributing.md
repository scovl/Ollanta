# Contributing

This guide explains how to contribute to Ollanta without having to infer the workflow from multiple modules.

If you want the short operational version, see [../CONTRIBUTIONS.md](../CONTRIBUTIONS.md).

## Before you start

Ollanta is split into multiple Go modules inside a single workspace. The project mixes:

- pure Go modules
- scanner and rules packages that depend on tree-sitter through CGo
- an embedded frontend bundle for the local scanner UI
- a separate server stack for centralized history, auth, and APIs

Because of that, the right validation depends on what you changed.

## Local setup

Minimum tooling:

- Go 1.21+
- C compiler available for CGo-backed modules
- Docker and Docker Compose
- Node.js for `ollantascanner/server/static`
- `golangci-lint`

Useful baseline commands:

```sh
make build
make test
make lint
```

## Which commands to run

### If you changed scanner, rules, parser, or shared Go code

Run:

```sh
make build
make test
make lint
```

### If you changed the local scanner frontend

Run:

```sh
cd ollantascanner/server/static
npm test
npm run build
```

Then, if you use Docker for the scanner UI:

```sh
cd ../../..
docker compose up -d --build --force-recreate serve
```

### If you changed the centralized server

Run the standard Go validation, then rebuild the server image when you need to test through Docker:

```sh
docker compose --profile server build ollantaweb
```

## Embedded scanner UI

The scanner UI is compiled from TypeScript into `ollantascanner/server/static/dist/app.js` and then embedded into the scanner binary.

That has two consequences:

1. Editing `src/*.ts` is not enough; you must regenerate `dist/app.js`
2. Rebuilding Docker `serve` is required if you want the containerized scanner to pick up the new frontend

If the browser keeps showing old behavior after a frontend change, confirm these steps in order:

1. `npm run build` was executed in `ollantascanner/server/static`
2. `docker compose up -d --build --force-recreate serve` was executed if you are using Docker
3. the browser was refreshed after the new container came up

## AI fix workflow

The local scanner UI includes a `Fix with AI` tab in issue details.

Supported local configuration patterns:

- simple OpenAI-compatible setup through `OPENAI_API_KEY` and `OLLANTA_AI_OPENAI_MODEL`
- explicit multi-agent JSON configuration through `OLLANTA_AI_AGENTS`
- local mock development through `OLLANTA_AI_ENABLE_MOCK=1`

The apply step is intentionally guarded: if the target file changed after preview generation, Ollanta rejects the apply request and asks for a fresh preview.

## Documentation expectations

Update documentation when you change:

- CLI flags
- Docker workflows
- environment variables
- server routes or auth requirements
- scanner UI behavior that users rely on

Relevant docs include:

- [architecture.md](architecture.md)
- [api.md](api.md)
- [quality-gates.md](quality-gates.md)
- [rules.md](rules.md)

## Review expectations

Good contributions in Ollanta usually have these properties:

- they solve the root cause instead of patching symptoms
- they preserve the hexagonal boundaries
- they do not duplicate types or data sources
- they leave the repo with passing validation for the touched area

If your change affects scanner UX, include the exact commands another contributor can run to verify the behavior.