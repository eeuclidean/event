package publisher

import (
	"encoding/json"
	"event/user/aggregates"
	"os"

	"github.com/eeuclidean/eventsourcing"
	"github.com/eeuclidean/eventsourcing/publisher"
)

const (
	ANTRIAN_CHANNEL = "ANTRIAN_CHANNEL"
)

var antrianRedisEventPublisher *publisher.RedisAdapterDomainEventPublisher

func NewAntrianEventPublisher(antrian aggregates.Antrian) EventPublisher {
	return antrianEventPublisher{
		Antrian:     antrian,
		EventName:   "antrian",
		DataName:    "antrian",
		ServiceName: "booking_service",
	}
}

type antrianEventPublisher struct {
	Antrian     aggregates.Antrian
	EventName   string
	DataName    string
	ServiceName string
}

func (eventPublisher *antrianEventPublisher) Publish() (err error) {
	if antrianRedisEventPublisher == nil {
		antrianRedisEventPublisher, err = eventPublisher.initRedisEventPublisher()
		if err != nil {
			return
		}
	}
	err = antrianRedisEventPublisher.Publish(eventPublisher.getDomainEvent())
	return
}

func (eventPublisher *antrianEventPublisher) initRedisEventPublisher() (*publisher.RedisAdapterDomainEventPublisher, error) {
	return &publisher.RedisAdapterDomainEventPublisher{
		RedisUrl:         os.Getenv(REDIS_ADDRESS),
		DBIndex:          os.Getenv(REDIS_DB),
		Password:         os.Getenv(REDIS_PASSWORD),
		Channel:          os.Getenv(ANTRIAN_CHANNEL),
		SaveToEventStore: true,
		EventStore:       newEventStore(),
	}, nil
}

func (eventPublisher *antrianEventPublisher) getDomainEvent() eventsourcing.Event {
	data, _ := json.Marshal(eventPublisher.Antrian)
	return eventsourcing.Event{
		EventType: eventPublisher.Antrian.Status,
		EventName: eventPublisher.EventName,
		DataID:    eventPublisher.Antrian.ID,
		DataName:  eventPublisher.DataName,
		Data:      string(data),
		CreateBy:  eventPublisher.ServiceName,
	}
}
