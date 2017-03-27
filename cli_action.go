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
	if len(args) != 5 {
		cli.ShowCommandHelp(c, "listen")
		return cli.NewExitError("Invalid usage", 10)
	}

	//Get arguments from commandline
	kafka := c.Args().Get(0)
	topic := c.Args().Get(1)
	group := c.Args().Get(2)
	username := c.Args().Get(3)
	password := c.Args().Get(4)

	go consumeFromKafka(kafka, topic, group, username, password)

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
		c.JSON(200, gin.H{ "status": mapPb[key]})
	})
	r.Run()
	return nil
}