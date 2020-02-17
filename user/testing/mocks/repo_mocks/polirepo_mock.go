package repo_mocks

import (
	"event/user/aggregates"

	"github.com/stretchr/testify/mock"
)

type PoliRepoMock struct {
	mock.Mock
}

func (mock *PoliRepoMock) Save(poli aggregates.Poli) error {
	return nil
}
func (mock *PoliRepoMock) Update(poli aggregates.Poli) error {
	return nil
}
func (mock *PoliRepoMock) Remove(poliid string) error {
	return nil
}
func (mock *PoliRepoMock) Get(poliid string) (aggregates.Poli, error) {
	args := mock.Called(poliid)
	return args.Get(0).(aggregates.Poli), args.Error(1)
}
