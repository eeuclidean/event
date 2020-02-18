package consumer

import (
	"encoding/json"
	"event/user/aggregates"
	"event/user/service"

	"github.com/eeuclidean/eventsourcing"
)

type branchEventConsumerHandler struct {
	Service service.Service
	Log     func(functionName string, msg string)
}

func (consumer branchEventConsumerHandler) Apply(event eventsourcing.Event) error {
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
