package service

import (
	"bookingsvr/booking/event"
	"errors"
	"event/user/aggregates"
	"event/user/commands"
	"event/user/event/publisher"
	"event/user/repositories/antrianrepo"
	"event/user/repositories/bookingrepo"
	"event/user/repositories/branchrepo"
	"event/user/repositories/polirepo"
	"event/user/repositories/schedulerepo"

	"strconv"

	"github.com/eeuclidean/eventsourcing/publisher"
)

const ERR_PREFIX = "400 "

func NewService() (Service, error) {
	bookingRepo, err := bookingrepo.NewMongoAdapterBookingRepo()
	if err != nil {
		return ServiceImpl{}, err
	}
	poliRepo, err := polirepo.NewMongoAdapterPoliRepo()
	if err != nil {
		return ServiceImpl{}, err
	}
	antrianRepo, err := antrianrepo.NewMongoAdapterAntrianRepo()
	if err != nil {
		return ServiceImpl{}, err
	}
	branchRepo, err := branchrepo.NewMongoAdapterBranchRepo()
	if err != nil {
		return ServiceImpl{}, err
	}
	scheduleRepo, err := schedulerepo.NewMongoAdapterScheduleRepo()
	if err != nil {
		return ServiceImpl{}, err
	}
	bookingEventPublisher, err := event.NewBookingRedisEventPublisher()
	if err != nil {
		return ServiceImpl{}, err
	}
	antrianEventPublisher, err := event.NewAntrianRedisEventPublisher()
	if err != nil {
		return ServiceImpl{}, err
	}
	return ServiceImpl{
		BookingRepo:           bookingRepo,
		PoliRepo:              poliRepo,
		AntrianRepo:           antrianRepo,
		ScheduleRepo:          scheduleRepo,
		BranchRepo:            branchRepo,
		BookingEventPublisher: bookingEventPublisher,
		AntrianEventPublisher: antrianEventPublisher,
	}, nil
}

type ServiceImpl struct {
	BookingRepo           bookingrepo.BookingRepository
	PoliRepo              polirepo.PoliRepository
	AntrianRepo           antrianrepo.AntrianRepository
	BranchRepo            branchrepo.BranchRepository
	ScheduleRepo          schedulerepo.ScheduleRepository
	BookingEventPublisher publisher.DomainEventPublisher
	AntrianEventPublisher publisher.DomainEventPublisher
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
	poli, err := svc.PoliRepo.Get(command.PoliID)
	if err != nil {
		return "", err
	}
	date, month, year, _ := command.GetDateMonthYear()
	_, err = svc.ScheduleRepo.GetByDate(command.PoliID, year, month, date)
	if err == nil {
		return "", errors.New("Tanggal " + command.Tanggal + " libur")
	}
	if command.SubPoliID != "" {
		_, err = svc.ScheduleRepo.GetByDate(command.SubPoliID, year, month, date)
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
	branch, err := svc.BranchRepo.Get(poli.BranchID)
	if err != nil {
		return "", err
	}
	bookings, err := svc.BookingRepo.GetManyByPatientIDAndDate(command.BranchID, command.PatientID, command.Tanggal)
	if err != nil {
		return "", err
	}
	if branch.MaxBookingPerDay > 0 && len(bookings) >= branch.MaxBookingPerDay {
		return "", errors.New("Limit Booking " + strconv.Itoa(branch.MaxBookingPerDay) + " kali")
	}
	antrianBranch := aggregates.NewAntrianBranch(poli, command.Tanggal)
	err = svc.AntrianRepo.Save(antrianBranch)
	if err != nil {
		if err.Error() == antrianrepo.ANTRIAN_EXIST {
			var err error
			antrianBranch, err = svc.AntrianRepo.GetAntrianBranch(command.GetAntrainBranchID(), antrianrepo.REGISTER_CONTEXT)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	if err := svc.AntrianEventPublisher.Publish(antrianBranch.GetDomainEvent()); err != nil {
		return "", err
	}
	antrianPoli := aggregates.NewAntrianPoli(poli, command.Tanggal)
	err = svc.AntrianRepo.Save(antrianPoli)
	if err != nil {
		if err.Error() == antrianrepo.ANTRIAN_EXIST {
			var err error
			antrianPoli, err = svc.AntrianRepo.GetAntrianPoli(command.GetAntrainPoliID(), antrianrepo.REGISTER_CONTEXT)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	if err := svc.AntrianEventPublisher.Publish(antrianPoli.GetDomainEvent()); err != nil {
		return "", err
	}
	booking := aggregates.NewBooking(command, antrianBranch.Terisi, poli.PayAmount)
	err = svc.BookingEventPublisher.Publish(booking.GetDomainEvent())
	if err != nil {
		return "", err
	}
	err = svc.BookingRepo.Save(booking)
	if err != nil {
		return "", err
	}
	return booking.ID, nil
}
func (svc ServiceImpl) CallBooking(command commands.CallBookingCommand) error {
	booking, err := svc.BookingRepo.Get(command.ID)
	if err != nil {
		return err
	}
	if booking.IsToday() {
		booking.Call()
		err = svc.BookingRepo.Update(booking)
		if err != nil {
			return err
		}
		return svc.BookingEventPublisher.Publish(booking.GetDomainEvent())
	}
	return errors.New("Bukan Hari Ini")
}
func (svc ServiceImpl) CheckInLoketBooking(command commands.LoketCheckinBookingCommand) error {
	booking, err := svc.BookingRepo.Get(command.ID)
	if err != nil {
		return err
	}
	if booking.IsToday() {
		if err := booking.SetStatusLoketCheckin(command); err != nil {
			return err
		}
		antrianBranch, err := svc.AntrianRepo.GetAntrianBranch(booking.AntrianBranchID, antrianrepo.CHECKIN_CONTEXT)
		if err != nil {
			return err
		}
		if err := svc.AntrianEventPublisher.Publish(antrianBranch.GetDomainEvent()); err != nil {
			return err
		}
		antrianPoli, err := svc.AntrianRepo.GetAntrianPoli(booking.AntrianPoliID, antrianrepo.CHECKIN_CONTEXT)
		if err != nil {
			return err
		}
		if err := svc.AntrianEventPublisher.Publish(antrianPoli.GetDomainEvent()); err != nil {
			return err
		}
		booking.NoAntrianPoli = antrianPoli.Checkin
		err = svc.BookingRepo.Update(booking)
		if err != nil {
			return err
		}
		return svc.BookingEventPublisher.Publish(booking.GetDomainEvent())
	}
	return errors.New("Bukan Hari Ini")
}
func (svc ServiceImpl) CheckInPoliBooking(command commands.PoliCheckinBookingCommand) error {
	booking, err := svc.BookingRepo.Get(command.ID)
	if err != nil {
		return err
	}
	if booking.IsToday() {
		if err := booking.SetStatusPoliCheckIn(command.By); err != nil {
			return err
		}
		antrianPoli, err := svc.AntrianRepo.GetAntrianPoli(booking.AntrianPoliID, antrianrepo.FINISH_CONTEXT)
		if err != nil {
			return err
		}
		if err := svc.AntrianEventPublisher.Publish(antrianPoli.GetDomainEvent()); err != nil {
			return err
		}
		err = svc.BookingRepo.Update(booking)
		if err != nil {
			return err
		}
		return svc.BookingEventPublisher.Publish(booking.GetDomainEvent())
	}
	return errors.New("Bukan Hari Ini")
}

func (svc ServiceImpl) PayBooking(command commands.PayBookingCommand) error {
	booking, err := svc.BookingRepo.Get(command.ID)
	if err != nil {
		return err
	}
	if booking.IsToday() {
		if err := booking.SetStatusPayed(command); err != nil {
			return err
		}
		err = svc.BookingRepo.Update(booking)
		if err != nil {
			return err
		}
		return svc.BookingEventPublisher.Publish(booking.GetDomainEvent())
	}
	return errors.New("Bukan Hari Ini")
}
func (svc ServiceImpl) CancelBooking(command commands.CancelBookingCommand) error {
	booking, err := svc.BookingRepo.Get(command.ID)
	if err != nil {
		return err
	}
	if err := booking.SetStatusCanceled(); err != nil {
		return err
	}
	antrianPoli, err := svc.AntrianRepo.GetAntrianPoli(booking.AntrianPoliID, antrianrepo.CANCEL_CONTEXT)
	if err != nil {
		return err
	}
	err = svc.AntrianEventPublisher.Publish(antrianPoli.GetDomainEvent())
	if err != nil {
		return err
	}
	antrianBranch, err := svc.AntrianRepo.GetAntrianBranch(booking.AntrianBranchID, antrianrepo.CANCEL_CONTEXT)
	if err != nil {
		return err
	}
	err = svc.AntrianEventPublisher.Publish(antrianBranch.GetDomainEvent())
	if err != nil {
		return err
	}
	err = svc.BookingRepo.Update(booking)
	if err != nil {
		return err
	}
	return svc.BookingEventPublisher.Publish(booking.GetDomainEvent())
}

func (svc ServiceImpl) UpdateAntrian(command commands.UpdateAntrianCommand) error {
	poli, err := svc.PoliRepo.Get(command.PoliID)
	if err != nil {
		return err
	}
	antrianPoli := aggregates.NewEmptyAntrianPoli(poli, command.Tanggal)
	err = svc.AntrianRepo.Save(antrianPoli)
	if err != nil {
		if err.Error() == antrianrepo.ANTRIAN_EXIST {
			var err error
			antrianPoli, err = svc.AntrianRepo.Get(command.GetAntrainPoliID())
			if err != nil {
				return err
			}
			if err := svc.AntrianRepo.UpdateAntrianKuota(antrianPoli.ID, command.Kuota); err != nil {
				return err
			}
			antrianPoli, err = svc.AntrianRepo.Get(command.GetAntrainPoliID())
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return svc.AntrianEventPublisher.Publish(antrianPoli.GetDomainEvent())

}
func (svc ServiceImpl) AddPoli(poli aggregates.Poli) error {
	return svc.PoliRepo.Save(poli)
}
func (svc ServiceImpl) UpdatePoli(poli aggregates.Poli) error {
	return svc.PoliRepo.Update(poli)
}

func (svc ServiceImpl) AddBranch(branch aggregates.Branch) error {
	return svc.BranchRepo.Save(branch)
}
func (svc ServiceImpl) UpdateBranch(branch aggregates.Branch) error {
	return svc.BranchRepo.Update(branch)
}

func (svc ServiceImpl) AddSchedule(schedule aggregates.Schedule) error {
	return svc.ScheduleRepo.Save(schedule)
}
func (svc ServiceImpl) UpdateSchedule(schedule aggregates.Schedule) error {
	return svc.ScheduleRepo.Update(schedule)
}
