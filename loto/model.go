package loto

import (
	"time"

	"github.com/winning-number/fdj-sdk/v2/model"
)

// Lottery version supported value.
const (
	LotteryV0 LotteryVersion = "v0"
	LotteryV1 LotteryVersion = "v1"
	LotteryV2 LotteryVersion = "v2"
	LotteryV3 LotteryVersion = "v3"
	LotteryV4 LotteryVersion = "v4"
)

const (
	// SuperLotto is a super lotto event.
	SuperLotto LottoType = "super-lotto"
	// GrandLotto is a big lotto event.
	GrandLotto LottoType = "grand-lotto"
	// XmasLotto is a christmas lotto event.
	XmasLotto LottoType = "xmas-lotto"
	// NewLotto is a common type for the lotto. It define a basic game.
	NewLotto LottoType = "new-lotto"
)

// NumBall set the expected balls in a Draw depends of the Old loto rules or the actual.
const (
	NumBallInOldVersion = 7
	NumBallInDraw       = 5
)

// Rank is the different positions to gain something.
// example with the new lotto.WinRank1: one ball + lucky ball.
// Rank1: 5 balls with lucky ball.
// Rank2: 5 balls without lucky ball.
// Rank3: 4 balls with lucky ball.
// Rank4: 4 balls without lucky ball.
// Rank5: 3 balls with lucky ball.
// Rank6: 3 balls without lucky ball.
// Rank7: 2 balls with lucky ball.
// Rank8: 2 balls without lucky ball.
// Rank9: 1 or 0 balls with lucky ball.
const (
	Rank1 WinRank = 1
	Rank2 WinRank = 2
	Rank3 WinRank = 3
	Rank4 WinRank = 4
	Rank5 WinRank = 5
	Rank6 WinRank = 6
	Rank7 WinRank = 7
	Rank8 WinRank = 8
	Rank9 WinRank = 9
)

// NumRank is the number rank available by lottery version.
const (
	NumRankInVersion0           = 7
	NumRankInVersion1           = 7
	NumRankInVersion2           = 6
	NumRankInVersion3           = 9
	NumRankInVersion4           = 9
	NumRankInVersion4SecondDraw = 4
)

// Draw is the result for a draw section in a LotteryDraw.
type Draw struct {
	NumBalls     int32
	Balls        []int32
	LuckyBall    int32
	HasLuckyBall bool
	NumRanks     int32
	WinStats     map[WinRank]WinStat
}

// LotteryDraw represents a lotto draw.
type LotteryDraw struct {
	Metadata      Metadata
	FirstDraw     Draw
	SecondDraw    Draw
	WinningCode   WinningCode
	HasSecondDraw bool
	JokerPlus     string
	JokerV1       string
	HasJokerV1    bool
}

// LotteryVersion setup the configuration for the LotteryDraw.
// Depends on the LotteryVersion some elements exist or not (ex: second draw).
type LotteryVersion string

// LottoType allow to make a diff between classic lotto (new-lotto), grand lotto, super lotto and xmas lotto.
type LottoType string

// Metadata represents the metadata of a LotteryDraw.
type Metadata struct {
	FDJID          string
	Date           time.Time
	ForclosureDate time.Time
	Version        LotteryVersion
	Type           LottoType
	Day            model.Day
	Currency       model.Currency
	TirageOrder    int32
	IsOldType      bool
}

// WinningCode is a result for the winning codes for each LotteryDraw.
type WinningCode struct {
	Codes    []string
	Price    float64
	NumCodes int
}

// WinRank type set the gain position.
type WinRank int8

// WinStat give the rate and number winner by rank.
type WinStat struct {
	Rate   float64
	Number int32
}
