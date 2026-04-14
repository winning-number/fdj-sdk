# FDJ-SDK

<!-- Badge section [ci - codecov - release - godoc - code quality - codebeat - license - fossa ]-->
[![Static Badge](https://img.shields.io/badge/project%20use%20codesystem-green?link=https%3A%2F%2Fgithub.com%2Fgofast-pkg%2Fcodesystem)](https://github.com/gofast-pkg/codesystem)
![Build status](https://github.com/winning-number/fdj-sdk/v2/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/winning-number/fdj-sdk/branch/main/graph/badge.svg?token=7TCE3QB21E)](https://codecov.io/gh/winning-number/fdj-sdk)
[![Release](https://img.shields.io/github/release/winning-number/fdj-sdk.svg?style=flat-square)](https://github.com/winning-number/fdj-sdk/v2/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/winning-number/fdj-sdk.svg)](https://pkg.go.dev/github.com/winning-number/fdj-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/winning-number/fdj-sdk)](https://goreportcard.com/report/github.com/winning-number/fdj-sdk)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwinning-number%2Ffdj-sdk.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwinning-number%2Ffdj-sdk?ref=badge_shield)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://raw.githubusercontent.com/winning-number/fdj-sdk/main/LICENSE.md)

<!-- Description section -->

This package is a fdj's sdk (FRANCAISE DES JEUX). It could get the full game's history. The actual implementation expose the loto game model and data history. But you could get data from any game by using custom decoder.

[Loto history csv files](https://www.fdj.fr/jeux-de-tirage/loto/historique).

> **Disclaimer:** This SDK uses FDJ's public-facing API endpoints (e.g. `sto.api.fdj.fr`). These endpoints are not officially documented or supported for third-party use. They may change or become unavailable without notice. This project is not affiliated with, endorsed by, or in any way officially connected to the Française des Jeux (FDJ). Use at your own risk and ensure your usage complies with FDJ's terms of service.

## Installation

Download fdj-sdk:

```sh
$ go get -u github.com/winning-number/fdj-sdk@latest
# go: downloading github.com/winning-number/fdj-sdk ...
# go: added github.com/winning-number/fdj-sdk ...
```

## Usage

[Read the godoc documentation](https://pkg.go.dev/github.com/winning-number/fdj-sdk)

They are two way for use this sdk. You can get directly the history from the `FRANCAISE DES JEUX` or directly by providing a CSV file (in this case, you need to know [all type of Draw](#type-of-draw))

* From the `FRANCAISE DES JEUX` api

```golang
package main

import (
    "context"
    fdj "github.com/winning-number/fdj-sdk"
)

func main() {
    var err error
    var driver lotto.API

    if driver, err = fdj.NewAPI(); err != nil {
        panic(err)
    }
    history, err := driver.DownloadHistory(context.Background(), loto.Loto201911UUID)
    if err != nil {
        panic(err)
    }
    dec, err := loto.NewDecoder(&history)
    if err != nil {
        panic(err)
    }
    draws, err := dec.Decode()
    if err != nil {
        panic(err)
    }
}
```

* From your custom zip history files

```golang
package main

import (
    "context"

    "github.com/gofast-pkg/zip"
    fdj "github.com/winning-number/fdj-sdk"
)

func main() {
    var err error
    var driver lotto.API

    zipReader, err := zip.NewReader("your_filepath.zip")
    if err != nil {
        panic(err)
    }
    dec, err := loto.NewDecoderWithDataset(
        &source.Source{
            Metadata: source.Metadata{
                Identifier: "your-id", // in this context could be ignore.
                Size: sizeFile // in this context could be ignore.
                RequestdAt: time.Now() // in this context could be ignore.
                FileName: "your-file-name"
            },
            Data: zipReader,
        },
        loto.DatasetInfo{
            Type: loto.NewLotto,
            Version: loto.LotteryV4,
        })
    draws, err := dec.Decode()
    if err != nil {
        panic(err)
    }
}
```

## Type of DRAW

From the begin of the lotto, the rules are been updated any time. So, to exploit the full history, each version are interpreting like a global `Draw` type.

### V0

* Concern only the `super lotto` type `before 2008 october`. They was `6` balls between `1` and `49` included and one lucky ball between `1` and `49` included. All balls are added inside the global slice of ball inside the Draw type. Result is provinding in a pick order and number order.
* Only one a `joker+` picked by draw.
* `7 Winners Rank` by Draw.

### V1

* Concern only the `classic lotto` type `before 2008 october`. They was `6` balls between `1` and `49` included and one lucky ball between `1` and `49` included. All balls are added inside the global slice of ball inside the Draw type. Result is provinding in a pick order and number order.
* They was one a `joker+` and one `joker number` picked by draw.
* `7` Winners Rank by Draw.

### V2

* Concern the `classic lotto` and the `super lotto` type `between 2008 october and 2017 march`. They was `5` balls between `1` and `49` included and one lucky ball between `1` and `9` included. these 5 balls are added inside the global slice of ball inside the Draw type and the lucky ball is added inside a LuckyBall field (int32).
* Only one a `joker+` picked by draw.
* `6` Winners Rank by Draw.

### V3

* Concern the `classic lotto`, the `super lotto`, the `grand lotto` and the `grand lotto (xmas)` type `between 2017 march and 2019 november`. They was `5` balls between `1` and `49` included and one lucky ball between `1` and `9` included. these 5 balls are added inside the global slice of ball inside the Draw type and the lucky ball is added inside a LuckyBall field (int32). Result is provinding in a pick order and number order.
* Only one a `joker+` picked by draw.
* `9` Winners Rank by Draw.
* Add `Winning number` (any by draw depends of type of draw).

### V4

* Concern the `classic lotto` type `from 2019 november`. They was `5` balls between `1` and `49` included and one lucky ball between `1` and `9` included. these 5 balls are added inside the global slice of ball inside the Draw type and the lucky ball is added inside a LuckyBall field (int32). Result is provinding in a pick order and number order.
* Only one a `joker+` picked by draw.
* `9` Winners Rank by Draw.
* Add `Winning number` (any by draw depends of type of draw).
* Add a `Second Roll` with `5` balls between `1` and `49` included (no lucky ball). Result is providing in only in a number order.
* `4` Winners Rank for the second Roll

## Draw Info Endpoint: `service-draw-info/v3/draws`

This endpoint provides **draw metadata** — not winning numbers. It is used to retrieve current and upcoming draw information for any FDJ game.

**Base URL:** `https://www.sto.api.fdj.fr/anonymous/service-draw-info/v3/draws`

### Query Parameters

| Parameter | Values | Description |
| --- | --- | --- |
| `game_name` | `fdj.Game` | Game identifier |
| `sort` | `fdj.Sort` | Sort order (ASC - DESC) |
| `current` | `bool` | `true` = current open draw |
| `from_planned_at` | `ISO date` or `now` | Filter draws starting from |
| `to_planned_at` | `ISO date` | Filter draws up to |
| `page_number` | `int` | Pagination (page-size: 20) |

**`current=true` + `game=loto` returns the current loto draw.**
**`current=false` + `game=loto` + `from_planned_at=now` returns the next loto draw(s) after the current.**
**`current=false` + `game=loto` + `to_planned_at=now` returns the prev loto draw(s) before the current.**

### Response Fields

| Field | Type | Description |
| --- | --- | --- |
| `id` | string | FDJ internal draw ID |
| `external_id` | string | Official draw number (e.g. `26042`) |
| `game_external_id` | `fdj.GameExternalID` | Game identifier (e.g. `draw-13` = Loto) |
| `game_version` | int | Rules version (`0` = current) |
| `cdc` | int | Internal prize grid identifier |
| `type` | *string | Draw type: `"1"` = classic, `"2"` = exceptional, `"9"` = super mega jackpot. **Optional** — absent on some draws |
| `is_current` | bool | `true` = next draw open for bets |
| `planned_at` | datetime | Scheduled draw date/time |
| `wagering_ends_at` | datetime | Betting deadline |
| `forcloses_at` | datetime | Prize claim deadline (~90 days) |
| `processed_at` | *datetime | Processing timestamp. **Absent = future draw** |
| `cycle_number` | *string | Draw number within the current jackpot cycle |
| `rolldown` | bool | Jackpot redistributed to lower rank |
| `flowdown` | *bool | Jackpot cascade redistribution variant |
| `theme_id` | string | Visual theme (`"0"` = standard) |
| `guaranteed_amounts` | []Amount | Guaranteed minimum jackpot (EUR + XPF) |
| `guaranteed_raffle_prizes` | []RafflePrize | Raffle: winners count and amount per winner |
| `videos` | []Video | YouTube draw video (past draws only) |

#### Amount Decoding

Monetary values use a `value` + `scale` encoding:

```text
real_amount = value / 10^scale
```

Examples:

* `value=1_300_000_000, scale=2` → **13 000 000 €** (13M)
* `value=2_000_000, scale=2` → **20 000 €**
* XPF uses `scale=0` (no decimals — French Polynesia currency)

In Go: `amount := float64(value) / math.Pow10(scale)`

> **Note on datetime parsing:** the format `2026-04-11T20:15:00.000+02:00` is not parsed natively by `time.RFC3339` in Go (milliseconds + offset). Use `time.Parse("2006-01-02T15:04:05.000-07:00", ...)` or a custom type.

## Contributing

&nbsp;:grey_exclamation:&nbsp; Use issues for everything

Read more informations with the [CONTRIBUTING_GUIDE](./.github/CONTRIBUTING.md)

For all changes, please update the CHANGELOG.txt file by replacing the existant content.

Thank you &nbsp;:pray:&nbsp;&nbsp;:+1:&nbsp;

<a href="https://github.com/winning-number/fdj-sdk/v2/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=winning-number/fdj-sdk" />
</a>

Made with [contrib.rocks](https://contrib.rocks).

## License

[GNU GPL v3](https://choosealicense.com/licenses/gpl-3.0/)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwinning-number%2Ffdj-sdk.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwinning-number%2Ffdj-sdk?ref=badge_large)
