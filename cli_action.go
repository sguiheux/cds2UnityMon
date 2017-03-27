package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/urfave/cli.v1"
)

func helpAction(c *cli.Context) error {
	fmt.Println("Coucou")
	return nil
}

func listenAction(c *cli.Context) error {
	args := []string{c.Args().First()}
	args = append(args, c.Args().Tail()...)
	if len(args) != 5 {
		cli.ShowCommandHelp(c, "listen")
		return cli.NewExitError("Invalid usage", 10)
	}

	kafkaHost := os.Getenv("kafka_host")
	topic := os.Getenv("kafka_topic")
	group := os.Getenv("kafka_group")
	username := os.Getenv("kafka_username")
	password := os.Getenv("kafka_password")

	if kafkaHost == "" || topic == "" || group == "" || username == "" || password == "" {
		return cli.NewExitError("Missing env variable", 11)
	}

	go consumeFromKafka(kafkaHost, topic, group, username, password)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/queue", func(c *gin.Context) {
		c.JSON(200, mapPb)
	})
	r.GET("/queue/:key", func(c *gin.Context) {
		key := c.Param("key")
		if _, ok := mapPb[key]; !ok {
			c.JSON(410, "")
			return
		}
		c.JSON(200, gin.H{"status": mapPb[key]})
	})
	r.Run()
	return nil
}
