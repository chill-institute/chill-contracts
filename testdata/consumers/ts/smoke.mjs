import assert from "node:assert/strict";
import { create, fromBinary, fromJson, toBinary, toJson } from "@bufbuild/protobuf";
import {
  ResolvePlaybackRequestSchema,
  ResolvePlaybackResponseSchema,
  AddTransferRequestSchema,
  CatalogOriginSchema,
  CatalogSettingsSchema,
  CatalogSort,
  DownloadSettingsSchema,
  ReleaseInfoSchema,
  SearchSettingsSchema,
  SearchResultSchema,
  SearchResultDisplayBehavior,
  UserSettingsSchema,
  UserService,
} from "@chill-institute/contracts/chill/v4/api_pb";

assert.equal(UserService.typeName, "chill.v4.UserService");
assert.equal(SearchResultDisplayBehavior.FASTEST, 2);
for (const [schema, typeName] of [
  [ReleaseInfoSchema, "chill.v4.ReleaseInfo"],
  [SearchResultSchema, "chill.v4.SearchResult"],
  [UserSettingsSchema, "chill.v4.UserSettings"],
  [SearchSettingsSchema, "chill.v4.SearchSettings"],
  [CatalogSettingsSchema, "chill.v4.CatalogSettings"],
  [DownloadSettingsSchema, "chill.v4.DownloadSettings"],
  [CatalogOriginSchema, "chill.v4.CatalogOrigin"],
  [AddTransferRequestSchema, "chill.v4.AddTransferRequest"],
]) {
  assert.equal(schema.typeName, typeName);
}

for (const sort of [
  undefined,
  CatalogSort.UNSPECIFIED,
  CatalogSort.POPULARITY,
  CatalogSort.RELEASE_DATE_DESC,
]) {
  const settings = create(CatalogSettingsSchema, { sort });
  for (const decoded of [
    fromBinary(CatalogSettingsSchema, toBinary(CatalogSettingsSchema, settings)),
    fromJson(CatalogSettingsSchema, toJson(CatalogSettingsSchema, settings)),
  ]) {
    assert.equal(decoded.sort, sort);
  }
}

const largeId = "9223372036854775807";
const request = fromJson(ResolvePlaybackRequestSchema, { fileId: largeId });
assert.equal(toJson(ResolvePlaybackRequestSchema, request).fileId, largeId);
for (const result of [
  { ready: { media: { url: "https://example.com/movie.mp4", expiryUnknown: true } } },
  { ready: { media: { url: "https://example.com/movie.mp4", expiresAt: "2030-01-01T00:00:00Z" } } },
  { pending: { reason: "PENDING_REASON_PROCESSING" } },
  { unavailable: { reason: "UNAVAILABLE_REASON_NOT_FOUND" } },
]) {
  const parsed = fromJson(ResolvePlaybackResponseSchema, result);
  assert.deepEqual(toJson(ResolvePlaybackResponseSchema, parsed), result);
}
assert.ok(UserService.method.resolvePlayback);

console.log("ts consumer import smoke passed");
