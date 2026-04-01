package loto

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winning-number/fdj-sdk/loto/csv"
	"github.com/winning-number/fdj-sdk/loto/testdata"
	"github.com/winning-number/fdj-sdk/model"
)

func ExpectedMetadata(t *testing.T, change *Metadata, isOld bool) Metadata {
	t.Helper()

	data := Metadata{
		FDJID:          testdata.ValidFDJIdentifier,
		ID:             fmt.Sprintf("%s-1-", testdata.ValidFDJIdentifier),
		Date:           time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
		ForclosureDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
		Version:        LotteryV0,
		Type:           NewLotto,
		Day:            model.Monday,
		Currency:       model.EUR,
		TirageOrder:    1,
		IsOldType:      isOld,
	}

	if change == nil {
		return data
	}

	if change.FDJID != "" {
		data.FDJID = change.FDJID
	}
	if change.ID != "" {
		data.ID = change.ID
	}
	if !change.Date.IsZero() {
		data.Date = change.Date
	}
	if !change.ForclosureDate.IsZero() {
		data.ForclosureDate = change.ForclosureDate
	}
	if change.Version != "" {
		data.Version = change.Version
	}
	if change.Type != "" {
		data.Type = change.Type
	}
	if change.Day != "" {
		data.Day = change.Day
	}
	if change.Currency != "" {
		data.Currency = change.Currency
	}
	if change.TirageOrder != 0 {
		data.TirageOrder = change.TirageOrder
	}

	return data
}

func TestCSVConverter(t *testing.T) {
	t.Run("should return an error", func(t *testing.T) {
		versions := map[LotteryVersion]struct {
			nilPtr      any
			newInstance func() any
		}{
			LotteryV0: {
				nilPtr:      (*csv.Version0)(nil),
				newInstance: func() any { return &csv.Version0{} },
			},
			LotteryV1: {
				nilPtr:      (*csv.Version1)(nil),
				newInstance: func() any { return &csv.Version1{} },
			},
			LotteryV2: {
				nilPtr:      (*csv.Version2)(nil),
				newInstance: func() any { return &csv.Version2{} },
			},
			LotteryV3: {
				nilPtr:      (*csv.Version3)(nil),
				newInstance: func() any { return &csv.Version3{} },
			},
			LotteryV4: {
				nilPtr:      (*csv.Version4)(nil),
				newInstance: func() any { return &csv.Version4{} },
			},
		}

		for version, data := range versions {
			t.Run(string(version), func(t *testing.T) {
				testCases := map[string]struct {
					common       csv.Common
					injectCommon bool
					expectedErr  error
					instance     any
				}{
					"when the instance is unexpected for the version": {
						injectCommon: false,
						expectedErr:  ErrInvalidOBJInstanceConverter,
						instance:     testdata.NotValidValue,
					},
					"when the instance is nil for the version": {
						injectCommon: false,
						expectedErr:  ErrNilOBJInstance,
						instance:     data.nilPtr,
					},
					"when the instance metadata has an invalid date": {
						common: csv.Common{
							Date: testdata.NotValidValue,
						},
						injectCommon: true,
						expectedErr:  ErrInvalidMetadata,
						instance:     data.newInstance(),
					},
					"when the instance metadata has an empty date": {
						common: csv.Common{
							Date: testdata.EmptyValue,
						},
						injectCommon: true,
						expectedErr:  ErrEmptyDate,
						instance:     data.newInstance(),
					},
					"when the instance metadata has an invalid foreclosure date": {
						common: csv.Common{
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.NotValidValue,
						},
						injectCommon: true,
						expectedErr:  ErrInvalidMetadata,
						instance:     data.newInstance(),
					},
					"when the instance metadata has an invalid day": {
						common: csv.Common{
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            testdata.NotValidValue,
						},
						injectCommon: true,
						expectedErr:  ErrUnknownDay,
						instance:     data.newInstance(),
					},
					"when the instance metadata has an empty day": {
						common: csv.Common{
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            testdata.EmptyValue,
						},
						injectCommon: true,
						expectedErr:  ErrEmptyDay,
						instance:     data.newInstance(),
					},
					"when the instance metadata has an invalid currency": {
						common: csv.Common{
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortMonday,
							Currency:       testdata.NotValidValue,
						},
						injectCommon: true,
						expectedErr:  ErrUnknownCurrency,
						instance:     data.newInstance(),
					},
					"when the instance metadata has an empty currency": {
						common: csv.Common{
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortMonday,
							Currency:       testdata.EmptyValue,
						},
						injectCommon: true,
						expectedErr:  ErrEmptyCurrency,
						instance:     data.newInstance(),
					},
				}

				for name, tt := range testCases {
					t.Run(name, func(t *testing.T) {
						if tt.injectCommon {
							testdata.CopyCommon(t, tt.instance, tt.common)
						}

						draw, err := CSVConverter(NewLotto, version, tt.instance)
						require.ErrorIs(t, err, tt.expectedErr)
						assert.Nil(t, draw)
					})
				}
			})
		}
	})
	t.Run("should return an error with an unknown version", func(t *testing.T) {
		draw, err := CSVConverter(NewLotto, testdata.NotValidValue, &csv.Version0{})
		require.ErrorIs(t, err, ErrUnknownLotteryVersion)
		assert.Nil(t, draw)
	})
	t.Run("should convert the metadata", func(t *testing.T) {
		versions := map[LotteryVersion]struct {
			newInstance func() any
			old         bool
		}{
			LotteryV0: {
				newInstance: func() any { return &csv.Version0{} },
				old:         true,
			},
			LotteryV1: {
				newInstance: func() any {
					return &csv.Version1{
						Tirage: 1,
					}
				},
				old: true,
			},
			LotteryV2: {
				newInstance: func() any { return &csv.Version2{} },
				old:         false,
			},
			LotteryV3: {
				newInstance: func() any { return &csv.Version3{} },
				old:         false,
			},
			LotteryV4: {
				newInstance: func() any { return &csv.Version4{} },
				old:         false,
			},
		}

		for version, data := range versions {
			t.Run(string(version), func(t *testing.T) {
				testCases := map[string]struct {
					common         csv.Common
					lottoType      LottoType
					metadataChange Metadata
				}{
					"with type lotto SuperLotto": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortMonday,
							Currency:       csv.Euro,
						},
						lottoType: SuperLotto,
						metadataChange: Metadata{
							Type:    SuperLotto,
							Version: version,
						},
					},
					"with type lotto GrandLotto": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortMonday,
							Currency:       csv.Euro,
						},
						lottoType: GrandLotto,
						metadataChange: Metadata{
							Type:    GrandLotto,
							Version: version,
						},
					},
					"with type lotto XmasLotto": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortMonday,
							Currency:       csv.Euro,
						},
						lottoType: XmasLotto,
						metadataChange: Metadata{
							Type:    XmasLotto,
							Version: version,
						},
					},
					"with type lotto NewLotto": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortMonday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Type:    NewLotto,
							Version: version,
						},
					},
					"with an other format date": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           "02/01/2006",
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortMonday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Version: version,
						},
					},
					"without foreclosure date": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: "04/01/2006",
							Day:            csv.ShortMonday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Version:        version,
							ForclosureDate: time.Date(2006, 1, 4, 0, 0, 0, 0, time.UTC),
						},
					},
					"with currency in FRANC": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortMonday,
							Currency:       csv.Franc,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Version:  version,
							Currency: model.FRANC,
						},
					},
					"with long Monday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.Monday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Monday,
							Version: version,
						},
					},
					"with long tuesday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.Tuesday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Tuesday,
							Version: version,
						},
					},
					"with short tuesday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortTuesday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Tuesday,
							Version: version,
						},
					},
					"with long wednesday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.Wednesday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Wednesday,
							Version: version,
						},
					},
					"with short wednesday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortWednesday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Wednesday,
							Version: version,
						},
					},
					"with long thursday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.Thursday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Thursday,
							Version: version,
						},
					},
					"with short thursday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortThursday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Thursday,
							Version: version,
						},
					},
					"with long friday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.Friday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Friday,
							Version: version,
						},
					},
					"with short friday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortFriday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Friday,
							Version: version,
						},
					},
					"with long saturday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.Saturday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Saturday,
							Version: version,
						},
					},
					"with short saturday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortSaturday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Saturday,
							Version: version,
						},
					},
					"with long sunday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.Sunday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Sunday,
							Version: version,
						},
					},
					"with short sunday value": {
						common: csv.Common{
							ID:             testdata.ValidFDJIdentifier,
							Date:           testdata.ValidDateValue,
							ForclosureDate: testdata.ValidDateValue,
							Day:            csv.ShortSunday,
							Currency:       csv.Euro,
						},
						lottoType: NewLotto,
						metadataChange: Metadata{
							Day:     model.Sunday,
							Version: version,
						},
					},
				}

				for name, tt := range testCases {
					t.Run(name, func(t *testing.T) {
						instance := data.newInstance()
						testdata.CopyCommon(t, instance, tt.common)

						draw, err := CSVConverter(tt.lottoType, version, instance)
						require.NoError(t, err)
						assert.Equal(t, ExpectedMetadata(t, &tt.metadataChange, data.old), draw.Metadata)
					})
				}
			})
		}
	})
	t.Run("should convert the CSV version 0", func(t *testing.T) {
		draw, err := CSVConverter(NewLotto, LotteryV0, testdata.CSVVersion0Data())
		require.NoError(t, err)

		assert.Equal(t, &LotteryDraw{
			Metadata: Metadata{
				FDJID:          testdata.ValidFDJIdentifier,
				ID:             fmt.Sprintf("%s-1-%s", testdata.ValidFDJIdentifier, testdata.ValidWinOrder),
				Date:           time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
				ForclosureDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
				Version:        LotteryV0,
				Type:           NewLotto,
				Day:            model.Monday,
				Currency:       model.EUR,
				TirageOrder:    1,
				IsOldType:      true,
			},
			FirstDraw: Draw{
				NumBalls: 7,
				Balls: []int32{
					testdata.ValidBall1,
					testdata.ValidBall2,
					testdata.ValidBall3,
					testdata.ValidBall4,
					testdata.ValidBall5,
					testdata.ValidBall6,
					testdata.ValidAdditionalBall,
				},
				HasLuckyBall: false,
				NumRanks:     7,
				WinStats: map[WinRank]WinStat{
					Rank1: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank2: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank3: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank4: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank5: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank6: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank7: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
				},
			},
			HasSecondDraw: false,
			JokerPlus:     testdata.ValidJokerPlus,
			HasJokerV1:    false,
		}, draw)
	})
	t.Run("should convert the CSV version 1", func(t *testing.T) {
		draw, err := CSVConverter(NewLotto, LotteryV1, testdata.CSVVersion1Data())
		require.NoError(t, err)

		assert.Equal(t, &LotteryDraw{
			Metadata: Metadata{
				FDJID: testdata.ValidFDJIdentifier,
				ID: fmt.Sprintf("%s-%d-%s",
					testdata.ValidFDJIdentifier,
					testdata.ValidTirageOrder,
					testdata.ValidWinOrder),
				Date:           time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
				ForclosureDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
				Version:        LotteryV1,
				Type:           NewLotto,
				Day:            model.Monday,
				Currency:       model.EUR,
				TirageOrder:    testdata.ValidTirageOrder,
				IsOldType:      true,
			},
			FirstDraw: Draw{
				NumBalls: 7,
				Balls: []int32{
					testdata.ValidBall1,
					testdata.ValidBall2,
					testdata.ValidBall3,
					testdata.ValidBall4,
					testdata.ValidBall5,
					testdata.ValidBall6,
					testdata.ValidAdditionalBall,
				},
				HasLuckyBall: false,
				NumRanks:     7,
				WinStats: map[WinRank]WinStat{
					Rank1: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank2: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank3: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank4: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank5: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank6: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank7: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
				},
			},
			HasSecondDraw: false,
			JokerPlus:     testdata.ValidJokerPlus,
			JokerV1:       testdata.ValidJokerV1,
			HasJokerV1:    true,
		}, draw)
	})
	t.Run("should convert the CSV version 2", func(t *testing.T) {
		draw, err := CSVConverter(NewLotto, LotteryV2, testdata.CSVVersion2Data())
		require.NoError(t, err)

		assert.Equal(t, &LotteryDraw{
			Metadata: Metadata{
				FDJID: testdata.ValidFDJIdentifier,
				ID: fmt.Sprintf("%s-1-%s",
					testdata.ValidFDJIdentifier,
					testdata.ValidWinOrder),
				Date:           time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
				ForclosureDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
				Version:        LotteryV2,
				Type:           NewLotto,
				Day:            model.Monday,
				Currency:       model.EUR,
				TirageOrder:    1,
				IsOldType:      false,
			},
			FirstDraw: Draw{
				NumBalls: 5,
				Balls: []int32{
					testdata.ValidBall1,
					testdata.ValidBall2,
					testdata.ValidBall3,
					testdata.ValidBall4,
					testdata.ValidBall5,
				},
				LuckyBall:    testdata.ValidLuckyBall,
				HasLuckyBall: true,
				NumRanks:     6,
				WinStats: map[WinRank]WinStat{
					Rank1: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank2: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank3: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank4: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank5: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank6: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
				},
			},
			HasSecondDraw: false,
			JokerPlus:     testdata.ValidJokerPlus,
			HasJokerV1:    false,
		}, draw)
	})
	t.Run("should convert the CSV version 3", func(t *testing.T) {
		draw, err := CSVConverter(NewLotto, LotteryV3, testdata.CSVVersion3Data())
		require.NoError(t, err)

		assert.Equal(t, &LotteryDraw{
			Metadata: Metadata{
				FDJID: testdata.ValidFDJIdentifier,
				ID: fmt.Sprintf("%s-1-%s",
					testdata.ValidFDJIdentifier,
					testdata.ValidWinOrder),
				Date:           time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
				ForclosureDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
				Version:        LotteryV3,
				Type:           NewLotto,
				Day:            model.Monday,
				Currency:       model.EUR,
				TirageOrder:    1,
				IsOldType:      false,
			},
			FirstDraw: Draw{
				NumBalls: 5,
				Balls: []int32{
					testdata.ValidBall1,
					testdata.ValidBall2,
					testdata.ValidBall3,
					testdata.ValidBall4,
					testdata.ValidBall5,
				},
				LuckyBall:    testdata.ValidLuckyBall,
				HasLuckyBall: true,
				NumRanks:     9,
				WinStats: map[WinRank]WinStat{
					Rank1: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank2: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank3: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank4: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank5: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank6: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank7: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank8: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank9: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
				},
			},
			WinningCode: WinningCode{
				Codes: []string{
					"423456",
					"345678",
					"0928474",
				},
				Price:    float64(testdata.ValidGainRank),
				NumCodes: testdata.ValidNumberWinCode,
			},
			HasSecondDraw: false,
			JokerPlus:     testdata.ValidJokerPlus,
			HasJokerV1:    false,
		}, draw)
	})
	t.Run("should convert the CSV version 4", func(t *testing.T) {
		draw, err := CSVConverter(NewLotto, LotteryV4, testdata.CSVVersion4Data())
		require.NoError(t, err)

		assert.Equal(t, &LotteryDraw{
			Metadata: Metadata{
				FDJID: testdata.ValidFDJIdentifier,
				ID: fmt.Sprintf("%s-1-%s",
					testdata.ValidFDJIdentifier,
					testdata.ValidWinOrder),
				Date:           time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
				ForclosureDate: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
				Version:        LotteryV4,
				Type:           NewLotto,
				Day:            model.Monday,
				Currency:       model.EUR,
				TirageOrder:    1,
				IsOldType:      false,
			},
			FirstDraw: Draw{
				NumBalls: 5,
				Balls: []int32{
					testdata.ValidBall1,
					testdata.ValidBall2,
					testdata.ValidBall3,
					testdata.ValidBall4,
					testdata.ValidBall5,
				},
				LuckyBall:    testdata.ValidLuckyBall,
				HasLuckyBall: true,
				NumRanks:     9,
				WinStats: map[WinRank]WinStat{
					Rank1: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank2: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank3: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank4: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank5: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank6: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank7: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank8: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank9: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
				},
			},
			WinningCode: WinningCode{
				Codes: []string{
					"423456",
					"345678",
					"0928474",
				},
				Price:    float64(testdata.ValidGainRank),
				NumCodes: testdata.ValidNumberWinCode,
			},
			HasSecondDraw: true,
			SecondDraw: Draw{
				NumBalls: 5,
				Balls: []int32{
					testdata.ValidBall1,
					testdata.ValidBall2,
					testdata.ValidBall3,
					testdata.ValidBall4,
					testdata.ValidBall5,
				},
				HasLuckyBall: false,
				NumRanks:     4,
				WinStats: map[WinRank]WinStat{
					Rank1: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank2: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank3: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
					Rank4: {
						Rate:   float64(testdata.ValidGainRank),
						Number: testdata.ValidWinnerRank,
					},
				},
			},
			JokerPlus:  testdata.ValidJokerPlus,
			HasJokerV1: false,
		}, draw)
	})
}
