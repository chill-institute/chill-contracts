# Architecture

`chill-contracts` turns one protobuf source tree into the artifacts consumed by
the API, web app, CLI, and external clients.

## System Context

```mermaid
graph LR
  Proto["protobuf schemas"] --> Generated["generated artifacts"]
  Generated --> Web["web"]
  Generated --> CLI["cli"]
  Generated --> Backend["backend"]
  Generated --> Consumers["other consumers"]
```

## Components

| Path | Owns |
| --- | --- |
| `proto/` | Canonical public API schemas |
| `gen/go/` | Go types and Connect bindings |
| `gen/ts/` | TypeScript types and Connect bindings |
| `gen/openapi/` | OpenAPI output |
| `testdata/consumers/` | Clean Go and TypeScript compile checks |

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

## Contract

- This repo owns the public schemas and generated contract artifacts.
- Consumer repos own the behavior built on those contracts.
- Compatibility decisions here are consumer-facing and should be treated like API changes.

## Delivery

Pull requests run `mise run verify`. Merges to `main` are reverified before
semantic-release publishes the npm package and immutable GitHub release. The
release tag is also the Go module version and must resolve under the `/v2`
module path. The release commit is created through GitHub with bot attribution
and verification.
