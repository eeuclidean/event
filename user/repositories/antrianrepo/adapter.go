package antrianrepo

import (
	"event/user/aggregates"
	"event/user/utils/utilsmongo"
	"errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

const (
	ANTRIAN_COLL = "ANTRIAN_COLL"
)

const (
	ANTRIAN_EXIST = "antrian_exist"
)

const (
	REGISTER_CONTEXT = "terisi"
	CHECKIN_CONTEXT  = "check_in"
	FINISH_CONTEXT   = "selesai"
	CANCEL_CONTEXT   = "cancel"
)

func NewMongoAdapterAntrianRepo() (AntrianRepository, error) {
	db, err := utilsmongo.MongoDBLogin()
	if err != nil {
		return MongoAdapterAntrianRepo{}, err
	}
	return MongoAdapterAntrianRepo{
		Collection:     db.C(os.Getenv(ANTRIAN_COLL)),
		CircuitBreaker: utilsmongo.NewCircuitBreaker(db),
	}, nil
}

type MongoAdapterAntrianRepo struct {
	Collection     *mgo.Collection
	CircuitBreaker *utilsmongo.MongoCircuitBreaker
}

func (adapter MongoAdapterAntrianRepo) Save(antrian aggregates.Antrian) error {
	return adapter.CircuitBreaker.Execute(func() error {
		err := adapter.Collection.Insert(antrian)
		if err != nil {
			if mgo.IsDup(err) {
				return errors.New(ANTRIAN_EXIST)
			}
		}
		return nil
	})
}

func (adapter MongoAdapterAntrianRepo) UpdateAntrianKuota(antrianid string, kuota int) error {
	return adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.UpdateId(antrianid, bson.M{"$set": bson.M{"kuota": kuota, "status": aggregates.ANTRIAN_UPDATED}})
	})
}

func (adapter MongoAdapterAntrianRepo) Get(antrianid string) (aggregates.Antrian, error) {
	var antrian aggregates.Antrian
	err := adapter.CircuitBreaker.Execute(func() error {
		return adapter.Collection.FindId(antrianid).One(&antrian)
	})
	return antrian, err
}
func (adapter MongoAdapterAntrianRepo) GetAntrianPoli(antrianid string, context string) (aggregates.Antrian, error) {
	var antrian aggregates.Antrian
	err := adapter.CircuitBreaker.Execute(func() error {
		if context == REGISTER_CONTEXT {
			var existingAntrian aggregates.Antrian
			err := adapter.Collection.FindId(antrianid).One(&existingAntrian)
			if err != nil {
				return err
			}
			err = adapter.Collection.Update(
				bson.M{"_id": antrianid, context: bson.M{"$lt": existingAntrian.Kuota}, "type": aggregates.ANTRIAN_POLI},
				bson.M{"$inc": bson.M{context: 1}, "$set": bson.M{"status": aggregates.ANTRIAN_UPDATED}},
			)
			if err != nil {
				if err.Error() == "not found" {
					return errors.New("Antrian Penuh")
				}
				return err
			}
		} else if context == CANCEL_CONTEXT {
			err := adapter.Collection.UpdateId(antrianid, bson.M{"$inc": bson.M{REGISTER_CONTEXT: -1}, "$set": bson.M{"status": aggregates.ANTRIAN_UPDATED}})
			if err != nil {
				return err
			}
		} else {
			err := adapter.Collection.UpdateId(antrianid, bson.M{"$inc": bson.M{context: 1}, "$set": bson.M{"status": aggregates.ANTRIAN_UPDATED}})
			if err != nil {
				return err
			}
		}
		return adapter.Collection.FindId(antrianid).One(&antrian)
	})
	return antrian, err
}

func (adapter MongoAdapterAntrianRepo) GetAntrianBranch(antrianid string, context string) (aggregates.Antrian, error) {
	var antrian aggregates.Antrian
	err := adapter.CircuitBreaker.Execute(func() error {
		if context == REGISTER_CONTEXT {
			err := adapter.Collection.UpdateId(antrianid, bson.M{"$inc": bson.M{context: 1}, "$set": bson.M{"status": aggregates.ANTRIAN_UPDATED}})
			if err != nil {
				return err
			}
		} else if context == CANCEL_CONTEXT {
			err := adapter.Collection.UpdateId(antrianid, bson.M{"$inc": bson.M{REGISTER_CONTEXT: -1}, "$set": bson.M{"status": aggregates.ANTRIAN_UPDATED}})
			if err != nil {
				return err
			}
		} else {
			err := adapter.Collection.UpdateId(antrianid, bson.M{"$inc": bson.M{FINISH_CONTEXT: 1, CHECKIN_CONTEXT: 1}, "$set": bson.M{"status": aggregates.ANTRIAN_UPDATED}})
			if err != nil {
				return err
			}
		}
		return adapter.Collection.FindId(antrianid).One(&antrian)
	})
	return antrian, err
}
