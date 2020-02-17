package branchrepo

import (
	"event/user/aggregates"
	"event/user/utils/utilsmongo"
	"os"

	mgo "gopkg.in/mgo.v2"
)

const (
	BRANCH_COLL = "BRANCH_COLL"
)

func NewMongoAdapterBranchRepo() (BranchRepository, error) {
	db, err := utilsmongo.MongoDBLogin()
	if err != nil {
		return MongoAdapterBranchRepo{}, err
	}
	return MongoAdapterBranchRepo{
		Collection:     db.C(os.Getenv(BRANCH_COLL)),
		CircuitBreaker: utilsmongo.NewCircuitBreaker(db),
	}, nil
}

type MongoAdapterBranchRepo struct {
	Collection     *mgo.Collection
	CircuitBreaker *utilsmongo.MongoCircuitBreaker
}

func (adapter MongoAdapterBranchRepo) Save(branch aggregates.Branch) error {
	return adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.Insert(branch)
	})
}

func (adapter MongoAdapterBranchRepo) Update(branch aggregates.Branch) error {
	return adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.UpdateId(branch.ID, branch)
	})
}

func (adapter MongoAdapterBranchRepo) Get(branchid string) (aggregates.Branch, error) {
	var branch aggregates.Branch
	err := adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.FindId(branchid).One(&branch)
	})
	return branch, err
}
