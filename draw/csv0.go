package draw

// CSV0 is the version 0 of the CSV file.
// It is the oldest version of the CSV file.
type CSV0 struct {
	coreCSV

	JokerPlus string `csv:"numero_jokerplus"`
	//nolint:misspell // This field is named like this in the CSV file.
	WinOrder       string `csv:"combinaison_gagnante_en_ordre_croissant"`
	GainR1         string `csv:"rapport_du_rang1"`
	GainR2         string `csv:"rapport_du_rang2"`
	GainR3         string `csv:"rapport_du_rang3"`
	GainR4         string `csv:"rapport_du_rang4"`
	GainR5         string `csv:"rapport_du_rang5"`
	GainR6         string `csv:"rapport_du_rang6"`
	GainR7         string `csv:"rapport_du_rang7"`
	AdditionalBall int32  `csv:"boule_complementaire"`
	B1             int32  `csv:"boule_1"`
	B2             int32  `csv:"boule_2"`
	B3             int32  `csv:"boule_3"`
	B4             int32  `csv:"boule_4"`
	B5             int32  `csv:"boule_5"`
	B6             int32  `csv:"boule_6"`
	WinnerR1       int32  `csv:"nombre_de_gagnant_au_rang1"`
	WinnerR2       int32  `csv:"nombre_de_gagnant_au_rang2"`
	WinnerR3       int32  `csv:"nombre_de_gagnant_au_rang3"`
	WinnerR4       int32  `csv:"nombre_de_gagnant_au_rang4"`
	WinnerR5       int32  `csv:"nombre_de_gagnant_au_rang5"`
	WinnerR6       int32  `csv:"nombre_de_gagnant_au_rang6"`
	WinnerR7       int32  `csv:"nombre_de_gagnant_au_rang7"`
}

// no winCode in this version.
//
//nolint:unparam // This winCode method satisfies the DrawConverter interface.
func (c CSV0) winCode() (WinCode, error) {
	return WinCode{
		Number: 0,
	}, nil
}

func (c CSV0) winStat() (WinStats, error) {
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
		},
		WinRate:    prices,
		SecondRoll: false,
	}, nil
}

func (c CSV0) joker() Joker {
	return Joker{
		Plus:    c.JokerPlus,
		HasBase: false,
	}
}

// roll returns the roll of the draw.
// in this version, the additional ball is a value between 1 and 49.
// so the additional ball is added into the first roll.
// the lucky ball is not present in this version.
// so the second roll is not present in this version.
func (c CSV0) roll() Roll {
	return Roll{
		First:     []int32{c.B1, c.B2, c.B3, c.B4, c.B5, c.B6, c.AdditionalBall},
		HasLucky:  false,
		HasSecond: false,
	}
}

func (c CSV0) metadata(drawT Type) (Metadata, error) {
	var meta Metadata
	var err error

	if meta, err = c.coreCSV.metadata(drawT); err != nil {
		return Metadata{}, err
	}
	meta.OldType = true
	meta.Version = V0
	meta.TirageOrder = 1
	meta.setID(c.WinOrder)

	return meta, nil
}
