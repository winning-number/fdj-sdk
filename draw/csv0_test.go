package draw

import (
	"io"
	"os"
	"testing"

	"github.com/gofast-pkg/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winning-number/fdj-sdk-lotto/draw/testdata"
)

func loadCSV[T any](t *testing.T, filepath string, data *[]T) {
	t.Helper()

	var err error
	var file *os.File
	var csvReader csv.CSV

	//nolint:gosec // This function is used to load a file from the local file system in tests context.
	if file, err = os.Open(filepath); err != nil {
		t.Fatal(err)
	}
	defer func() { require.NoError(t, file.Close()) }()

	if csvReader, err = csv.New(file, ';'); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(*data); i++ {
		var warn csv.Warning

		warn, err = csvReader.Decode(&(*data)[i])
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		if len(warn) > 0 {
			t.Fatal(warn)
		}
	}
}

func TestCSV0_Convert(t *testing.T) {
	t.Run("Should convert csv0 to draw", func(t *testing.T) {
		data := make([]CSV0, 2)
		loadCSV(t, testdata.FileSuperLotoV0, &data)
		expected := DataSetSuperLottoV0()

		for i, csv := range data {
			draw, err := Convert(csv, SuperLottoType)
			if assert.NoError(t, err) {
				assert.Equal(t, expected[i], draw)
			}
		}
	})
	t.Run("Should return an error on the first draw because metadata has failed", func(t *testing.T) {
		data := make([]CSV0, 2)
		loadCSV(t, testdata.FileSuperLotoV0, &data)
		expected := DataSetSuperLottoV0()

		data[0].Date = invalidDate
		for i, csv := range data {
			draw, err := Convert(csv, SuperLottoType)
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
		data := make([]CSV0, 2)
		loadCSV(t, testdata.FileSuperLotoV0, &data)
		expected := DataSetSuperLottoV0()

		data[0].GainR1 = invalidGain
		for i, csv := range data {
			draw, err := Convert(csv, SuperLottoType)
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
