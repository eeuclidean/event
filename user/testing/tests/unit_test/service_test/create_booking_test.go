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

func TestCreateBooking(t *testing.T) {
	assert := assert.New(t)
	bookingRepoMock := new(repo_mocks.BookingRepoMock)
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
	t.Run("Create Booking Success", func(t *testing.T) {
		poliRepoMock := new(repo_mocks.PoliRepoMock)
		poliRepoMock.On("Get", poli.ID).Return(poli, nil)
		bookingService := service.ServiceImpl{
			EventPublisher: &publisher.PublisherSuccessStub{},
			BookingRepo:    bookingRepoMock,
			PoliRepo:       poliRepoMock,
		}

		err := bookingService.CreateBooking(command)
		assert.Nil(err)
	})

	t.Run("Create Booking Error because Event Publisher Error", func(t *testing.T) {
		poliRepoMock := new(repo_mocks.PoliRepoMock)
		poliRepoMock.On("Get", poli.ID).Return(poli, nil)
		bookingService := service.ServiceImpl{
			EventPublisher: &publisher.PublisherErrorStub{},
			BookingRepo:    bookingRepoMock,
			PoliRepo:       poliRepoMock,
		}

		err := bookingService.CreateBooking(command)
		assert.NotNil(err)
	})

	t.Run("Create Booking Error because Poli Not Found", func(t *testing.T) {
		poliRepoMock := new(repo_mocks.PoliRepoMock)
		poliRepoMock.On("Get", poli.ID).Return(aggregates.Poli{}, errors.New("not found"))
		bookingService := service.ServiceImpl{
			EventPublisher: &publisher.PublisherSuccessStub{},
			BookingRepo:    bookingRepoMock,
			PoliRepo:       poliRepoMock,
		}

		err := bookingService.CreateBooking(command)
		assert.NotNil(err)
	})
}
