package main

import (
	"fmt"
	"os"

	"io/ioutil"

	"os/exec"

	"github.com/jaredallard/containerenv/pkg/containerenv"
	log "github.com/sirupsen/logrus"
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
		log.Fatalln("Disabling systemd is not supported at this time.")
		os.Exit(1)
	}

	if _, _, err := containerenv.GetConfig(conf.Environment.Name); err == nil {
		return fmt.Errorf("environment '%s' already exists. use 'containerenv start %s' instead", conf.Environment.Name, conf.Environment.Name)
	}

	x11Conf := containerenv.X11Settings{}
	if conf.Environment.Options.X11 {

		// determine which mode to enable x11 under
		if _, err := os.Stat("/tmp/.X11-unix"); os.IsNotExist(err) {
			x11Conf.Containerized = true
		} else {
			log.Infof("running xhost ...")
			com := exec.Command("xhost", "local:root")
			com.Env = os.Environ()
			com.Stderr = os.Stderr
			com.Stdin = os.Stdin
			com.Stdout = os.Stdout

			if err := com.Run(); err != nil {
				log.Warnf("Failed to run xhost, running X applications may fail")
			}

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

	log.Infoln("Pulling docker image ...")
	err = containerenv.PullImage(image)
	if err != nil {
		log.Warnf("Failed to pull docker image.")
	}

	id, err := containerenv.CreateContainer(env)
	if err != nil {
		return err
	}

	log.Infof("Starting container: %s\n", id)

	err = containerenv.StartContainer(id)
	return err
}
