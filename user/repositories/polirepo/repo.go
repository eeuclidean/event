package polirepo

import (
	"event/user/aggregates"
)

type PoliRepository interface {
	Save(poli aggregates.Poli) error
	Update(poli aggregates.Poli) error
	Remove(poliid string) error
	Get(poliid string) (aggregates.Poli, error)
}
