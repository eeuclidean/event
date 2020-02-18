package service

import (
	"errors"
	"event/user/aggregates"
	"event/user/commands"
	"event/user/event/publisher"
	"event/user/repositories"

	"strconv"
)

const ERR_PREFIX = "400 "

func NewService() (Service, error) {
	repo, err := repositories.NewRepositories()
	if err != nil {
		return ServiceImpl{}, err
	}
	eventPublisher := publisher.NewEventPublisher()
	return ServiceImpl{
		Repositories:   repo,
		EventPublisher: eventPublisher,
	}, nil
}

type ServiceImpl struct {
	Repositories   repositories.Repositories
	EventPublisher publisher.EventPublisher
}

func (svc ServiceImpl) CreateBooking(command commands.AddBookingCommand) (string, error) {
	if ok, err := command.IsInWeekEnd(); err != nil {
		return "", err
	} else if ok {
		return "", errors.New("Sabtu Minggu Libur")
	}
	daysLeft, err := command.GetDaysLeft()
	if err != nil {
		return "", err
	}
	poli, err := svc.Repositories.GetPoliRepository().Get(command.PoliID)
	if err != nil {
		return "", err
	}
	date, month, year, _ := command.GetDateMonthYear()
	_, err = svc.Repositories.GetScheduleRepository().GetByDate(command.PoliID, year, month, date)
	if err == nil {
		return "", errors.New("Tanggal " + command.Tanggal + " libur")
	}
	if command.SubPoliID != "" {
		_, err = svc.Repositories.GetScheduleRepository().GetByDate(command.SubPoliID, year, month, date)
		if err == nil {
			return "", errors.New("Tanggal " + command.Tanggal + " libur")
		}
	}
	if poli.PolicyMaxDayBooking > 0 && daysLeft > poli.PolicyMaxDayBooking {
		return "", errors.New("Maksimum Rentang " + strconv.Itoa(poli.PolicyMaxDayBooking) + " hari")
	}
	if daysLeft == 0 && poli.CloseTime > 0 {
		if command.GetTodayHour() >= poli.CloseTime {
			return "", errors.New("Maksimum Waktu Booking Jam " + strconv.Itoa(poli.CloseTime))
		}
	}
	branch, err := svc.Repositories.GetBranchRepository().Get(poli.BranchID)
	if err != nil {
		return "", err
	}
	bookings, err := svc.Repositories.GetBookingResository().GetManyByPatientIDAndDate(command.BranchID, command.PatientID, command.Tanggal)
	if err != nil {
		return "", err
	}
	if branch.MaxBookingPerDay > 0 && len(bookings) >= branch.MaxBookingPerDay {
		return "", errors.New("Limit Booking " + strconv.Itoa(branch.MaxBookingPerDay) + " kali")
	}
	antrianBranch := aggregates.NewAntrianBranch(poli, command.Tanggal)
	err = svc.Repositories.GetAntrianRepository().Save(antrianBranch)
	if err != nil {
		if err.Error() == repositories.ANTRIAN_EXIST {
			var err error
			antrianBranch, err = svc.Repositories.GetAntrianRepository().GetAntrianBranch(command.GetAntrainBranchID(), repositories.REGISTER_CONTEXT)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	if err := svc.EventPublisher.PublishAntrianEvent(antrianBranch); err != nil {
		return "", err
	}
	antrianPoli := aggregates.NewAntrianPoli(poli, command.Tanggal)
	err = svc.Repositories.GetAntrianRepository().Save(antrianPoli)
	if err != nil {
		if err.Error() == repositories.ANTRIAN_EXIST {
			var err error
			antrianPoli, err = svc.Repositories.GetAntrianRepository().GetAntrianPoli(command.GetAntrainPoliID(), repositories.REGISTER_CONTEXT)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	if err := svc.EventPublisher.PublishAntrianEvent(antrianPoli); err != nil {
		return "", err
	}
	booking := aggregates.NewBooking(command, antrianBranch.Terisi, poli.PayAmount)
	err = svc.EventPublisher.PublishBookingEvent(booking)
	if err != nil {
		return "", err
	}
	err = svc.Repositories.GetBookingResository().Save(booking)
	if err != nil {
		return "", err
	}
	return booking.ID, nil
}

func (svc ServiceImpl) CallBooking(command commands.CallBookingCommand) error {
	booking, err := svc.Repositories.GetBookingResository().Get(command.ID)
	if err != nil {
		return err
	}
	if booking.IsToday() {
		booking.Call()
		err = svc.Repositories.GetBookingResository().Update(booking)
		if err != nil {
			return err
		}
		return svc.EventPublisher.PublishBookingEvent(booking)
	}
	return errors.New("Bukan Hari Ini")
}

func (svc ServiceImpl) CheckInLoketBooking(command commands.LoketCheckinBookingCommand) error {
	booking, err := svc.Repositories.GetBookingResository().Get(command.ID)
	if err != nil {
		return err
	}
	if booking.IsToday() {
		if err := booking.SetStatusLoketCheckin(command); err != nil {
			return err
		}
		antrianBranch, err := svc.Repositories.GetAntrianRepository().GetAntrianBranch(booking.AntrianBranchID, repositories.CHECKIN_CONTEXT)
		if err != nil {
			return err
		}
		if err := svc.EventPublisher.PublishAntrianEvent(antrianBranch); err != nil {
			return err
		}
		antrianPoli, err := svc.Repositories.GetAntrianRepository().GetAntrianPoli(booking.AntrianPoliID, repositories.CHECKIN_CONTEXT)
		if err != nil {
			return err
		}
		if err := svc.EventPublisher.PublishAntrianEvent(antrianPoli); err != nil {
			return err
		}
		booking.NoAntrianPoli = antrianPoli.Checkin
		err = svc.Repositories.GetBookingResository().Update(booking)
		if err != nil {
			return err
		}
		return svc.EventPublisher.PublishBookingEvent(booking)
	}
	return errors.New("Bukan Hari Ini")
}
func (svc ServiceImpl) CheckInPoliBooking(command commands.PoliCheckinBookingCommand) error {
	booking, err := svc.Repositories.GetBookingResository().Get(command.ID)
	if err != nil {
		return err
	}
	if booking.IsToday() {
		if err := booking.SetStatusPoliCheckIn(command.By); err != nil {
			return err
		}
		antrianPoli, err := svc.Repositories.GetAntrianRepository().GetAntrianPoli(booking.AntrianPoliID, repositories.FINISH_CONTEXT)
		if err != nil {
			return err
		}
		if err := svc.EventPublisher.PublishAntrianEvent(antrianPoli); err != nil {
			return err
		}
		err = svc.Repositories.GetBookingResository().Update(booking)
		if err != nil {
			return err
		}
		return svc.EventPublisher.PublishBookingEvent(booking)
	}
	return errors.New("Bukan Hari Ini")
}

func (svc ServiceImpl) PayBooking(command commands.PayBookingCommand) error {
	booking, err := svc.Repositories.GetBookingResository().Get(command.ID)
	if err != nil {
		return err
	}
	if booking.IsToday() {
		if err := booking.SetStatusPayed(command); err != nil {
			return err
		}
		err = svc.Repositories.GetBookingResository().Update(booking)
		if err != nil {
			return err
		}
		return svc.EventPublisher.PublishBookingEvent(booking)
	}
	return errors.New("Bukan Hari Ini")
}
func (svc ServiceImpl) CancelBooking(command commands.CancelBookingCommand) error {
	booking, err := svc.Repositories.GetBookingResository().Get(command.ID)
	if err != nil {
		return err
	}
	if err := booking.SetStatusCanceled(); err != nil {
		return err
	}
	antrianPoli, err := svc.Repositories.GetAntrianRepository().GetAntrianPoli(booking.AntrianPoliID, repositories.CANCEL_CONTEXT)
	if err != nil {
		return err
	}
	err = svc.EventPublisher.PublishAntrianEvent(antrianPoli)
	if err != nil {
		return err
	}
	antrianBranch, err := svc.Repositories.GetAntrianRepository().GetAntrianBranch(booking.AntrianBranchID, repositories.CANCEL_CONTEXT)
	if err != nil {
		return err
	}
	err = svc.EventPublisher.PublishAntrianEvent(antrianBranch)
	if err != nil {
		return err
	}
	err = svc.Repositories.GetBookingResository().Update(booking)
	if err != nil {
		return err
	}
	return svc.EventPublisher.PublishBookingEvent(booking)
}

func (svc ServiceImpl) UpdateAntrian(command commands.UpdateAntrianCommand) error {
	poli, err := svc.Repositories.GetPoliRepository().Get(command.PoliID)
	if err != nil {
		return err
	}
	antrianPoli := aggregates.NewEmptyAntrianPoli(poli, command.Tanggal)
	err = svc.Repositories.GetAntrianRepository().Save(antrianPoli)
	if err != nil {
		if err.Error() == repositories.ANTRIAN_EXIST {
			var err error
			antrianPoli, err = svc.Repositories.GetAntrianRepository().Get(command.GetAntrainPoliID())
			if err != nil {
				return err
			}
			if err := svc.Repositories.GetAntrianRepository().UpdateAntrianKuota(antrianPoli.ID, command.Kuota); err != nil {
				return err
			}
			antrianPoli, err = svc.Repositories.GetAntrianRepository().Get(command.GetAntrainPoliID())
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return svc.EventPublisher.PublishAntrianEvent(antrianPoli)

}
func (svc ServiceImpl) AddPoli(poli aggregates.Poli) error {
	return svc.Repositories.GetPoliRepository().Save(poli)
}
func (svc ServiceImpl) UpdatePoli(poli aggregates.Poli) error {
	return svc.Repositories.GetPoliRepository().Update(poli)
}

func (svc ServiceImpl) AddBranch(branch aggregates.Branch) error {
	return svc.Repositories.GetBranchRepository().Save(branch)
}
func (svc ServiceImpl) UpdateBranch(branch aggregates.Branch) error {
	return svc.Repositories.GetBranchRepository().Update(branch)
}

func (svc ServiceImpl) AddSchedule(schedule aggregates.Schedule) error {
	return svc.Repositories.GetScheduleRepository().Save(schedule)
}
func (svc ServiceImpl) UpdateSchedule(schedule aggregates.Schedule) error {
	return svc.Repositories.GetScheduleRepository().Update(schedule)
}
