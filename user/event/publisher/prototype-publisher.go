package publisher

import (
	"os"

	"github.com/eeuclidean/eventsourcing/publisher"
)

const (
	REDIS_ADDRESS  = "REDIS_ADDRESS"
	REDIS_PASSWORD = "REDIS_PASSWORD"
	REDIS_DB       = "REDIS_DB"
)

var publisherAdapterPrototype = publisher.RedisAdapterDomainEventPublisher{
	RedisUrl:         os.Getenv(REDIS_ADDRESS),
	DBIndex:          os.Getenv(REDIS_DB),
	Password:         os.Getenv(REDIS_PASSWORD),
	SaveToEventStore: false,
}
