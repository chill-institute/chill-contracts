# Contributing

Edit schemas in `proto/` and commit their generated Go, TypeScript, and OpenAPI
artifacts in the same change.

## Start

```bash
mise install
mise run generate
mise run verify
```

`mise run verify` runs the steps in [`tasks.verify`](./mise.toml). It does not
type-check the TypeScript consumer. CI separately rejects uncommitted generated
changes.

The TypeScript fixture keeps its lockfile directory local so pnpm resolves the
file dependency independently of the root workspace. Its pnpm bootstrap lockfile
is ignored; dependency resolution stays lockfile-free. Update both manifests'
pnpm versions together; Renovate groups these pins.

## Review Contract Changes

- Treat every schema change as consumer-facing.
- Review `gen/go/`, `gen/ts/`, and `gen/openapi/` after generation.
- Explain compatibility or migration risk in the pull request.
- Consumers should depend on releases, not `main`.
- `CONTRACTS_BASE_REF=<git-ref> mise run compatibility:check` checks an explicit
  baseline. Both CI paths use this policy: pull requests compare with their base
  commit; main compares with the highest stable release tag reachable before
  HEAD. A source break requires a Conventional Commit breaking marker in that
  baseline-to-HEAD range and must still preserve wire compatibility.

## Release

Merges to `main` are released from Conventional Commits. Automation creates the
tag, npm package, GitHub release, and signed version commit.
