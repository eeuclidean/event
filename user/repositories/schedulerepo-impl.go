package repositories

import (
	"event/user/aggregates"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoAdapterScheduleRepo struct {
	Collection  *mgo.Collection
	ConnChecker *connectionChecker
}

func (adapter mongoAdapterScheduleRepo) Save(schedule aggregates.Schedule) error {
	return adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.Insert(schedule)
	})
}

func (adapter mongoAdapterScheduleRepo) Update(schedule aggregates.Schedule) error {
	return adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.UpdateId(schedule.ID, schedule)
	})
}

func (adapter mongoAdapterScheduleRepo) GetByDate(poliid string, year, month, date int) (aggregates.Schedule, error) {
	var schedule aggregates.Schedule
	err := adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.Find(bson.M{"poli_id": poliid, "type": "except", "year": year, "month_of_year": month, "day_of_month": date}).One(&schedule)
	})
	return schedule, err
}
