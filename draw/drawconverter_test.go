package draw

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoreCSV_metadata(t *testing.T) {
	c := coreCSV{
		ID:             "123456",
		Date:           "invalid",
		ForclosureDate: "invalid",
		Day:            "invalid",
		Currency:       "eur",
	}
	t.Run("Should get an error from dateConverter with an invalid date", func(t *testing.T) {
		m, err := c.metadata(LottoType)
		if assert.Error(t, err) {
			assert.Empty(t, m)
			require.ErrorIs(t, err, ErrCSVDate)
		}
	})
	t.Run("Should get an error from dateConverter with an invalid forclosureDate", func(t *testing.T) {
		c.Date = "01/01/2020"

		m, err := c.metadata(LottoType)
		if assert.Error(t, err) {
			assert.Empty(t, m)
			require.ErrorIs(t, err, ErrCSVDate)
		}
	})
	t.Run("Should get an error from dayConverter with an invalid day", func(t *testing.T) {
		c.ForclosureDate = "01/01/2020"

		m, err := c.metadata(LottoType)
		if assert.Error(t, err) {
			assert.Empty(t, m)
			require.ErrorIs(t, err, ErrCSVDay)
		}
	})
	t.Run("Should convert the metadata", func(t *testing.T) {
		c.Day = "LU"

		expectedMetadata := Metadata{
			Version:        "",
			OldType:        false,
			DrawType:       LottoType,
			TirageOrder:    0,
			FDJID:          c.ID,
			ID:             "",
			Date:           time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			ForclosureDate: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			Day:            DayMonday,
			Currency:       CurrencyEur,
		}
		m, err := c.metadata(LottoType)
		if assert.NoError(t, err) {
			assert.Equal(t, expectedMetadata, m)
		}
	})
}
