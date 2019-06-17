package main

import (
	"fmt"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/urfave/cli"
)

func logsCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return fmt.Errorf("Missing environment name")
	}

	envName := c.Args().First()

	_, id, err := containerenv.GetConfig(envName)
	if err != nil {
		return err
	}

	// TODO: actually parse logs
	err = containerenv.Logs(id, []string{
		"-f",
	})
	return err
}
