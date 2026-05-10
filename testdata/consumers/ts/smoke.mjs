import {
  ReleaseInfoSchema,
  SearchResultSchema,
  SearchResultDisplayBehavior,
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

console.log("ts consumer import smoke passed");
