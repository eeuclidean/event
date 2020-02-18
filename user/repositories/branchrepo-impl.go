package repositories

import (
	"event/user/aggregates"

	mgo "gopkg.in/mgo.v2"
)

type mongoAdapterBranchRepo struct {
	Collection  *mgo.Collection
	ConnChecker *connectionChecker
}

func (adapter mongoAdapterBranchRepo) Save(branch aggregates.Branch) error {
	return adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.Insert(branch)
	})
}

func (adapter mongoAdapterBranchRepo) Update(branch aggregates.Branch) error {
	return adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.UpdateId(branch.ID, branch)
	})
}

func (adapter mongoAdapterBranchRepo) Get(branchid string) (aggregates.Branch, error) {
	var branch aggregates.Branch
	err := adapter.ConnChecker.Execute(func() error {
		return adapter.Collection.FindId(branchid).One(&branch)
	})
	return branch, err
}
