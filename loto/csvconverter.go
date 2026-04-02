//nolint:dupl // the conversion should catch some duplications but it's ok for exhaustive mapping.
package loto

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/winning-number/fdj-sdk/v2/loto/csv"
	"github.com/winning-number/fdj-sdk/v2/model"
)

// Public error used to convert csv data type to loto model type.
var (
	ErrInvalidOBJInstanceConverter = errors.New("invalid csv instance to conversion")
	ErrUnknownLotteryVersion       = errors.New("invalid lottery version")
	ErrNilOBJInstance              = errors.New("obj instance can't be nil")
	ErrInvalidMetadata             = errors.New("invalid metadata to conversion into LotteryDraw")
	ErrEmptyDate                   = errors.New("date is empty")
	ErrEmptyDay                    = errors.New("day is empty")
	ErrUnknownDay                  = errors.New("day is unknown")
	ErrEmptyCurrency               = errors.New("currency is empty")
	ErrUnknownCurrency             = errors.New("currency is unknown")
)

type metadataParameter struct {
	lottoType   LottoType
	version     LotteryVersion
	tirageOrder int32
	isOldType   bool
	winOrder    string
}

// CSVConverter to LotteryDraw model.
func CSVConverter(lottoType LottoType, version LotteryVersion, csvOBJ any) (*LotteryDraw, error) {
	switch version {
	case LotteryV0:
		obj, ok := csvOBJ.(*csv.Version0)
		if !ok {
			return nil, fmt.Errorf("version0: %w", ErrInvalidOBJInstanceConverter)
		}

		return convertCSVVersion0(lottoType, obj)
	case LotteryV1:
		obj, ok := csvOBJ.(*csv.Version1)
		if !ok {
			return nil, fmt.Errorf("version1: %w", ErrInvalidOBJInstanceConverter)
		}

		return convertCSVVersion1(lottoType, obj)
	case LotteryV2:
		obj, ok := csvOBJ.(*csv.Version2)
		if !ok {
			return nil, fmt.Errorf("version2: %w", ErrInvalidOBJInstanceConverter)
		}

		return convertCSVVersion2(lottoType, obj)
	case LotteryV3:
		obj, ok := csvOBJ.(*csv.Version3)
		if !ok {
			return nil, fmt.Errorf("version3: %w", ErrInvalidOBJInstanceConverter)
		}

		return convertCSVVersion3(lottoType, obj)
	case LotteryV4:
		obj, ok := csvOBJ.(*csv.Version4)
		if !ok {
			return nil, fmt.Errorf("version4: %w", ErrInvalidOBJInstanceConverter)
		}

		return convertCSVVersion4(lottoType, obj)
	default:
		return nil, fmt.Errorf("version %s: %w", version, ErrUnknownLotteryVersion)
	}
}

// Version 0 is the oldest version of lottery game.
// In this version 6 balls are picks + one additional ball. All between 1 and 49.
// To normalize LotteryDraw, all balls are includes in the draw (include the additional ball).
// There is no second draw.
// There is 7 winner ranks.
// There is a Joker+ and the Joker v1 is not implemented yet.
func convertCSVVersion0(lottoType LottoType, obj *csv.Version0) (*LotteryDraw, error) {
	if obj == nil {
		return nil, ErrNilOBJInstance
	}

	metadata, err := convertMetadata(&obj.Common, metadataParameter{
		lottoType:   lottoType,
		version:     LotteryV0,
		tirageOrder: 1,
		isOldType:   true,
		winOrder:    obj.WinOrder,
	})
	if err != nil {
		return nil, errors.Join(err, ErrInvalidMetadata)
	}

	balls := []int32{
		obj.B1,
		obj.B2,
		obj.B3,
		obj.B4,
		obj.B5,
		obj.B6,
		obj.AdditionalBall,
	}

	draw := LotteryDraw{
		Metadata: metadata,
		FirstDraw: Draw{
			NumBalls:     NumBallInOldVersion,
			Balls:        balls,
			HasLuckyBall: false,
			WinStats: map[WinRank]WinStat{
				Rank1: {
					Rate:   float64(obj.GainR1),
					Number: obj.WinnerR1,
				},
				Rank2: {
					Rate:   float64(obj.GainR2),
					Number: obj.WinnerR2,
				},
				Rank3: {
					Rate:   float64(obj.GainR3),
					Number: obj.WinnerR3,
				},
				Rank4: {
					Rate:   float64(obj.GainR4),
					Number: obj.WinnerR4,
				},
				Rank5: {
					Rate:   float64(obj.GainR5),
					Number: obj.WinnerR5,
				},
				Rank6: {
					Rate:   float64(obj.GainR6),
					Number: obj.WinnerR6,
				},
				Rank7: {
					Rate:   float64(obj.GainR7),
					Number: obj.WinnerR7,
				},
			},
			NumRanks: NumRankInVersion0,
		},
		HasSecondDraw: false,
		JokerPlus:     obj.JokerPlus,
		HasJokerV1:    false,
	}

	return &draw, nil
}

// Version 1 is an old version of lottery game.
// In this version 6 balls are picks + one additional ball. All between 1 and 49.
// To normalize LotteryDraw, all balls are includes in the draw (include the additional ball).
// There is a second draw but record in an other raw from the history differentiable by the tirageOrder.
// There is 7 winner ranks.
// There is a Joker+ and the Joker v1 is added in this version.
func convertCSVVersion1(lottoType LottoType, obj *csv.Version1) (*LotteryDraw, error) {
	if obj == nil {
		return nil, ErrNilOBJInstance
	}

	metadata, err := convertMetadata(&obj.Common, metadataParameter{
		lottoType:   lottoType,
		version:     LotteryV1,
		tirageOrder: obj.Tirage,
		isOldType:   true,
		winOrder:    obj.WinOrder,
	})
	if err != nil {
		return nil, errors.Join(err, ErrInvalidMetadata)
	}

	balls := []int32{
		obj.B1,
		obj.B2,
		obj.B3,
		obj.B4,
		obj.B5,
		obj.B6,
		obj.AdditionalBall,
	}

	draw := LotteryDraw{
		Metadata: metadata,
		FirstDraw: Draw{
			NumBalls:     NumBallInOldVersion,
			Balls:        balls,
			HasLuckyBall: false,
			WinStats: map[WinRank]WinStat{
				Rank1: {
					Rate:   float64(obj.GainR1),
					Number: obj.WinnerR1,
				},
				Rank2: {
					Rate:   float64(obj.GainR2),
					Number: obj.WinnerR2,
				},
				Rank3: {
					Rate:   float64(obj.GainR3),
					Number: obj.WinnerR3,
				},
				Rank4: {
					Rate:   float64(obj.GainR4),
					Number: obj.WinnerR4,
				},
				Rank5: {
					Rate:   float64(obj.GainR5),
					Number: obj.WinnerR5,
				},
				Rank6: {
					Rate:   float64(obj.GainR6),
					Number: obj.WinnerR6,
				},
				Rank7: {
					Rate:   float64(obj.GainR7),
					Number: obj.WinnerR7,
				},
			},
			NumRanks: NumRankInVersion1,
		},
		HasSecondDraw: false,
		JokerPlus:     obj.JokerPlus,
		JokerV1:       obj.Joker,
		HasJokerV1:    true,
	}

	return &draw, nil
}

// Version 2 is not an old version of lottery game because the final rules are implemented.
// In this version 5 balls are picks (between 1 and 49) + one lucky ball (between 1 and 10).
// There are not second draw.
// There is 6 winner ranks.
// There is a Joker+ and the Joker v1 is removed in this version.
func convertCSVVersion2(lottoType LottoType, obj *csv.Version2) (*LotteryDraw, error) {
	if obj == nil {
		return nil, ErrNilOBJInstance
	}

	metadata, err := convertMetadata(&obj.Common, metadataParameter{
		lottoType:   lottoType,
		version:     LotteryV2,
		tirageOrder: 1,
		isOldType:   false,
		winOrder:    obj.WinOrder,
	})
	if err != nil {
		return nil, errors.Join(err, ErrInvalidMetadata)
	}

	balls := []int32{
		obj.B1,
		obj.B2,
		obj.B3,
		obj.B4,
		obj.B5,
	}

	draw := LotteryDraw{
		Metadata: metadata,
		FirstDraw: Draw{
			NumBalls:     NumBallInDraw,
			Balls:        balls,
			LuckyBall:    obj.LuckyBall,
			HasLuckyBall: true,
			WinStats: map[WinRank]WinStat{
				Rank1: {
					Rate:   float64(obj.GainR1),
					Number: obj.WinnerR1,
				},
				Rank2: {
					Rate:   float64(obj.GainR2),
					Number: obj.WinnerR2,
				},
				Rank3: {
					Rate:   float64(obj.GainR3),
					Number: obj.WinnerR3,
				},
				Rank4: {
					Rate:   float64(obj.GainR4),
					Number: obj.WinnerR4,
				},
				Rank5: {
					Rate:   float64(obj.GainR5),
					Number: obj.WinnerR5,
				},
				Rank6: {
					Rate:   float64(obj.GainR6),
					Number: obj.WinnerR6,
				},
			},
			NumRanks: NumRankInVersion2,
		},
		HasSecondDraw: false,
		JokerPlus:     obj.JokerPlus,
		HasJokerV1:    false,
	}

	return &draw, nil
}

// Version 3 use the final rules of the game and implement the win codes.
// In this version 5 balls are picks (between 1 and 49) + one lucky ball (between 1 and 10).
// There are not second draw.
// There is 9 winner ranks.
// There is a Joker+.
// Winning codes are codes automatically generated on each lottery ticket,
// at the time of the draw, a certain number is randomly selected to win a prize.
//
//nolint:dupl,funlen // the conversion should catch some duplications but it's ok for exhaustive mapping.
func convertCSVVersion3(lottoType LottoType, obj *csv.Version3) (*LotteryDraw, error) {
	if obj == nil {
		return nil, ErrNilOBJInstance
	}

	metadata, err := convertMetadata(&obj.Common, metadataParameter{
		lottoType:   lottoType,
		version:     LotteryV3,
		tirageOrder: 1,
		isOldType:   false,
		winOrder:    obj.WinOrder,
	})
	if err != nil {
		return nil, errors.Join(err, ErrInvalidMetadata)
	}

	balls := []int32{
		obj.B1,
		obj.B2,
		obj.B3,
		obj.B4,
		obj.B5,
	}

	draw := LotteryDraw{
		Metadata: metadata,
		FirstDraw: Draw{
			NumBalls:     NumBallInDraw,
			Balls:        balls,
			LuckyBall:    obj.LuckyBall,
			HasLuckyBall: true,
			WinStats: map[WinRank]WinStat{
				Rank1: {
					Rate:   float64(obj.GainR1),
					Number: obj.WinnerR1,
				},
				Rank2: {
					Rate:   float64(obj.GainR2),
					Number: obj.WinnerR2,
				},
				Rank3: {
					Rate:   float64(obj.GainR3),
					Number: obj.WinnerR3,
				},
				Rank4: {
					Rate:   float64(obj.GainR4),
					Number: obj.WinnerR4,
				},
				Rank5: {
					Rate:   float64(obj.GainR5),
					Number: obj.WinnerR5,
				},
				Rank6: {
					Rate:   float64(obj.GainR6),
					Number: obj.WinnerR6,
				},
				Rank7: {
					Rate:   float64(obj.GainR7),
					Number: obj.WinnerR7,
				},
				Rank8: {
					Rate:   float64(obj.GainR8),
					Number: obj.WinnerR8,
				},
				Rank9: {
					Rate:   float64(obj.GainR9),
					Number: obj.WinnerR9,
				},
			},
			NumRanks: NumRankInVersion3,
		},
		WinningCode: WinningCode{
			Codes:    normalizeCodes(obj.WinCodes),
			Price:    float64(obj.GainCode),
			NumCodes: int(obj.NumberWinCodes),
		},
		HasSecondDraw: false,
		JokerPlus:     obj.JokerPlus,
		HasJokerV1:    false,
	}

	return &draw, nil
}

// Version 4 use the final rules of the game and implement the second draw.
// In this version 5 balls are picks (between 1 and 49) + one lucky ball (between 1 and 10) on the first draw.
// 5 balls are picks for the second draws (between 1 and 49).
// There is 9 winner ranks for the first draw and 4 winner rank for the second draw.
// There is a Joker+.
// Winning codes are codes automatically generated on each lottery ticket,
// at the time of the draw, a certain number is randomly selected to win a prize.
// A new field 'promotion_second_tirage' is present in the schema but is already empty.
//
//nolint:dupl,funlen // the conversion should catch some duplications but it's ok for exhaustive mapping.
func convertCSVVersion4(lottoType LottoType, obj *csv.Version4) (*LotteryDraw, error) {
	if obj == nil {
		return nil, ErrNilOBJInstance
	}

	metadata, err := convertMetadata(&obj.Common, metadataParameter{
		lottoType:   lottoType,
		version:     LotteryV4,
		tirageOrder: 1,
		isOldType:   false,
		winOrder:    obj.WinOrder,
	})
	if err != nil {
		return nil, errors.Join(err, ErrInvalidMetadata)
	}

	ballsFirstDraw := []int32{
		obj.B1,
		obj.B2,
		obj.B3,
		obj.B4,
		obj.B5,
	}

	ballsSecondDraw := []int32{
		obj.B1SecondRoll,
		obj.B2SecondRoll,
		obj.B3SecondRoll,
		obj.B4SecondRoll,
		obj.B5SecondRoll,
	}

	draw := LotteryDraw{
		Metadata: metadata,
		FirstDraw: Draw{
			NumBalls:     NumBallInDraw,
			Balls:        ballsFirstDraw,
			LuckyBall:    obj.LuckyBall,
			HasLuckyBall: true,
			WinStats: map[WinRank]WinStat{
				Rank1: {
					Rate:   float64(obj.GainR1),
					Number: obj.WinnerR1,
				},
				Rank2: {
					Rate:   float64(obj.GainR2),
					Number: obj.WinnerR2,
				},
				Rank3: {
					Rate:   float64(obj.GainR3),
					Number: obj.WinnerR3,
				},
				Rank4: {
					Rate:   float64(obj.GainR4),
					Number: obj.WinnerR4,
				},
				Rank5: {
					Rate:   float64(obj.GainR5),
					Number: obj.WinnerR5,
				},
				Rank6: {
					Rate:   float64(obj.GainR6),
					Number: obj.WinnerR6,
				},
				Rank7: {
					Rate:   float64(obj.GainR7),
					Number: obj.WinnerR7,
				},
				Rank8: {
					Rate:   float64(obj.GainR8),
					Number: obj.WinnerR8,
				},
				Rank9: {
					Rate:   float64(obj.GainR9),
					Number: obj.WinnerR9,
				},
			},
			NumRanks: NumRankInVersion4,
		},
		SecondDraw: Draw{
			NumBalls:     NumBallInDraw,
			Balls:        ballsSecondDraw,
			HasLuckyBall: false,
			WinStats: map[WinRank]WinStat{
				Rank1: {
					Rate:   float64(obj.GainR1SecondRoll),
					Number: obj.WinnerR1SecondRoll,
				},
				Rank2: {
					Rate:   float64(obj.GainR2SecondRoll),
					Number: obj.WinnerR2SecondRoll,
				},
				Rank3: {
					Rate:   float64(obj.GainR3SecondRoll),
					Number: obj.WinnerR3SecondRoll,
				},
				Rank4: {
					Rate:   float64(obj.GainR4SecondRoll),
					Number: obj.WinnerR4SecondRoll,
				},
			},
			NumRanks: NumRankInVersion4SecondDraw,
		},
		WinningCode: WinningCode{
			Codes:    normalizeCodes(obj.WinCodes),
			Price:    float64(obj.GainCode),
			NumCodes: int(obj.NumberWinCodes),
		},
		HasSecondDraw: true,
		JokerPlus:     obj.JokerPlus,
		HasJokerV1:    false,
	}

	return &draw, nil
}

// Metadata conversion:
//
//	FDJID: is the official ID from the fdj API represented by the year and the draw number. It's not necessary uniq.
//	ID: is a composition reproducible to detect if draw match and
//		is composed by FDJID + draw number and the result number in ascendant order.
//	Date: is a draw date.
//	ForclosureDate: is the forclosure date, if empty, the date is set to the draw date + 2 days.
//	Version: represent the origin version of the lotto.
//	Type: represent the kind of lotto.
//	Day: represent the draw day.
//	Currency: represent the currency of the gain.
//	TirageOrder: represent the tirage order. Depends of the history version,
//		the second draw is not include in the same raw of the csv and this field.
//	allow to detect if it's the first or the second draw.
//	IsOldType: define it's this draw is a oldest Lottery game.
func convertMetadata(obj *csv.Common, params metadataParameter) (Metadata, error) {
	drawDate, err := dateConverter(obj.Date)
	if err != nil {
		return Metadata{}, fmt.Errorf("invalid date: %w", err)
	}
	forclosureDate, err := dateConverter(obj.ForclosureDate)
	if err != nil {
		if !errors.Is(err, ErrEmptyDate) {
			return Metadata{}, fmt.Errorf("invalid forclosure date: %w", err)
		}

		forclosureDate = drawDate.Add(model.DurationTwoDays)
	}
	day, err := dayConverter(obj.Day)
	if err != nil {
		return Metadata{}, fmt.Errorf("invalid day: %w", err)
	}
	currency, err := currencyConverter(obj.Currency)
	if err != nil {
		return Metadata{}, fmt.Errorf("invalid currency: %w", err)
	}

	return Metadata{
		FDJID: obj.ID,
		ID: fmt.Sprintf(
			"%s-%s-%s",
			obj.ID,
			strconv.Itoa(int(params.tirageOrder)),
			params.winOrder,
		),
		Date:           drawDate,
		ForclosureDate: forclosureDate,
		Version:        params.version,
		Type:           params.lottoType,
		Day:            day,
		Currency:       currency,
		TirageOrder:    params.tirageOrder,
		IsOldType:      params.isOldType,
	}, nil
}

// dateConverter detect the format date and convert it to a date type [time.Time].
// The input date support two format:
//   - 20060102 default
//   - 02/01/2006 (need to be reverse).
func dateConverter(date string) (time.Time, error) {
	var t time.Time
	var err error

	if date == "" {
		return time.Time{}, ErrEmptyDate
	}

	format := "20060102"
	if strings.ContainsRune(date, '/') {
		format = "02/01/2006"
	}
	if t, err = time.Parse(format, date); err != nil {
		return time.Time{}, err
	}

	return t, nil
}

// dayConverter detect the day format and convert it to a model.Day.
func dayConverter(day string) (model.Day, error) {
	normalizedDay := strings.TrimSpace(day)

	if normalizedDay == "" {
		return "", ErrEmptyDay
	}

	switch normalizedDay {
	case csv.Monday, csv.ShortMonday:
		return model.Monday, nil
	case csv.Tuesday, csv.ShortTuesday:
		return model.Tuesday, nil
	case csv.Wednesday, csv.ShortWednesday:
		return model.Wednesday, nil
	case csv.Thursday, csv.ShortThursday:
		return model.Thursday, nil
	case csv.Friday, csv.ShortFriday:
		return model.Friday, nil
	case csv.Saturday, csv.ShortSaturday:
		return model.Saturday, nil
	case csv.Sunday, csv.ShortSunday:
		return model.Sunday, nil
	default:
		return "", fmt.Errorf("%s: %w", normalizedDay, ErrUnknownDay)
	}
}

// currencyConverter convert a csv currency into a model.Currency.
func currencyConverter(currency string) (model.Currency, error) {
	if currency == "" {
		return "", ErrEmptyCurrency
	}

	switch currency {
	case csv.Euro:
		return model.EUR, nil
	case csv.Franc:
		return model.FRANC, nil
	default:
		return "", fmt.Errorf("%s: %w", currency, ErrUnknownCurrency)
	}
}

// normalizeCodes transform the input string in list of code.
// It format them to remove leading and trailing space char.
func normalizeCodes(input string) []string {
	codes := strings.Split(input, csv.WinCodeSeparator)

	for i, code := range codes {
		codes[i] = strings.Trim(code, " ")
	}

	return codes
}
