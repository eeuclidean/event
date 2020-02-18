package aggregates

import "gopkg.in/mgo.v2/bson"

func generateID() string {
	return bson.NewObjectId().Hex()
}
