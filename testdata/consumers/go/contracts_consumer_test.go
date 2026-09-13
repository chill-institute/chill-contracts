package consumer_test

import (
	"google.golang.org/protobuf/encoding/protojson"
	"math"
	"net/http"
	"testing"

	chillv4 "github.com/chill-institute/chill-contracts/v2/gen/go/chill/v4"
	chillv4connect "github.com/chill-institute/chill-contracts/v2/gen/go/chill/v4/chillv4connect"
	"google.golang.org/protobuf/proto"
)

func TestGeneratedContractsCompileForConsumer(t *testing.T) {
	request := &chillv4.GetUserProfileRequest{}
	if request == nil {
		t.Fatal("expected generated request type")
	}

	folderID := int64(123)
	settings := &chillv4.UserSettings{
		Search: &chillv4.SearchSettings{
			FilterNastyResults:          true,
			SearchResultDisplayBehavior: chillv4.SearchResultDisplayBehavior_SEARCH_RESULT_DISPLAY_BEHAVIOR_FASTEST,
			SortBy:                      chillv4.SortBy_SORT_BY_SEEDERS,
			SortDirection:               chillv4.SortDirection_SORT_DIRECTION_DESC,
		},
		Catalog: &chillv4.CatalogSettings{
			MoviesSource:  chillv4.MoviesSource_MOVIES_SOURCE_IMDB_MOVIEMETER,
			TvShowsSource: chillv4.TVShowsSource_TV_SHOWS_SOURCE_NETFLIX,
		},
		Download: &chillv4.DownloadSettings{
			FolderId: &folderID,
		},
	}
	if settings.GetDownload().GetFolderId() != folderID {
		t.Fatal("expected nested settings type to expose download folder")
	}
	transfer := &chillv4.AddTransferRequest{
		Url: "magnet:?xt=urn:btih:example",
		CatalogOrigin: &chillv4.CatalogOrigin{
			Catalog: &chillv4.CatalogOrigin_MoviesSource{
				MoviesSource: chillv4.MoviesSource_MOVIES_SOURCE_TRAKT,
			},
		},
	}
	if transfer.GetCatalogOrigin().GetMoviesSource() != chillv4.MoviesSource_MOVIES_SOURCE_TRAKT {
		t.Fatal("expected transfer request to expose catalog origin")
	}

	client := chillv4connect.NewUserServiceClient(http.DefaultClient, "https://example.com")
	if client == nil {
		t.Fatal("expected generated connect client")
	}
}

func TestCatalogSortPresenceSurvivesWireRoundTrip(t *testing.T) {
	for _, sort := range []*chillv4.CatalogSort{
		nil,
		chillv4.CatalogSort_CATALOG_SORT_UNSPECIFIED.Enum(),
		chillv4.CatalogSort_CATALOG_SORT_POPULARITY.Enum(),
		chillv4.CatalogSort_CATALOG_SORT_RELEASE_DATE_ASC.Enum(),
	} {
		original := &chillv4.CatalogSettings{Sort: sort}
		encoded, err := proto.Marshal(original)
		if err != nil {
			t.Fatal(err)
		}
		decoded := &chillv4.CatalogSettings{}
		if err := proto.Unmarshal(encoded, decoded); err != nil {
			t.Fatal(err)
		}
		if !proto.Equal(original, decoded) || (decoded.Sort == nil) != (sort == nil) {
			t.Fatalf("sort presence changed: original=%v decoded=%v", original, decoded)
		}
	}
}

func TestPlaybackContractsPreserveIDsAndStates(t *testing.T) {
	request := &chillv4.ResolvePlaybackRequest{FileId: math.MaxInt64}
	encoded, err := protojson.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	var decoded chillv4.ResolvePlaybackRequest
	if err := protojson.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.GetFileId() != math.MaxInt64 {
		t.Fatal("file ID lost precision")
	}
	for _, payload := range []string{
		`{"ready":{"media":{"url":"https://example.com/movie.mp4","expiryUnknown":true}}}`,
		`{"ready":{"media":{"url":"https://example.com/movie.mp4","expiresAt":"2030-01-01T00:00:00Z"}}}`,
		`{"pending":{"reason":"PENDING_REASON_PROCESSING"}}`,
		`{"unavailable":{"reason":"UNAVAILABLE_REASON_NOT_FOUND"}}`,
	} {
		var response chillv4.ResolvePlaybackResponse
		if err := protojson.Unmarshal([]byte(payload), &response); err != nil {
			t.Fatal(err)
		}
		if response.GetResult() == nil {
			t.Fatal("missing playback result")
		}
		if ready := response.GetReady(); ready != nil && ready.GetMedia().GetExpiry() == nil {
			t.Fatal("missing media expiry state")
		}
	}
}
