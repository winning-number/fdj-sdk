package csv

// Version0 is the first version of the FDJ's lotto history.
// At this time 6 balls are picks + an additional ball between 1 and 49.
type Version0 struct {
	Common

	JokerPlus string `csv:"numero_jokerplus"`
	//nolint:misspell // This field is named like this in the CSV file.
	WinOrder       string      `csv:"combinaison_gagnante_en_ordre_croissant"`
	GainR1         FrenchFloat `csv:"rapport_du_rang1"`
	GainR2         FrenchFloat `csv:"rapport_du_rang2"`
	GainR3         FrenchFloat `csv:"rapport_du_rang3"`
	GainR4         FrenchFloat `csv:"rapport_du_rang4"`
	GainR5         FrenchFloat `csv:"rapport_du_rang5"`
	GainR6         FrenchFloat `csv:"rapport_du_rang6"`
	GainR7         FrenchFloat `csv:"rapport_du_rang7"`
	AdditionalBall int32       `csv:"boule_complementaire"`
	B1             int32       `csv:"boule_1"`
	B2             int32       `csv:"boule_2"`
	B3             int32       `csv:"boule_3"`
	B4             int32       `csv:"boule_4"`
	B5             int32       `csv:"boule_5"`
	B6             int32       `csv:"boule_6"`
	WinnerR1       int32       `csv:"nombre_de_gagnant_au_rang1"`
	WinnerR2       int32       `csv:"nombre_de_gagnant_au_rang2"`
	WinnerR3       int32       `csv:"nombre_de_gagnant_au_rang3"`
	WinnerR4       int32       `csv:"nombre_de_gagnant_au_rang4"`
	WinnerR5       int32       `csv:"nombre_de_gagnant_au_rang5"`
	WinnerR6       int32       `csv:"nombre_de_gagnant_au_rang6"`
	WinnerR7       int32       `csv:"nombre_de_gagnant_au_rang7"`
}
