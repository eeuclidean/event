package aggregates

const (
	BRANCH_CREATED = "branch_created"
	BRANCH_UPDATED = "branch_updated"
)

type Branch struct {
	ID               string `json:"id" bson:"_id"`
	MaxBookingPerDay int    `json:"maxbooking_perpatient_perday" bson:"maxbooking_perpatient_perday"`
}
