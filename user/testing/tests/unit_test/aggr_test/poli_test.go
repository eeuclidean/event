package aggr_test

import (
	"event/user/aggregates"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoli(t *testing.T) {
	assert := assert.New(t)
	poli := aggregates.Poli{
		Max:   100,
		Count: 0,
	}

	t.Run("Add Poli Count", func(t *testing.T) {
		assert.Equal(0, poli.Count)
		poli.AddCount()
		assert.Equal(1, poli.Count)
		poli.AddCount()
		assert.Equal(2, poli.Count)
	})

	t.Run("Max Count", func(t *testing.T) {
		assert.Equal(false, poli.IsMax())
		poli.Count = 100
		assert.Equal(true, poli.IsMax())
	})
}
