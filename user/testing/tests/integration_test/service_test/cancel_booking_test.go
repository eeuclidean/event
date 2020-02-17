package service_test

import (
	"os"
	"event/user/aggregates"
	"event/user/commands"
	"event/user/repositories/bookingrepo"
	"event/user/service"
	"event/user/testing/config"
	"event/user/utils/utilsgenerator"
	"event/user/utils/utilsmongo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCancelBooking(t *testing.T) {
	assert := assert.New(t)
	config.RunConfig()
	collection, err := utilsmongo.GetCollection(os.Getenv(bookingrepo.BOOKING_COLL))
	assert.Nil(err)
	t.Run("Cancel Booking Success", func(t *testing.T) {
		booking := aggregates.Booking{
			ID:          utilsgenerator.NewID(),
			PatientID:   utilsgenerator.NewID(),
			BranchID:    utilsgenerator.NewID(),
			PoliID:      utilsgenerator.NewID(),
			SubPoliID:   utilsgenerator.NewID(),
			InsuranceID: utilsgenerator.NewID(),
			Tanggal:     "12-12-2019",
			CreateBy:    utilsgenerator.NewID(),
			Status:      aggregates.BOOKING_CREATED,
		}
		err := collection.Insert(booking)
		assert.Nil(err)

		command := commands.CancelBookingCommand{
			ID: booking.ID,
		}
		bookingService, err := service.NewService()
		assert.Nil(err)
		err = bookingService.CancelBooking(command)
		assert.Nil(err)
	})

	t.Run("Cancel Booking Error Not Found", func(t *testing.T) {
		command := commands.CancelBookingCommand{
			ID: utilsgenerator.NewID(),
		}
		bookingService, err := service.NewService()
		assert.Nil(err)
		err = bookingService.CancelBooking(command)
		assert.NotNil(err)
	})

}
