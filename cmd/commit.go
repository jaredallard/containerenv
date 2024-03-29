package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"strings"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/urfave/cli"
)

func commitCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return fmt.Errorf("Missing environment name")
	}

	envName := c.Args().First()

	var commitImage string

	e, _, err := containerenv.GetConfig(envName)
	if err != nil {
		return err
	}

	if c.String("image") != "" {
		commitImage = c.String("image")
	} else {
		if e.CommitImage == "" {
			return fmt.Errorf("Environment did not specify a reference format. Please supply one with --image")
		}

		commitImage = strings.Split(e.Image, ":")[0]
	}

	log.Infof("creating a new version of environment '%s' and publishing to '%s'", envName, commitImage)

	imageName, err := containerenv.Commit(envName, commitImage)
	if err != nil {
		return err
	}

	if !c.Bool("no-push") {
		log.Infof("pushing image ...")
		err = containerenv.Push(imageName)
		if err != nil {
			return err
		}
	}

	log.Infof("recreating container")
	err = containerenv.StopContainer(envName)
	if err != nil {
		log.Warnf("failed to stop container: %v", err)
	}

	err = containerenv.RemoveContainer(envName)
	if err != nil {
		return err
	}

	e.Image = imageName
	id, err := containerenv.CreateContainer(e)
	if err != nil {
		return err
	}

	err = containerenv.StartContainer(id)
	return err
}
