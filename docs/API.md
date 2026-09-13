# API Reference

Entry point for the public `chill.institute` API docs.

## Canonical URLs

| Surface | URL |
|---------|-----|
| API host | https://api.chill.institute |
| v4 OpenAPI doc | https://api.chill.institute/v4 |

## Auth

| Surface | Header | Format |
|---------|--------|--------|
| `v4.CoreService` | `X-API-Key` | Service key |
| `v4.UserService` | `Authorization` | User Token: `Bearer <token>` |

## Source Of Truth

The schema lives in:

- `proto/chill/v4/api.proto`
- `gen/openapi/chill/v4/api.swagger.yaml`
- `gen/go/chill/v4/`
- `gen/ts/chill/v4/`

Use the generated OpenAPI and generated client artifacts for field-level request and response details.

## Catalog dates and sorting

`Movie.release_date` and `TVShow.first_air_date` contain date-only `YYYY-MM-DD`
values, or an empty string when unknown. TV dates refer to the series premiere;
`TVShowDetail.first_air_date` has the same meaning.

`CatalogSettings.sort` is shared across movies, TV shows, and providers, independently
of search-result sorting. `POPULARITY` is the default and preserves the selected provider's order.
Other choices sort by rating or release date in either direction.

On `SaveUserSettings`, omitting the optional sort preserves its stored value.
Missing stored preferences default to `POPULARITY`. Sending `POPULARITY` explicitly
resets a preference; sending `UNSPECIFIED` explicitly is invalid.

The deprecated `movies_sort` and `tv_shows_sort` fields remain in the wire contract;
new consumers use `sort`.

## Playback resolution

`UserService.ResolvePlayback` uses the same regular chill bearer token as the
other user RPCs. Reuse `UserService.GetFolder` for browsing; no separate library
service, delegated credential or client-specific authorization is introduced.
Schema availability does not establish that Engine implements resolution yet.

Resolution checks the authenticated user's current access to the requested
file and returns ready, pending or unavailable. Missing and inaccessible files
share the not-found result. Invalid credentials use `unauthenticated`;
transport/provider outages remain RPC errors. Resolution never starts transfers
or conversion. Pending describes existing provider work; a later explicit
request can observe readiness changes.

Ready contains a direct HTTPS media URL usable without forwarding provider
credentials, known format information, and optional subtitles. Unspecified
format values mean unknown; each client decides compatibility. Media/subtitle
URLs declare known or explicitly unknown expiry. Provider int64 file IDs retain
precision in generated clients and JSON strings.

Personalized responses require `Cache-Control: no-store`. Bearer tokens and
signed URLs must stay out of logs, traces and shared caches. Protobuf does not
enforce these semantics; Engine and its consumers must validate them.
