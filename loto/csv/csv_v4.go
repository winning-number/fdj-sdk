package csv

// Version4 is the fifth versions of the FDJ's lotto history and the lastest version
// as we know it toaday.
// It add win codes and integrate a second draw without lucky ball.
// The new field 'promotion_second_tirage' is already empty and so never used.
type Version4 struct {
	Common

	JokerPlus string      `csv:"numero_jokerplus"`
	GainCode  FrenchFloat `csv:"rapport_codes_gagnants"`
	WinCodes  string      `csv:"codes_gagnants"`
	//nolint:misspell // This field is named like this in the CSV file.
	WinOrder            string      `csv:"combinaison_gagnante_en_ordre_croissant"`
	GainR1              FrenchFloat `csv:"rapport_du_rang1"`
	GainR2              FrenchFloat `csv:"rapport_du_rang2"`
	GainR3              FrenchFloat `csv:"rapport_du_rang3"`
	GainR4              FrenchFloat `csv:"rapport_du_rang4"`
	GainR5              FrenchFloat `csv:"rapport_du_rang5"`
	GainR6              FrenchFloat `csv:"rapport_du_rang6"`
	GainR7              FrenchFloat `csv:"rapport_du_rang7"`
	GainR8              FrenchFloat `csv:"rapport_du_rang8"`
	GainR9              FrenchFloat `csv:"rapport_du_rang9"`
	PromotionSecondRoll string      `csv:"promotion_second_tirage"`
	//nolint:misspell // This field is named like this in the CSV file.
	WinOrderSecondRoll string      `csv:"combinaison_gagnant_second_tirage_en_ordre_croissant"`
	GainR1SecondRoll   FrenchFloat `csv:"rapport_du_rang1_second_tirage"`
	GainR2SecondRoll   FrenchFloat `csv:"rapport_du_rang2_second_tirage"`
	GainR3SecondRoll   FrenchFloat `csv:"rapport_du_rang3_second_tirage"`
	GainR4SecondRoll   FrenchFloat `csv:"rapport_du_rang4_second_tirage"`
	NumberWinCodes     int32       `csv:"nombre_de_codes_gagnants"`
	B1                 int32       `csv:"boule_1"`
	B2                 int32       `csv:"boule_2"`
	B3                 int32       `csv:"boule_3"`
	B4                 int32       `csv:"boule_4"`
	B5                 int32       `csv:"boule_5"`
	LuckyBall          int32       `csv:"numero_chance"`
	WinnerR1           int32       `csv:"nombre_de_gagnant_au_rang1"`
	WinnerR2           int32       `csv:"nombre_de_gagnant_au_rang2"`
	WinnerR3           int32       `csv:"nombre_de_gagnant_au_rang3"`
	WinnerR4           int32       `csv:"nombre_de_gagnant_au_rang4"`
	WinnerR5           int32       `csv:"nombre_de_gagnant_au_rang5"`
	WinnerR6           int32       `csv:"nombre_de_gagnant_au_rang6"`
	WinnerR7           int32       `csv:"nombre_de_gagnant_au_rang7"`
	WinnerR8           int32       `csv:"nombre_de_gagnant_au_rang8"`
	WinnerR9           int32       `csv:"nombre_de_gagnant_au_rang9"`
	B1SecondRoll       int32       `csv:"boule_1_second_tirage"`
	B2SecondRoll       int32       `csv:"boule_2_second_tirage"`
	B3SecondRoll       int32       `csv:"boule_3_second_tirage"`
	B4SecondRoll       int32       `csv:"boule_4_second_tirage"`
	B5SecondRoll       int32       `csv:"boule_5_second_tirage"`
	WinnerR1SecondRoll int32       `csv:"nombre_de_gagnant_au_rang_1_second_tirage"`
	WinnerR2SecondRoll int32       `csv:"nombre_de_gagnant_au_rang_2_second_tirage"`
	WinnerR3SecondRoll int32       `csv:"nombre_de_gagnant_au_rang_3_second_tirage"`
	WinnerR4SecondRoll int32       `csv:"nombre_de_gagnant_au_rang_4_second_tirage"`
}
