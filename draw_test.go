package fdj

import (
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
)

const (
	drawInvalidNumber string = "invalid-number"
	drawValidNumber   string = "1"
)

var testBody = []byte(`[{
    "rolldown": false,
    "flowdown": false,
    "external_id": "26038",
    "cycle_number": "8",
    "game_external_id": "draw-13",
    "game_version": 0,
    "cdc": 9495,
    "wagering_ends_at": "2026-03-30T20:15:00.000+02:00",
    "theme_id": "0",
    "planned_at": "2026-03-30T20:55:00.000+02:00",
    "processed_at": "2026-03-30T20:21:57.000+02:00",
    "forcloses_at": "2026-06-29T00:00:00.000+02:00",
    "guaranteed_amounts": [
        {
            "value": 900000000,
            "currency": "EUR",
            "scale": 2
        },
        {
            "value": 1073985680,
            "currency": "XPF",
            "scale": 0
        }
    ],
    "guaranteed_raffle_prizes": [
        {
            "count": 10,
            "amount": 2000000,
            "currency": "EUR",
            "scale": 2,
            "annuity_amount": null,
            "annuity_period": null,
            "annuity_duration": null
        },
        {
            "count": 10,
            "amount": 2500000,
            "currency": "XPF",
            "scale": 0,
            "annuity_amount": null,
            "annuity_period": null,
            "annuity_duration": null
        }
    ],
    "is_current": false,
    "videos": [
        {
            "provider_name": "youtube",
            "address": "7yBL0EUL1Fg",
            "id": "7f704099-1ce6-4b5e-b4fa-14f2c84d6235",
            "broadcastable_at": "2026-03-30T19:05:05.000+00:00"
        }
    ],
    "id": "1003"
}]`)

func TestAPI_Draws(t *testing.T) {
	t.Run("should return an error without filter", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)
		api := &api{httpClient: httpClient.Client()}

		data, err := api.Draws(context.Background(), nil)
		require.ErrorIs(t, err, ErrDrawNoFilter)
		assert.Nil(t, data)
	})
	t.Run("should return an error with a nil context", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)
		api := &api{httpClient: httpClient.Client()}

		data, err := api.Draws(context.Context(nil), &DrawFilter{})
		require.ErrorIs(t, err, ErrNilContext)
		assert.Nil(t, data)
	})
	t.Run("should return an error when calling the fdj's api", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)

		ctx := context.Background()
		api := &api{httpClient: httpClient.Client()}

		url := fmt.Sprintf("%s/%s", defaultBaseURL, drawEndpointPath)
		expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

		httpClient.AddCall(testify.Caller{
			ExpectedRequest: expectedRequest,
			Err:             assert.AnError,
		})

		data, err := api.Draws(ctx, &DrawFilter{})
		require.ErrorIs(t, err, ErrHTTPClient)
		assert.Nil(t, data)
	})
	t.Run("should return an error to decode the json", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)

		ctx := context.Background()
		api := &api{httpClient: httpClient.Client()}

		url := fmt.Sprintf("%s/%s", defaultBaseURL, drawEndpointPath)
		expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

		expectedResponse := &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("1234"))),
		}
		httpClient.AddCall(testify.Caller{
			ExpectedRequest: expectedRequest,
			Response:        expectedResponse,
		})

		data, err := api.Draws(ctx, &DrawFilter{})
		require.ErrorIs(t, err, ErrDrawBodyResponse)
		assert.Nil(t, data)
	})
	t.Run("should return when parsing the response header", func(t *testing.T) {
		testCases := map[string]struct {
			headerTotalPageValue  string
			headerTotalPageIsSet  bool
			headerPageNumberValue string
			headerPageNumberIsSet bool
			errHas                string
		}{
			"with missing total page": {
				headerTotalPageIsSet:  false,
				headerTotalPageValue:  "",
				headerPageNumberIsSet: true,
				headerPageNumberValue: drawValidNumber,
				errHas:                drawResponseHeaderTotalPage,
			},
			"with invalid total page": {
				headerTotalPageIsSet:  true,
				headerTotalPageValue:  drawInvalidNumber,
				headerPageNumberIsSet: true,
				headerPageNumberValue: drawValidNumber,
				errHas:                drawResponseHeaderTotalPage,
			},
			"with missing page number": {
				headerTotalPageIsSet:  true,
				headerTotalPageValue:  drawValidNumber,
				headerPageNumberIsSet: false,
				headerPageNumberValue: "",
				errHas:                drawResponseHeaderPageNumber,
			},
			"with invalid page number": {
				headerTotalPageIsSet:  true,
				headerTotalPageValue:  drawValidNumber,
				headerPageNumberIsSet: true,
				headerPageNumberValue: drawInvalidNumber,
				errHas:                drawResponseHeaderPageNumber,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				httpClient := testify.NewHTTPClient(t)

				ctx := context.Background()
				api := &api{httpClient: httpClient.Client()}

				url := fmt.Sprintf("%s/%s", defaultBaseURL, drawEndpointPath)
				expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)

				header := http.Header{}
				if tt.headerPageNumberIsSet {
					header.Add(drawResponseHeaderPageNumber, tt.headerPageNumberValue)
				}
				if tt.headerTotalPageIsSet {
					header.Add(drawResponseHeaderTotalPage, tt.headerTotalPageValue)
				}

				expectedResponse := &http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(testBody)),
					Header:     header,
				}
				httpClient.AddCall(testify.Caller{
					ExpectedRequest: expectedRequest,
					Response:        expectedResponse,
				})

				data, err := api.Draws(ctx, &DrawFilter{})
				require.ErrorIs(t, err, ErrDrawHeaderResponse)
				require.ErrorContains(t, err, tt.errHas)
				assert.Nil(t, data)
			})
		}
	})
	t.Run("should return unexpected status code", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)

		ctx := context.Background()
		api := &api{httpClient: httpClient.Client()}

		url := fmt.Sprintf("%s/%s", defaultBaseURL, drawEndpointPath)
		expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
		values := expectedRequest.URL.Query()
		values.Add(drawHeaderGameName, "loto")

		expectedRequest.URL.RawQuery = values.Encode()

		expectedResponse := &http.Response{
			Status:     http.StatusText(http.StatusBadRequest),
			StatusCode: http.StatusBadRequest,
			Body:       http.NoBody,
		}
		httpClient.AddCall(testify.Caller{
			ExpectedRequest: expectedRequest,
			Response:        expectedResponse,
		})

		data, err := api.Draws(ctx, &DrawFilter{
			Games: []Game{GameLoto},
		})
		require.ErrorIs(t, err, ErrHTTPResponse)
		assert.Empty(t, data)
	})
	t.Run("should be done with content", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)

		ctx := context.Background()
		api := &api{httpClient: httpClient.Client()}

		url := fmt.Sprintf("%s/%s", defaultBaseURL, drawEndpointPath)
		expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
		values := expectedRequest.URL.Query()
		values.Add(drawHeaderGameName, "loto")
		values.Add(drawHeaderCurrent, "true")
		values.Add(drawHeaderFromPlannedAt, "2026-04-11T20:55:00.000+02:00")
		values.Add(drawHeaderToPlannedAt, "2026-04-11T20:55:00.000+02:00")
		values.Add(drawHeaderSort, "planned_at:asc")
		values.Add(drawHeaderPageNumber, "1")

		expectedRequest.URL.RawQuery = values.Encode()

		header := http.Header{}
		header.Add(drawResponseHeaderPageNumber, drawValidNumber)
		header.Add(drawResponseHeaderTotalPage, drawValidNumber)

		expectedResponse := &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(testBody)),
			Header:     header,
		}
		httpClient.AddCall(testify.Caller{
			ExpectedRequest: expectedRequest,
			Response:        expectedResponse,
		})

		expectedData := &DrawData{
			Draws: []*Draw{
				{
					ID:             "1003",
					ExternalID:     "26038",
					GameExternalID: GameExternalIDLoto,
					GameVersion:    0,
					CDC:            9495,
					ThemeID:        "0",
					Rolldown:       false,
					IsCurrent:      false,
					PlannedAt:      "2026-03-30T20:55:00.000+02:00",
					WageringEndsAt: "2026-03-30T20:15:00.000+02:00",
					ForclosesAt:    "2026-06-29T00:00:00.000+02:00",
					GuaranteedAmounts: []Amount{
						{
							Value:    900000000,
							Currency: "EUR",
							Scale:    2,
						},
						{
							Value:    1073985680,
							Currency: "XPF",
							Scale:    0,
						},
					},
					GuaranteedRafflePrizes: []RafflePrize{
						{
							Count:    10,
							Amount:   2000000,
							Currency: "EUR",
							Scale:    2,
						},
						{
							Count:    10,
							Amount:   2500000,
							Currency: "XPF",
							Scale:    0,
						},
					},
					CycleNumber: "8",
					Flowdown:    false,
					ProcessedAt: "2026-03-30T20:21:57.000+02:00",
					Videos: []Video{
						{
							ID:              "7f704099-1ce6-4b5e-b4fa-14f2c84d6235",
							ProviderName:    "youtube",
							Address:         "7yBL0EUL1Fg",
							BroadcastableAt: "2026-03-30T19:05:05.000+00:00",
						},
					},
				},
			},
			CurrentPage: 1,
			NumberPage:  1,
		}

		from, err := time.Parse(drawTimeFormat, "2026-04-11T20:55:00.000+02:00")
		require.NoError(t, err)
		to, err := time.Parse(drawTimeFormat, "2026-04-11T20:55:00.000+02:00")
		require.NoError(t, err)

		data, err := api.Draws(ctx, &DrawFilter{
			Games:        []Game{GameLoto},
			Current:      true,
			CurrentIsSet: true,
			From:         from,
			To:           to,
			Sort:         SortASC,
			SortIsSet:    true,
			PageNumber:   1,
		})
		require.NoError(t, err)
		assert.Equal(t, expectedData, data)
	})
	t.Run("should be done without content", func(t *testing.T) {
		httpClient := testify.NewHTTPClient(t)

		ctx := context.Background()
		api := &api{httpClient: httpClient.Client()}

		url := fmt.Sprintf("%s/%s", defaultBaseURL, drawEndpointPath)
		expectedRequest := httptest.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
		values := expectedRequest.URL.Query()
		values.Add(drawHeaderGameName, "loto")
		values.Add(drawHeaderCurrent, "true")
		values.Add(drawHeaderFromPlannedAt, "2026-04-11T20:55:00.000+02:00")
		values.Add(drawHeaderToPlannedAt, "2026-04-11T20:55:00.000+02:00")
		values.Add(drawHeaderSort, "planned_at:asc")
		values.Add(drawHeaderPageNumber, "1")

		expectedRequest.URL.RawQuery = values.Encode()

		expectedResponse := &http.Response{
			Status:     http.StatusText(http.StatusNoContent),
			StatusCode: http.StatusNoContent,
			Body:       http.NoBody,
		}
		httpClient.AddCall(testify.Caller{
			ExpectedRequest: expectedRequest,
			Response:        expectedResponse,
		})

		from, err := time.Parse(drawTimeFormat, "2026-04-11T20:55:00.000+02:00")
		require.NoError(t, err)
		to, err := time.Parse(drawTimeFormat, "2026-04-11T20:55:00.000+02:00")
		require.NoError(t, err)

		data, err := api.Draws(ctx, &DrawFilter{
			Games:        []Game{GameLoto},
			Current:      true,
			CurrentIsSet: true,
			From:         from,
			To:           to,
			Sort:         SortASC,
			SortIsSet:    true,
			PageNumber:   1,
		})
		require.NoError(t, err)
		assert.Empty(t, data)
	})
}
