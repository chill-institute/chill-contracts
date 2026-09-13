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

`CatalogSettings.movies_sort` and `tv_shows_sort` are independent of search-result
sorting. `POPULARITY` is the default and preserves the selected provider's order.
Other choices sort by rating or release date in either direction.

On `SaveUserSettings`, omitting either optional sort preserves its stored value.
Missing stored preferences default to `POPULARITY`. Sending `POPULARITY` explicitly
resets a preference; sending `UNSPECIFIED` explicitly is invalid.
