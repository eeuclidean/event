package repo_test

import (
	"os"
	"event/user/aggregates"
	"event/user/repositories/bookingrepo"
	"event/user/testing/config"
	"event/user/utils/utilsgenerator"
	"event/user/utils/utilsmongo"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBookingRepo(t *testing.T) {
	assert := assert.New(t)
	config.RunConfig()

	var bookingRepo bookingrepo.BookingRepository
	{
		var err error
		bookingRepo, err = bookingrepo.NewMongoAdapterBookingRepo()
		assert.Nil(err)
	}

	var booking aggregates.Booking
	{
		booking = aggregates.Booking{
			ID:          utilsgenerator.NewID(),
			PatientID:   utilsgenerator.NewID(),
			BranchID:    utilsgenerator.NewID(),
			PoliID:      utilsgenerator.NewID(),
			SubPoliID:   utilsgenerator.NewID(),
			InsuranceID: utilsgenerator.NewID(),
			NoAntrian:   1,
			TotalCalls:  0,
			Tanggal:     "12-12-2019",
			Created:     time.Now().Format(time.RFC3339),
			Status:      aggregates.BOOKING_CREATED,
		}
	}
	collection, err := utilsmongo.GetCollection(os.Getenv(bookingrepo.BOOKING_COLL))
	assert.Nil(err)

	t.Run("Save Booking", func(t *testing.T) {
		err := bookingRepo.Save(booking)
		assert.Nil(err)
	})

	t.Run("Get Booking", func(t *testing.T) {
		savedBooking, err := bookingRepo.Get(booking.ID)
		assert.Nil(err)
		assert.Equal(booking, savedBooking)
	})

	t.Run("Update Booking", func(t *testing.T) {
		newBooking := aggregates.Booking{
			ID:          booking.ID,
			PatientID:   booking.PatientID,
			BranchID:    booking.BranchID,
			PoliID:      booking.PoliID,
			SubPoliID:   booking.SubPoliID,
			InsuranceID: booking.InsuranceID,
			NoAntrian:   booking.NoAntrian,
			TotalCalls:  booking.TotalCalls,
			Tanggal:     booking.Tanggal,
			Created:     booking.Created,
			Status:      aggregates.BOOKING_CANCELED,
		}
		err := bookingRepo.Update(newBooking)
		assert.Nil(err)

		var savedBooking aggregates.Booking
		err = collection.FindId(newBooking.ID).One(&savedBooking)
		assert.Nil(err)
		assert.Equal(newBooking, savedBooking)

	})

	t.Run("Delete Booking", func(t *testing.T) {
		err := bookingRepo.Remove(booking.ID)
		assert.Nil(err)
	})
}
