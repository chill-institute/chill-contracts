# API

The v4 API is served at <https://api.chill.institute>; browsable docs are at
<https://api.chill.institute/v4> and its OpenAPI document at
<https://api.chill.institute/v4/docs>. Field-level semantics, including
deprecations, are comments in [`api.proto`](../proto/chill/v4/api.proto);
generated details are in
[`gen/openapi`](../gen/openapi/chill/v4/api.swagger.yaml).

| Service | Auth header |
| --- | --- |
| `CoreService` | `X-API-Key: <service key>` |
| `UserService` | `Authorization: Bearer <user token>` |

Rules the schema cannot enforce:

- Call `StartPlaybackConversion` only when a user chooses to play, never to list
  or preview; it is the only playback call that changes provider state.
- Every `CoreService` and `UserService` response sends `Cache-Control: no-store`.
  Bearer tokens and signed URLs stay out of logs, traces, and shared caches.
