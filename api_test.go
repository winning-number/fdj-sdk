package lotto

import (
	"archive/zip"
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	xcsv "github.com/gofast-pkg/csv"
	"github.com/gofast-pkg/http/testify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winning-number/fdj-sdk-lotto/draw"
)

func TestNewAPI(t *testing.T) {
	t.Run("Should return an API", func(t *testing.T) {
		api := NewAPI()
		assert.NotNil(t, api)
	})
}

func TestAPI_LoadSource(t *testing.T) {
	httpTester := testify.NewHTTPClient(t)
	ctx := context.Background()

	t.Run("Should nothing to do because source list is empty or nil", func(t *testing.T) {
		var err error

		a := NewAPI()
		err = a.LoadSource(context.Background(), []SourceInfo{})
		if assert.NoError(t, err) {
			assert.Empty(t, a.(*api).draws)
		}

		err = a.LoadSource(context.Background(), nil)
		if assert.NoError(t, err) {
			assert.Empty(t, a.(*api).draws)
		}
	})
	t.Run("Should return an error because context is nil", func(t *testing.T) {
		a := NewAPI()
		//nolint:staticcheck // ignore SA1012 because we want to test nil context
		err := a.LoadSource(nil, []SourceInfo{})

		if assert.Error(t, err) {
			require.ErrorIs(t, err, ErrNoContext)
		}
	})
	t.Run("Should return an error because new request has failed", func(t *testing.T) {
		var err error

		a := &api{httpClient: httpTester.Client()}
		err = a.LoadSource(ctx, []SourceInfo{{APIBase: APIBasePath, Name: 9999, APIPath: "%"}})

		if assert.Error(t, err) {
			require.ErrorIs(t, err, ErrInvalidHTTPRequest)
		}
	})
	t.Run("Should return an error because httpclient.Do() has failed", func(t *testing.T) {
		var req *http.Request
		var err error

		// setup http client tester
		func() {
			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Err:             ErrHTTPClientDo,
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		err = a.LoadSource(ctx, []SourceInfo{{APIBase: APIBasePath, Name: 9999, APIPath: "test"}})

		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			require.ErrorIs(t, err, ErrHTTPClientDo)
		}
	})
	t.Run("Should return an error because client return a bad status code", func(t *testing.T) {
		var req *http.Request
		var err error

		// setup http client tester
		func() {
			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Response:        &http.Response{StatusCode: http.StatusBadRequest},
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		err = a.LoadSource(ctx, []SourceInfo{{APIBase: APIBasePath, Name: 9999, APIPath: "test"}})

		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			require.ErrorIs(t, err, ErrInvalidResponseHTTP)
		}
	})
	t.Run("Should return an error because zip new reader has failed", func(t *testing.T) {
		var req *http.Request
		var err error

		// setup http client tester
		func() {
			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Response:        &http.Response{StatusCode: http.StatusOK},
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		err = a.LoadSource(ctx, []SourceInfo{{APIBase: APIBasePath, Name: 9999, APIPath: "test"}})

		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			require.ErrorIs(t, err, zip.ErrFormat)
		}
	})
	t.Run("Should return an error because zip file is empty", func(t *testing.T) {
		var req *http.Request
		var err error
		var f *os.File

		// setup http client tester
		func() {
			if f, err = os.Open("draw/testdata/emptyfile.zip"); err != nil {
				t.Error(err)
			}
			// do not close file because the request reader should will do it
			// defer f.Close()

			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Response:        &http.Response{StatusCode: http.StatusOK, Body: f},
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		err = a.LoadSource(ctx, []SourceInfo{{
			APIBase: APIBasePath,
			Name:    9999,
			APIPath: "test",
			Version: draw.V1,
			Type:    draw.LottoType,
		}})
		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			assert.ErrorContains(t, err, "error to create a csvutil.NewDecoder")
		}
	})
	t.Run(
		"Should return an error because csvReader.DecodeWithDecoder fail to parse the csv content",
		func(t *testing.T) {
			var req *http.Request
			var err error
			var f *os.File

			// setup http client tester
			func() {
				if f, err = os.Open("draw/testdata/classic-loto-v1.csv.zip"); err != nil {
					t.Error(err)
				}
				// do not close file because the request reader should will do it
				// defer f.Close()

				url := fmt.Sprintf("%s/%s", APIBasePath, "test")
				if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
					t.Error(err)
				}

				httpTester.AddCall(testify.Caller{
					ExpectedRequest: req,
					Response:        &http.Response{StatusCode: http.StatusOK, Body: f},
				})
			}()

			a := &api{httpClient: httpTester.Client()}
			// use the v2 draw version to have a complex object compare to the csv content
			err = a.LoadSource(ctx, []SourceInfo{{
				APIBase: APIBasePath,
				Name:    9999,
				APIPath: "test",
				Version: draw.V2,
				Type:    draw.LottoType,
			}})

			if assert.Error(t, err) {
				httpTester.ExpectedCalls()
				require.ErrorIs(t, err, xcsv.ErrOBJDecode)
			}
		})
	t.Run(
		"Should return an error because version and type doesn't match the csv content",
		func(t *testing.T) {
			var req *http.Request
			var err error
			var f *os.File

			// setup http client tester
			func() {
				if f, err = os.Open("draw/testdata/classic-loto-v1.csv.zip"); err != nil {
					t.Error(err)
				}
				// do not close file because the request reader should will do it
				// defer f.Close()

				url := fmt.Sprintf("%s/%s", APIBasePath, "test")
				if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
					t.Error(err)
				}

				httpTester.AddCall(testify.Caller{
					ExpectedRequest: req,
					Response:        &http.Response{StatusCode: http.StatusOK, Body: f},
				})
			}()

			a := &api{httpClient: httpTester.Client()}
			// use the v0 draw version to an incomplete object compare to the csv content
			err = a.LoadSource(ctx, []SourceInfo{{
				APIBase: APIBasePath,
				Name:    9999,
				APIPath: "test",
				Version: draw.V0,
				Type:    draw.LottoType,
			}})

			if assert.Error(t, err) {
				httpTester.ExpectedCalls()
				require.ErrorIs(t, err, ErrInvalidCSVType)
			}
		})
	t.Run("Should load the source", func(t *testing.T) {
		var req *http.Request
		var err error
		var f *os.File

		// setup http client tester
		func() {
			if f, err = os.Open("draw/testdata/classic-loto-v1.csv.zip"); err != nil {
				t.Error(err)
			}
			// do not close file because the request reader should will do it
			// defer f.Close()

			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Response:        &http.Response{StatusCode: http.StatusOK, Body: f},
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		err = a.LoadSource(ctx, []SourceInfo{{
			APIBase: APIBasePath,
			Name:    9999,
			APIPath: "test",
			Version: draw.V1,
			Type:    draw.LottoType,
		}})

		if assert.NoError(t, err) {
			httpTester.ExpectedCalls()
			assert.Len(t, a.draws, 2)
		}
	})
}

func TestAPI_DownloadSource(t *testing.T) {
	httpTester := testify.NewHTTPClient(t)
	ctx := context.Background()
	t.Run("Should return an error because the context is nil", func(t *testing.T) {
		a := &api{}
		//nolint:staticcheck // ignore SA1012 because we want to test nil context
		err := a.DownloadSource(nil, "/tmp", []SourceInfo{{}})
		if assert.Error(t, err) {
			require.ErrorIs(t, err, ErrNoContext)
		}
	})
	t.Run("Should nothing to do because source list is empty or nil", func(t *testing.T) {
		var err error

		a := NewAPI()
		err = a.DownloadSource(context.Background(), "./tmp", []SourceInfo{})
		if assert.NoError(t, err) {
			assert.Empty(t, a.(*api).draws)
		}

		err = a.DownloadSource(context.Background(), "./tmp", nil)
		if assert.NoError(t, err) {
			assert.Empty(t, a.(*api).draws)
		}
	})
	t.Run("Should return an error because new request has failed", func(t *testing.T) {
		var err error

		a := &api{httpClient: httpTester.Client()}
		err = a.DownloadSource(ctx, "./tmp", []SourceInfo{{Name: 9999, APIPath: "%"}})

		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			require.ErrorIs(t, err, ErrInvalidHTTPRequest)
		}
	})
	t.Run("Should return an error because httpclient.Do() has failed", func(t *testing.T) {
		var req *http.Request
		var err error

		// setup http client tester
		func() {
			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Err:             ErrHTTPClientDo,
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		err = a.DownloadSource(ctx, "./tmp", []SourceInfo{{APIBase: APIBasePath, Name: 9999, APIPath: "test"}})

		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			require.ErrorIs(t, err, ErrHTTPClientDo)
		}
	})
	t.Run("Should return an error because client return a bad status code", func(t *testing.T) {
		var req *http.Request
		var err error

		// setup http client tester
		func() {
			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Response:        &http.Response{StatusCode: http.StatusBadRequest},
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		err = a.DownloadSource(ctx, "./tmp", []SourceInfo{{APIBase: APIBasePath, Name: 9999, APIPath: "test"}})

		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			require.ErrorIs(t, err, ErrInvalidResponseHTTP)
		}
	})
	t.Run("Should return an error because zip new reader has failed", func(t *testing.T) {
		var req *http.Request
		var err error

		// setup http client tester
		func() {
			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Response:        &http.Response{StatusCode: http.StatusOK},
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		err = a.DownloadSource(ctx, "./tmp", []SourceInfo{{APIBase: APIBasePath, Name: 9999, APIPath: "test"}})

		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			require.ErrorIs(t, err, zip.ErrFormat)
		}
	})
	t.Run("Should return an error because path to create file do not exist", func(t *testing.T) {
		var req *http.Request
		var err error
		var f *os.File

		// setup http client tester
		func() {
			if f, err = os.Open("draw/testdata/emptyfile.zip"); err != nil {
				t.Error(err)
			}
			// do not close file because the request reader should will do it
			// defer f.Close()

			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Response:        &http.Response{StatusCode: http.StatusOK, Body: f},
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		err = a.DownloadSource(ctx, "./tmp", []SourceInfo{{
			APIBase: APIBasePath,
			Name:    9999,
			APIPath: "test",
			Version: draw.V1,
			Type:    draw.LottoType,
		}})
		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			require.ErrorIs(t, err, os.ErrNotExist)
		}
	})
	t.Run("Should download the source", func(t *testing.T) {
		var req *http.Request
		var err error
		var f *os.File

		if err = os.Mkdir("./tmp", 0750); err != nil {
			t.Error(err)
		}
		defer func() {
			require.NoError(t, os.RemoveAll("./tmp"))
		}()

		// setup http client tester
		func() {
			if f, err = os.Open("draw/testdata/classic-loto-v1.csv.zip"); err != nil {
				t.Error(err)
			}
			// do not close file because the request reader should will do it
			// defer f.Close()

			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Response:        &http.Response{StatusCode: http.StatusOK, Body: f},
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		err = a.DownloadSource(ctx, "./tmp", []SourceInfo{{
			APIBase: APIBasePath,
			Name:    9999,
			APIPath: "test",
			Version: draw.V1,
			Type:    draw.LottoType,
		}})

		if assert.NoError(t, err) {
			httpTester.ExpectedCalls()
			assert.FileExists(t, "./tmp/classic-loto-v1.csv")
		}
	})
}

func TestAPI_SourceUpdatedAt(t *testing.T) {
	httpTester := testify.NewHTTPClient(t)
	ctx := context.Background()

	t.Run("Should return an error because the context is nil", func(t *testing.T) {
		a := &api{}

		//nolint:staticcheck // ignore SA1012 because we want to test nil context
		updatedAt, err := a.SourceUpdatedAt(nil, SourceInfo{})
		if assert.Error(t, err) {
			require.ErrorIs(t, err, ErrNoContext)
			assert.Empty(t, updatedAt)
		}
	})
	t.Run("Should nothing to do because source list is empty or nil", func(t *testing.T) {
		var err error

		a := NewAPI()
		updatedAt, err := a.SourceUpdatedAt(ctx, SourceInfo{APIPath: "%"})
		if assert.Error(t, err) {
			require.ErrorIs(t, err, ErrInvalidHTTPRequest)
			assert.Empty(t, updatedAt)
		}
	})
	t.Run("Should return an error because httpclient.Do() has failed", func(t *testing.T) {
		var req *http.Request
		var err error
		var updatedAt time.Time

		// setup http client tester
		func() {
			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodHead, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Err:             ErrHTTPClientDo,
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		updatedAt, err = a.SourceUpdatedAt(ctx, SourceInfo{APIBase: APIBasePath, Name: 9999, APIPath: "test"})

		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			require.ErrorIs(t, err, ErrHTTPClientDo)
			assert.Empty(t, updatedAt)
		}
	})
	t.Run("Should return an error because client return a bad status code", func(t *testing.T) {
		var req *http.Request
		var err error
		var updatedAt time.Time

		// setup http client tester
		func() {
			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodHead, url, nil); err != nil {
				t.Error(err)
			}

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Response:        &http.Response{StatusCode: http.StatusBadRequest},
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		updatedAt, err = a.SourceUpdatedAt(ctx, SourceInfo{APIBase: APIBasePath, Name: 9999, APIPath: "test"})

		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			require.ErrorIs(t, err, ErrInvalidResponseHTTP)
			assert.Empty(t, updatedAt)
		}
	})
	t.Run("Should return an error because time.Parse fail with a bad format", func(t *testing.T) {
		var req *http.Request
		var err error
		var updatedAt time.Time

		// setup http client tester
		func() {
			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodHead, url, nil); err != nil {
				t.Error(err)
			}
			header := make(http.Header)
			header.Set("Last-Modified", "bad format")

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Response:        &http.Response{StatusCode: http.StatusOK, Header: header},
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		updatedAt, err = a.SourceUpdatedAt(ctx, SourceInfo{APIBase: APIBasePath, Name: 9999, APIPath: "test"})

		if assert.Error(t, err) {
			httpTester.ExpectedCalls()
			assert.ErrorContains(t, err, "cannot parse \"bad format\" as")
			assert.Empty(t, updatedAt)
		}
	})
	t.Run("Should get the source update at time", func(t *testing.T) {
		var req *http.Request
		var err error
		var updatedAt time.Time

		// setup http client tester
		func() {
			url := fmt.Sprintf("%s/%s", APIBasePath, "test")
			if req, err = http.NewRequestWithContext(ctx, http.MethodHead, url, nil); err != nil {
				t.Error(err)
			}
			header := make(http.Header)
			header.Set("Last-Modified", "Mon, 02 Jan 2006 15:04:05 GMT")

			httpTester.AddCall(testify.Caller{
				ExpectedRequest: req,
				Response:        &http.Response{StatusCode: http.StatusOK, Header: header},
			})
		}()

		a := &api{httpClient: httpTester.Client()}
		updatedAt, err = a.SourceUpdatedAt(ctx, SourceInfo{APIBase: APIBasePath, Name: 9999, APIPath: "test"})

		if assert.NoError(t, err) {
			httpTester.ExpectedCalls()
			assert.NotEmpty(t, updatedAt)
			assert.Equal(t, time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), updatedAt.UTC())
		}
	})
}

func TestAPI_LoadFile(t *testing.T) {
	t.Run("Should receive file not found", func(t *testing.T) {
		a := &api{}

		err := a.LoadFile("", SourceInfo{})
		if assert.Error(t, err) {
			require.ErrorIs(t, err, os.ErrNotExist)
		}
	})
	t.Run("Should receive an error because parseCSV has failed", func(t *testing.T) {
		a := &api{}

		err := a.LoadFile(
			"./draw/testdata/classic-loto-v1.csv",
			SourceInfo{Version: draw.V0, Type: draw.XmasLottoType})
		if assert.Error(t, err) {
			require.ErrorIs(t, err, ErrInvalidCSVType)
		}
	})
	t.Run("Should load the file", func(t *testing.T) {
		a := &api{}

		err := a.LoadFile(
			"./draw/testdata/classic-loto-v1.csv",
			SourceInfo{Version: draw.V1, Type: draw.LottoType})
		if assert.NoError(t, err) {
			assert.Len(t, a.draws, 2)
		}
	})
}

func TestAPI_NDraws(t *testing.T) {
	t.Run("Should not return the draws count", func(t *testing.T) {
		a := &api{draws: []draw.Draw{{
			Metadata: draw.Metadata{
				DrawType: draw.LottoType,
			},
		}, {
			Metadata: draw.Metadata{
				DrawType: draw.LottoType,
			},
		}}}

		assert.EqualValues(t, 0, a.NDraws(Filter{}))
	})
	t.Run("Should return the draws count", func(t *testing.T) {
		a := &api{draws: []draw.Draw{{
			Metadata: draw.Metadata{
				DrawType: draw.LottoType,
			},
		}, {
			Metadata: draw.Metadata{
				DrawType: draw.LottoType,
			},
		}}}

		assert.EqualValues(t, 2, a.NDraws(Filter{
			SuperLotto:   true,
			GrandLotto:   true,
			XmasLotto:    true,
			ClassicLotto: true,
			OldLotto:     true,
		}))
	})
}

func TestAPI_Draws(t *testing.T) {
	t.Run("Should not return the draws list", func(t *testing.T) {
		a := &api{draws: []draw.Draw{{
			Metadata: draw.Metadata{
				DrawType: draw.LottoType,
			},
		}, {
			Metadata: draw.Metadata{
				DrawType: draw.LottoType,
			},
		}}}

		assert.Equal(t, []draw.Draw{}, a.Draws(Filter{}, draw.OrderNone))
	})
	t.Run("Should return the draws list", func(t *testing.T) {
		a := &api{draws: []draw.Draw{{
			Metadata: draw.Metadata{
				DrawType: draw.LottoType,
			},
		}, {
			Metadata: draw.Metadata{
				DrawType: draw.LottoType,
			},
		}}}

		a.Draws(Filter{
			SuperLotto:   true,
			GrandLotto:   true,
			XmasLotto:    true,
			ClassicLotto: true,
			OldLotto:     true,
		}, draw.OrderNone)
		assert.Len(t, a.draws, 2)
	})
}

func TestAPI_Reset(t *testing.T) {
	t.Run("Should reset the draws list", func(t *testing.T) {
		a := &api{draws: []draw.Draw{{
			Metadata: draw.Metadata{
				DrawType: draw.LottoType,
			},
		}, {
			Metadata: draw.Metadata{
				DrawType: draw.LottoType,
			},
		}}}

		a.Reset()
		assert.Empty(t, a.draws)
	})
}
