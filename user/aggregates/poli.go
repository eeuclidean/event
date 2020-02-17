package aggregates

const (
	POLI_CREATED = "poli_created"
	POLI_UPDATED = "poli_updated"
)

const (
	POLI_EVENT_NAME = "poli"
)

type Poli struct {
	ID                  string `json:"id" bson:"_id"`
	TenantID            string `json:"tenant_id" bson:"tenant_id"`
	BranchID            string `json:"branch_id" bson:"branch_id"`
	Name                string `json:"name" bson:"name"`
	Max                 int    `json:"max" bson:"max"`
	PolicyMaxDayBooking int    `json:"policy_max_day_booking" bson:"policy_max_day_booking"`
	CloseTime           int    `json:"close_time" bson:"close_time"`
	PayAmount           int    `json:"pay_amount" bson:"pay_amount"`
}
