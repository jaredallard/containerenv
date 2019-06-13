package main

import (
	"fmt"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/urfave/cli"
)

func execCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return fmt.Errorf("Missing environment name")
	}

	err := containerenv.Exec(c.Args().First())
	return err
}
