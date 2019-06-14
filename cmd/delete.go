package main

import (
	"fmt"
	"log"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/urfave/cli"
)

func deleteCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return fmt.Errorf("Missing environment name")
	}

	envName := c.Args().First()
	_, err := containerenv.GetConfig(envName)
	if err != nil {
		return err
	}

	log.Printf("destroying environment")
	err = containerenv.StopContainer(envName)
	if err != nil {
		log.Printf("WARNING: failed to stop container: %v", err)
	}

	err = containerenv.RemoveContainer(envName)
	return err
}
