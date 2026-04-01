# FDJ-SDK

<!-- Badge section [ci - codecov - release - godoc - code quality - codebeat - license - fossa ]-->
[![Static Badge](https://img.shields.io/badge/project%20use%20codesystem-green?link=https%3A%2F%2Fgithub.com%2Fgofast-pkg%2Fcodesystem)](https://github.com/gofast-pkg/codesystem)
![Build status](https://github.com/winning-number/fdj-sdk/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/winning-number/fdj-sdk/branch/main/graph/badge.svg?token=7TCE3QB21E)](https://codecov.io/gh/winning-number/fdj-sdk)
[![Release](https://img.shields.io/github/release/winning-number/fdj-sdk.svg?style=flat-square)](https://github.com/winning-number/fdj-sdk/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/winning-number/fdj-sdk.svg)](https://pkg.go.dev/github.com/winning-number/fdj-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/winning-number/fdj-sdk)](https://goreportcard.com/report/github.com/winning-number/fdj-sdk)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwinning-number%2Ffdj-sdk.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwinning-number%2Ffdj-sdk?ref=badge_shield)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://raw.githubusercontent.com/winning-number/fdj-sdk/main/LICENSE.md)

<!-- Description section -->

This package is a fdj's sdk (FRANCAISE DES JEUX). It could get the full game's history. The actual implementation expose the loto game model and data history. But you could get data from any game by using custom decoder.

[Loto history csv files](https://www.fdj.fr/jeux-de-tirage/loto/historique).

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

## Contributing

&nbsp;:grey_exclamation:&nbsp; Use issues for everything

Read more informations with the [CONTRIBUTING_GUIDE](./.github/CONTRIBUTING.md)

For all changes, please update the CHANGELOG.txt file by replacing the existant content.

Thank you &nbsp;:pray:&nbsp;&nbsp;:+1:&nbsp;

<a href="https://github.com/winning-number/fdj-sdk/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=winning-number/fdj-sdk" />
</a>

Made with [contrib.rocks](https://contrib.rocks).

## License

[GNU GPL v3](https://choosealicense.com/licenses/gpl-3.0/)

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwinning-number%2Ffdj-sdk.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwinning-number%2Ffdj-sdk?ref=badge_large)
