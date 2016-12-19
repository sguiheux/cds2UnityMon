package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
	"github.com/gin-gonic/gin"
)

func helpAction(c *cli.Context) error {
	fmt.Println("Coucou")
	return nil
}

func listenAction(c *cli.Context) error {
	args := []string{c.Args().First()}
	args = append(args, c.Args().Tail()...)
	if len(args) != 4 {
		cli.ShowCommandHelp(c, "listen")
		return cli.NewExitError("Invalid usage", 10)
	}

	//Get arguments from commandline
	kafka := c.Args().Get(0)
	topic := c.Args().Get(1)
	group := c.Args().Get(2)
	key := c.Args().Get(3)

	if err := consumeFromKafka(kafka, topic, group, key); err != nil {
		return cli.NewExitError(err.Error(), 13)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
	return nil
}