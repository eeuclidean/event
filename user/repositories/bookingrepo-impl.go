package repositories

import (
	"event/user/aggregates"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

type mongoAdapterBookingRepo struct {
	Collection  *mgo.Collection
	ConnChecker *connectionChecker
}

func (adapter mongoAdapterBookingRepo) Save(booking aggregates.Booking) error {
	return adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.Insert(booking)
	})
}

func (adapter mongoAdapterBookingRepo) Update(booking aggregates.Booking) error {
	return adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.UpdateId(booking.ID, booking)
	})
}

func (adapter mongoAdapterBookingRepo) Remove(bookingid string) error {
	return adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.RemoveId(bookingid)
	})
}

func (adapter mongoAdapterBookingRepo) Get(bookingid string) (aggregates.Booking, error) {
	var booking aggregates.Booking
	err := adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.FindId(bookingid).One(&booking)
	})
	return booking, err
}
func (adapter mongoAdapterBookingRepo) GetManyByPatientIDAndDate(branchid, patientid, tanggal string) ([]aggregates.Booking, error) {
	var bookings []aggregates.Booking
	err := adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.Find(bson.M{"branch_id": branchid, "patient_id": patientid, "tanggal": tanggal}).All(&bookings)
	})
	return bookings, err
}
