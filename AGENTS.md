# Contracts

Protobuf schemas and generated Go, TypeScript, and OpenAPI artifacts for the
public `chill.institute` API. Consumers own behavior built on them.

```bash
mise install
mise run generate   # rebuild gen/ from proto/
mise run verify
```

## Rules

- Every schema change is a consumer-facing API change; keep wire compatibility.
- Edit `proto/`, never `gen/`; commit schema and regenerated `gen/` together.
- Source breaks need a Conventional Commit breaking marker; release types come
  from `releaseRules` in [`.releaserc.json`](./.releaserc.json), shared by every
  chill.institute package repo and checked by [`smoke.mjs`](./.github/release/smoke.mjs).
- Field-level semantics live in proto comments; cross-cutting API rules in
  [docs/api.md](./docs/api.md).

## Proof map

| Change | Check |
| --- | --- |
| Protobuf, `buf*.yaml` | `mise run verify`, then `CONTRACTS_BASE_REF=<git-ref> mise run compatibility:check` |
| Package metadata, generated surface | `mise run verify-consumers` (in `verify`) |
| Workflows | `mise run actions`, `go test ./internal/workflowpolicy` (in `verify`) |
| Release rules | `mise run release:smoke` |

CI also runs `mise run generated:check`, which fails when `gen/` differs from the
commit. Gaps: the TypeScript consumer fixture is not type-checked; docs have no
link check.
