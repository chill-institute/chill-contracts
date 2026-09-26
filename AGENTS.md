# Contracts

`contracts` is the public contract boundary for `chill.institute`.

## Work

- `mise install`
- `mise run generate` regenerates Go, TypeScript, and OpenAPI output.

## Conventions

- Treat changes here as consumer-facing API changes.
- Keep protobuf sources, generated artifacts, and package metadata in sync.
- Commit schema and generated output together.
- Keep application behavior in consumers; this repo owns only schemas and
  generated artifacts.
- Only `feat` (minor), `fix`, `perf`, `refactor`, `revert` (patch) and breaking
  (major) commits release from `main`; `docs`, `test`, `build`, `ci`, `chore`
  and `deps` do not. The rule set is the commit-analyzer `releaseRules` in
  [`.releaserc.json`](./.releaserc.json), identical in every chill.institute
  package repo; [`smoke.mjs`](./.github/release/smoke.mjs) fails when it drifts.

## Proof map

| Change | Check | Runs | Leaves |
| --- | --- | --- | --- |
| Docs | `mise run verify` | local, [PR](./.github/workflows/verify.yml) and [main](./.github/workflows/main.yml) `verify` jobs | exit status |
| Protobuf schemas, `buf*.yaml` | `mise run verify`, then `CONTRACTS_BASE_REF=<git-ref> mise run compatibility:check` | local, [PR](./.github/workflows/verify.yml) against the base commit, [main](./.github/workflows/main.yml) against the preceding release | exit status, regenerated `gen/` diff; CI fails on uncommitted output |
| Package metadata, generated consumer surface | `mise run verify-consumers` (in `verify`) | local, CI `verify` jobs | exit status; Go compile in `testdata/consumers/go`, JS import and round-trip in `testdata/consumers/ts` |
| Release rules, locked semantic-release | `npm ci --prefix .github/release --ignore-scripts`, then `node .github/release/smoke.mjs` | local, PR [Verify](./.github/workflows/verify.yml) | exit status |
| Workflows | `mise run actions` and `go test ./internal/workflowpolicy` (in `verify`) | local, CI `verify` jobs | exit status |
| Pushed workflow changes | [shared scan](https://github.com/chill-institute/.github/tree/main/.github/actions/scan), last step of the [Main](./.github/workflows/main.yml) `verify` job: Actionlint and Zizmor when the pushed range touches workflows; secrets rely on GitHub secret scanning | CI on push to `main` (pushed range) and Main dispatch (full history) | failed run |
| Release | [main](./.github/workflows/main.yml) `release` job after `verify` | CI only, `release` Environment | npm package, tag, GitHub release, `[skip ci]` version commit; `gh release verify` and `go mod download` of the tag |

Gaps:

- No TypeScript type-check of the consumer fixture. Owner: chill-institute/chill-contracts.
- No Markdown or link check for docs. Owner: chill-institute/chill-contracts.

## Read More

- public-boundary model and generation flow: [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)
- setup, consumer usage, and release flow: [CONTRIBUTING.md](./CONTRIBUTING.md)
