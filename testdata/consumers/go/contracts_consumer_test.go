package consumer_test

import (
	"context"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"connectrpc.com/connect/v2/connectproto"
	chillv4 "github.com/chill-institute/chill-contracts/v3/gen/go/chill/v4"
	chillv4connect "github.com/chill-institute/chill-contracts/v3/gen/go/chill/v4/chillv4connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestGeneratedContractsCompileForConsumer(t *testing.T) {
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

	var _ chillv4connect.UserServiceClient = chillv4connect.NewUserServiceClient(connect.NewClient(connecthttp.NewTransport(http.DefaultClient, "https://example.com")))
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

type unresolvedBingeTitles struct {
	chillv4connect.UnimplementedUserServiceHandler
}

func (unresolvedBingeTitles) CreateBinge(context.Context, *chillv4.CreateBingeRequest) (*chillv4.CreateBingeResponse, error) {
	detail, err := connectproto.NewErrorDetail(&chillv4.BingeItemsUnresolved{ImdbIds: []string{"tt0000001", "tt0000002"}})
	if err != nil {
		return nil, err
	}
	return nil, connect.NewError(connect.CodeInvalidArgument, "binge items could not be resolved").WithDetail(detail)
}

func TestUnresolvedBingeTitlesReachClientAsErrorDetail(t *testing.T) {
	server := connect.NewServer()
	chillv4connect.RegisterUserServiceHandler(server, unresolvedBingeTitles{})
	mux := http.NewServeMux()
	connecthttp.Mount(mux, server)
	httpServer := httptest.NewServer(mux)
	t.Cleanup(httpServer.Close)
	client := chillv4connect.NewUserServiceClient(connect.NewClient(connecthttp.NewTransport(httpServer.Client(), httpServer.URL)))

	_, err := client.CreateBinge(t.Context(), &chillv4.CreateBingeRequest{Title: "Heists"})

	var connectErr *connect.Error
	if !errors.As(err, &connectErr) || connectErr.Code() != connect.CodeInvalidArgument {
		t.Fatalf("expected invalid_argument, got %v", err)
	}
	var unresolved []string
	for _, detail := range connectErr.Details() {
		message, err := connectproto.UnmarshalErrorDetail(detail)
		if err != nil {
			t.Fatal(err)
		}
		if items, ok := message.(*chillv4.BingeItemsUnresolved); ok {
			unresolved = items.GetImdbIds()
		}
	}
	if !slices.Equal(unresolved, []string{"tt0000001", "tt0000002"}) {
		t.Fatalf("unresolved IMDb IDs = %v", unresolved)
	}
}
