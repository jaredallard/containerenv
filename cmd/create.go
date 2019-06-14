package main

import (
	"fmt"
	"log"
	"os"

	"io/ioutil"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func createCommand(c *cli.Context) error {
	if c.Args().First() == "" {
		return fmt.Errorf("Missing environment config path")
	}

	b, err := ioutil.ReadFile(c.Args().First())
	if err != nil {
		return fmt.Errorf("Failed to read config file: %v", err)
	}

	var conf containerenv.ConfigFileV1
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return fmt.Errorf("Failed to read config file: %v", err)
	}

	if !conf.Environment.Options.SystemD {
		log.Fatalf("Disabling systemd is not supported at this time.")
		os.Exit(1)
	}

	x11Conf := containerenv.X11Settings{}
	if conf.Environment.Options.X11 {

		// determine which mode to enable x11 under
		if _, err := os.Stat("/tmp/.X11-unix"); os.IsNotExist(err) {
			x11Conf.Containerized = true
		} else {
			x11Conf.Host = true
		}
	}

	var image string
	if conf.Environment.Base != "" {
		image = fmt.Sprintf("jaredallard/containerenv-%s", conf.Environment.Base)
	} else if conf.Environment.Image != "" {
		image = conf.Environment.Image
	} else {
		return fmt.Errorf("Missing a Base or Image declaration")
	}

	env := &containerenv.Environment{
		Name:     conf.Environment.Name,
		SystemD:  conf.Environment.Options.PulseAudio,
		Image:    image,
		Username: conf.Environment.Username,
		PulseAudio: containerenv.PulseAudioSettings{
			Host: true,
		},
		X11: x11Conf,
	}

	if conf.Environment.CommitOptions.Image != "" {
		env.CommitImage = conf.Environment.CommitOptions.Image
	}

	if conf.Environment.Binds != nil {
		env.Binds = conf.Environment.Binds
	}

	id, err := containerenv.CreateContainer(env)
	if err != nil {
		return err
	}

	log.Printf("starting container: %s\n", id)

	err = containerenv.StartContainer(id)
	return err
}
