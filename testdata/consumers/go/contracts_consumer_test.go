package consumer_test

import (
	"net/http"
	"testing"

	chillv4 "github.com/chill-institute/chill-contracts/v2/gen/go/chill/v4"
	chillv4connect "github.com/chill-institute/chill-contracts/v2/gen/go/chill/v4/chillv4connect"
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
