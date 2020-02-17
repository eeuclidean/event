package schedulerepo

import (
	"event/user/aggregates"
	"event/user/utils/utilsmongo"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	SCHEDULE_COLL = "SCHEDULE_COLL"
)

func NewMongoAdapterScheduleRepo() (ScheduleRepository, error) {
	db, err := utilsmongo.MongoDBLogin()
	if err != nil {
		return MongoAdapterScheduleRepo{}, err
	}
	return MongoAdapterScheduleRepo{
		Collection:     db.C(os.Getenv(SCHEDULE_COLL)),
		CircuitBreaker: utilsmongo.NewCircuitBreaker(db),
	}, nil
}

type MongoAdapterScheduleRepo struct {
	Collection     *mgo.Collection
	CircuitBreaker *utilsmongo.MongoCircuitBreaker
}

func (adapter MongoAdapterScheduleRepo) Save(schedule aggregates.Schedule) error {
	return adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.Insert(schedule)
	})
}

func (adapter MongoAdapterScheduleRepo) Update(schedule aggregates.Schedule) error {
	return adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.UpdateId(schedule.ID, schedule)
	})
}

func (adapter MongoAdapterScheduleRepo) GetByDate(poliid string, year, month, date int) (aggregates.Schedule, error) {
	var schedule aggregates.Schedule
	err := adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.Find(bson.M{"poli_id": poliid, "type": "except", "year": year, "month_of_year": month, "day_of_month": date}).One(&schedule)
	})
	return schedule, err
}
