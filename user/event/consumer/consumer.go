package consumer

import (
	"os"

	"event/user/service"

	"github.com/eeuclidean/eventsourcing/consumer"
)

const (
	REDIS_ADDRESS  = "REDIS_ADDRESS"
	REDIS_PASSWORD = "REDIS_PASSWORD"
	REDIS_DB       = "REDIS_DB"
)

const (
	REDIS_STREAM_GROUP = "REDIS_STREAM_GROUP"
)

const (
	POLI_CHANNEL   = "POLI_CHANNEL"
	BRANCH_CHANNEL = "BRANCH_CHANNEL"
)

func NewEventConsumer(svc service.Service, log func(functionName string, msg string)) (consumer.EventConsumer, error) {
	consumerName, err := os.Hostname()
	if err != nil {
		return consumer.RedisEventConsumer{}, err
	}
	handlers := map[string]consumer.EventConsumerHandler{
		os.Getenv(POLI_CHANNEL):   poliEventConsumerHandler{Service: svc, Log: log},
		os.Getenv(BRANCH_CHANNEL): branchEventConsumerHandler{Service: svc, Log: log},
	}
	return consumer.RedisEventConsumer{
		RedisURL:         os.Getenv(REDIS_ADDRESS),
		DBIndex:          os.Getenv(REDIS_DB),
		Password:         os.Getenv(REDIS_PASSWORD),
		Group:            os.Getenv(REDIS_STREAM_GROUP),
		Consumer:         consumerName,
		NewMessage:       false,
		Foreground:       false,
		HandlerConsumers: handlers,
		LogFunction:      log,
	}, nil
}
