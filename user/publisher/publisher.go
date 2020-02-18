package publisher

import (
	"encoding/json"
	"event/user/aggregates"
	"os"

	"github.com/eeuclidean/eventsourcing/publisher"
)

const (
	BOOKING_CHANNEL = "BOOKING_CHANNEL"
	ANTRIAN_CHANNEL = "ANTRIAN_CHANNEL"
)

func NewEventPublisher() EventPublisher {
	antrianRedisEventPublisher := publisherAdapterPrototype
	antrianRedisEventPublisher.Channel = os.Getenv(ANTRIAN_CHANNEL)
	bookingRedisEventPublisher := publisherAdapterPrototype
	bookingRedisEventPublisher.Channel = os.Getenv(BOOKING_CHANNEL)
	return &eventPublisherImpl{
		AntrianPublisherAdapter: &antrianRedisEventPublisher,
		BookingPublisherAdapter: &bookingRedisEventPublisher,
	}
}

type EventPublisher interface {
	PublishBookingEvent(booking aggregates.Booking) error
	PublishAntrianEvent(antrian aggregates.Antrian) error
}

type eventPublisherImpl struct {
	AntrianPublisherAdapter *publisher.RedisAdapterDomainEventPublisher
	BookingPublisherAdapter *publisher.RedisAdapterDomainEventPublisher
}

func (factory *eventPublisherImpl) PublishBookingEvent(booking aggregates.Booking) error {
	data, _ := json.Marshal(booking)
	publisher := &domainPublisher{
		EventName:           "booking",
		DataName:            "booking",
		ServiceName:         "booking_service",
		DataID:              booking.ID,
		Channel:             BOOKING_CHANNEL,
		EventType:           booking.Status,
		Data:                string(data),
		RedisEventPublisher: factory.BookingPublisherAdapter,
	}
	return publisher.Publish()
}

func (factory *eventPublisherImpl) PublishAntrianEvent(antrian aggregates.Antrian) error {
	data, _ := json.Marshal(antrian)
	publisher := &domainPublisher{
		EventName:           "antrian",
		DataName:            "antrian",
		ServiceName:         "booking_service",
		DataID:              antrian.ID,
		Channel:             ANTRIAN_CHANNEL,
		EventType:           antrian.Status,
		Data:                string(data),
		RedisEventPublisher: factory.AntrianPublisherAdapter,
	}
	return publisher.Publish()
}
