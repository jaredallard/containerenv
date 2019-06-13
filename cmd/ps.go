package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"time"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/urfave/cli"
)

func psCommand(c *cli.Context) error {
	envs, err := containerenv.ListContainers()
	if err != nil {
		return err
	}

	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 5, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tCONTAINER ID\tCREATED AT\t")
	for _, c := range *envs {
		idRunes := []rune(c.ID)
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%v\t", c.Labels["jaredallard.containerenv/environment-name"], string(idRunes[0:12]), time.Unix(c.Created, 0)))
	}
	w.Flush()

	return nil
}
