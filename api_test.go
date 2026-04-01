package fdj

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofast-pkg/http/testify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winning-number/fdj-sdk/v2/loto"
	"github.com/winning-number/fdj-sdk/v2/source"
)

const (
	invalidHeaderValue      = "invalid-value"
	contentDispositionValue = "attachment; filename=file.zip;"
	contentSizeValue        = "1024"
	dateReceiveValue        = "2026-04-01T07:40:44.514496Z"
	filename                = "file.zip"
)

func TestNewAPI(t *testing.T) {
	t.Run("Should return an API", func(t *testing.T) {
		api := NewAPI()
		assert.NotNil(t, api)
	})
}

func TestAPI_DownloadHistory(t *testing.T) {
	t.Run("should return an error with a nil context", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)
		api := &api{domain: APIBaseURL, httpClient: httpClient.Client()}

		history, err := api.DownloadHistory(context.Context(nil), string(loto.Loto201911UUID))
		require.ErrorIs(t, err, ErrNilContext)
		assert.Empty(t, history)
	})
	t.Run("should return an error when create new request", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)
		api := &api{domain: "%%invalid-url%%", httpClient: httpClient.Client()}

		history, err := api.DownloadHistory(context.Background(), string(loto.Loto201911UUID))
		require.ErrorIs(t, err, ErrHTTPRequest)
		assert.Empty(t, history)
	})
	t.Run("should return an error when calling the fdj's api", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)

		ctx := context.Background()
		api := &api{domain: APIBaseURL, httpClient: httpClient.Client()}

		identifier := string(loto.Loto201911UUID)
		url := fmt.Sprintf("%s/%s", api.domain, identifier)
		expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

		httpClient.AddCall(testify.Caller{
			ExpectedRequest: expectedRequest,
			Err:             assert.AnError,
		})

		history, err := api.DownloadHistory(ctx, identifier)
		require.ErrorIs(t, err, ErrHTTPClient)
		assert.Empty(t, history)
	})
	t.Run("should return an error when parsing the header", func(t *testing.T) {
		testCases := map[string]struct {
			header      http.Header
			expectedErr error
		}{
			"when content-disposition is missing": {
				header:      http.Header{},
				expectedErr: ErrMissingContentDisposition,
			},
			"when content-disposition is invalid": {
				header: func() http.Header {
					header := http.Header{}
					header.Set(headerContentDisposition, `something; "thing=any"`)

					return header
				}(),
				expectedErr: ErrInvalidContentDisposition,
			},
			"when filename is missing": {
				header: func() http.Header {
					header := http.Header{}
					header.Add(headerContentDisposition, "media-type; param1=value1;")

					return header
				}(),
				expectedErr: ErrMissingFilename,
			},
			"when content-length is missing": {
				header: func() http.Header {
					header := http.Header{}
					header.Add(headerContentDisposition, contentDispositionValue)

					return header
				}(),
				expectedErr: ErrMissingContentLength,
			},
			"when content-length is invalid": {
				header: func() http.Header {
					header := http.Header{}
					header.Add(headerContentDisposition, contentDispositionValue)
					header.Set(headerContentLength, invalidHeaderValue)

					return header
				}(),
				expectedErr: ErrInvalidContentLength,
			},
			"when date-received is missing": {
				header: func() http.Header {
					header := http.Header{}
					header.Add(headerContentDisposition, contentDispositionValue)
					header.Set(headerContentLength, contentSizeValue)

					return header
				}(),
				expectedErr: ErrMissingDateReceived,
			},
			"when date-received is invalid": {
				header: func() http.Header {
					header := http.Header{}
					header.Add(headerContentDisposition, contentDispositionValue)
					header.Set(headerContentLength, contentSizeValue)
					header.Add(headerDateReceived, invalidHeaderValue)

					return header
				}(),
				expectedErr: ErrInvalidDateReceived,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				httpClient := testify.NewHTTPClient(t)

				ctx := context.Background()
				api := &api{domain: APIBaseURL, httpClient: httpClient.Client()}

				identifier := string(loto.Loto201911UUID)
				url := fmt.Sprintf("%s/%s", api.domain, identifier)
				expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

				httpClient.AddCall(testify.Caller{
					Response: &http.Response{
						Status:     http.StatusText(http.StatusOK),
						StatusCode: http.StatusOK,
						Header:     tt.header,
					},
					ExpectedRequest: expectedRequest,
					Err:             nil,
				})

				history, err := api.DownloadHistory(ctx, identifier)
				require.ErrorIs(t, err, tt.expectedErr)
				assert.Empty(t, history)
			})
		}
	})
	t.Run("should return an error when creating the zip reader", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)

		ctx := context.Background()
		api := &api{domain: APIBaseURL, httpClient: httpClient.Client()}

		identifier := string(loto.Loto201911UUID)
		url := fmt.Sprintf("%s/%s", api.domain, identifier)
		expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

		header := http.Header{}
		header.Add(headerContentDisposition, fmt.Sprintf("attachment; filename=%s;", filename))
		header.Add(headerContentLength, contentSizeValue)
		header.Add(headerDateReceived, dateReceiveValue)

		httpClient.AddCall(testify.Caller{
			Response: &http.Response{
				Status:     http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Header:     header,
			},
			ExpectedRequest: expectedRequest,
			Err:             nil,
		})

		history, err := api.DownloadHistory(ctx, identifier)
		require.ErrorIs(t, err, ErrCreateZipReader)
		assert.Empty(t, history)
	})
	t.Run("should be ok", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)

		ctx := context.Background()
		api := &api{domain: APIBaseURL, httpClient: httpClient.Client()}

		identifier := string(loto.Loto201911UUID)
		url := fmt.Sprintf("%s/%s", api.domain, identifier)
		expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

		header := http.Header{}
		header.Add(headerContentDisposition, contentDispositionValue)
		header.Add(headerContentLength, contentSizeValue)
		header.Add(headerDateReceived, dateReceiveValue)

		buf := new(bytes.Buffer)
		zw := zip.NewWriter(buf)
		f, err := zw.Create(filename)
		require.NoError(t, err)
		_, err = f.Write([]byte("something"))
		require.NoError(t, err)
		require.NoError(t, zw.Close())

		httpClient.AddCall(testify.Caller{
			Response: &http.Response{
				Status:     http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Header:     header,
				Body:       io.NopCloser(buf),
			},
			ExpectedRequest: expectedRequest,
			Err:             nil,
		})

		expected := source.Metadata{
			Identifier:  string(loto.Loto201911UUID),
			Size:        1024,
			RequestedAt: time.Date(2026, time.April, 1, 7, 40, 44, 514496000, time.UTC),
			FileName:    filename,
		}

		history, err := api.DownloadHistory(ctx, identifier)
		require.NoError(t, err)
		assert.Equal(t, expected, history.Metadata)
		assert.NotNil(t, history.Data)
	})
}

func TestAPI_HistoryMetada(t *testing.T) {
	t.Run("should return an error with a nil context", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)
		api := &api{domain: APIBaseURL, httpClient: httpClient.Client()}

		meta, err := api.HistoryMetadata(context.Context(nil), string(loto.Loto201911UUID))
		require.ErrorIs(t, err, ErrNilContext)
		assert.Empty(t, meta)
	})
	t.Run("should return an error when create new request", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)
		api := &api{domain: "%%invalid-url%%", httpClient: httpClient.Client()}

		meta, err := api.HistoryMetadata(context.Background(), string(loto.Loto201911UUID))
		require.ErrorIs(t, err, ErrHTTPRequest)
		assert.Empty(t, meta)
	})
	t.Run("should return an error when calling the fdj's api", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)

		ctx := context.Background()
		api := &api{domain: APIBaseURL, httpClient: httpClient.Client()}

		identifier := string(loto.Loto201911UUID)
		url := fmt.Sprintf("%s/%s", api.domain, identifier)
		expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

		httpClient.AddCall(testify.Caller{
			ExpectedRequest: expectedRequest,
			Err:             assert.AnError,
		})

		meta, err := api.HistoryMetadata(ctx, string(loto.Loto201911UUID))
		require.ErrorIs(t, err, ErrHTTPClient)
		assert.Empty(t, meta)
	})
	t.Run("should return an error when parsing the header", func(t *testing.T) {
		testCases := map[string]struct {
			header      http.Header
			expectedErr error
		}{
			"when content-disposition is missing": {
				header:      http.Header{},
				expectedErr: ErrMissingContentDisposition,
			},
			"when content-disposition is invalid": {
				header: func() http.Header {
					header := http.Header{}
					header.Set(headerContentDisposition, `something; "thing=any"`)

					return header
				}(),
				expectedErr: ErrInvalidContentDisposition,
			},
			"when filename is missing": {
				header: func() http.Header {
					header := http.Header{}
					header.Add(headerContentDisposition, "media-type; param1=value1;")

					return header
				}(),
				expectedErr: ErrMissingFilename,
			},
			"when content-length is missing": {
				header: func() http.Header {
					header := http.Header{}
					header.Add(headerContentDisposition, contentDispositionValue)

					return header
				}(),
				expectedErr: ErrMissingContentLength,
			},
			"when content-length is invalid": {
				header: func() http.Header {
					header := http.Header{}
					header.Add(headerContentDisposition, contentDispositionValue)
					header.Set(headerContentLength, invalidHeaderValue)

					return header
				}(),
				expectedErr: ErrInvalidContentLength,
			},
			"when date-received is missing": {
				header: func() http.Header {
					header := http.Header{}
					header.Add(headerContentDisposition, contentDispositionValue)
					header.Set(headerContentLength, contentSizeValue)

					return header
				}(),
				expectedErr: ErrMissingDateReceived,
			},
			"when date-received is invalid": {
				header: func() http.Header {
					header := http.Header{}
					header.Add(headerContentDisposition, contentDispositionValue)
					header.Set(headerContentLength, contentSizeValue)
					header.Add(headerDateReceived, invalidHeaderValue)

					return header
				}(),
				expectedErr: ErrInvalidDateReceived,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				httpClient := testify.NewHTTPClient(t)

				ctx := context.Background()
				api := &api{domain: APIBaseURL, httpClient: httpClient.Client()}

				identifier := string(loto.Loto201911UUID)
				url := fmt.Sprintf("%s/%s", api.domain, identifier)
				expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

				httpClient.AddCall(testify.Caller{
					Response: &http.Response{
						Status:     http.StatusText(http.StatusOK),
						StatusCode: http.StatusOK,
						Header:     tt.header,
					},
					ExpectedRequest: expectedRequest,
					Err:             nil,
				})

				meta, err := api.HistoryMetadata(ctx, string(loto.Loto201911UUID))
				require.ErrorIs(t, err, tt.expectedErr)
				assert.Empty(t, meta)
			})
		}
	})
	t.Run("should be ok", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)

		ctx := context.Background()
		api := &api{domain: APIBaseURL, httpClient: httpClient.Client()}

		identifier := string(loto.Loto201911UUID)
		url := fmt.Sprintf("%s/%s", api.domain, identifier)
		expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

		header := http.Header{}
		header.Add(headerContentDisposition, fmt.Sprintf("attachment; filename=%s;", filename))
		header.Add(headerContentLength, contentSizeValue)
		header.Add(headerDateReceived, dateReceiveValue)

		httpClient.AddCall(testify.Caller{
			Response: &http.Response{
				Status:     http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Header:     header,
				Body:       http.NoBody,
			},
			ExpectedRequest: expectedRequest,
			Err:             nil,
		})

		expected := source.Metadata{
			Identifier:  string(loto.Loto201911UUID),
			Size:        1024,
			RequestedAt: time.Date(2026, time.April, 1, 7, 40, 44, 514496000, time.UTC),
			FileName:    filename,
		}

		meta, err := api.HistoryMetadata(ctx, string(loto.Loto201911UUID))
		require.NoError(t, err)
		assert.Equal(t, expected, meta)
	})
}
