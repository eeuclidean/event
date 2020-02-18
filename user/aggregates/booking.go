package aggregates

import (
	"errors"
	"event/user/commands"
	"strconv"
	"strings"
	"time"
)

const (
	BOOKING_CREATED         = "booking_created"
	BOOKING_CANCELED        = "booking_canceled"
	BOOKING_PAYED           = "booking_payed"
	BOOKING_LOKET_CHECKEDIN = "booking_loket_checkedin"
	BOOKING_POLI_CHECKEDIN  = "booking_poli_checkedin"
)

func NewBooking(command commands.AddBookingCommand, noantrian int, amount int) Booking {
	return Booking{
		ID:              generateID(),
		PatientID:       command.PatientID,
		BranchID:        command.BranchID,
		PoliID:          command.PoliID,
		NoAntrianBranch: noantrian,
		SubPoliID:       command.SubPoliID,
		InsuranceID:     command.InsuranceID,
		AntrianPoliID:   command.GetAntrainPoliID(),
		AntrianBranchID: command.GetAntrainBranchID(),
		Amount:          amount,
		Tanggal:         command.Tanggal,
		Status:          BOOKING_CREATED,
		CreateBy:        command.CreateBy,
		Created:         time.Now().Format(time.RFC3339),
	}
}

type Booking struct {
	ID              string `json:"id,omitempty" bson:"_id,omitempty"`
	PatientID       string `json:"patient_id,omitempty" bson:"patient_id,omitempty"`
	BranchID        string `json:"branch_id,omitempty" bson:"branch_id,omitempty"`
	PoliID          string `json:"poli_id,omitempty" bson:"poli_id,omitempty"`
	AntrianBranchID string `json:"antrian_branch_id,omitempty" bson:"antrin_branch_id,omitempty"`
	AntrianPoliID   string `json:"antrian_poli_id,omitempty" bson:"antrin_poli_id,omitempty"`
	SubPoliID       string `json:"subpoli_id,omitempty" bson:"subpoli_id,omitempty"`
	NoAntrianBranch int    `json:"no_antrian_branch,omitempty" bson:"no_antrian_branch,omitempty"`
	NoAntrianPoli   int    `json:"no_antrian_poli,omitempty" bson:"no_antrian_poli,omitempty"`
	TotalPoliCalls  int    `json:"total_poli_calls,omitempty" bson:"total_poli_calls,omitempty"`
	IsPoliCall      bool   `json:"ispoli_calls,omitempty" bson:"ispoli_calls,omitempty"`
	TotalLoketCalls int    `json:"total_loket_calls,omitempty" bson:"total_loket_calls,omitempty"`
	IsLoketCall     bool   `json:"isloket_calls,omitempty" bson:"isloket_calls,omitempty"`
	InsuranceID     string `json:"asuransi_id,omitempty" bson:"asuransi_id,omitempty"`
	Tanggal         string `json:"tanggal,omitempty" bson:"tanggal,omitempty"`
	Status          string `json:"status,omitempty" bson:"status,omitempty"`
	CreateBy        string `json:"create_by,omitempty" bson:"create_by,omitempty"`
	LoketCheckInBy  string `json:"loket_check_in_by,omitempty" bson:"loket_check_in_by,omitempty"`
	PoliCheckInBy   string `json:"poli_check_in_by,omitempty" bson:"poli_check_in_by,omitempty"`
	Amount          int    `json:"amount,omitempty" bson:"amount,omitempty"`
	AmountPayed     int    `json:"amount_payed,omitempty" bson:"amount_payed,omitempty"`
	StatusBayar     bool   `json:"statusbayar,omitempty" bson:"statusbayar,omitempty"`
	Created         string `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}

func (booking *Booking) Call() error {
	if booking.Status == BOOKING_CREATED {
		booking.TotalLoketCalls += 1
		booking.IsLoketCall = true
		return nil
	} else if booking.Status == BOOKING_LOKET_CHECKEDIN {
		booking.TotalPoliCalls += 1
		booking.IsPoliCall = true
	}
	return errors.New(booking.getStatusMsg())
}

func (booking *Booking) SetStatusCanceled() error {
	if booking.Status == BOOKING_CREATED {
		booking.Status = BOOKING_CANCELED
		return nil
	}
	return errors.New(booking.getStatusMsg())
}

func (booking *Booking) SetStatusLoketCheckin(command commands.LoketCheckinBookingCommand) error {
	if booking.Status == BOOKING_CREATED {
		booking.Status = BOOKING_LOKET_CHECKEDIN
		booking.LoketCheckInBy = command.By
		booking.IsLoketCall = false
		return nil
	}
	return errors.New(booking.getStatusMsg())
}

func (booking *Booking) SetStatusPayed(command commands.PayBookingCommand) error {
	if booking.Status == BOOKING_LOKET_CHECKEDIN {
		booking.Status = BOOKING_PAYED
		booking.LoketCheckInBy = command.By
		booking.StatusBayar = command.StatusBayar
		booking.AmountPayed = command.Amount
		return nil
	}
	return errors.New(booking.getStatusMsg())
}

func (booking *Booking) SetStatusPoliCheckIn(by string) error {
	if booking.Status == BOOKING_PAYED {
		booking.Status = BOOKING_POLI_CHECKEDIN
		booking.PoliCheckInBy = by
		booking.IsPoliCall = false
		return nil
	}
	return errors.New(booking.getStatusMsg())
}

func (booking *Booking) IsToday() bool {
	values := strings.Split(booking.Tanggal, "-")
	if len(values) != 3 {
		return false
	}
	date, _ := strconv.Atoi(values[0])
	month, _ := strconv.Atoi(values[1])
	year, _ := strconv.Atoi(values[2])
	today := time.Now().UTC().Add(7 * time.Hour)
	return (date == today.Day()) && (month == int(today.Month())) && (year == today.Year())
}

func (booking *Booking) getStatusMsg() string {
	switch booking.Status {
	case BOOKING_CREATED:
		return "Pasien Belum Checkin Di Loket"
	case BOOKING_CANCELED:
		return "Pasien Membatalkan Booking"
	case BOOKING_LOKET_CHECKEDIN:
		return "Pasien Sudah Checkin Di Loket"
	case BOOKING_POLI_CHECKEDIN:
		return "Pasien Sudah Checkin Di Poli"
	default:
		return "Terjadi Kesalahan"
	}
}
