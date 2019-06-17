package main

import (
	"fmt"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func deleteCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return fmt.Errorf("Missing environment name")
	}

	envName := c.Args().First()
	_, id, err := containerenv.GetConfig(envName)
	if err != nil {
		return err
	}

	log.Infof("destroying environment")
	err = containerenv.StopContainer(id)
	if err != nil {
		log.Warnf("failed to stop container: %v", err)
	}

	err = containerenv.RemoveContainer(id)
	return err
}
