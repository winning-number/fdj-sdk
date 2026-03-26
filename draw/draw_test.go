package draw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func messyDraws() []Draw {
	return []Draw{
		DataSetSuperLottoV0()[0],
		DataSetClassicLottoV1()[1],
		DataSetClassicLottoV1()[0],
		DataSetSuperLottoV0()[1],
		DataSetClassicLottoV1()[1],
	}
}

func TestOrderDraws(t *testing.T) {
	t.Run("Should be ok to order draws by DESC order", func(t *testing.T) {
		expectedDraw := []Draw{
			DataSetClassicLottoV1()[0],
			DataSetClassicLottoV1()[1],
			DataSetClassicLottoV1()[1],
			DataSetSuperLottoV0()[0],
			DataSetSuperLottoV0()[1],
		}

		draws := messyDraws()
		OrderDraws(&draws, OrderASC)
		assert.Equal(t, expectedDraw, draws)
	})
	t.Run("Should be ok to order draws by DESC order", func(t *testing.T) {
		expectedDraw := []Draw{
			DataSetSuperLottoV0()[1],
			DataSetSuperLottoV0()[0],
			DataSetClassicLottoV1()[1],
			DataSetClassicLottoV1()[1],
			DataSetClassicLottoV1()[0],
		}

		draws := messyDraws()
		OrderDraws(&draws, OrderDESC)
		assert.Equal(t, expectedDraw, draws)
	})
	t.Run("Should be ok to order draws by NONE order", func(t *testing.T) {
		expectedDraw := messyDraws()

		draws := messyDraws()
		OrderDraws(&draws, OrderNone)
		assert.Equal(t, expectedDraw, draws)
	})
}

func TestFinder(t *testing.T) {
	t.Run("Should be ok to find draw", func(t *testing.T) {
		expectedDraw := DataSetClassicLottoV1()[0]

		draws := messyDraws()
		ok := Finder(&draws, expectedDraw)
		assert.True(t, ok)
	})
	t.Run("Should not find the draw", func(t *testing.T) {
		expectedDraw := DataSetClassicLottoV1()[0]

		expectedDraw.Metadata.ID = "different ID"
		draws := messyDraws()
		ok := Finder(&draws, expectedDraw)
		assert.False(t, ok)
	})
}
