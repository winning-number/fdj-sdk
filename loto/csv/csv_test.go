package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const precisionEpsilon = 0.01

func TestFrenchFloat_UnmarshalCSV(t *testing.T) {
	t.Run("should return an error because data is invalid", func(t *testing.T) {
		var f FrenchFloat

		err := f.UnmarshalCSV([]byte("not-a-float"))
		require.ErrorIs(t, err, ErrToParseFloat64)
		assert.Zero(t, f)
	})
	t.Run("should return 0.0 when data is empty", func(t *testing.T) {
		var f FrenchFloat

		err := f.UnmarshalCSV([]byte(""))
		require.NoError(t, err)
		assert.Zero(t, f)
	})
	t.Run("should be ok with success conversion in FrenchFloat", func(t *testing.T) {
		var f FrenchFloat

		err := f.UnmarshalCSV([]byte("2000,3"))
		require.NoError(t, err)
		assert.InEpsilon(t, 2000.3, float64(f), precisionEpsilon)
	})
}
