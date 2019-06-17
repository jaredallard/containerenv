package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/urfave/cli"
)

func startCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return fmt.Errorf("Missing environment name")
	}

	envName := c.Args().First()

	cs, err := containerenv.ListContainers()

	var cont *types.Container
	for _, c := range *cs {
		if c.Labels["jaredallard.containerenv/environment-name"] == envName {
			cont = &c
		}

		continue
	}

	if cont == nil {
		return fmt.Errorf("Failed to find a container")
	}

	log.Infof("Running status: %s", cont.State)

	if cont.State == "running" {
		return fmt.Errorf("container is already running. Run 'containerenv exec %s' to exec into it", envName)
	}

	err = containerenv.StartContainer(cont.ID)
	return err
}
