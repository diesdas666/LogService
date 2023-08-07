package main

import (
	"example_consumer/internal/kafka/configkafka"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
)

func main() {

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
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt)
	run := true
	for run {
		select {
		case sig := <-signals:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		case ev := <-consumer.Events():
			switch e := ev.(type) {
			case *kafka.Message:
				// Process the message
				fmt.Printf("Received message on topic %s: %s\n", *e.TopicPartition.Topic, string(e.Value))
			}
		}
	}
	// Close the consumer
	consumer.Close()
}
