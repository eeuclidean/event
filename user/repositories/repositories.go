package repositories

type Repositories interface {
	GetAntrianRepository() AntrianRepository
	GetBookingResository() BookingRepository
	GetBranchRepository() BranchRepository
	GetPoliRepository() PoliRepository
	GetScheduleRepository() ScheduleRepository
}
