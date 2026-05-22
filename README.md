# chill-contracts

![chill.institute contracts](https://chill.institute/banner.png)

Public API contracts for [chill.institute](https://chill.institute), your favorite [put.io](https://put.io) extension since 2018.

This repo is the boundary between the backend, web app, CLI, and downstream consumers. It owns protobuf schemas plus generated Go, TypeScript, and OpenAPI artifacts.

## Install

TypeScript consumers:

```bash
npm install @chill-institute/contracts
```

Go consumers:

```bash
go get github.com/chill-institute/chill-contracts
```

## Develop

```bash
mise install
mise run generate
mise run verify
```

Regenerate artifacts after changing protobuf definitions and commit schema plus generated output together.

## Docs

- [API reference](./docs/API.md): generated contract surface
- [Architecture](./docs/ARCHITECTURE.md): package layout and generation flow
- [Security](./SECURITY.md): reporting and contract safety notes

## Contributing

Please read the [contributing guide](./CONTRIBUTING.md).

## License

Licensed under the [MIT License](./LICENSE).
