package event_test

import (
	"context"
	"os"
	"event/user/aggregates"
	"event/user/event"
	"event/user/event/handlers"
	"event/user/repositories/polirepo"
	"event/user/service"
	"event/user/testing/config"
	"event/user/utils/utilsgenerator"
	"event/user/utils/utilsmongo"
	"testing"
	"time"

	"github.com/eeuclidean/eventsourcing"
	"github.com/eeuclidean/eventsourcing/consumer"

	"github.com/stretchr/testify/assert"
)

func TestEventPoliUpdated(t *testing.T) {
	assert := assert.New(t)
	config.RunConfig()
	var err error

	var svc service.Service
	{
		svc, err = service.NewService()
		assert.Nil(err)
	}
	var eventConsumers consumer.EventConsumer
	{
		eventConsumers, err = event.NewRedisEventConsumer(handlers.NewEventHandlers(svc, Log), Log)
		assert.Nil(err)
	}
	var cancel context.CancelFunc
	{
		var ctx context.Context
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
		err = eventConsumers.Run(ctx)
		assert.Nil(err, "error run event consumers")
	}
	domainEventPublisher, err := event.NewRedisEventPublisher()
	assert.Nil(err)
	collection, err := utilsmongo.GetCollection(os.Getenv(polirepo.POLI_COLL))
	assert.Nil(err)
	t.Run("Create Poli Success", func(t *testing.T) {
		id := utilsgenerator.NewID()
		poli := aggregates.Poli{
			ID:    id,
			Max:   100,
			Count: 0,
		}
		err := collection.Insert(poli)
		assert.Nil(err)
		poliEvent := eventsourcing.Event{
			Data:      `{"id":"` + id + `","max":120}`,
			DataID:    id,
			EventType: aggregates.POLI_UPDATED,
			DataName:  "poli",
			CreateBy:  "poli_service",
		}
		domainEventPublisher.Channel = os.Getenv("POLI_CHANNEL")
		err = domainEventPublisher.Publish(poliEvent)
		assert.Nil(err)

		time.Sleep(1 * time.Second)
		var savedPoli aggregates.Poli
		err = collection.FindId(id).One(&savedPoli)
		assert.Nil(err)
		assert.Equal(120, savedPoli.Max)
	})
}
