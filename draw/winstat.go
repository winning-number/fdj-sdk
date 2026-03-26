package draw

// WinStats is the win stats of a draw for each rank.
type WinStats struct {
	// WinRateSecondRoll present amount money win by each winner for the second roll.
	// Index 0 is the Rank1.
	WinRateSecondRoll []float64
	// WinRate present amount money win by each winner. Index 0 is the Rank1.
	WinRate []float64
	// WinNumber present number of winner by Rank. Index 0 is the Rank1.
	WinNumber []int32
	// WinNumberSecondRoll present number of winner by Rank for the second roll.
	// Index 0 is the Rank1.
	WinNumberSecondRoll []int32
	SecondRoll          bool
}
