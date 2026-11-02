# Contributing

Edit schemas in `proto/`, run `mise run generate`, and commit the regenerated
`gen/` in the same change.

```bash
mise install
mise run generate
mise run verify
CONTRACTS_BASE_REF=origin/main mise run compatibility:check
```

Pull requests explain compatibility or migration risk.
[compatibility.go](./internal/workflowpolicy/compatibility.go) owns the breaking
marker and wire-compatibility checks.

TypeScript consumer setup: [verify-consumers-ts](./mise.toml).

Release rules: [`.releaserc.json`](./.releaserc.json). Publication:
[Main](./.github/workflows/main.yml).
