package events

import (
	"context"
	"encoding/json"
	"github.com/c0olix/goChan"
	"github.com/c0olix/goChan/kafka"
	"github.com/go-playground/validator/v10"

	"github.com/pkg/errors"
	kafkaGo "github.com/segmentio/kafka-go"
	"reflect"
	"time"
)

const (
	CustomerDataCreateTopicName   = "CUSTOMER_DATA_CREATE"
	CustomerDataUpdateTopicName   = "CUSTOMER_DATA_UPDATE"
	CustomerDeactivationTopicName = "CUSTOMER_DEACTIVATION"
	CustomerDeletionTopicName     = "CUSTOMER_DELETION"
	CustomerStatusUpdateTopicName = "CUSTOMER_STATUS_UPDATE"
	TemplateMessageTopicName      = "TEMPLATE_MESSAGE"
)

var validate *validator.Validate = validator.New()

func UnmarshalEvent(proto interface{}, messageArray []byte) (interface{}, error) {
	eventType := reflect.TypeOf(proto)
	value := reflect.New(eventType).Interface()
	err := json.Unmarshal(messageArray, value)
	return value, err
}

// CustomerDeleteEventBody Event that holds the information needed to delete a customer
type CustomerDeleteEventBody struct {
	CustomerNumber string `json:"customerNumber" validate:"required"`
}

// Validate validates the CustomerDeleteEventBody) struct for its requirements
func (thiz CustomerDeleteEventBody) Validate() error {
	return validate.Struct(thiz)
}

// NewCustomerDeleteEventBody creates a new struct of the CustomerDeleteEventBody type and validates its content
func NewCustomerDeleteEventBody(customerNumber string) (*CustomerDeleteEventBody, error) {
	customerDeleteEventBody := CustomerDeleteEventBody{
		CustomerNumber: customerNumber,
	}
	err := customerDeleteEventBody.Validate()
	if err != nil {
		return nil, err
	}
	return &customerDeleteEventBody, nil
}

// CustomerStatusEventBody Event that holds information about the customer status, for example if the customer is allowed to order
type CustomerStatusEventBody struct {
	CustomerNumber string `json:"customerNumber" validate:"required"`
	Email          string `json:"email" validate:"required"`
	Locked         bool   `json:"locked"`
}

// Validate validates the CustomerStatusEventBody) struct for its requirements
func (thiz CustomerStatusEventBody) Validate() error {
	return validate.Struct(thiz)
}

// NewCustomerStatusEventBody creates a new struct of the CustomerStatusEventBody type and validates its content
func NewCustomerStatusEventBody(customerNumber string, email string, locked bool) (*CustomerStatusEventBody, error) {
	customerStatusEventBody := CustomerStatusEventBody{
		CustomerNumber: customerNumber,
		Email:          email,
		Locked:         locked,
	}
	err := customerStatusEventBody.Validate()
	if err != nil {
		return nil, err
	}
	return &customerStatusEventBody, nil
}

// CustomerRegistrationEventBody Event which holds information about a newly registered customer
type CustomerRegistrationEventBody struct {
	Context        *string    `json:"context,omitempty"`
	CustomerNumber string     `json:"customerNumber" validate:"required"`
	DateOfBirth    *time.Time `json:"dateOfBirth,omitempty"`
	Email          string     `json:"email" validate:"required"`
	FirstName      *string    `json:"firstName,omitempty"`
	LastName       *string    `json:"lastName,omitempty"`
	Salutation     *string    `json:"salutation,omitempty"`
}

// Validate validates the CustomerRegistrationEventBody) struct for its requirements
func (thiz CustomerRegistrationEventBody) Validate() error {
	return validate.Struct(thiz)
}

// NewCustomerRegistrationEventBody creates a new struct of the CustomerRegistrationEventBody type and validates its content
func NewCustomerRegistrationEventBody(context *string, customerNumber string, dateOfBirth *time.Time, email string, firstName *string, lastName *string, salutation *string) (*CustomerRegistrationEventBody, error) {
	customerRegistrationEventBody := CustomerRegistrationEventBody{
		Context:        context,
		CustomerNumber: customerNumber,
		DateOfBirth:    dateOfBirth,
		Email:          email,
		FirstName:      firstName,
		LastName:       lastName,
		Salutation:     salutation,
	}
	err := customerRegistrationEventBody.Validate()
	if err != nil {
		return nil, err
	}
	return &customerRegistrationEventBody, nil
}

// TemplateMailMessageBody Event that holds information needed to fill a template mail message
type TemplateMailMessageBody struct {
	Files        []TemplateMessageBodyFile `json:"files" validate:"required"`
	Recipient    string                    `json:"recipient" validate:"required"`
	TemplateData map[string]string         `json:"templateData" validate:"required"`
	TemplateId   string                    `json:"templateId" validate:"required"`
}

// Validate validates the TemplateMailMessageBody) struct for its requirements
func (thiz TemplateMailMessageBody) Validate() error {
	return validate.Struct(thiz)
}

// NewTemplateMailMessageBody creates a new struct of the TemplateMailMessageBody type and validates its content
func NewTemplateMailMessageBody(files []TemplateMessageBodyFile, recipient string, templateData map[string]string, templateId string) (*TemplateMailMessageBody, error) {
	templateMailMessageBody := TemplateMailMessageBody{
		Files:        files,
		Recipient:    recipient,
		TemplateData: templateData,
		TemplateId:   templateId,
	}
	err := templateMailMessageBody.Validate()
	if err != nil {
		return nil, err
	}
	return &templateMailMessageBody, nil
}

// CustomerUpdateEventBody Event which holds the customer update data
type CustomerUpdateEventBody struct {
	DsgvoRelevant   bool        `json:"dsgvoRelevant"`
	UpdatedCustomer CustomerDTO `json:"updatedCustomer" validate:"required"`
}

// Validate validates the CustomerUpdateEventBody) struct for its requirements
func (thiz CustomerUpdateEventBody) Validate() error {
	return validate.Struct(thiz)
}

// NewCustomerUpdateEventBody creates a new struct of the CustomerUpdateEventBody type and validates its content
func NewCustomerUpdateEventBody(dsgvoRelevant bool, updatedCustomer CustomerDTO) (*CustomerUpdateEventBody, error) {
	customerUpdateEventBody := CustomerUpdateEventBody{
		DsgvoRelevant:   dsgvoRelevant,
		UpdatedCustomer: updatedCustomer,
	}
	err := customerUpdateEventBody.Validate()
	if err != nil {
		return nil, err
	}
	return &customerUpdateEventBody, nil
}

// CustomerDeactivationEventBody Event that holds the information needed to deactivate a customer
type CustomerDeactivationEventBody struct {
	CustomerNumber string `json:"customerNumber" validate:"required"`
}

// Validate validates the CustomerDeactivationEventBody) struct for its requirements
func (thiz CustomerDeactivationEventBody) Validate() error {
	return validate.Struct(thiz)
}

// NewCustomerDeactivationEventBody creates a new struct of the CustomerDeactivationEventBody type and validates its content
func NewCustomerDeactivationEventBody(customerNumber string) (*CustomerDeactivationEventBody, error) {
	customerDeactivationEventBody := CustomerDeactivationEventBody{
		CustomerNumber: customerNumber,
	}
	err := customerDeactivationEventBody.Validate()
	if err != nil {
		return nil, err
	}
	return &customerDeactivationEventBody, nil
}

// CustomerDTO Nested object for CustomerUpdateEventBody type
type CustomerDTO struct {
	CustomerNumber string  `json:"customerNumber" validate:"required"`
	Email          string  `json:"email" validate:"required"`
	FirstName      *string `json:"firstName,omitempty"`
	LastName       *string `json:"lastName,omitempty"`
}

// Validate validates the CustomerDTO struct for its requirements
func (thiz CustomerDTO) Validate() error {
	return validate.Struct(thiz)
}

// NewCustomerDTO creates a new struct of the CustomerDTO type and validates its content
func NewCustomerDTO(customerNumber string, email string, firstName *string, lastName *string) (*CustomerDTO, error) {
	customerDTO := CustomerDTO{
		CustomerNumber: customerNumber,
		Email:          email,
		FirstName:      firstName,
		LastName:       lastName,
	}
	err := customerDTO.Validate()
	if err != nil {
		return nil, err
	}
	return &customerDTO, nil
}

// TemplateMessageBodyFile Nested object for TemplateMailMessageBody type
type TemplateMessageBodyFile struct {
	Content []byte `json:"content" validate:"required"`
	Name    string `json:"name" validate:"required"`
}

// Validate validates the TemplateMessageBodyFile struct for its requirements
func (thiz TemplateMessageBodyFile) Validate() error {
	return validate.Struct(thiz)
}

// NewTemplateMessageBodyFile creates a new struct of the TemplateMessageBodyFile type and validates its content
func NewTemplateMessageBodyFile(content []byte, name string) (*TemplateMessageBodyFile, error) {
	templateMessageBodyFile := TemplateMessageBodyFile{
		Content: content,
		Name:    name,
	}
	err := templateMessageBodyFile.Validate()
	if err != nil {
		return nil, err
	}
	return &templateMessageBodyFile, nil
}

// ConsumerInterface Interface for all events to be consumed by application
type ConsumerInterface interface {
	ConsumeCustomerUpdateEvent(handler goChan.Handler) chan error
	ConsumeCustomerDeactivationEvent(handler goChan.Handler) chan error
	ConsumeCustomerDeleteEvent(handler goChan.Handler) chan error
	ConsumeCustomerStatusEvent(handler goChan.Handler) chan error
}

// ProducerInterface Interface for all events to be produced by application
type ProducerInterface interface {
	ProduceCustomerRegistrationEvent(ctx context.Context, event CustomerRegistrationEventBody) error
	ProduceTemplateMailMessageEvent(ctx context.Context, event TemplateMailMessageBody) error
}

// DefaultConsumer implements ConsumerInterface and consumes events with go kafka mosaic style flavor
type DefaultConsumer struct {
	CustomerDataUpdateTopic   goChan.ChannelInterface
	CustomerDeactivationTopic goChan.ChannelInterface
	CustomerDeletionTopic     goChan.ChannelInterface
	CustomerStatusUpdateTopic goChan.ChannelInterface
}

// DefaultProducer implements ProducerInterface and produces events with go kafka mosaic style flavor
type DefaultProducer struct {
	CustomerDataCreateTopic goChan.ChannelInterface
	TemplateMessageTopic    goChan.ChannelInterface
}

// NewDefaultConsumer wires all needed dependencies to create a DefaultConsumer
func NewDefaultConsumer(manager goChan.ManagerInterface, config kafka.ChannelConfig, mw ...goChan.Middleware) (*DefaultConsumer, error) {
	CustomerDataUpdateTopic, err := manager.CreateChannel(CustomerDataUpdateTopicName, config)
	if err != nil {
		return nil, err
	}
	CustomerDataUpdateTopic.SetReaderMiddleWares(mw...)
	CustomerDeactivationTopic, err := manager.CreateChannel(CustomerDeactivationTopicName, config)
	if err != nil {
		return nil, err
	}
	CustomerDeactivationTopic.SetReaderMiddleWares(mw...)
	CustomerDeletionTopic, err := manager.CreateChannel(CustomerDeletionTopicName, config)
	if err != nil {
		return nil, err
	}
	CustomerDeletionTopic.SetReaderMiddleWares(mw...)
	CustomerStatusUpdateTopic, err := manager.CreateChannel(CustomerStatusUpdateTopicName, config)
	if err != nil {
		return nil, err
	}
	CustomerStatusUpdateTopic.SetReaderMiddleWares(mw...)
	return &DefaultConsumer{
		CustomerDataUpdateTopic:   CustomerDataUpdateTopic,
		CustomerDeactivationTopic: CustomerDeactivationTopic,
		CustomerDeletionTopic:     CustomerDeletionTopic,
		CustomerStatusUpdateTopic: CustomerStatusUpdateTopic,
	}, nil
}

// NewDefaultProducer wires all needed dependencies to create a DefaultProducer
func NewDefaultProducer(manager goChan.ManagerInterface, config kafka.ChannelConfig, mw ...goChan.Middleware) (*DefaultProducer, error) {
	CustomerDataCreateTopic, err := manager.CreateChannel(CustomerDataCreateTopicName, config)
	if err != nil {
		return nil, err
	}
	CustomerDataCreateTopic.SetWriterMiddleWares(mw...)
	TemplateMessageTopic, err := manager.CreateChannel(TemplateMessageTopicName, config)
	if err != nil {
		return nil, err
	}
	TemplateMessageTopic.SetWriterMiddleWares(mw...)
	return &DefaultProducer{
		CustomerDataCreateTopic: CustomerDataCreateTopic,
		TemplateMessageTopic:    TemplateMessageTopic,
	}, nil
}

// ConsumeCustomerUpdateEvent is the go kafka mosaic style flavored implementation of the ConsumerInterface registered on DefaultConsumer
// to consume the CustomerUpdateEventBody Event.
func (d DefaultConsumer) ConsumeCustomerUpdateEvent(handler goChan.Handler) chan error {
	return d.CustomerDataUpdateTopic.Consume(handler)
}

// ConsumeCustomerDeactivationEvent is the go kafka mosaic style flavored implementation of the ConsumerInterface registered on DefaultConsumer
// to consume the CustomerDeactivationEventBody Event.
func (d DefaultConsumer) ConsumeCustomerDeactivationEvent(handler goChan.Handler) chan error {
	return d.CustomerDeactivationTopic.Consume(handler)
}

// ConsumeCustomerDeleteEvent is the go kafka mosaic style flavored implementation of the ConsumerInterface registered on DefaultConsumer
// to consume the CustomerDeleteEventBody Event.
func (d DefaultConsumer) ConsumeCustomerDeleteEvent(handler goChan.Handler) chan error {
	return d.CustomerDeletionTopic.Consume(handler)
}

// ConsumeCustomerStatusEvent is the go kafka mosaic style flavored implementation of the ConsumerInterface registered on DefaultConsumer
// to consume the CustomerStatusEventBody Event.
func (d DefaultConsumer) ConsumeCustomerStatusEvent(handler goChan.Handler) chan error {
	return d.CustomerStatusUpdateTopic.Consume(handler)
}

// ProduceCustomerRegistrationEvent is the go kafka mosaic style flavored implementation of the ConsumerInterface registered on DefaultConsumer
// to produce the CustomerRegistrationEventBody Event.
func (d DefaultProducer) ProduceCustomerRegistrationEvent(ctx context.Context, event CustomerRegistrationEventBody) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "unable to marshal eventdata")
	}
	msg := kafkaGo.Message{
		Value: eventData,
	}
	err = d.CustomerDataCreateTopic.Produce(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

// ProduceTemplateMailMessageEvent is the go kafka mosaic style flavored implementation of the ConsumerInterface registered on DefaultConsumer
// to produce the TemplateMailMessageBody Event.
func (d DefaultProducer) ProduceTemplateMailMessageEvent(ctx context.Context, event TemplateMailMessageBody) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "unable to marshal eventdata")
	}
	msg := kafkaGo.Message{
		Value: eventData,
	}
	err = d.TemplateMessageTopic.Produce(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}
