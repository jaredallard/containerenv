package main

import (
	"fmt"
	"log"

	"strings"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/urfave/cli"
)

func updateCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return fmt.Errorf("Missing environment name")
	}

	envName := c.Args().First()

	e, err := containerenv.GetConfig(envName)
	if err != nil {
		return err
	}

	imageName := e.Image
	if strings.Count(imageName, ":") == 0 {
		imageName = imageName + ":latest"
	} else {
		imageSplit := strings.Split(imageName, ":")
		imageName = imageSplit[0] + ":latest"
	}

	containerenv.PullImage(e.Image)

	log.Printf("recreating container")
	err = containerenv.StopContainer(envName)
	if err != nil {
		log.Printf("WARNING: failed to stop container: %v", err)
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
