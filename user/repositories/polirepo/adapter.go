package polirepo

import (
	"event/user/aggregates"
	"event/user/utils/utilsmongo"
	"os"

	mgo "gopkg.in/mgo.v2"
)

const (
	POLI_COLL = "POLI_COLL"
)

func NewMongoAdapterPoliRepo() (PoliRepository, error) {
	db, err := utilsmongo.MongoDBLogin()
	if err != nil {
		return MongoAdapterPoliRepo{}, err
	}
	return MongoAdapterPoliRepo{
		Collection:     db.C(os.Getenv(POLI_COLL)),
		CircuitBreaker: utilsmongo.NewCircuitBreaker(db),
	}, nil
}

type MongoAdapterPoliRepo struct {
	Collection     *mgo.Collection
	CircuitBreaker *utilsmongo.MongoCircuitBreaker
}

func (adapter MongoAdapterPoliRepo) Save(poli aggregates.Poli) error {
	return adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.Insert(poli)
	})
}

func (adapter MongoAdapterPoliRepo) Update(poli aggregates.Poli) error {
	return adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.UpdateId(poli.ID, poli)
	})
}

func (adapter MongoAdapterPoliRepo) Remove(poliid string) error {
	return adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.RemoveId(poliid)
	})
}

func (adapter MongoAdapterPoliRepo) Get(poliid string) (aggregates.Poli, error) {
	var poli aggregates.Poli
	err := adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.FindId(poliid).One(&poli)
	})
	return poli, err
}
