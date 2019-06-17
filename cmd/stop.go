package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/urfave/cli"
)

func stopCommand(c *cli.Context) error {
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

	if cont.State != "running" {
		return fmt.Errorf("container for environment '%s' is not running", envName)
	}

	log.Infof("stopping environment '%s' ...", envName)
	err = containerenv.StopContainer(cont.ID)
	return err
}
