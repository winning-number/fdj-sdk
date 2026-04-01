package csv

// Version1 is the second version of the FDJ's lotto history.
// It's exactly identical to the Version0.
// This version add the second draw and a new Joker draw.
type Version1 struct {
	Version0

	Joker  string `csv:"numero_joker"`
	Tirage int32  `csv:"1er_ou_2eme_tirage"`
}
