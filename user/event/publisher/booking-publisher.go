package publisher

import (
	"os"

	"github.com/eeuclidean/eventsourcing/publisher"
)

const (
	BOOKING_CHANNEL = "BOOKING_CHANNEL"
)

func newBookingRedisEventPublisher() (*publisher.RedisAdapterDomainEventPublisher, error) {
	return &publisher.RedisAdapterDomainEventPublisher{
		RedisUrl:         os.Getenv(REDIS_ADDRESS),
		DBIndex:          os.Getenv(REDIS_DB),
		Password:         os.Getenv(REDIS_PASSWORD),
		Channel:          os.Getenv(BOOKING_CHANNEL),
		SaveToEventStore: true,
		EventStore:       newEventStore(),
	}, nil
}
