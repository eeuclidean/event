package service_test

import (
	"event/user/aggregates"
	"event/user/service"
	"event/user/testing/config"
	"event/user/utils/utilsgenerator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPoli(t *testing.T) {
	assert := assert.New(t)
	config.RunConfig()
	t.Run("Add Poli Success", func(t *testing.T) {
		poli := aggregates.Poli{
			ID:    utilsgenerator.NewID(),
			Max:   100,
			Count: 0,
		}
		bookingService, err := service.NewService()
		assert.Nil(err)
		err = bookingService.AddPoli(poli)
		assert.Nil(err)
	})
}
