package aggr_test

import (
	"event/user/aggregates"
	"event/user/commands"
	"event/user/utils/utilsgenerator"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestBooking(t *testing.T) {
	assert := assert.New(t)
	command := commands.AddBookingCommand{
		PatientID:   utilsgenerator.NewID(),
		BranchID:    utilsgenerator.NewID(),
		PoliID:      utilsgenerator.NewID(),
		SubPoliID:   utilsgenerator.NewID(),
		InsuranceID: utilsgenerator.NewID(),
		Tanggal:     "12-12-2019",
		CreateBy:    utilsgenerator.NewID(),
	}
	noantrian := 2
	booking := aggregates.NewBooking(command, noantrian)
	t.Run("NewBooking function", func(t *testing.T) {
		assert.NotEqual("", booking.ID)
		assert.Equal(booking.PatientID, command.PatientID)
		assert.Equal(aggregates.BOOKING_CREATED, booking.Status)
		assert.Equal(command.CreateBy, booking.CreateBy)
		assert.Equal(command.PoliID, booking.PoliID)
		assert.Equal(command.BranchID, booking.BranchID)
		assert.Equal(command.SubPoliID, booking.SubPoliID)
		assert.Equal(command.InsuranceID, booking.InsuranceID)
	})

	t.Run("GetDomainEvent", func(t *testing.T) {
		event := booking.GetDomainEvent()
		data, _ := json.Marshal(booking)
		assert.Equal(booking.Status, event.EventType)
		assert.Equal(aggregates.BOOKING_EVENT_NAME, event.EventName)
		assert.Equal(booking.ID, event.DataID)
		assert.Equal(aggregates.BOOKING_DATA_NAME, event.DataName)
		assert.Equal(aggregates.BOOKING_SERVICE_NAME, event.CreateBy)
		assert.Equal(string(data), event.Data)
	})

	t.Run("Booking Call", func(t *testing.T) {
		booking.SetStatusCalled()
		assert.Equal(aggregates.BOOKING_CALLED, booking.Status)
		assert.Equal(1, booking.TotalCalls)
		booking.SetStatusCalled()
		assert.Equal(aggregates.BOOKING_CALLED, booking.Status)
		assert.Equal(2, booking.TotalCalls)
	})

	t.Run("Booking Loket Checkin", func(t *testing.T) {
		booking.SetStatusLoketCheckin()
		assert.Equal(aggregates.BOOKING_LOKET_CHECKEDIN, booking.Status)
	})  

	t.Run("Booking Poli Checkin", func(t *testing.T) {
		booking.SetStatusPoliCheckIn()
		assert.Equal(aggregates.BOOKING_POLI_CHECKEDIN, booking.Status)
	})

	t.Run("Booking Poli Canceled", func(t *testing.T) {
		booking.SetStatusCanceled()
		assert.Equal(aggregates.BOOKING_CANCELED, booking.Status)
	})
}
