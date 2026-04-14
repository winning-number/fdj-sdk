// Package fdj provide a FDJ's sdk to call the history endpoint to get game's data.
package fdj

import (
	"context"
	"errors"
	"net/http"

	xhttp "github.com/gofast-pkg/http"
	"github.com/winning-number/fdj-sdk/v2/source"
)

const (
	defaultBaseURL = "https://www.sto.api.fdj.fr"
)

// Errors types specific to the APIV3 implementation.
var (
	// HTTP errors.
	ErrNilContext   = errors.New("context is nil")
	ErrHTTPRequest  = errors.New("http request error")
	ErrHTTPClient   = errors.New("http client error")
	ErrHTTPResponse = errors.New("http response error")
	ErrHTTPDomain   = errors.New("http invalid domain")
)

// API is the interface to load the history source from the FDJ API.
// It can load the history source from the FDJ API.
//
//mockery:generate: true
type API interface {
	// DownloadHistory downloads the history source from the FDJ API.
	DownloadHistory(ctx context.Context, datasetID string) (source.Source, error)
	// HistoryMetadata fetches the metadata of the history source from the FDJ API.
	HistoryMetadata(ctx context.Context, datasetID string) (source.Metadata, error)
	// Draws return past or future draws bases depending on the filters selected on all FDJ games.
	// The return is paginated and the current and total page are set in the response.
	Draws(ctx context.Context, filter *DrawFilter) (*DrawData, error)
}

type api struct {
	httpClient *http.Client
}

// NewAPI returns a new API instance.
func NewAPI() API {
	return &api{
		httpClient: xhttp.NewClient(),
	}
}
