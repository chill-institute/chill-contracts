# Contracts

`contracts` is the public contract boundary for `chill.institute`.

## Stack

- protobuf source schemas in `proto/`
- generated Go, TypeScript, and OpenAPI artifacts
- `mise` for generation and verification

## Lifecycle

- `mise install`
- `mise run generate` regenerates all published artifacts.
- `mise run verify` is the canonical gate and compiles the generated package
  through clean Go and TypeScript consumers.

## Conventions

- Treat changes here as consumer-facing API changes.
- Keep protobuf sources, generated artifacts, and package metadata in sync.
- Commit schema and generated output together.
- Keep application behavior in the consuming repos; this repo owns the public schema and published artifacts.

## Read More

- public-boundary model and generation flow: [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)
- setup, consumer usage, and release flow: [CONTRIBUTING.md](./CONTRIBUTING.md)
