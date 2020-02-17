package service

import (
	"event/user/aggregates"
	"event/user/commands"
)

type Service interface {
	CreateBooking(command commands.AddBookingCommand) (string, error)
	PayBooking(command commands.PayBookingCommand) error
	CallBooking(command commands.CallBookingCommand) error
	CheckInLoketBooking(command commands.LoketCheckinBookingCommand) error
	CheckInPoliBooking(command commands.PoliCheckinBookingCommand) error
	CancelBooking(command commands.CancelBookingCommand) error
	UpdateAntrian(command commands.UpdateAntrianCommand) error
	AddPoli(poli aggregates.Poli) error
	UpdatePoli(poli aggregates.Poli) error
	AddBranch(branch aggregates.Branch) error
	UpdateBranch(branch aggregates.Branch) error
	AddSchedule(schedule aggregates.Schedule) error
	UpdateSchedule(schedule aggregates.Schedule) error
}
