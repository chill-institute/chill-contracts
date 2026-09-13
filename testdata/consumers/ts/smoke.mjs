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

if (UserService.typeName !== "chill.v4.UserService") {
  throw new Error(`unexpected service type name: ${UserService.typeName}`);
}

if (SearchResultDisplayBehavior.FASTEST !== 2) {
  throw new Error("unexpected enum export value for FASTEST display behavior");
}

if (ReleaseInfoSchema.typeName !== "chill.v4.ReleaseInfo") {
  throw new Error(`unexpected release info type name: ${ReleaseInfoSchema.typeName}`);
}

if (SearchResultSchema.typeName !== "chill.v4.SearchResult") {
  throw new Error(`unexpected search result type name: ${SearchResultSchema.typeName}`);
}

if (UserSettingsSchema.typeName !== "chill.v4.UserSettings") {
  throw new Error(`unexpected user settings type name: ${UserSettingsSchema.typeName}`);
}

if (SearchSettingsSchema.typeName !== "chill.v4.SearchSettings") {
  throw new Error(`unexpected search settings type name: ${SearchSettingsSchema.typeName}`);
}

if (CatalogSettingsSchema.typeName !== "chill.v4.CatalogSettings") {
  throw new Error(`unexpected catalog settings type name: ${CatalogSettingsSchema.typeName}`);
}

if (DownloadSettingsSchema.typeName !== "chill.v4.DownloadSettings") {
  throw new Error(`unexpected download settings type name: ${DownloadSettingsSchema.typeName}`);
}

if (CatalogOriginSchema.typeName !== "chill.v4.CatalogOrigin") {
  throw new Error(`unexpected catalog origin type name: ${CatalogOriginSchema.typeName}`);
}

if (AddTransferRequestSchema.typeName !== "chill.v4.AddTransferRequest") {
  throw new Error(`unexpected add transfer request type name: ${AddTransferRequestSchema.typeName}`);
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
if (toJson(ResolvePlaybackRequestSchema, request).fileId !== largeId) {
  throw new Error("playback file ID lost int64 precision");
}
for (const result of [
  { ready: { media: { url: "https://example.com/movie.mp4", expiryUnknown: true } } },
  { ready: { media: { url: "https://example.com/movie.mp4", expiresAt: "2030-01-01T00:00:00Z" } } },
  { pending: { reason: "PENDING_REASON_PROCESSING" } },
  { unavailable: { reason: "UNAVAILABLE_REASON_NOT_FOUND" } },
]) {
  const parsed = fromJson(ResolvePlaybackResponseSchema, result);
  const encoded = toJson(ResolvePlaybackResponseSchema, parsed);
  if (JSON.stringify(encoded) !== JSON.stringify(result)) {
    throw new Error("Playback result or expiry state changed on JSON round trip");
  }
}
if (!UserService.method.resolvePlayback) {
  throw new Error("UserService playback method is missing");
}

console.log("ts consumer import smoke passed");
