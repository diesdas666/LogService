package service

import (
	"encoding/json"
	"example_consumer/internal/asyncApi/events"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
)

type LogServiceInterface interface {
	UpdateStatusOfLoadingStation(msg interface{})
}
type LogService struct {
}

func (logService LogService) UpdateStatusOfLoadingStation(msg interface{}) error {
	message, ok := msg.(kafka.Message)
	if !ok {
		return errors.New("type is not kafka message")
	}
	eventBody := events.EventBodyMessage{}
	err := json.Unmarshal(message.Value, &eventBody)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal event")
	}
	err = Validate(eventBody)
	if err != nil {
		return errors.Wrap(err, "validation of eventBody results in failure")
	}
	// todo persist
	//err = api.identityService.DeactivateUserByCustomerNumber(ctx, eventBody.CustomerNumber)
	//if err != nil {
	//	return errors.Wrap(err, "unable to deactivate identity")
	//}
	return nil

}

func Validate(body events.EventBodyMessage) error {
	payload := body.Payload
	println(payload)
	if true {
		return nil
	} else {
		return errors.New("type is not kafka message")
	}
}
