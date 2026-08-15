# Contributing

Edit schemas in `proto/` and commit their generated Go, TypeScript, and OpenAPI
artifacts in the same change.

## Start

```bash
mise install
mise run generate
mise run verify
```

`mise run verify` checks formatting, generation drift, and clean Go and
TypeScript consumers under `testdata/consumers/`.

## Review Contract Changes

- Treat every schema change as consumer-facing.
- Review `gen/go/`, `gen/ts/`, and `gen/openapi/` after generation.
- Explain compatibility or migration risk in the pull request.
- Consumers should depend on releases, not `main`.

## Release

Merges to `main` are released from Conventional Commits. Automation creates the
tag, npm package, GitHub release, and signed version commit.
