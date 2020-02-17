package commands

type PayBookingCommand struct {
	ID          string `json:"id,omitempty"`
	Amount      int    `json:"amount,omitempty"`
	By          string `json:"by,omitempty"`
	StatusBayar bool   `json:"statusbayar,omitempty"`
}
