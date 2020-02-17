package publisher

import (
	"os"

	"github.com/eeuclidean/eventsourcing"
	"github.com/eeuclidean/eventsourcing/publisher"
)

const (
	REDIS_ADDRESS  = "REDIS_ADDRESS"
	REDIS_PASSWORD = "REDIS_PASSWORD"
	REDIS_DB       = "REDIS_DB"
)

type publisherUtils struct {
	EventName   string
	DataName    string
	ServiceName string
	Channel     string
	Status      string
	Data        string
	DataID      string
}

func (utils *publisherUtils) initRedisEventPublisher() (*publisher.RedisAdapterDomainEventPublisher, error) {
	return &publisher.RedisAdapterDomainEventPublisher{
		RedisUrl:         os.Getenv(REDIS_ADDRESS),
		DBIndex:          os.Getenv(REDIS_DB),
		Password:         os.Getenv(REDIS_PASSWORD),
		Channel:          os.Getenv(utils.Channel),
		SaveToEventStore: true,
		EventStore:       newEventStore(),
	}, nil
}

func (utils *publisherUtils) getDomainEvent() eventsourcing.Event {
	return eventsourcing.Event{
		EventType: utils.Status,
		EventName: utils.EventName,
		DataID:    utils.DataID,
		DataName:  utils.DataName,
		Data:      string(utils.Data),
		CreateBy:  utils.ServiceName,
	}
}
