package main

import (
	"fmt"

	"gopkg.in/bsm/sarama-cluster.v2"
	"github.com/Shopify/sarama"
)

func consumeFromKafka(kafka, topic, group, username, password string) error {

	var config = sarama.NewConfig()
	config.Net.TLS.Enable = true
	config.Net.SASL.Enable = true
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Version = sarama.V0_10_0_1
	config.ClientID = username

	clusterConfig := cluster.NewConfig()
	clusterConfig.Config = *config
	clusterConfig.Consumer.Return.Errors = true

	var errConsumer error
	consumer, errConsumer := cluster.NewConsumer([]string{kafka}, group, []string{topic}, clusterConfig)
	if errConsumer != nil {
		fmt.Printf("Error creating consumer: %s\n", errConsumer)
		return errConsumer
	}

	// Consume errors
	go func() {
		for err := range consumer.Errors() {
			fmt.Printf("Error during consumption: %s\n", err)
		}
	}()

	fmt.Println("Ready to consume messages...")
	// Consume messages
	for msg := range consumer.Messages() {
		fmt.Println(string(msg.Value))
	}
	return nil
}