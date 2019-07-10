package main

import (
	"fmt"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/urfave/cli"
)

func execCommand(c *cli.Context) error {
	envName := c.Args().First()
	if envName == "" {
		return fmt.Errorf("Missing environment name")
	}

	conf, id, err := containerenv.GetConfig(envName)
	if err != nil {
		return err
	}

	err = containerenv.Exec(id, conf)
	return err
}
