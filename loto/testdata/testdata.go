package testdata

import (
	"testing"

	"github.com/winning-number/fdj-sdk/v2/loto/csv"
)

const (
	NotValidValue       = "not-a-valid-value"
	EmptyValue          = ""
	ValidDateValue      = "20060102"
	ValidFDJIdentifier  = "fdj-identifier"
	ValidJokerPlus      = "joker-plus"
	ValidJokerV1        = "joker-v1"
	ValidGainRank       = csv.FrenchFloat(42.0)
	ValidAdditionalBall = 45
	ValidBall1          = 20
	ValidBall2          = 31
	ValidBall3          = 8
	ValidBall4          = 48
	ValidBall5          = 2
	ValidBall6          = 17
	ValidLuckyBall      = 7
	ValidWinnerRank     = 14
	ValidWinOrder       = "win-order"
	ValidTirageOrder    = 2
	ValidWinCode        = "423456,345678,0928474"
	ValidNumberWinCode  = 3
)

func CopyCommon(t *testing.T, obj any, common csv.Common) {
	t.Helper()

	switch v := obj.(type) {
	case *csv.Version0:
		v.Common = common
	case *csv.Version1:
		v.Common = common
	case *csv.Version2:
		v.Common = common
	case *csv.Version3:
		v.Common = common
	case *csv.Version4:
		v.Common = common
	}
}

func CSVVersion0Data() *csv.Version0 {
	return &csv.Version0{
		Common: csv.Common{
			ID:             ValidFDJIdentifier,
			Date:           ValidDateValue,
			ForclosureDate: ValidDateValue,
			Day:            csv.ShortMonday,
			Currency:       csv.Euro,
		},
		JokerPlus:      ValidJokerPlus,
		WinOrder:       ValidWinOrder,
		GainR1:         ValidGainRank,
		GainR2:         ValidGainRank,
		GainR3:         ValidGainRank,
		GainR4:         ValidGainRank,
		GainR5:         ValidGainRank,
		GainR6:         ValidGainRank,
		GainR7:         ValidGainRank,
		AdditionalBall: ValidAdditionalBall,
		B1:             ValidBall1,
		B2:             ValidBall2,
		B3:             ValidBall3,
		B4:             ValidBall4,
		B5:             ValidBall5,
		B6:             ValidBall6,
		WinnerR1:       ValidWinnerRank,
		WinnerR2:       ValidWinnerRank,
		WinnerR3:       ValidWinnerRank,
		WinnerR4:       ValidWinnerRank,
		WinnerR5:       ValidWinnerRank,
		WinnerR6:       ValidWinnerRank,
		WinnerR7:       ValidWinnerRank,
	}
}

func CSVVersion1Data() *csv.Version1 {
	return &csv.Version1{
		Version0: csv.Version0{
			Common: csv.Common{
				ID:             ValidFDJIdentifier,
				Date:           ValidDateValue,
				ForclosureDate: ValidDateValue,
				Day:            csv.ShortMonday,
				Currency:       csv.Euro,
			},
			JokerPlus:      ValidJokerPlus,
			WinOrder:       ValidWinOrder,
			GainR1:         ValidGainRank,
			GainR2:         ValidGainRank,
			GainR3:         ValidGainRank,
			GainR4:         ValidGainRank,
			GainR5:         ValidGainRank,
			GainR6:         ValidGainRank,
			GainR7:         ValidGainRank,
			AdditionalBall: ValidAdditionalBall,
			B1:             ValidBall1,
			B2:             ValidBall2,
			B3:             ValidBall3,
			B4:             ValidBall4,
			B5:             ValidBall5,
			B6:             ValidBall6,
			WinnerR1:       ValidWinnerRank,
			WinnerR2:       ValidWinnerRank,
			WinnerR3:       ValidWinnerRank,
			WinnerR4:       ValidWinnerRank,
			WinnerR5:       ValidWinnerRank,
			WinnerR6:       ValidWinnerRank,
			WinnerR7:       ValidWinnerRank,
		},
		Joker:  ValidJokerV1,
		Tirage: ValidTirageOrder,
	}
}

func CSVVersion2Data() *csv.Version2 {
	return &csv.Version2{
		Common: csv.Common{
			ID:             ValidFDJIdentifier,
			Date:           ValidDateValue,
			ForclosureDate: ValidDateValue,
			Day:            csv.ShortMonday,
			Currency:       csv.Euro,
		},
		JokerPlus: ValidJokerPlus,
		WinOrder:  ValidWinOrder,
		GainR1:    ValidGainRank,
		GainR2:    ValidGainRank,
		GainR3:    ValidGainRank,
		GainR4:    ValidGainRank,
		GainR5:    ValidGainRank,
		GainR6:    ValidGainRank,
		LuckyBall: ValidLuckyBall,
		B1:        ValidBall1,
		B2:        ValidBall2,
		B3:        ValidBall3,
		B4:        ValidBall4,
		B5:        ValidBall5,
		WinnerR1:  ValidWinnerRank,
		WinnerR2:  ValidWinnerRank,
		WinnerR3:  ValidWinnerRank,
		WinnerR4:  ValidWinnerRank,
		WinnerR5:  ValidWinnerRank,
		WinnerR6:  ValidWinnerRank,
	}
}

func CSVVersion3Data() *csv.Version3 {
	return &csv.Version3{
		Version2: csv.Version2{
			Common: csv.Common{
				ID:             ValidFDJIdentifier,
				Date:           ValidDateValue,
				ForclosureDate: ValidDateValue,
				Day:            csv.ShortMonday,
				Currency:       csv.Euro,
			},
			JokerPlus: ValidJokerPlus,
			WinOrder:  ValidWinOrder,
			GainR1:    ValidGainRank,
			GainR2:    ValidGainRank,
			GainR3:    ValidGainRank,
			GainR4:    ValidGainRank,
			GainR5:    ValidGainRank,
			GainR6:    ValidGainRank,
			LuckyBall: ValidLuckyBall,
			B1:        ValidBall1,
			B2:        ValidBall2,
			B3:        ValidBall3,
			B4:        ValidBall4,
			B5:        ValidBall5,
			WinnerR1:  ValidWinnerRank,
			WinnerR2:  ValidWinnerRank,
			WinnerR3:  ValidWinnerRank,
			WinnerR4:  ValidWinnerRank,
			WinnerR5:  ValidWinnerRank,
			WinnerR6:  ValidWinnerRank,
		},
		WinCodes:       ValidWinCode,
		GainCode:       ValidGainRank,
		GainR7:         ValidGainRank,
		GainR8:         ValidGainRank,
		GainR9:         ValidGainRank,
		NumberWinCodes: ValidNumberWinCode,
		WinnerR7:       ValidWinnerRank,
		WinnerR8:       ValidWinnerRank,
		WinnerR9:       ValidWinnerRank,
	}
}

func CSVVersion4Data() *csv.Version4 {
	return &csv.Version4{
		Common: csv.Common{
			ID:             ValidFDJIdentifier,
			Date:           ValidDateValue,
			ForclosureDate: ValidDateValue,
			Day:            csv.ShortMonday,
			Currency:       csv.Euro,
		},
		JokerPlus:           ValidJokerPlus,
		WinOrder:            ValidWinOrder,
		GainR1:              ValidGainRank,
		GainR2:              ValidGainRank,
		GainR3:              ValidGainRank,
		GainR4:              ValidGainRank,
		GainR5:              ValidGainRank,
		GainR6:              ValidGainRank,
		GainR7:              ValidGainRank,
		GainR8:              ValidGainRank,
		GainR9:              ValidGainRank,
		LuckyBall:           ValidLuckyBall,
		B1:                  ValidBall1,
		B2:                  ValidBall2,
		B3:                  ValidBall3,
		B4:                  ValidBall4,
		B5:                  ValidBall5,
		WinnerR1:            ValidWinnerRank,
		WinnerR2:            ValidWinnerRank,
		WinnerR3:            ValidWinnerRank,
		WinnerR4:            ValidWinnerRank,
		WinnerR5:            ValidWinnerRank,
		WinnerR6:            ValidWinnerRank,
		WinnerR7:            ValidWinnerRank,
		WinnerR8:            ValidWinnerRank,
		WinnerR9:            ValidWinnerRank,
		WinCodes:            ValidWinCode,
		GainCode:            ValidGainRank,
		NumberWinCodes:      ValidNumberWinCode,
		B1SecondRoll:        ValidBall1,
		B2SecondRoll:        ValidBall2,
		B3SecondRoll:        ValidBall3,
		B4SecondRoll:        ValidBall4,
		B5SecondRoll:        ValidBall5,
		WinnerR1SecondRoll:  ValidWinnerRank,
		WinnerR2SecondRoll:  ValidWinnerRank,
		WinnerR3SecondRoll:  ValidWinnerRank,
		WinnerR4SecondRoll:  ValidWinnerRank,
		GainR1SecondRoll:    ValidGainRank,
		GainR2SecondRoll:    ValidGainRank,
		GainR3SecondRoll:    ValidGainRank,
		GainR4SecondRoll:    ValidGainRank,
		PromotionSecondRoll: "",
		WinOrderSecondRoll:  ValidWinOrder,
	}
}
