package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"event/user/aggregates"
	"event/user/commands"
	"event/user/gokit"
	"event/user/repositories/bookingrepo"
	"event/user/repositories/polirepo"
	"event/user/service"
	"event/user/testing/config"
	"event/user/utils/utilsgenerator"
	"event/user/utils/utilsmongo"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestBookingAPI(t *testing.T) {
	assert := assert.New(t)
	config.RunConfig()
	svc, err := service.NewService()
	if err != nil {
		panic(err)
	}
	endpoints := gokit.NewEndPoints(svc, gokit.Logger(), gokit.RequestDuration())
	httpHandler := gokit.NewHTTPServer(context.Background(), endpoints)

	poliCollection, err := utilsmongo.GetCollection(os.Getenv(polirepo.POLI_COLL))
	assert.Nil(err)

	bookingCollection, err := utilsmongo.GetCollection(os.Getenv(bookingrepo.BOOKING_COLL))
	assert.Nil(err)

	poli := aggregates.Poli{
		ID:    utilsgenerator.NewID(),
		Max:   100,
		Count: 0,
	}
	assert.Nil(poliCollection.Insert(poli))

	addBookingCommand := commands.AddBookingCommand{
		PatientID:   utilsgenerator.NewID(),
		BranchID:    utilsgenerator.NewID(),
		PoliID:      poli.ID,
		SubPoliID:   utilsgenerator.NewID(),
		InsuranceID: utilsgenerator.NewID(),
		Tanggal:     "12-12-2019",
		CreateBy:    utilsgenerator.NewID(),
	}
	var booking aggregates.Booking
	t.Run("Add Booking Success", func(t *testing.T) {
		data, _ := json.Marshal(addBookingCommand)
		req, err := http.NewRequest("POST", "/api/v1/user/create", bytes.NewReader(data))
		assert.Nil(err)

		rr := httptest.NewRecorder()
		httpHandler.ServeHTTP(rr, req)
		assert.Equal(http.StatusOK, rr.Code)

		err = bookingCollection.Find(bson.M{"patient_id": addBookingCommand.PatientID}).One(&booking)
		assert.Nil(err)
		assert.NotEqual(aggregates.Booking{}, booking)
	})
	t.Run("Cancel Booking Success", func(t *testing.T) {
		command := commands.CancelBookingCommand{
			ID: booking.ID,
		}
		data, _ := json.Marshal(command)
		req, err := http.NewRequest("PUT", "/api/v1/user/cancel", bytes.NewReader(data))
		assert.Nil(err)

		rr := httptest.NewRecorder()
		httpHandler.ServeHTTP(rr, req)
		assert.Equal(http.StatusOK, rr.Code)

		var savedBooking aggregates.Booking
		err = bookingCollection.FindId(booking.ID).One(&savedBooking)
		assert.Nil(err)
		assert.Equal(aggregates.BOOKING_CANCELED, savedBooking.Status)
	})
	t.Run("Call Booking Success", func(t *testing.T) {
		command := commands.CallBookingCommand{
			ID: booking.ID,
		}
		data, _ := json.Marshal(command)
		req, err := http.NewRequest("PUT", "/api/v1/user/call", bytes.NewReader(data))
		assert.Nil(err)

		rr := httptest.NewRecorder()
		httpHandler.ServeHTTP(rr, req)
		assert.Equal(http.StatusOK, rr.Code)

		var savedBooking aggregates.Booking
		err = bookingCollection.FindId(booking.ID).One(&savedBooking)
		assert.Nil(err)
		assert.Equal(aggregates.BOOKING_CALLED, savedBooking.Status)
		assert.Equal(1, savedBooking.TotalCalls)
	})
	t.Run("Checkin Loket Booking Success", func(t *testing.T) {
		command := commands.LoketCheckinBookingCommand{
			ID: booking.ID,
		}
		data, _ := json.Marshal(command)
		req, err := http.NewRequest("PUT", "/api/v1/user/checkin_loket", bytes.NewReader(data))
		assert.Nil(err)

		rr := httptest.NewRecorder()
		httpHandler.ServeHTTP(rr, req)
		assert.Equal(http.StatusOK, rr.Code)

		var savedBooking aggregates.Booking
		err = bookingCollection.FindId(booking.ID).One(&savedBooking)
		assert.Nil(err)
		assert.Equal(aggregates.BOOKING_LOKET_CHECKEDIN, savedBooking.Status)
	})
	t.Run("Checkin poli Booking Success", func(t *testing.T) {
		command := commands.PoliCheckinBookingCommand{
			ID: booking.ID,
		}
		data, _ := json.Marshal(command)
		req, err := http.NewRequest("PUT", "/api/v1/user/checkin_poli", bytes.NewReader(data))
		assert.Nil(err)

		rr := httptest.NewRecorder()
		httpHandler.ServeHTTP(rr, req)
		assert.Equal(http.StatusOK, rr.Code)

		var savedBooking aggregates.Booking
		err = bookingCollection.FindId(booking.ID).One(&savedBooking)
		assert.Nil(err)
		assert.Equal(aggregates.BOOKING_POLI_CHECKEDIN, savedBooking.Status)
	})
}
