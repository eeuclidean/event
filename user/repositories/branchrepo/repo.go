package branchrepo

import (
	"event/user/aggregates"
)

type BranchRepository interface {
	Save(branch aggregates.Branch) error
	Update(branch aggregates.Branch) error
	Get(branchid string) (aggregates.Branch, error)
}
