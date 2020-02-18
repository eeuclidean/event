package repositories

import "os"

const (
	SCHEDULE_COLL = "SCHEDULE_COLL"
	ANTRIAN_COLL  = "ANTRIAN_COLL"
	BOOKING_COLL  = "BOOKING_COLL"
	BRANCH_COLL   = "BRANCH_COLL"
	POLI_COLL     = "POLI_COLL"
)

func NewRepositories() (Repositories, error) {
	db, err := mongoDBLogin()
	if err != nil {
		return &repositoriesImpl{}, err
	}
	return &repositoriesImpl{
		AntrianRepo: mongoAdapterAntrianRepo{
			Collection:  db.C(os.Getenv(ANTRIAN_COLL)),
			ConnChecker: &connectionChecker{Database: db},
		},
		BookingRepo: mongoAdapterBookingRepo{
			Collection:  db.C(os.Getenv(BOOKING_COLL)),
			ConnChecker: &connectionChecker{Database: db},
		},
		BranchRepo: mongoAdapterBranchRepo{
			Collection:  db.C(os.Getenv(BRANCH_COLL)),
			ConnChecker: &connectionChecker{Database: db},
		},
		PoliRepo: mongoAdapterPoliRepo{
			Collection:  db.C(os.Getenv(POLI_COLL)),
			ConnChecker: &connectionChecker{Database: db},
		},
		ScheduleRepo: mongoAdapterScheduleRepo{
			Collection:  db.C(os.Getenv(SCHEDULE_COLL)),
			ConnChecker: &connectionChecker{Database: db},
		},
	}, nil
}

type repositoriesImpl struct {
	AntrianRepo  AntrianRepository
	BookingRepo  BookingRepository
	BranchRepo   BranchRepository
	PoliRepo     PoliRepository
	ScheduleRepo ScheduleRepository
}

func (repo *repositoriesImpl) GetAntrianRepository() AntrianRepository {
	return repo.AntrianRepo
}
func (repo *repositoriesImpl) GetBookingResository() BookingRepository {
	return repo.BookingRepo
}
func (repo *repositoriesImpl) GetBranchRepository() BranchRepository {
	return repo.BranchRepo
}
func (repo *repositoriesImpl) GetPoliRepository() PoliRepository {
	return repo.PoliRepo
}
func (repo *repositoriesImpl) GetScheduleRepository() ScheduleRepository {
	return repo.ScheduleRepo
}
