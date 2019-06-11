package containerenv

import (
	"fmt"
	"log"

	"os"

	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// CreateContainer creates a docker container and returns the id of the container
func CreateContainer(e *Environment) (string, error) {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		return "", fmt.Errorf("ERROR: couldn't create docker client\n%+v", err)
	}

	config := &container.Config{
		Image: e.Image,
		Env:   []string{fmt.Sprintf("USERNAME_CONFIG=%s", e.Username)},
	}

	hostconfig := &container.HostConfig{
		Tmpfs:  make(map[string]string),
		CapAdd: make([]string, 0),
	}

	if e.SystemD {
		log.Println("enabling systemd")
		systemdTmpfsMounts := []string{
			"/run",
			"/run/lock",
			"/tmp",
			"/sys/fs/cgroup/systemd",
			"/var/lib/journal",
		}

		for _, mount := range systemdTmpfsMounts {
			hostconfig.Tmpfs[mount] = ""
		}

		hostconfig.CapAdd = append(hostconfig.CapAdd, "SYS_ADMIN")
		hostconfig.Binds = append(hostconfig.Binds, "/sys/fs/cgroup:/sys/fs/cgroup:ro")

		// systemd expects this signal for graceful shutdown
		config.StopSignal = "SIGRTMIN+3"
		config.Env = append(config.Env, "SYSTEMD_CONFIG=enabled")
	}

	if e.PulseAudio.Host {
		log.Println("using pulseaudio from host")
		uid := os.Getuid()

		hostconfig.Binds = append(hostconfig.Binds, fmt.Sprintf("/run/user/%d/pulse:/run/user/1000/pulse", uid))

		config.Env = append(config.Env, "PULSEAUDIO_CONFIG=HOST")
	} else if e.PulseAudio.Containerized {
		log.Println("setting up pulseaudio in the container")

		hostconfig.Devices = append(hostconfig.Devices, container.DeviceMapping{
			PathOnHost:        "/dev/snd",
			PathInContainer:   "/dev/snd",
			CgroupPermissions: "rwm",
		})

		config.Env = append(config.Env, "PULSEAUDIO_CONFIG=CONTAINER")
	}

	if e.X11.Host {
		log.Println("using X11 from the host")
		hostconfig.Binds = append(hostconfig.Binds, "/tmp/.X11-unix:/tmp/.X11-unix:ro")
		hostconfig.Devices = append(hostconfig.Devices, container.DeviceMapping{
			PathOnHost:        "/dev/dri",
			PathInContainer:   "/dev/dri",
			CgroupPermissions: "rwm",
		})

		config.Env = append(config.Env, []string{"X11_CONFIG=HOST", "DISPLAY=:0", "PULSE_SERVER"}...)
	} else if e.X11.Containerized {
		log.Println("running xorg in the container")
		hostconfig.Privileged = true

		// TODO: collapse w/ the above definition
		hostconfig.Devices = append(hostconfig.Devices, container.DeviceMapping{
			PathOnHost:        "/dev/dri",
			PathInContainer:   "/dev/dri",
			CgroupPermissions: "rwm",
		})

		config.Env = append(config.Env, "X11_CONFIG=CONTAINER")
	}

	resp, err := docker.ContainerCreate(ctx, config, hostconfig, &network.NetworkingConfig{}, "environment")
	if client.IsErrImageNotFound(err) {

	} else if err != nil {
		return "", err
	}

	return resp.ID, err
}

// StartContainer starts a container
func StartContainer(id string) error {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("ERROR: couldn't create docker client\n%+v", err)
	}

	return docker.ContainerStart(ctx, id, types.ContainerStartOptions{})
}

// StopContainer stops a container
func StopContainer(id string) error {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("ERROR: couldn't create docker client\n %+v", err)
	}

	dur := time.Minute * 5
	return docker.ContainerStop(ctx, id, &dur)
}
