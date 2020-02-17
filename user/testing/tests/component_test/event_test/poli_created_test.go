package event_test

import (
	"context"
	"fmt"
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

func Log(functionName string, msg string) {
	fmt.Println(functionName + " : " + msg)
}

func TestEventPoliCreated(t *testing.T) {
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
		poliEvent := eventsourcing.Event{
			Data:      `{"id":"` + id + `"}`,
			DataID:    id,
			EventType: aggregates.POLI_CREATED,
			DataName:  "poli",
			CreateBy:  "poli_service",
		}
		domainEventPublisher.Channel = os.Getenv("POLI_CHANNEL")
		err = domainEventPublisher.Publish(poliEvent)
		assert.Nil(err)

		time.Sleep(1 * time.Second)
		var poli aggregates.Poli
		err = collection.FindId(id).One(&poli)
		assert.Nil(err)
		assert.NotEqual(aggregates.Poli{}, poli)
	})
}
