package main

import (
	"fmt"

	"io/ioutil"

	"github.com/jaredallard/containerenv/pkg/cliutils"
	"github.com/jaredallard/containerenv/pkg/containerenv"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func initCommand(c *cli.Context) error {
	images := []string{
		"archlinux",
	}

	e := containerenv.ConfigFileV1{
		Version: 1,
	}

	fmt.Println("Supported Distros:")
	fmt.Println()
	for _, i := range images {
		fmt.Printf(" * %s\n", i)
	}
	fmt.Println()
	fmt.Printf("Distro to base off of: ")
	base, err := cliutils.GetUserInput()
	if err != nil {
		return err
	}

	found := false
	for _, i := range images {
		if i == base {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("distro '%s' isn't currently supported", base)
	}

	defaultName := namesgenerator.GetRandomName(0)
	fmt.Printf("Name of this environment (default will be %s): ", defaultName)
	envName, err := cliutils.GetUserInput()
	if err != nil {
		return err
	}

	if envName == "" {
		// fmt.Println(defaultName)
		envName = defaultName
	}

	e.Environment.Name = envName

	fmt.Printf("Username to be created in the image: ")
	username, err := cliutils.GetUserInput()
	if err != nil {
		return err
	}
	if username == "" {
		return fmt.Errorf("username is empty")
	}

	e.Environment.Username = username

	defaultImage := fmt.Sprintf("%s/%s", e.Environment.Username, e.Environment.Name)
	fmt.Printf("Docker image we should create (i.e jaredallard/env, defaults to '%s'): ", defaultImage)
	imageName, err := cliutils.GetUserInput()
	if err != nil {
		return err
	}
	if imageName == "" {
		// fmt.Println(imageName)
		imageName = defaultImage
	}

	e.Environment.Image = imageName
	e.Environment.CommitOptions.Image = imageName

	fmt.Printf("Enable pulseaudio? [Y/n]: ")
	v, err := cliutils.GetYesOrNoInput()
	if err != nil {
		return err
	}

	e.Environment.Options.PulseAudio = v

	fmt.Printf("Enable X11? [Y/n]: ")
	v, err = cliutils.GetYesOrNoInput()
	if err != nil {
		return err
	}

	e.Environment.Options.X11 = v
	e.Environment.Options.SystemD = true

	fmt.Println(" --> Building image ... ")
	fqbn := fmt.Sprintf("jaredallard/containerenv-%s", base)
	err = containerenv.PullImage(fqbn)
	if err != nil {
		return err
	}

	fmt.Println(" --> Tagging image")
	err = containerenv.Tag(fqbn, imageName)
	if err != nil {
		return err
	}

	fmt.Println(" --> Generating Config")
	b, err := yaml.Marshal(e)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./env.yaml", b, 0755)
	if err != nil {
		return err
	}

	fmt.Println(" --> Environment saved at 'env.yaml' !!")

	return nil
}
