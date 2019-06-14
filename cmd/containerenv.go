package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jaredallard/containerenv/pkg/version"
	"github.com/urfave/cli"
)

// App instance
var App cli.App

func main() {
	App := cli.NewApp()
	App.Name = "containerenv"
	App.Usage = "reproducible and shareable operating system environments"
	App.Version = version.GetVersion()

	App.Commands = []cli.Command{
		cli.Command{
			Name:      "create",
			Usage:     "Create and start an environment",
			UsageText: "create <container.yaml>",
			Action:    createCommand,
		},
		cli.Command{
			Name:      "exec",
			Usage:     "Exec into an environment",
			UsageText: "exec <environment-name>",
			Action:    execCommand,
		},
		cli.Command{
			Name:      "commit",
			Usage:     "Create a new version of an environment and then replace it.",
			UsageText: "commit <environment-name>",
			Action:    commitCommand,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "image",
					Usage: "--image jaredallard/myenv",
				},
				cli.BoolFlag{
					Name:  "no-push",
					Usage: "--no-push disables pushing the image (results in images only existing locally)",
				},
			},
		},
		cli.Command{
			Name:      "delete",
			Usage:     "Delete an environment",
			UsageText: "delete <environment-name>",
			Action:    deleteCommand,
		},
		cli.Command{
			Name:      "ps",
			Usage:     "List running environments",
			UsageText: "ps",
			Action:    psCommand,
		},
		cli.Command{
			Name:      "version",
			Usage:     "Return the version of the application",
			UsageText: "version",
			Action: func(c *cli.Context) error {
				fmt.Printf("%s version %s\n", App.Name, App.Version)
				return nil
			},
		},
	}

	err := App.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
