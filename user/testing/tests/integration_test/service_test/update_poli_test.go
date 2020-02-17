package service_test

import (
	"os"
	"event/user/aggregates"
	"event/user/repositories/polirepo"
	"event/user/service"
	"event/user/testing/config"
	"event/user/utils/utilsgenerator"
	"event/user/utils/utilsmongo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdatePoli(t *testing.T) {
	assert := assert.New(t)
	config.RunConfig()
	collection, err := utilsmongo.GetCollection(os.Getenv(polirepo.POLI_COLL))
	assert.Nil(err)
	t.Run("Update Poli Success", func(t *testing.T) {
		poli := aggregates.Poli{
			ID:    utilsgenerator.NewID(),
			Max:   100,
			Count: 0,
		}
		err := collection.Insert(poli)
		assert.Nil(err)
		poli.Max = 150
		bookingService, err := service.NewService()
		assert.Nil(err)
		err = bookingService.UpdatePoli(poli)
		assert.Nil(err)
	})
}
