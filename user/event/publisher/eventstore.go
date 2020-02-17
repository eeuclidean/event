package publisher

import (
	"github.com/eeuclidean/eventsourcing/eventstore"
)

const (
	MONGO_EVENTSTORE_URL        = "MONGO_EVENTSTORE_URL"
	MONGO_EVENTSTORE_NAME       = "MONGO_EVENTSTORE_NAME"
	MONGO_EVENTSTORE_COLLECTION = "MONGO_EVENTSTORE_COLLECTION"
	MONGO_EVENTSTORE_USERNAME   = "MONGO_EVENTSTORE_USERNAME"
	MONGO_EVENTSTORE_PASSWORD   = "MONGO_EVENTSTORE_PASSWORD"
)

func newEventStore() *eventstore.MongoAdapterEventStore {
	return &eventstore.MongoAdapterEventStore{
		URL:        os.Getenv(MONGO_EVENTSTORE_URL),
		DB:         os.Getenv(MONGO_EVENTSTORE_NAME),
		Collection: os.Getenv(MONGO_EVENTSTORE_COLLECTION),
		Username:   os.Getenv(MONGO_EVENTSTORE_USERNAME),
		Password:   os.Getenv(MONGO_EVENTSTORE_PASSWORD),
	}
}
