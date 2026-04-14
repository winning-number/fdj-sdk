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
	historyEndpointPath = "anonymous/service-draw-info/v3/documentations"

	headerContentDisposition = "content-disposition"
	headerContentLength      = "content-length"
	headerDateReceived       = "date-received"
	headerFileName           = "filename"
)

// Error history list.
var (
	ErrHistoryMissingHeader         = errors.New("header data is missing")
	ErrHistoryInvalidHeader         = errors.New("header data is invalid")
	ErrHistoryFailedToParseMetadata = errors.New("failed to parse metadata from response headers")
	ErrHistoryInvalidIdentifier     = errors.New("identifier must contain only digits and hyphens")

	ErrHistoryCreateZipReader = errors.New("failed to create zip reader")
)

func (a *api) DownloadHistory(ctx context.Context, identifier string) (source.Source, error) {
	var ret source.Source
	var resp *http.Response
	var err error

	if resp, err = a.historyDo(ctx, identifier); err != nil {
		return source.Source{}, err
	}
	defer func() { err = errors.Join(err, resp.Body.Close()) }()

	ret.Metadata, err = buildMetadata(resp.Header, identifier)
	if err != nil {
		return source.Source{}, errors.Join(err, ErrHistoryFailedToParseMetadata)
	}

	ret.Data, err = zip.NewReader(resp.Body)
	if err != nil {
		return source.Source{}, errors.Join(err, ErrHistoryCreateZipReader)
	}

	return ret, nil
}

func (a *api) HistoryMetadata(ctx context.Context, identifier string) (source.Metadata, error) {
	var ret source.Metadata
	var resp *http.Response
	var err error

	if resp, err = a.historyDo(ctx, identifier); err != nil {
		return source.Metadata{}, err
	}
	defer func() { err = errors.Join(err, resp.Body.Close()) }()

	ret, err = buildMetadata(resp.Header, identifier)
	if err != nil {
		return source.Metadata{}, errors.Join(err, ErrHistoryFailedToParseMetadata)
	}

	return ret, nil
}

func (a *api) historyDo(ctx context.Context, identifier string) (*http.Response, error) {
	var err error
	var req *http.Request
	var resp *http.Response

	if ctx == nil {
		return nil, ErrNilContext
	}

	url := fmt.Sprintf("%s/%s/%s", defaultBaseURL, historyEndpointPath, identifier)
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
		return source.Metadata{}, fmt.Errorf("%s: %w", headerContentDisposition, ErrHistoryMissingHeader)
	}
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return source.Metadata{}, errors.Join(
			err,
			fmt.Errorf("%s: %w", headerContentDisposition, ErrHistoryInvalidHeader))
	}

	var ok bool
	if ret.FileName, ok = params[headerFileName]; !ok {
		return source.Metadata{}, fmt.Errorf("%s: %w", headerFileName, ErrHistoryMissingHeader)
	}

	contentLength := header.Get(headerContentLength)
	if contentLength == "" {
		return source.Metadata{}, fmt.Errorf("%s: %w", headerContentLength, ErrHistoryMissingHeader)
	}
	ret.Size, err = strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return source.Metadata{}, errors.Join(
			err,
			fmt.Errorf("%s: %w", headerContentLength, ErrHistoryInvalidHeader))
	}

	dateReceived := header.Get(headerDateReceived)
	if dateReceived == "" {
		return source.Metadata{}, fmt.Errorf("%s: %w", headerDateReceived, ErrHistoryMissingHeader)
	}
	ret.RequestedAt, err = time.Parse(time.RFC3339Nano, dateReceived)
	if err != nil {
		return source.Metadata{}, errors.Join(
			err,
			fmt.Errorf("%s: %w", headerDateReceived, ErrHistoryInvalidHeader))
	}

	return ret, nil
}
