package repositories

import (
	"event/user/aggregates"
)

type AntrianRepository interface {
	Save(antrian aggregates.Antrian) error
	UpdateAntrianKuota(antrianid string, kuota int) error
	Get(antrianid string) (aggregates.Antrian, error)
	GetAntrianPoli(antrianid string, context string) (aggregates.Antrian, error)
	GetAntrianBranch(antrianid string, context string) (aggregates.Antrian, error)
}
