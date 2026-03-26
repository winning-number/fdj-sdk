package draw

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// bitSize is the number of bits in the mantissa of a float64.
const bitSize = 64

const (
	csvCurrencyEur = "eur"
	csvCurrencyFRF = "frf"
)

// Currency list.
const (
	CurrencyEur Currency = "EUR"
	CurrencyFR  Currency = "FRF"
)

// Currency is a type to represent the currency of a draw's gain.
type Currency string

// currencyConverter convert a currency [string] from a csv data to a [Currency].
// If the input parameter is unknow, the default currency [CurrencyEur] should be return.
func currencyConverter(name string) Currency {
	switch name {
	case csvCurrencyEur:
		return CurrencyEur
	case csvCurrencyFRF:
		return CurrencyFR
	default:
		return CurrencyEur
	}
}

// priceConverter convert any string to a slice of float64.
func priceConverter(prices ...string) ([]float64, error) {
	var ret []float64
	var err error

	ret = make([]float64, len(prices))
	for i, price := range prices {
		var v float64
		if v, err = priceToFloat64(price); err != nil {
			return nil, err
		}
		ret[i] = v
	}

	return ret, nil
}

// priceToFloat64 convert a price type [string] into a [float64].
// Support two type of input parameter:
//   - 42.0 (basic)
//   - 42,0 (specific from csv history)
func priceToFloat64(price string) (float64, error) {
	var ret float64
	var err error

	if price == "" {
		return 0.0, ErrEmptyMoney
	}
	price = strings.Replace(price, ",", ".", 1)
	if ret, err = strconv.ParseFloat(price, bitSize); err != nil {
		return 0.0, errors.Wrap(ErrCSVPrice, err.Error())
	}

	return ret, nil
}
