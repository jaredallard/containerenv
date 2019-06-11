package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "containerenv"
	app.Usage = "reproducible and shareable operating system environments"

	app.Commands = []cli.Command{
		cli.Command{
			Name:      "create",
			Usage:     "Create and start an environment",
			UsageText: "create <container.yaml>",
			Action:    createCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
