package repo_test

import (
	"os"
	"event/user/aggregates"
	"event/user/repositories/polirepo"
	"event/user/testing/config"
	"event/user/utils/utilsgenerator"
	"event/user/utils/utilsmongo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoliRepo(t *testing.T) {
	assert := assert.New(t)
	config.RunConfig()

	var poliRepo polirepo.PoliRepository
	{

		var err error
		poliRepo, err = polirepo.NewMongoAdapterPoliRepo()
		assert.Nil(err)
	}

	var poli aggregates.Poli
	{
		poli = aggregates.Poli{
			ID:       utilsgenerator.NewID(),
			TenantID: utilsgenerator.NewID(),
			BranchID: utilsgenerator.NewID(),
			Name:     utilsgenerator.NewID(),
			Max:      100,
			Count:    10,
		}
	}

	collection, err := utilsmongo.GetCollection(os.Getenv(polirepo.POLI_COLL))
	assert.Nil(err)

	t.Run("Save Poli", func(t *testing.T) {
		err := poliRepo.Save(poli)
		assert.Nil(err)
	})

	t.Run("Get Poli", func(t *testing.T) {
		savedPoli, err := poliRepo.Get(poli.ID)
		assert.Nil(err)
		assert.Equal(poli, savedPoli)
	})

	t.Run("Update Poli", func(t *testing.T) {
		newPoli := aggregates.Poli{
			ID:       poli.ID,
			TenantID: poli.TenantID,
			BranchID: poli.BranchID,
			Name:     poli.Name,
			Max:      poli.Max,
			Count:    122,
		}
		err := poliRepo.Update(newPoli)
		assert.Nil(err)

		var savedPoli aggregates.Poli
		err = collection.FindId(newPoli.ID).One(&savedPoli)
		assert.Nil(err)
		assert.Equal(newPoli, savedPoli)

	})

	t.Run("Delete Poli", func(t *testing.T) {
		err := poliRepo.Remove(poli.ID)
		assert.Nil(err)
	})
}
