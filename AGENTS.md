# Contracts

`contracts` is the public contract boundary for `chill.institute`.

## Work

- `mise install`
- `mise run generate` regenerates Go, TypeScript, and OpenAPI output.
- `mise run verify` is the canonical gate, including clean consumer builds.

## Conventions

- Treat changes here as consumer-facing API changes.
- Keep protobuf sources, generated artifacts, and package metadata in sync.
- Commit schema and generated output together.
- Keep application behavior in consumers; this repo owns only schemas and
  generated artifacts.

## Read More

- public-boundary model and generation flow: [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)
- setup, consumer usage, and release flow: [CONTRIBUTING.md](./CONTRIBUTING.md)
