// Package csv describe all formats available for the FDJ's lotto history.
// Each csv version match a csv format and rules.
package csv

import (
	"errors"
	"strconv"
	"strings"
)

// ErrToParseFloat64 is used when the parsing of FrenchFloat in csv file fails.
var ErrToParseFloat64 = errors.New("value is not convertible in FrenchFloat - float64")

// Day representation foundable in the lotto history source.
const (
	Monday         = "LUNDI"
	Tuesday        = "MARDI"
	Wednesday      = "MERCREDI"
	Thursday       = "JEUDI"
	Friday         = "VENDREDI"
	Saturday       = "SAMEDI"
	Sunday         = "DIMANCHE"
	ShortMonday    = "LU"
	ShortTuesday   = "MA"
	ShortWednesday = "ME"
	ShortThursday  = "JE"
	ShortFriday    = "VE"
	ShortSaturday  = "SA"
	ShortSunday    = "DI"
)

// Currencies supported.
const (
	Euro  = "eur"
	Franc = "frf"
)

// Separator to extract the winCodes inside a csv history file.
const (
	WinCodeSeparator = ","
)

// Common data available for all csv version.
type Common struct {
	ID             string `csv:"annee_numero_de_tirage"`
	Date           string `csv:"date_de_tirage"`
	ForclosureDate string `csv:"date_de_forclusion"`
	Day            string `csv:"jour_de_tirage"`
	Currency       string `csv:"devise"`
}

// FrenchFloat is a float64 which use ',' instead of '.' in csv file.
type FrenchFloat float64

// UnmarshalCSV override the method used by the csv unmarhsaller.
func (f *FrenchFloat) UnmarshalCSV(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	str := strings.ReplaceAll(string(data), ",", ".")
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return errors.Join(ErrToParseFloat64, err)
	}
	*f = FrenchFloat(val)

	return nil
}
