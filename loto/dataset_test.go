package loto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatasetIdentifiers(t *testing.T) {
	t.Run("should matches", func(t *testing.T) {
		assert.Equal(t, []DatasetIdentifier{
			GrandLotoUUID,
			GrandLotoNoelUUID,
			SuperLoto199605UUID,
			SuperLoto200810UUID,
			SuperLoto201703UUID,
			SuperLoto201907UUID,
			Loto197605UUID,
			Loto200810UUID,
			Loto201703UUID,
			Loto201902UUID,
			Loto201911UUID,
		}, DatasetIdentifiers())
	})
}
