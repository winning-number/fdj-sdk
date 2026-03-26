package draw

// Type lists the different types of draw available.
const (
	SuperLottoType Type = "super-lotto"
	GrandLottoType Type = "grand-lotto"
	XmasLottoType  Type = "xmas-lotto"
	LottoType      Type = "new-lotto"
)

// Type is the type of draw.
type Type string
