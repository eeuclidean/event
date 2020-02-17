package handlers

import (
	"event/user/service"
	"os"

	"github.com/eeuclidean/eventsourcing/consumer"
)

func NewEventHandlers(service service.Service, log func(functionName string, msg string)) map[string]consumer.EventConsumerHandler {
	return map[string]consumer.EventConsumerHandler{
		os.Getenv(POLI_CHANNEL):   NewPoliEventConsumerHandler(service, log),
		os.Getenv(BRANCH_CHANNEL): NewBranchEventConsumerHandler(service, log),
	}
}
