package draw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winning-number/fdj-sdk-lotto/draw/testdata"
)

func TestCSV1_Convert(t *testing.T) {
	t.Run("Should convert csv1 to draw", func(t *testing.T) {
		data := make([]CSV1, 2)
		loadCSV(t, testdata.FileClassicLotoV1, &data)
		expected := DataSetClassicLottoV1()

		for i, csv := range data {
			draw, err := Convert(csv, LottoType)
			if assert.NoError(t, err) {
				assert.Equal(t, expected[i], draw)
			}
		}
	})
	t.Run("Should return an error on the first draw because metadata has failed", func(t *testing.T) {
		data := make([]CSV1, 2)
		loadCSV(t, testdata.FileClassicLotoV1, &data)
		expected := DataSetClassicLottoV1()

		data[0].Date = invalidDate
		for i, csv := range data {
			draw, err := Convert(csv, LottoType)
			if i == 0 && assert.Error(t, err) {
				require.ErrorIs(t, err, ErrCSVDate)
				assert.Empty(t, draw)

				continue
			}

			if assert.NoError(t, err) {
				assert.Equal(t, expected[i], draw)
			}
		}
	})
	t.Run("Should return an error on the first draw because winStat has failed", func(t *testing.T) {
		data := make([]CSV1, 2)
		loadCSV(t, testdata.FileClassicLotoV1, &data)
		expected := DataSetClassicLottoV1()

		data[0].GainR1 = invalidGain
		for i, csv := range data {
			draw, err := Convert(csv, LottoType)
			if i == 0 && assert.Error(t, err) {
				require.ErrorIs(t, err, ErrCSVPrice)
				assert.Empty(t, draw)

				continue
			}

			if assert.NoError(t, err) {
				assert.Equal(t, expected[i], draw)
			}
		}
	})
}
