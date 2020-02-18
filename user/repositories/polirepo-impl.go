package repositories

import (
	"event/user/aggregates"

	mgo "gopkg.in/mgo.v2"
)

type mongoAdapterPoliRepo struct {
	Collection  *mgo.Collection
	ConnChecker *connectionChecker
}

func (adapter mongoAdapterPoliRepo) Save(poli aggregates.Poli) error {
	return adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.Insert(poli)
	})
}

func (adapter mongoAdapterPoliRepo) Update(poli aggregates.Poli) error {
	return adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.UpdateId(poli.ID, poli)
	})
}

func (adapter mongoAdapterPoliRepo) Remove(poliid string) error {
	return adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.RemoveId(poliid)
	})
}

func (adapter mongoAdapterPoliRepo) Get(poliid string) (aggregates.Poli, error) {
	var poli aggregates.Poli
	err := adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.FindId(poliid).One(&poli)
	})
	return poli, err
}
