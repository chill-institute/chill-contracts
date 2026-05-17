# Architecture

This document describes how `chill-contracts` is built.

## System Context

```mermaid
graph LR
  Proto["protobuf schemas"] --> Generated["generated artifacts"]
  Generated --> Web["web"]
  Generated --> CLI["cli"]
  Generated --> Backend["backend"]
  Generated --> Future["future clients"]
```

## Components

| Component | Responsibility | Talks to |
|-----------|----------------|----------|
| `proto/` | Canonical public API schemas | Buf, generators |
| `gen/go/` | Generated Go contract types | backend consumers |
| `gen/ts/` | Generated TypeScript contract types | `web`, future JS clients |
| `gen/openapi/` | Generated OpenAPI output | docs and tooling |
| `testdata/consumers/` | Tiny downstream fixtures proving Go and TypeScript consumers still compile/import | generated artifacts, `mise` |
| package metadata | Publish the TypeScript package and release artifacts | npm, GitHub releases |

## Generation Model

```mermaid
graph TD
  Author["edit proto"] --> Generate["mise run generate"]
  Generate --> Go["gen/go"]
  Generate --> TS["gen/ts"]
  Generate --> OpenAPI["gen/openapi"]
  Go --> Backend["backend consumers"]
  TS --> Web["web consumers"]
```

## Boundaries

- This repo owns public contracts only.
- Backend behavior, web UX, and CLI command surfaces do not live here.
- Compatibility decisions here are consumer-facing and should be treated like API changes.

## Release Model

- `main` is the release branch.
- `Verify` runs on pull requests.
- `Main` runs on pushes to `main`.
- verification includes tiny downstream consumer checks for Go and TypeScript fixtures under `testdata/consumers/`
- `Main` re-verifies the repo and then publishes the npm package and GitHub release from `main`
- npm publishing uses trusted publishing from the `Main` workflow with the `release` environment
- release jobs opt into only the GitHub token scopes they need, run only from `main`, and use pinned workflow actions
- `Main` remains available as a manual fallback for the trusted publish path from `main`
