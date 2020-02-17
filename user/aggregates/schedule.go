package aggregates

const (
	SCHEDULE_CREATED = "schedule_created"
	SCHEDULE_UPDATED = "schedule_updated"
)

const (
	SCHEDULE_EVENT_NAME = "schedule"
)

type Schedule struct {
	ID          string `json:"id" bson:"_id"`
	PoliID      string `json:"poli_id" bson:"poli_id"`
	PoliName    string `json:"poli_name" bson:"poli_name"`
	Name        string `json:"name" bson:"name"`
	Type        string `json:"type" bson:"type"`
	Status      int    `json:"status" bson:"status"`
	Setting     string `json:"setting" bson:"setting"`
	DayOfMonth  int    `json:"day_of_month" bson:"day_of_month"`
	MonthOfYear int    `json:"month_of_year" bson:"month_of_year"`
	Year        int    `json:"year" bson:"year"`
}
