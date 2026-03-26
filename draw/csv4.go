package draw

import "strings"

// CSV4 is the CSV format for the draw type 4.
// It is the most recent format.
// It implement the second roll.
type CSV4 struct {
	coreCSV

	JokerPlus string `csv:"numero_jokerplus"`
	GainCode  string `csv:"rapport_codes_gagnants"`
	WinCodes  string `csv:"codes_gagnants"`
	//nolint:misspell // This field is named like this in the CSV file.
	WinOrder            string `csv:"combinaison_gagnante_en_ordre_croissant"`
	GainR1              string `csv:"rapport_du_rang1"`
	GainR2              string `csv:"rapport_du_rang2"`
	GainR3              string `csv:"rapport_du_rang3"`
	GainR4              string `csv:"rapport_du_rang4"`
	GainR5              string `csv:"rapport_du_rang5"`
	GainR6              string `csv:"rapport_du_rang6"`
	GainR7              string `csv:"rapport_du_rang7"`
	GainR8              string `csv:"rapport_du_rang8"`
	GainR9              string `csv:"rapport_du_rang9"`
	PromotionSecondRoll string `csv:"promotion_second_tirage"`
	//nolint:misspell // This field is named like this in the CSV file.
	WinOrderSecondRoll string `csv:"combinaison_gagnant_second_tirage_en_ordre_croissant"`
	GainR1SecondRoll   string `csv:"rapport_du_rang1_second_tirage"`
	GainR2SecondRoll   string `csv:"rapport_du_rang2_second_tirage"`
	GainR3SecondRoll   string `csv:"rapport_du_rang3_second_tirage"`
	GainR4SecondRoll   string `csv:"rapport_du_rang4_second_tirage"`
	NumberWinCodes     int32  `csv:"nombre_de_codes_gagnants"`
	B1                 int32  `csv:"boule_1"`
	B2                 int32  `csv:"boule_2"`
	B3                 int32  `csv:"boule_3"`
	B4                 int32  `csv:"boule_4"`
	B5                 int32  `csv:"boule_5"`
	LuckyBall          int32  `csv:"numero_chance"`
	WinnerR1           int32  `csv:"nombre_de_gagnant_au_rang1"`
	WinnerR2           int32  `csv:"nombre_de_gagnant_au_rang2"`
	WinnerR3           int32  `csv:"nombre_de_gagnant_au_rang3"`
	WinnerR4           int32  `csv:"nombre_de_gagnant_au_rang4"`
	WinnerR5           int32  `csv:"nombre_de_gagnant_au_rang5"`
	WinnerR6           int32  `csv:"nombre_de_gagnant_au_rang6"`
	WinnerR7           int32  `csv:"nombre_de_gagnant_au_rang7"`
	WinnerR8           int32  `csv:"nombre_de_gagnant_au_rang8"`
	WinnerR9           int32  `csv:"nombre_de_gagnant_au_rang9"`
	B1SecondRoll       int32  `csv:"boule_1_second_tirage"`
	B2SecondRoll       int32  `csv:"boule_2_second_tirage"`
	B3SecondRoll       int32  `csv:"boule_3_second_tirage"`
	B4SecondRoll       int32  `csv:"boule_4_second_tirage"`
	B5SecondRoll       int32  `csv:"boule_5_second_tirage"`
	WinnerR1SecondRoll int32  `csv:"nombre_de_gagnant_au_rang_1_second_tirage"`
	WinnerR2SecondRoll int32  `csv:"nombre_de_gagnant_au_rang_2_second_tirage"`
	WinnerR3SecondRoll int32  `csv:"nombre_de_gagnant_au_rang_3_second_tirage"`
	WinnerR4SecondRoll int32  `csv:"nombre_de_gagnant_au_rang_4_second_tirage"`
}

func (c CSV4) winCode() (WinCode, error) {
	price, err := priceToFloat64(c.GainCode)
	if err != nil {
		return WinCode{}, err
	}

	return WinCode{
		Number: int(c.NumberWinCodes),
		Codes:  strings.Split(c.WinCodes, ","),
		Price:  price,
	}, nil
}

func (c CSV4) winStat() (WinStats, error) {
	var err error
	var firstPrices []float64
	var secondPrices []float64

	if firstPrices, err = priceConverter(
		c.GainR1,
		c.GainR2,
		c.GainR3,
		c.GainR4,
		c.GainR5,
		c.GainR6,
		c.GainR7,
		c.GainR8,
		c.GainR9,
	); err != nil {
		return WinStats{}, err
	}
	if secondPrices, err = priceConverter(
		c.GainR1SecondRoll,
		c.GainR2SecondRoll,
		c.GainR3SecondRoll,
		c.GainR4SecondRoll,
	); err != nil {
		return WinStats{}, err
	}

	return WinStats{
		WinNumber: []int32{
			c.WinnerR1,
			c.WinnerR2,
			c.WinnerR3,
			c.WinnerR4,
			c.WinnerR5,
			c.WinnerR6,
			c.WinnerR7,
			c.WinnerR8,
			c.WinnerR9,
		},
		WinRate:    firstPrices,
		SecondRoll: true,
		WinNumberSecondRoll: []int32{
			c.WinnerR1SecondRoll,
			c.WinnerR2SecondRoll,
			c.WinnerR3SecondRoll,
			c.WinnerR4SecondRoll,
		},
		WinRateSecondRoll: secondPrices,
	}, nil
}

func (c CSV4) joker() Joker {
	return Joker{
		Plus:    c.JokerPlus,
		HasBase: false,
	}
}

// roll returns the roll of the draw
// this version implements the second roll.
func (c CSV4) roll() Roll {
	return Roll{
		First:     []int32{c.B1, c.B2, c.B3, c.B4, c.B5},
		HasLucky:  true,
		LuckyBall: c.LuckyBall,
		HasSecond: true,
		Second: []int32{
			c.B1SecondRoll,
			c.B2SecondRoll,
			c.B3SecondRoll,
			c.B4SecondRoll,
			c.B5SecondRoll,
		},
	}
}

func (c CSV4) metadata(drawT Type) (Metadata, error) {
	var meta Metadata
	var err error

	if meta, err = c.coreCSV.metadata(drawT); err != nil {
		return Metadata{}, err
	}
	meta.OldType = false
	meta.Version = V4
	meta.TirageOrder = 1
	meta.setID(c.WinOrder)

	return meta, nil
}
