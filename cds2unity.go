package main

import (
	"os"
)

func main() {
	app := initCli()
	app.Run(os.Args)
}