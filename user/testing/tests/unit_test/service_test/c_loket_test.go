package service_test

import (
	"errors"
	"event/user/aggregates"
	"event/user/commands"
	"event/user/service"
	"event/user/testing/mocks/repo_mocks"
	"event/user/utils/utilsgenerator"
	"testing"

	"github.com/eeuclidean/eventsourcing/publisher"

	"github.com/stretchr/testify/assert"
)

func TestLoketCheckinBooking(t *testing.T) {
	assert := assert.New(t)
	booking := aggregates.Booking{
		ID:          utilsgenerator.NewID(),
		PatientID:   utilsgenerator.NewID(),
		BranchID:    utilsgenerator.NewID(),
		PoliID:      utilsgenerator.NewID(),
		SubPoliID:   utilsgenerator.NewID(),
		InsuranceID: utilsgenerator.NewID(),
		Tanggal:     "12-12-2019",
		CreateBy:    utilsgenerator.NewID(),
	}
	command := commands.LoketCheckinBookingCommand{
		ID: booking.ID,
	}
	t.Run("Loket Checkin Success", func(t *testing.T) {
		bookingRepoMock := new(repo_mocks.BookingRepoMock)
		bookingRepoMock.On("Get", booking.ID).Return(booking, nil)
		bookingService := service.ServiceImpl{
			EventPublisher: &publisher.PublisherSuccessStub{},
			BookingRepo:    bookingRepoMock,
		}
		err := bookingService.CheckInLoketBooking(command)
		assert.Nil(err)
	})

	t.Run("Loket Checkin Error because Event Publisher Error", func(t *testing.T) {
		bookingRepoMock := new(repo_mocks.BookingRepoMock)
		bookingRepoMock.On("Get", booking.ID).Return(booking, nil)
		bookingService := service.ServiceImpl{
			EventPublisher: &publisher.PublisherErrorStub{},
			BookingRepo:    bookingRepoMock,
		}
		err := bookingService.CheckInLoketBooking(command)
		assert.NotNil(err)
	})

	t.Run("Loket Checkin Error because Booking Not Found", func(t *testing.T) {
		bookingRepoMock := new(repo_mocks.BookingRepoMock)
		bookingRepoMock.On("Get", booking.ID).Return(aggregates.Booking{}, errors.New("not found"))
		bookingService := service.ServiceImpl{
			EventPublisher: &publisher.PublisherSuccessStub{},
			BookingRepo:    bookingRepoMock,
		}
		err := bookingService.CheckInLoketBooking(command)
		assert.NotNil(err)
	})
}
