package csv

// Version3 is the fourth version of the FDJ's lotto history.
// It just extends the number win ranks for a draw.
type Version3 struct {
	Version2

	WinCodes       string      `csv:"codes_gagnants"`
	GainCode       FrenchFloat `csv:"rapport_codes_gagnants"`
	GainR7         FrenchFloat `csv:"rapport_du_rang7"`
	GainR8         FrenchFloat `csv:"rapport_du_rang8"`
	GainR9         FrenchFloat `csv:"rapport_du_rang9"`
	NumberWinCodes int32       `csv:"nombre_de_codes_gagnants"`
	WinnerR7       int32       `csv:"nombre_de_gagnant_au_rang7"`
	WinnerR8       int32       `csv:"nombre_de_gagnant_au_rang8"`
	WinnerR9       int32       `csv:"nombre_de_gagnant_au_rang9"`
}
