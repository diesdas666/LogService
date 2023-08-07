package main

import (
	"example_consumer/internal/asyncApi/configkafka"
	"example_consumer/internal/service"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
	"time"
)

func main() {

	logService := service.LogService{
		//Database konfig or address
	}

	// Set up Kafka consumer configuration
	config := &kafka.ConfigMap{
		configkafka.Host:   configkafka.Kafka.Host, // Update with your broker's address
		configkafka.Group:  configkafka.Kafka.Group,
		configkafka.Offset: configkafka.Kafka.Offset,
	}
	// Create Kafka consumer instance
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("Error creating consumer: %v\n", err)
		return
	}
	// Subscribe to the topic
	topic := configkafka.Topic
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Error subscribing to topic: %v\n", err)
		return
	}

	// Handle OS signals to gracefully close the consumer
	sigchan := make(chan os.Signal, 1)

	signal.Notify(sigchan, os.Interrupt)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		//case ev := <-consumer.Events():
		//	switch e := ev.(type) {
		//	case *kafka.Message:
		//		// Process the message
		//		fmt.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, string(e.Value))
		//	}
		default:
			ev, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			err = logService.UpdateStatusOfLoadingStation(ev)
			if err != nil {
				fmt.Printf("Error procesingevent: %v\n", err)
			}
			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}
	// Close the consumer
	consumer.Close()
}
