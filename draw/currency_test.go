package draw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrencyConverter(t *testing.T) {
	t.Run("Should convert currency with csv format: eur", func(t *testing.T) {
		var result Currency

		currency := csvCurrencyEur
		expected := CurrencyEur

		result = currencyConverter(currency)
		assert.Equal(t, expected, result)
	})
	t.Run("Should convert currency with csv format: frf", func(t *testing.T) {
		var result Currency

		currency := csvCurrencyFRF
		expected := CurrencyFR

		result = currencyConverter(currency)
		assert.Equal(t, expected, result)
	})
	t.Run("Should return the default currency if the input is unknown", func(t *testing.T) {
		var result Currency

		currency := "unknown"
		expected := CurrencyEur

		result = currencyConverter(currency)
		assert.Equal(t, expected, result)
	})
}

func TestPriceConverter(t *testing.T) {
	t.Run("Should convert price with csv format: 42.0", func(t *testing.T) {
		var result []float64
		var err error

		price := "42.0"
		expected := []float64{42.0}

		result, err = priceConverter(price)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
	})
	t.Run("Should convert price with csv format: 42,0", func(t *testing.T) {
		var result []float64
		var err error

		price := "42,0"
		expected := []float64{42.0}

		result, err = priceConverter(price)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
	})
	t.Run("Should return an error with en empty price", func(t *testing.T) {
		var result []float64
		var err error

		price := ""
		expected := []float64(nil)

		result, err = priceConverter(price)
		if assert.Error(t, err) {
			require.ErrorIs(t, err, ErrEmptyMoney)
			assert.Equal(t, expected, result)
		}
	})
	t.Run("Should return an error with a price that contains a letter", func(t *testing.T) {
		var result []float64
		var err error

		price := "42.0a"
		expected := []float64(nil)

		result, err = priceConverter(price)
		if assert.Error(t, err) {
			require.ErrorIs(t, err, ErrCSVPrice)
			assert.Equal(t, expected, result)
		}
	})
}
