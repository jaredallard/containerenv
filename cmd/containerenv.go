package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jaredallard/containerenv/pkg/version"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "containerenv"
	app.Usage = "reproducible and shareable operating system environments"
	app.Version = version.GetVersion()

	app.Commands = []cli.Command{
		cli.Command{
			Name:      "create",
			Usage:     "Create and start an environment",
			UsageText: "create <container.yaml>",
			Action:    createCommand,
		},
		cli.Command{
			Name:      "version",
			Usage:     "Return the version of the application",
			UsageText: "version",
			Action: func(c *cli.Context) error {
				fmt.Printf("%s version %s\n", app.Name, app.Version)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
