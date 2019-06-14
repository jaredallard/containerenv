package main

import (
	"fmt"
	"log"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/urfave/cli"
)

func commitCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return fmt.Errorf("Missing environment name")
	}

	envName := c.Args().First()

	var commitImage string

	e, err := containerenv.GetConfig(envName)
	if err != nil {
		return err
	}

	if c.String("image") != "" {
		commitImage = c.String("image")
	} else {
		if e.CommitImage == "" {
			return fmt.Errorf("Image did not specify a reference format. Please supply one with --image")
		}
	}

	log.Printf("creating a new version of environment '%s' and publishing to '%s'", envName, commitImage)

	imageName, err := containerenv.Commit(envName, commitImage)
	if err != nil {
		return err
	}

	if !c.Bool("no-push") {
		log.Printf("pushing image ...")
		err = containerenv.Push(imageName)
		if err != nil {
			return err
		}
	}

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
