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

| Change | Check | Runs |
| --- | --- | --- |
| Protobuf, `buf*.yaml` | `mise run verify`, then `CONTRACTS_BASE_REF=<git-ref> mise run compatibility:check` | local; [PR](./.github/workflows/verify.yml) vs base, [main](./.github/workflows/main.yml) vs preceding release; CI fails on uncommitted `gen/` |
| Package metadata, generated surface | `mise run verify-consumers` (in `verify`) | local, CI |
| Workflows | `mise run actions`, `go test ./internal/workflowpolicy` (in `verify`) | local, CI |
| Release rules | `npm ci --prefix .github/release --ignore-scripts && node .github/release/smoke.mjs` | local, PR |
| Release | [main](./.github/workflows/main.yml) `release` job | CI only: npm package, tag, GitHub release |

Gaps: the TypeScript consumer fixture is not type-checked; docs have no link
check.
