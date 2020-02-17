package service_test

import (
	"os"
	"event/user/aggregates"
	"event/user/commands"
	"event/user/repositories/polirepo"
	"event/user/service"
	"event/user/testing/config"
	"event/user/utils/utilsgenerator"
	"event/user/utils/utilsmongo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBooking(t *testing.T) {
	assert := assert.New(t)
	config.RunConfig()
	collection, err := utilsmongo.GetCollection(os.Getenv(polirepo.POLI_COLL))
	assert.Nil(err)
	t.Run("Create Booking Success", func(t *testing.T) {
		command := commands.AddBookingCommand{
			PatientID:   utilsgenerator.NewID(),
			BranchID:    utilsgenerator.NewID(),
			PoliID:      utilsgenerator.NewID(),
			SubPoliID:   utilsgenerator.NewID(),
			InsuranceID: utilsgenerator.NewID(),
			Tanggal:     "12-12-2019",
			CreateBy:    utilsgenerator.NewID(),
		}
		poli := aggregates.Poli{
			ID:    command.PoliID,
			Max:   100,
			Count: 0,
		}
		err := collection.Insert(poli)
		assert.Nil(err)
		bookingService, err := service.NewService()
		assert.Nil(err)
		err = bookingService.CreateBooking(command)
		assert.Nil(err)
	})

	t.Run("Create Booking Error Not Found", func(t *testing.T) {
		command := commands.AddBookingCommand{
			PatientID:   utilsgenerator.NewID(),
			BranchID:    utilsgenerator.NewID(),
			PoliID:      utilsgenerator.NewID(),
			SubPoliID:   utilsgenerator.NewID(),
			InsuranceID: utilsgenerator.NewID(),
			Tanggal:     "12-12-2019",
			CreateBy:    utilsgenerator.NewID(),
		}
		bookingService, err := service.NewService()
		assert.Nil(err)
		err = bookingService.CreateBooking(command)
		assert.NotNil(err)
	})

}
