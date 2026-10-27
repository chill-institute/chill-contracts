# Contributing

Edit schemas in `proto/`, run `mise run generate`, and commit the regenerated
`gen/` in the same change.

```bash
mise install
mise run generate
mise run verify
CONTRACTS_BASE_REF=origin/main mise run compatibility:check
```

Pull requests explain compatibility or migration risk. A source break needs a
Conventional Commit breaking marker and must keep wire compatibility.

The TypeScript consumer fixture resolves without a lockfile; bump its pnpm
version with the root `packageManager` and mise's `pnpm` (`verify-consumers-ts` fails
on a mismatch).

Merges to `main` release from Conventional Commits: tag, npm package, GitHub
release, and a Go module version under `/v3`.
