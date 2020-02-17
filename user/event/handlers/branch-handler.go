package handlers

import (
	"event/user/aggregates"
	"event/user/service"
	"encoding/json"

	"github.com/eeuclidean/eventsourcing"
)

const (
	BRANCH_CHANNEL = "BRANCH_CHANNEL"
)

func NewBranchEventConsumerHandler(service service.Service, log func(functionName string, msg string)) BranchEventConsumerHandler {
	return BranchEventConsumerHandler{
		Service: service,
		Log:     log,
	}
}

type BranchEventConsumerHandler struct {
	Service service.Service
	Log     func(functionName string, msg string)
}

func (consumer BranchEventConsumerHandler) Apply(event eventsourcing.Event) error {
	consumer.Log("Apply", "start")
	var branch aggregates.Branch
	err := json.Unmarshal([]byte(event.Data), &branch)
	if err != nil {
		return err
	}
	switch event.EventType {
	case aggregates.BRANCH_CREATED:
		if err := consumer.Service.AddBranch(branch); err != nil {
			consumer.Log("AddBranch", err.Error())
		}
		return err
	case aggregates.BRANCH_UPDATED:
		if err := consumer.Service.UpdateBranch(branch); err != nil {
			consumer.Log("UpdateBranch", err.Error())
		}
		return err
	default:
		return nil
	}
}
