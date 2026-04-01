// Package fdj provide a FDJ's sdk to call the history endpoint to get game's data.
package fdj

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"time"

	xhttp "github.com/gofast-pkg/http"
	"github.com/gofast-pkg/zip"
	"github.com/winning-number/fdj-sdk/v2/source"
)

const (
	// APIBaseURL is the base URL for the FDJ API.
	APIBaseURL = "https://www.sto.api.fdj.fr/anonymous/service-draw-info/v3/documentations"

	headerContentDisposition = "content-disposition"
	headerContentLength      = "content-length"
	headerDateReceived       = "date-received"
	headerFileName           = "filename"
)

// Errors types specific to the APIV3 implementation.
var (
	// HTTP errors.
	ErrNilContext   = errors.New("context is nil")
	ErrHTTPRequest  = errors.New("http request error")
	ErrHTTPClient   = errors.New("http client error")
	ErrHTTPResponse = errors.New("http response error")

	// Metadata errors.
	ErrMissingContentDisposition = errors.New("Content-Disposition header is missing")
	ErrMissingContentLength      = errors.New("missing Content-Length header")
	ErrMissingDateReceived       = errors.New("missing date-received header")
	ErrInvalidContentDisposition = errors.New("invalid Content-Disposition header")
	ErrMissingFilename           = errors.New("filename parameter is missing in Content-Disposition header")
	ErrInvalidContentLength      = errors.New("invalid Content-Length header")
	ErrInvalidDateReceived       = errors.New("invalid date-received header")
	ErrFailedToParseMetadata     = errors.New("failed to parse metadata from response headers")

	// Zip reader errors.
	ErrCreateZipReader = errors.New("failed to create zip reader")
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
}

type api struct {
	domain     string
	httpClient *http.Client
}

// NewAPI returns a new API instance.
func NewAPI() API {
	return &api{
		domain:     APIBaseURL,
		httpClient: xhttp.NewClient(),
	}
}

func (a *api) DownloadHistory(ctx context.Context, identifier string) (source.Source, error) {
	var ret source.Source
	var resp *http.Response
	var err error

	if resp, err = a.do(ctx, identifier); err != nil {
		return source.Source{}, err
	}
	defer func() { err = errors.Join(err, resp.Body.Close()) }()

	ret.Metadata, err = buildMetadata(resp.Header, identifier)
	if err != nil {
		return source.Source{}, errors.Join(err, ErrFailedToParseMetadata)
	}

	ret.Data, err = zip.NewReader(resp.Body)
	if err != nil {
		return source.Source{}, errors.Join(err, ErrCreateZipReader)
	}

	return ret, nil
}

func (a *api) HistoryMetadata(ctx context.Context, identifier string) (source.Metadata, error) {
	var ret source.Metadata
	var resp *http.Response
	var err error

	if resp, err = a.do(ctx, identifier); err != nil {
		return source.Metadata{}, err
	}
	defer func() { err = errors.Join(err, resp.Body.Close()) }()

	ret, err = buildMetadata(resp.Header, identifier)
	if err != nil {
		return source.Metadata{}, errors.Join(err, ErrFailedToParseMetadata)
	}

	return ret, nil
}

func (a *api) do(ctx context.Context, identifier string) (*http.Response, error) {
	var err error
	var resp *http.Response
	var req *http.Request

	if ctx == nil {
		return nil, ErrNilContext
	}

	url := fmt.Sprintf("%s/%s", a.domain, identifier)
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody); err != nil {
		return nil, errors.Join(err, ErrHTTPRequest)
	}
	if resp, err = a.httpClient.Do(req); err != nil {
		return nil, errors.Join(ErrHTTPClient, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w",
			xhttp.UnexpectedResponseBody(resp),
			ErrHTTPResponse)
	}

	return resp, nil
}

func buildMetadata(header http.Header, identifier string) (source.Metadata, error) {
	var ret source.Metadata

	ret.Identifier = identifier

	contentDisposition := header.Get(headerContentDisposition)
	if contentDisposition == "" {
		return source.Metadata{}, ErrMissingContentDisposition
	}
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return source.Metadata{}, errors.Join(err, ErrInvalidContentDisposition)
	}

	var ok bool
	if ret.FileName, ok = params[headerFileName]; !ok {
		return source.Metadata{}, ErrMissingFilename
	}

	contentLength := header.Get(headerContentLength)
	if contentLength == "" {
		return source.Metadata{}, ErrMissingContentLength
	}
	ret.Size, err = strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return source.Metadata{}, errors.Join(err, ErrInvalidContentLength)
	}

	dateReceived := header.Get(headerDateReceived)
	if dateReceived == "" {
		return source.Metadata{}, ErrMissingDateReceived
	}
	ret.RequestedAt, err = time.Parse(time.RFC3339Nano, dateReceived)
	if err != nil {
		return source.Metadata{}, errors.Join(err, ErrInvalidDateReceived)
	}

	return ret, nil
}
