package draw

// CSV2 is the draw type for the V2 format.
type CSV2 struct {
	coreCSV

	JokerPlus string `csv:"numero_jokerplus"`
	//nolint:misspell // This field is named like this in the CSV file.
	WinOrder  string `csv:"combinaison_gagnante_en_ordre_croissant"`
	GainR1    string `csv:"rapport_du_rang1"`
	GainR2    string `csv:"rapport_du_rang2"`
	GainR3    string `csv:"rapport_du_rang3"`
	GainR4    string `csv:"rapport_du_rang4"`
	GainR5    string `csv:"rapport_du_rang5"`
	GainR6    string `csv:"rapport_du_rang6"`
	B1        int32  `csv:"boule_1"`
	B2        int32  `csv:"boule_2"`
	B3        int32  `csv:"boule_3"`
	B4        int32  `csv:"boule_4"`
	B5        int32  `csv:"boule_5"`
	LuckyBall int32  `csv:"numero_chance"`
	WinnerR1  int32  `csv:"nombre_de_gagnant_au_rang1"`
	WinnerR2  int32  `csv:"nombre_de_gagnant_au_rang2"`
	WinnerR3  int32  `csv:"nombre_de_gagnant_au_rang3"`
	WinnerR4  int32  `csv:"nombre_de_gagnant_au_rang4"`
	WinnerR5  int32  `csv:"nombre_de_gagnant_au_rang5"`
	WinnerR6  int32  `csv:"nombre_de_gagnant_au_rang6"`
}

// no winCode in this version.
//
//nolint:unparam // This winCode method satisfies the DrawConverter interface.
func (c CSV2) winCode() (WinCode, error) {
	return WinCode{
		Number: 0,
	}, nil
}

func (c CSV2) winStat() (WinStats, error) {
	var err error
	var prices []float64

	if prices, err = priceConverter(
		c.GainR1,
		c.GainR2,
		c.GainR3,
		c.GainR4,
		c.GainR5,
		c.GainR6,
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
		},
		WinRate:    prices,
		SecondRoll: false,
	}, nil
}

func (c CSV2) joker() Joker {
	return Joker{
		Plus:    c.JokerPlus,
		HasBase: false,
	}
}

// roll returns the roll of the draw
// roll become common with 5 balls and 1 lucky ball
// so the second roll is not present in this version.
func (c CSV2) roll() Roll {
	return Roll{
		First:     []int32{c.B1, c.B2, c.B3, c.B4, c.B5},
		HasLucky:  true,
		LuckyBall: c.LuckyBall,
		HasSecond: false,
	}
}

func (c CSV2) metadata(drawT Type) (Metadata, error) {
	var meta Metadata
	var err error

	if meta, err = c.coreCSV.metadata(drawT); err != nil {
		return Metadata{}, err
	}
	meta.OldType = false
	meta.Version = V2
	meta.TirageOrder = 1
	meta.setID(c.WinOrder)

	return meta, nil
}
