package repositories

import (
	"event/user/aggregates"
)

type BookingRepository interface {
	Save(booking aggregates.Booking) error
	Update(booking aggregates.Booking) error
	Remove(bookingid string) error
	Get(bookingid string) (aggregates.Booking, error)
	GetManyByPatientIDAndDate(branchid, patientid, tanggal string) ([]aggregates.Booking, error)
}
