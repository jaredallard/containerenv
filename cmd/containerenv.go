package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"strings"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/jaredallard/containerenv/pkg/version"
	"github.com/urfave/cli"
)

// App instance
var App cli.App

func listEnvs(cli *cli.Context) {
	cs, err := containerenv.ListContainers()
	if err != nil {
		return
	}

	for _, c := range *cs {
		if c.Labels["jaredallard.containerenv/environment-name"] == "" {
			continue
		}

		fmt.Println(c.Labels["jaredallard.containerenv/environment-name"])
	}

	return
}

func main() {
	if strings.ToUpper(os.Getenv("LOG_LEVEL")) == "TRACE" {
		log.SetLevel(log.TraceLevel)
	}

	if strings.ToUpper(os.Getenv("LOG_OUTPUT")) == "JSON" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	App := cli.NewApp()
	App.Name = "containerenv"
	App.Usage = "reproducible and shareable operating system environments"
	App.Version = version.GetVersion()

	App.EnableBashCompletion = true

	App.Commands = []cli.Command{
		cli.Command{
			Name:      "create",
			Usage:     "Create and start an environment",
			UsageText: "create <container.yaml>",
			Action:    createCommand,
		},
		cli.Command{
			Name:         "exec",
			Usage:        "Exec into an environment",
			UsageText:    "exec <environment-name>",
			Action:       execCommand,
			BashComplete: listEnvs,
		},
		cli.Command{
			Name:         "update",
			Usage:        "Update a running environment",
			UsageText:    "update <environment-name>",
			Action:       updateCommand,
			BashComplete: listEnvs,
		},
		cli.Command{
			Name:         "start",
			Usage:        "Start an already created environment",
			UsageText:    "start <environment-name>",
			Action:       startCommand,
			BashComplete: listEnvs,
		},
		cli.Command{
			Name:        "completion",
			Usage:       "Generate shell completion",
			Description: "Add source <(containerenv completion YOUR_SHELL_HERE) to your shell rc (e.g ~/.bashrc)",
			Action:      generateCompletion,
		},
		cli.Command{
			Name:         "commit",
			Usage:        "Create a new version of an environment and then replace it.",
			UsageText:    "commit <environment-name>",
			Action:       commitCommand,
			BashComplete: listEnvs,
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
			Name:      "init",
			Usage:     "Create an environment configuration",
			UsageText: "init",
			Action:    initCommand,
		},
		cli.Command{
			Name:         "delete",
			Usage:        "Delete an environment",
			UsageText:    "delete <environment-name>",
			Action:       deleteCommand,
			BashComplete: listEnvs,
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
		log.Fatalln(err.Error())
		os.Exit(1)
	}
}
