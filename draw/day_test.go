package draw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateConverter(t *testing.T) {
	t.Run("Should convert date with csv format: 20060102", func(t *testing.T) {
		date := "20190101"
		expected := "2019-01-01 00:00:00 +0000 UTC"

		result, err := dateConverter(date)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result.String())
		}
	})
	t.Run("Should convert date with csv format: 02/01/2006", func(t *testing.T) {
		date := "01/01/2019"
		expected := "2019-01-01 00:00:00 +0000 UTC"

		result, err := dateConverter(date)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result.String())
		}
	})
	t.Run("Should return an error if the date format is unknown", func(t *testing.T) {
		date := "2019-01-01"

		result, err := dateConverter(date)
		if assert.Error(t, err) {
			require.ErrorIs(t, err, ErrCSVDate)
			assert.Empty(t, result)
		}
	})
}

func TestDayConverter(t *testing.T) {
	t.Run("Should convert day with csv format: LUNDI | LU", func(t *testing.T) {
		var result Day
		var err error

		day1 := csvDayMonday
		day2 := csvDayShortMonday
		expected := DayMonday

		result, err = dayConverter(day1)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
		result, err = dayConverter(day2)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
	})
	t.Run("Should convert day with csv format: MARDI | MA", func(t *testing.T) {
		var result Day
		var err error

		day1 := csvDayTuesday
		day2 := csvDayShortTuesday
		expected := DayTuesday

		result, err = dayConverter(day1)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
		result, err = dayConverter(day2)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
	})
	t.Run("Should convert day with csv format: MERCREDI | ME", func(t *testing.T) {
		var result Day
		var err error

		day1 := csvDayWednesday
		day2 := csvDayShortWednesday
		expected := DayWednesday

		result, err = dayConverter(day1)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
		result, err = dayConverter(day2)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
	})
	t.Run("Should convert day with csv format: JEUDI | JE", func(t *testing.T) {
		var result Day
		var err error

		day1 := csvDayThursday
		day2 := csvDayShortThursday
		expected := DayThursday

		result, err = dayConverter(day1)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
		result, err = dayConverter(day2)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
	})
	t.Run("Should convert day with csv format: VENDREDI | VE", func(t *testing.T) {
		var result Day
		var err error

		day1 := csvDayFriday
		day2 := csvDayShortFriday
		expected := DayFriday

		result, err = dayConverter(day1)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
		result, err = dayConverter(day2)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
	})
	t.Run("Should convert day with csv format: SAMEDI | SA", func(t *testing.T) {
		var result Day
		var err error

		day1 := csvDaySaturday
		day2 := csvDayShortSaturday
		expected := DaySaturday

		result, err = dayConverter(day1)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
		result, err = dayConverter(day2)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
	})
	t.Run("Should convert day with csv format: DIMANCHE | DI", func(t *testing.T) {
		var result Day
		var err error

		day1 := csvDaySunday
		day2 := csvDayShortSunday
		expected := DaySunday

		result, err = dayConverter(day1)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
		result, err = dayConverter(day2)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, result)
		}
	})
	t.Run("Should return an error if the day format is unknown", func(t *testing.T) {
		day := "MONDAY"

		result, err := dayConverter(day)
		if assert.Error(t, err) {
			require.ErrorIs(t, err, ErrCSVDay)
			assert.Empty(t, result)
		}
	})
}
