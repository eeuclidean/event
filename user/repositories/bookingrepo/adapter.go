package bookingrepo

import (
	"event/user/aggregates"
	"event/user/utils/utilsmongo"
	"gopkg.in/mgo.v2/bson"
	"os"

	mgo "gopkg.in/mgo.v2"
)

const (
	BOOKING_COLL = "BOOKING_COLL"
)

func NewMongoAdapterBookingRepo() (BookingRepository, error) {
	db, err := utilsmongo.MongoDBLogin()
	if err != nil {
		return MongoAdapterBookingRepo{}, err
	}
	return MongoAdapterBookingRepo{
		Collection:     db.C(os.Getenv(BOOKING_COLL)),
		CircuitBreaker: utilsmongo.NewCircuitBreaker(db),
	}, nil
}

type MongoAdapterBookingRepo struct {
	Collection     *mgo.Collection
	CircuitBreaker *utilsmongo.MongoCircuitBreaker
}

func (adapter MongoAdapterBookingRepo) Save(booking aggregates.Booking) error {
	return adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.Insert(booking)
	})
}

func (adapter MongoAdapterBookingRepo) Update(booking aggregates.Booking) error {
	return adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.UpdateId(booking.ID, booking)
	})
}

func (adapter MongoAdapterBookingRepo) Remove(bookingid string) error {
	return adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.RemoveId(bookingid)
	})
}

func (adapter MongoAdapterBookingRepo) Get(bookingid string) (aggregates.Booking, error) {
	var booking aggregates.Booking
	err := adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.FindId(bookingid).One(&booking)
	})
	return booking, err
}
func (adapter MongoAdapterBookingRepo) GetManyByPatientIDAndDate(branchid, patientid, tanggal string) ([]aggregates.Booking, error) {
	var bookings []aggregates.Booking
	err := adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.Find(bson.M{"branch_id": branchid, "patient_id": patientid, "tanggal": tanggal}).All(&bookings)
	})
	return bookings, err
}
