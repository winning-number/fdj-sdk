package lotto

import (
	"testing"

	"github.com/gofast-pkg/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winning-number/fdj-sdk-lotto/draw"
)

// TestNewInstanceFunc tests the newInstanceFunc function.
// Private function tested for robustness.
func TestNewInstanceFunc(t *testing.T) {
	t.Run("Should be ok", func(t *testing.T) {
		var instance draw.Converter

		instance = newInstanceFunc(draw.V0)
		if assert.NotNil(t, instance) {
			assert.IsType(t, &draw.CSV0{}, instance)
		}
		instance = newInstanceFunc(draw.V1)
		if assert.NotNil(t, instance) {
			assert.IsType(t, &draw.CSV1{}, instance)
		}
		instance = newInstanceFunc(draw.V2)
		if assert.NotNil(t, instance) {
			assert.IsType(t, &draw.CSV2{}, instance)
		}
		instance = newInstanceFunc(draw.V3)
		if assert.NotNil(t, instance) {
			assert.IsType(t, &draw.CSV3{}, instance)
		}
		instance = newInstanceFunc(draw.V4)
		if assert.NotNil(t, instance) {
			assert.IsType(t, &draw.CSV4{}, instance)
		}
		instance = newInstanceFunc("unknown")
		if assert.NotNil(t, instance) {
			assert.IsType(t, &draw.CSV4{}, instance)
		}
	})
}

func TestSaveInstanceFunc(t *testing.T) {
	var dec csv.Decoder
	var err error

	if dec, err = csv.NewDecoder(csv.ConfigDecoder{
		NewInstanceFunc:  func() any { return nil },
		SaveInstanceFunc: func(any, csv.Decoder) error { return nil },
	}); err != nil {
		t.Fatal(err)
	}

	t.Run("Should return an error because the context is not set", func(t *testing.T) {
		var d draw.Draw

		d, err = saveInstanceFunc(&draw.CSV0{}, dec)
		if assert.Error(t, err) {
			assert.Equal(t, draw.Draw{}, d)
			require.ErrorIs(t, err, ErrDrawTypeKeyNotFound)
		}
	})
	t.Run(
		"Should return an error because the instance doesn't support the Converter type",
		func(t *testing.T) {
			var d draw.Draw

			dec.ContextSet(keyDrawType, string(draw.LottoType))
			d, err = saveInstanceFunc(&Filter{}, dec)
			if assert.Error(t, err) {
				assert.Equal(t, draw.Draw{}, d)
				require.ErrorIs(t, err, ErrInvalidDrawType)
			}
		})
	t.Run("Should return an error because converter has failed", func(t *testing.T) {
		var d draw.Draw

		dec.ContextSet(keyDrawType, string(draw.LottoType))
		d, err = saveInstanceFunc(&draw.CSV1{}, dec)
		if assert.Error(t, err) {
			assert.Equal(t, draw.Draw{}, d)
			require.ErrorIs(t, err, draw.ErrCSVDate)
		}
	})
}
