import {
  CatalogSettingsSchema,
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

console.log("ts consumer import smoke passed");
