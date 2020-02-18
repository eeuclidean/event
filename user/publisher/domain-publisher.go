package publisher

import (
	"github.com/eeuclidean/eventsourcing"
	"github.com/eeuclidean/eventsourcing/publisher"
)

type domainPublisher struct {
	EventName           string
	DataName            string
	ServiceName         string
	Channel             string
	EventType           string
	Data                string
	DataID              string
	RedisEventPublisher *publisher.RedisAdapterDomainEventPublisher
}

func (eventPublisher *domainPublisher) Publish() error {
	return eventPublisher.RedisEventPublisher.Publish(eventPublisher.getDomainEvent())
}

func (eventPublisher *domainPublisher) getDomainEvent() eventsourcing.Event {
	return eventsourcing.Event{
		EventType: eventPublisher.EventType,
		EventName: eventPublisher.EventName,
		DataID:    eventPublisher.DataID,
		DataName:  eventPublisher.DataName,
		Data:      eventPublisher.Data,
		CreateBy:  eventPublisher.ServiceName,
	}
}
