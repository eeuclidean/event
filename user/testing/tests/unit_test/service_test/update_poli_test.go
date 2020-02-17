package service_test

import (
	"event/user/aggregates"
	"event/user/service"
	"event/user/testing/mocks/repo_mocks"
	"event/user/utils/utilsgenerator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdatePoliBooking(t *testing.T) {
	assert := assert.New(t)
	poli := aggregates.Poli{
		ID:    utilsgenerator.NewID(),
		Max:   100,
		Count: 0,
	}
	poliRepoMock := new(repo_mocks.PoliRepoMock)
	t.Run("Update Poli Success", func(t *testing.T) {
		bookingService := service.ServiceImpl{
			PoliRepo: poliRepoMock,
		}
		err := bookingService.UpdatePoli(poli)
		assert.Nil(err)
	})
}
