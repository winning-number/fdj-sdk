package fdj

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPI(t *testing.T) {
	t.Run("Should return an API", func(t *testing.T) {
		api := NewAPI()
		assert.NotNil(t, api)
	})
}
