package draw

import "strings"

// CSV3 is the draw for the version 3 of the CSV.
// It is really near from the draw type v2.
// this version introduce the win code and extend the win rate to 9 ranks.
// So, CSV2 is embedded.
type CSV3 struct {
	CSV2

	WinCodes       string `csv:"codes_gagnants"`
	GainCode       string `csv:"rapport_codes_gagnants"`
	GainR7         string `csv:"rapport_du_rang7"`
	GainR8         string `csv:"rapport_du_rang8"`
	GainR9         string `csv:"rapport_du_rang9"`
	NumberWinCodes int32  `csv:"nombre_de_codes_gagnants"`
	WinnerR7       int32  `csv:"nombre_de_gagnant_au_rang7"`
	WinnerR8       int32  `csv:"nombre_de_gagnant_au_rang8"`
	WinnerR9       int32  `csv:"nombre_de_gagnant_au_rang9"`
}

func (c CSV3) winCode() (WinCode, error) {
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

func (c CSV3) winStat() (WinStats, error) {
	var err error
	var prices []float64

	if prices, err = priceConverter(
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
		WinRate:    prices,
		SecondRoll: false,
	}, nil
}

func (c CSV3) metadata(drawT Type) (Metadata, error) {
	var meta Metadata
	var err error

	if meta, err = c.coreCSV.metadata(drawT); err != nil {
		return Metadata{}, err
	}
	meta.OldType = false
	meta.Version = V3
	meta.TirageOrder = 1
	meta.setID(c.WinOrder)

	return meta, nil
}
