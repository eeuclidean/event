package handlers

import (
	"event/user/aggregates"
	"event/user/service"
	"encoding/json"

	"github.com/eeuclidean/eventsourcing"
)

const (
	POLI_CHANNEL = "POLI_CHANNEL"
)

func NewPoliEventConsumerHandler(service service.Service, log func(functionName string, msg string)) PoliEventConsumerHandler {
	return PoliEventConsumerHandler{
		Service: service,
		Log:     log,
	}
}

type PoliEventConsumerHandler struct {
	Service service.Service
	Log     func(functionName string, msg string)
}

func (consumer PoliEventConsumerHandler) Apply(event eventsourcing.Event) error {
	switch event.EventName {
	case aggregates.POLI_EVENT_NAME:
		var poli aggregates.Poli
		err := json.Unmarshal([]byte(event.Data), &poli)
		if err != nil {
			return err
		}
		switch event.EventType {
		case aggregates.POLI_CREATED:
			if err := consumer.Service.AddPoli(poli); err != nil {
				consumer.Log("AddPoli", err.Error())
			}
			return err
		case aggregates.POLI_UPDATED:
			if err := consumer.Service.UpdatePoli(poli); err != nil {
				consumer.Log("UpdatePoli", err.Error())
			}
			return err
		default:
			return nil
		}
	case aggregates.SCHEDULE_EVENT_NAME:
		var schedule aggregates.Schedule
		err := json.Unmarshal([]byte(event.Data), &schedule)
		if err != nil {
			return err
		}
		switch event.EventType {
		case aggregates.SCHEDULE_CREATED:
			if err := consumer.Service.AddSchedule(schedule); err != nil {
				consumer.Log("AddSchedule", err.Error())
			}
			return err
		case aggregates.SCHEDULE_UPDATED:
			if err := consumer.Service.UpdateSchedule(schedule); err != nil {
				consumer.Log("UpdateSchedule", err.Error())
			}
			return err
		default:
			return nil
		}
	default:
		return nil
	}
}
