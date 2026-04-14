package fdj

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	xhttp "github.com/gofast-pkg/http"
)

const (
	drawHeaderGameName           = "game_name"
	drawHeaderCurrent            = "current"
	drawHeaderSort               = "sort"
	drawHeaderToPlannedAt        = "to_planned_at"
	drawHeaderFromPlannedAt      = "from_planned_at"
	drawHeaderPageNumber         = "page"
	drawResponseHeaderTotalPage  = "total-count"
	drawResponseHeaderPageNumber = "page-number"
)

const (
	drawEndpointPath = "anonymous/service-draw-info/v3/draws"
	drawTimeFormat   = "2006-01-02T15:04:05.000-07:00"
)

// FDJ Games supported by the FDJ's API.
const (
	GameLoto           Game = "loto"
	GameSuperLoto      Game = "superloto"
	GameGrandLoto      Game = "grandloto"
	GameEuroMillions   Game = "euromillions"
	GameEuroMyMillions Game = "mymillion"
	GameKeno           Game = "keno"
	GameKeno2025       Game = "keno2025"
	GameCrescendo      Game = "crescendo"
	GameEuroDream      Game = "eurodreams"
	GameJoker          Game = "joker"
	GameAmigo          Game = "amigo"
)

// FDJ Games Extenal ID supported by the FDJ's API.
const (
	GameExternalIDLoto           GameExternalID = "draw-13"
	GameExternalIDSuperLoto      GameExternalID = "draw-14"
	GameExternalIDGrandLoto      GameExternalID = "draw-15"
	GameExternalIDEuroMillions   GameExternalID = "draw-20"
	GameExternalIDEuroMyMillions GameExternalID = "draw-30"
	GameExternalIDKeno           GameExternalID = "draw-26"
	GameExternalIDKeno2025       GameExternalID = "draw-28"
	GameExternalIDCrescendo      GameExternalID = "draw-25"
	GameExternalIDEuroDream      GameExternalID = "draw-35"
	GameExternalIDJoker          GameExternalID = "draw-11"
	GameExternalIDAmigo          GameExternalID = "draw-32"
)

// Sorting supported values.
const (
	SortASC  Sort = "planned_at:asc"
	SortDESC Sort = "planned_at:desc"
)

// Errors draw list.
var (
	ErrDrawHeaderResponse = errors.New("fails to parse the header data")
	ErrDrawBodyResponse   = errors.New("fails to parse body response")
	ErrDrawNoFilter       = errors.New("fails to call the draw endpoint without filter")
)

// Draw JSON data from the FDJ's API.
type Draw struct {
	ID                     string         `json:"id"`
	ExternalID             string         `json:"external_id"`
	GameExternalID         GameExternalID `json:"game_external_id"`
	GameVersion            int            `json:"game_version"`
	CDC                    int64          `json:"cdc"`
	ThemeID                string         `json:"theme_id"`
	Rolldown               bool           `json:"rolldown"`
	IsCurrent              bool           `json:"is_current"`
	PlannedAt              string         `json:"planned_at"`
	WageringEndsAt         string         `json:"wagering_ends_at"`
	ForclosesAt            string         `json:"forcloses_at"`
	GuaranteedAmounts      []Amount       `json:"guaranteed_amounts"`
	GuaranteedRafflePrizes []RafflePrize  `json:"guaranteed_raffle_prizes"`
	CycleNumber            string         `json:"cycle_number,omitempty"`
	Flowdown               bool           `json:"flowdown,omitempty"`
	ProcessedAt            string         `json:"processed_at,omitempty"`
	Videos                 []Video        `json:"videos,omitempty"`
}

// Amount JSON data from the FDJ's API.
type Amount struct {
	Value    int64  `json:"value"`
	Currency string `json:"currency"`
	Scale    int    `json:"scale"`
}

// RafflePrize JSON data from the FDJ's API.
type RafflePrize struct {
	Count           int     `json:"count"`
	Amount          int64   `json:"amount"`
	Currency        string  `json:"currency"`
	Scale           int     `json:"scale"`
	AnnuityAmount   *int64  `json:"annuity_amount"`
	AnnuityPeriod   *string `json:"annuity_period"`
	AnnuityDuration *int    `json:"annuity_duration"`
}

// Video JSON data from the FDJ's API.
type Video struct {
	ID              string `json:"id"`
	ProviderName    string `json:"provider_name"`
	Address         string `json:"address"`
	BroadcastableAt string `json:"broadcastable_at"`
}

// DrawFilter JSON data from the FDJ's API.
type DrawFilter struct {
	Games        []Game
	Current      bool
	CurrentIsSet bool
	From         time.Time
	To           time.Time
	Sort         Sort
	SortIsSet    bool
	PageNumber   int
}

// DrawData represent the returned data after calling the Draws method.
type DrawData struct {
	Draws       []*Draw
	CurrentPage int
	NumberPage  int
}

// Game type for the FDJ's API.
type Game string

// GameExternalID for the FDJ's API.
type GameExternalID string

// Sort type for the FDJ's API.
type Sort string

func (a *api) Draws(ctx context.Context, filter *DrawFilter) (*DrawData, error) {
	var err error
	var resp *http.Response

	if filter == nil {
		return nil, ErrDrawNoFilter
	}

	if resp, err = a.drawsDo(ctx, filter); err != nil {
		return nil, err
	}
	defer func() { err = errors.Join(err, resp.Body.Close()) }()

	if resp.StatusCode == http.StatusNoContent {
		return &DrawData{}, nil
	}

	var draws []*Draw
	if err = json.NewDecoder(resp.Body).Decode(&draws); err != nil {
		return nil, errors.Join(ErrDrawBodyResponse, err)
	}

	value := resp.Header.Get(drawResponseHeaderPageNumber)
	currentPage, err := strconv.Atoi(value)
	if err != nil {
		return nil, errors.Join(
			ErrDrawHeaderResponse,
			fmt.Errorf("%s: %w", drawResponseHeaderPageNumber, err))
	}

	value = resp.Header.Get(drawResponseHeaderTotalPage)
	totalPage, err := strconv.Atoi(value)
	if err != nil {
		return nil, errors.Join(
			ErrDrawHeaderResponse,
			fmt.Errorf("%s: %w", drawResponseHeaderTotalPage, err))
	}

	return &DrawData{
		Draws:       draws,
		CurrentPage: currentPage,
		NumberPage:  totalPage,
	}, nil
}

func (a *api) drawsDo(ctx context.Context, filter *DrawFilter) (*http.Response, error) {
	var err error
	var req *http.Request
	var resp *http.Response

	if ctx == nil {
		return nil, ErrNilContext
	}

	url := fmt.Sprintf("%s/%s", defaultBaseURL, drawEndpointPath)
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, errors.Join(err, ErrHTTPRequest)
	}

	setDrawsURLParam(req, filter)

	if resp, err = a.httpClient.Do(req); err != nil {
		return nil, errors.Join(ErrHTTPClient, err)
	}
	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("%s: %w",
			xhttp.UnexpectedResponseBody(resp),
			ErrHTTPResponse)
	}

	return resp, nil
}

func setDrawsURLParam(req *http.Request, filter *DrawFilter) {
	values := req.URL.Query()
	for _, g := range filter.Games {
		values.Add(drawHeaderGameName, string(g))
	}
	if filter.CurrentIsSet {
		values.Set(drawHeaderCurrent, fmt.Sprintf("%t", filter.Current))
	}
	if !filter.From.IsZero() {
		values.Set(drawHeaderFromPlannedAt, filter.From.Format(drawTimeFormat))
	}
	if !filter.To.IsZero() {
		values.Set(drawHeaderToPlannedAt, filter.To.Format(drawTimeFormat))
	}
	if filter.SortIsSet {
		values.Set(drawHeaderSort, string(filter.Sort))
	}
	if filter.PageNumber > 0 {
		values.Set(drawHeaderPageNumber, fmt.Sprintf("%d", filter.PageNumber))
	}

	req.URL.RawQuery = values.Encode()
}
