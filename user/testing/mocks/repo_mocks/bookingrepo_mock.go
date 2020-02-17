package repo_mocks

import (
	"event/user/aggregates"

	"github.com/stretchr/testify/mock"
)

type BookingRepoMock struct {
	mock.Mock
}

func (mock *BookingRepoMock) Save(booking aggregates.Booking) error {
	return nil
}
func (mock *BookingRepoMock) Update(booking aggregates.Booking) error {
	return nil
}
func (mock *BookingRepoMock) Remove(bookingid string) error {
	return nil
}
func (mock *BookingRepoMock) Get(bookingid string) (aggregates.Booking, error) {
	args := mock.Called(bookingid)
	return args.Get(0).(aggregates.Booking), args.Error(1)
}
