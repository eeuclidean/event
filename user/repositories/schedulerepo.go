package repositories

import (
	"event/user/aggregates"
)

type ScheduleRepository interface {
	Save(schedule aggregates.Schedule) error
	Update(schedule aggregates.Schedule) error
	GetByDate(scheduleid string, year, month, date int) (aggregates.Schedule, error)
}
