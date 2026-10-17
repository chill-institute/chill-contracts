# API

The v4 API is served at <https://api.chill.institute>; its OpenAPI document is
<https://api.chill.institute/v4>. Field-level semantics are comments in
[`api.proto`](../proto/chill/v4/api.proto); generated details are in
[`gen/openapi`](../gen/openapi/chill/v4/api.swagger.yaml).

| Service | Auth header |
| --- | --- |
| `CoreService` | `X-API-Key: <service key>` |
| `UserService` | `Authorization: Bearer <user token>` |

Rules the schema cannot enforce:

- `ResolvePlayback` answers not-found for both missing and inaccessible files.
- Call `StartPlaybackConversion` when a user chooses to play, not to list or
  preview; it is the only playback call that changes provider state.
- Personalized responses send `Cache-Control: no-store`. Bearer tokens and
  signed URLs stay out of logs, traces, and shared caches.
- New consumers use `CatalogSettings.sort`; `movies_sort` and `tv_shows_sort`
  are deprecated.
