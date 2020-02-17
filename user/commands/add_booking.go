package commands

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type AddBookingCommand struct {
	PatientID   string `json:"patient_id,omitempty"`
	BranchID    string `json:"branch_id,omitempty"`
	PoliID      string `json:"poli_id,omitempty"`
	SubPoliID   string `json:"subpoli_id,omitempty"`
	InsuranceID string `json:"asuransi_id,omitempty"`
	Tanggal     string `json:"tanggal,omitempty"`
	CreateBy    string `json:"create_by,omitempty"`
}

func (command AddBookingCommand) GetAntrainPoliID() string {
	return command.PoliID + strings.ReplaceAll(command.Tanggal, "-", "")
}

func (command AddBookingCommand) GetAntrainBranchID() string {
	return command.BranchID + strings.ReplaceAll(command.Tanggal, "-", "")
}

func (command AddBookingCommand) GetDaysLeft() (int, error) {
	date, month, year, err := command.GetDateMonthYear()
	if err != nil {
		return 0, err
	}
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).UTC().Unix()
	bookingTime := time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.Local).UTC().Unix()
	duration, err := time.ParseDuration(strconv.Itoa(int(bookingTime-now)) + "s")
	if err != nil {
		return 0, errors.New("Wrong Date Format")
	}
	days := int(duration.Hours() / 24)
	return days, nil
}

func (command AddBookingCommand) GetTodayHour() int {
	today := time.Now().UTC().Add(7 * time.Hour)
	hour, _, _ := today.Clock()
	return hour
}

func (command AddBookingCommand) IsInWeekEnd() (bool, error) {
	date, month, year, err := command.GetDateMonthYear()
	if err != nil {
		return false, err
	}
	dateTime := time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.Local).UTC().Add(7 * time.Hour)
	if dateTime.Weekday() == time.Sunday || dateTime.Weekday() == time.Saturday {
		return true, nil
	}
	return false, nil
}

func (command AddBookingCommand) GetDateMonthYear() (int, int, int, error) {
	values := strings.Split(command.Tanggal, "-")
	if len(values) != 3 {
		return 0, 0, 0, errors.New("Wrong Date Format")
	}
	date, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, 0, 0, errors.New("Wrong Date Format")
	}
	month, err := strconv.Atoi(values[1])
	if err != nil {
		return 0, 0, 0, errors.New("Wrong Date Format")
	}
	year, err := strconv.Atoi(values[2])
	if err != nil {
		return 0, 0, 0, errors.New("Wrong Date Format")
	}
	return date, month, year, nil
}
