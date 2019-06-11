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
	b, err := ioutil.ReadFile(c.Args().First())
	if err != nil {
		return fmt.Errorf("Failed to read config file: %v", err)
	}

	var conf containerenv.ConfigFileV1
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return fmt.Errorf("Failed to read config file: %v", err)
	}

	if !conf.Environment.Options.PulseAudio {
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

	id, err := containerenv.CreateContainer(&containerenv.Environment{
		StartOnBoot: false,
		SystemD:     conf.Environment.Options.PulseAudio,
		Image:       fmt.Sprintf("jaredallard/containerenv-%s", conf.Environment.Base),
		Username:    conf.Environment.Username,
		PulseAudio: containerenv.PulseAudioSettings{
			Host: true,
		},
		X11: x11Conf,
	})
	if err != nil {
		return err
	}

	log.Printf("starting container: %s\n", id)

	err = containerenv.StartContainer(id)
	return err
}
