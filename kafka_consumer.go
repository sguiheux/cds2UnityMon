package main

import (
	"fmt"
	"os"
	"github.com/Shopify/sarama"
	"os/signal"
)

func consumeFromKafka(kafka, topic, group, key string) error {

	//Create a new client
	var config = sarama.NewConfig()
	// Set key as the client id for authentication
	config.ClientID = key
	client, err := sarama.NewClient([]string{kafka}, config)
	if err != nil {
		return err
	}

	// Create an offsetManager
	offsetManager, err := sarama.NewOffsetManagerFromClient(group, client)
	if err != nil {
		return err
	}

	// Create a client
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return err
	}

	fmt.Printf("Listening Kafka %s on topic %s...\n", kafka, topic)

	// Create the message chan, that will receive the queue
	messagesChan := make(chan []byte)
	// Create the error chan, that will receive the queue
	errorsChan := make(chan error)

	// read the number of partition for the given topic
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}

	// Create a consumer for each partition
	if len(partitions) > 1 {
		return fmt.Errorf("Multiple partition not supported")
	}
	p := partitions[0]
	partitionOffsetManager, err := offsetManager.ManagePartition(topic, p)
	if err != nil {
		return err
	}
	defer partitionOffsetManager.AsyncClose()

	// Start a consumer at next offset
	offset, _ := partitionOffsetManager.NextOffset()
	partitionConsumer, err := consumer.ConsumePartition(topic, p, offset)
	if err != nil {
		return err
	}
	defer partitionConsumer.AsyncClose()

	// Asynchronously handle message
	go consumptionHandler(partitionConsumer, topic, partitionOffsetManager, messagesChan, errorsChan)

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)


	for {
		select {
		case msg := <-messagesChan:
			// Process msg
			fmt.Println(msg)
		case err := <-errorsChan:
			fmt.Printf("%s\n", err)
			return err
		case <-signals:
			return nil
		}
	}
}

// ConsumptionHandler pipes the handled messages and push them to a chan
func consumptionHandler(pc sarama.PartitionConsumer, topic string, po sarama.PartitionOffsetManager, messagesChan chan<- []byte, errorsChan chan<- error) {
	for {
		select {
		case msg := <-pc.Messages():
			messagesChan <- msg.Value
			po.MarkOffset(msg.Offset+1, topic)
		case err := <-pc.Errors():
			fmt.Println(err)
			errorsChan <- err
		case offsetErr := <-po.Errors():
			fmt.Println(offsetErr)
			errorsChan <- offsetErr
		}
	}
}