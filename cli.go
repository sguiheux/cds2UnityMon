package main

import (
	"os"
	"path/filepath"

	"gopkg.in/urfave/cli.v1"
)

func initCli() *cli.App {
	return &cli.App{
		Name:         filepath.Base(os.Args[0]),
		HelpName:     filepath.Base(os.Args[0]),
		Usage:        "CDS Event Listener & Cache for unity monitoring",
		UsageText:    "",
		Author:       "Steven GUIHEUX <steven.guiheux@gmail.com>",
		Version:      "0.1",
		BashComplete: cli.DefaultAppComplete,
		Writer:       os.Stdout,
		Commands: []cli.Command{
			{
				Name:      "help",
				Aliases:   []string{"h"},
				Usage:     "Shows a list of commands or help for one command",
				ArgsUsage: "[command]",
				Action:    helpAction,
			},
			{
				Name:      "listen",
				Aliases:   []string{"l"},
				Usage:     "Listen a Kafka topic and wait",
				ArgsUsage: "<kafka> <topic> <group> <username> <password>",
				Action: listenAction,
			},
		},
	}
}
