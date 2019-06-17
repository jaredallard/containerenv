package containerenv

import (
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"os"

	"time"

	"encoding/base64"
	"encoding/json"

	"strings"

	"os/user"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"k8s.io/kubernetes/pkg/kubectl/util/term"
)

// TODO: don't invoke docker

// CreateContainer creates a docker container and returns the id of the container
func CreateContainer(e *Environment) (string, error) {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		return "", fmt.Errorf("couldn't create docker client\n%+v", err)
	}

	b, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal config: %v", err)
	}

	config := &container.Config{
		Image: e.Image,
		Labels: map[string]string{
			"jaredallard.containerenv/environment":      "true",
			"jaredallard.containerenv/environment-name": e.Name,
			"jaredallard.containerenv/config":           base64.StdEncoding.EncodeToString(b),
		},
		Env: []string{fmt.Sprintf("USERNAME_CONFIG=%s", e.Username)},
	}

	hostconfig := &container.HostConfig{
		Tmpfs:       make(map[string]string),
		CapAdd:      make([]string, 0),
		NetworkMode: "host",
		Binds: []string{
			"/var/run/docker.sock:/var/run/docker.sock",
		},
	}

	g, err := user.LookupGroup("docker")
	if err == nil {
		config.Env = append(config.Env, fmt.Sprintf("DOCKER_GID=%s", g.Gid))
	}

	if e.Binds != nil {
		hostconfig.Binds = append(hostconfig.Binds, e.Binds...)
	}

	if e.SystemD {
		log.Tracef("SYSTEMD_CONFIG=enabled")
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
		log.Tracef("PULSEAUDIO_CONFIG=HOST")
		uid := os.Getuid()

		hostconfig.Binds = append(hostconfig.Binds, fmt.Sprintf("/run/user/%d/pulse:/run/user/1000/pulse", uid))

		config.Env = append(config.Env, "PULSEAUDIO_CONFIG=HOST")
	} else if e.PulseAudio.Containerized {
		log.Tracef("PULSEAUDIO_CONFIG=CONTAINER")

		hostconfig.Devices = append(hostconfig.Devices, container.DeviceMapping{
			PathOnHost:        "/dev/snd",
			PathInContainer:   "/dev/snd",
			CgroupPermissions: "rwm",
		})

		config.Env = append(config.Env, "PULSEAUDIO_CONFIG=CONTAINER")
	}

	if e.X11.Host {
		log.Tracef("X11_CONFIG=HOST")
		hostconfig.Binds = append(hostconfig.Binds, "/tmp/.X11-unix:/tmp/.X11-unix:ro")
		hostconfig.Devices = append(hostconfig.Devices, container.DeviceMapping{
			PathOnHost:        "/dev/dri",
			PathInContainer:   "/dev/dri",
			CgroupPermissions: "rwm",
		})

		config.Env = append(config.Env, []string{"X11_CONFIG=HOST", "DISPLAY=:0", "PULSE_SERVER"}...)
	} else if e.X11.Containerized {
		log.Tracef("X11_CONFIG=CONTAINER")
		hostconfig.Privileged = true

		// TODO: collapse w/ the above definition
		hostconfig.Devices = append(hostconfig.Devices, container.DeviceMapping{
			PathOnHost:        "/dev/dri",
			PathInContainer:   "/dev/dri",
			CgroupPermissions: "rwm",
		})

		config.Env = append(config.Env, "X11_CONFIG=CONTAINER")
	}

	bytesConfig, err := json.MarshalIndent(config, "", "  ")
	bytesHostConf, err := json.MarshalIndent(hostconfig, "", "  ")

	log.Tracef("config = '%s'", bytesConfig)
	log.Tracef("hostconfig = '%s'", bytesHostConf)

	resp, err := docker.ContainerCreate(ctx, config, hostconfig, &network.NetworkingConfig{}, e.Name)
	if client.IsErrImageNotFound(err) {
		PullImage(config.Image)
		resp, err = docker.ContainerCreate(ctx, config, hostconfig, &network.NetworkingConfig{}, e.Name)
		return resp.ID, err
	} else if err != nil {
		return "", err
	}

	return resp.ID, err
}

// PullImage pulls a docker image
// invokes docker to not reimplement this
func PullImage(image string) error {
	com := exec.Command("/usr/bin/docker", []string{
		"pull",
		image,
	}...)
	com.Env = os.Environ()
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout

	err := (term.TTY{In: os.Stdin, TryDev: true}).Safe(com.Run)
	return err
}

// StartContainer starts a container
func StartContainer(id string) error {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("couldn't create docker client\n%+v", err)
	}

	return docker.ContainerStart(ctx, id, types.ContainerStartOptions{})
}

// StopContainer stops a container
func StopContainer(id string) error {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("couldn't create docker client\n %+v", err)
	}

	dur := time.Second * 5
	return docker.ContainerStop(ctx, id, &dur)
}

// RemoveContainer removes a container
func RemoveContainer(id string) error {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("couldn't create docker client\n %+v", err)
	}

	return docker.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
}

// ListContainers returns a list of environment containers running
func ListContainers() (*[]types.Container, error) {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("couldn't create docker client\n%+v", err)
	}

	f := filters.NewArgs()
	f.Add("label", "jaredallard.containerenv/environment=true")
	f.Add("status", "running")
	f.Add("status", "exited")
	cnts, err := docker.ContainerList(ctx, types.ContainerListOptions{
		Filters: f,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search for environments: %s", err)
	}

	log.Tracef("ListContainers(): got len(%d) container(s)", len(cnts))

	names := make([]types.Container, len(cnts))
	for i, cnt := range cnts {
		names[i] = cnt
	}

	return &names, nil
}

// GetConfig returns the config of an environment
func GetConfig(name string) (*Environment, string, error) {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		return nil, "", fmt.Errorf("couldn't create docker client\n%+v", err)
	}

	f := filters.NewArgs()
	f.Add("label", fmt.Sprintf("jaredallard.containerenv/environment-name=%s", name))
	f.Add("status", "running")
	f.Add("status", "exited")
	cnts, err := docker.ContainerList(ctx, types.ContainerListOptions{
		Filters: f,
	})
	if err != nil {
		return nil, "", err
	}
	if len(cnts) != 1 {
		return nil, "", fmt.Errorf("failed to find container for environment (len %d)", len(cnts))
	}

	config := cnts[0].Labels["jaredallard.containerenv/config"]
	var env Environment

	decodedConfig, err := base64.StdEncoding.DecodeString(config)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base64 config: %v", err)
	}

	err = json.Unmarshal(decodedConfig, &env)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshall to json: %v", err)
	}

	return &env, cnts[0].ID, nil
}

// Exec opens a shell into an environment
// This just wraps docker exec due to not wanting to reimplement that.
func Exec(name string) error {
	com := exec.Command("/usr/bin/docker", []string{
		"exec",
		"-it",
		"--user",
		"1000",
		name,
		"bash",
		"--login",
	}...)
	com.Env = os.Environ()
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout

	err := (term.TTY{In: os.Stdin, TryDev: true}).Safe(com.Run)
	return err
}

// Commit creates a new version of an environment and pushes
// it to a docker registry
func Commit(name, image string) (string, error) {
	ctx := context.Background()

	docker, err := client.NewEnvClient()
	if err != nil {
		return "", fmt.Errorf("couldn't create docker client\n%+v", err)
	}

	imageName := fmt.Sprintf("%s:%s", image, time.Now().UTC().Format("2006-01-02T15-04-05"))

	if strings.Count(imageName, "/") == 1 {
		imageName = fmt.Sprintf("docker.io/%s", imageName)
	}

	_, err = docker.ContainerCommit(ctx, name, types.ContainerCommitOptions{
		Pause:     false,
		Reference: imageName,
	})

	return imageName, err
}

// Tag tags an image
func Tag(src, dst string) error {
	com := exec.Command("/usr/bin/docker", []string{
		"tag",
		src,
		dst,
	}...)
	com.Env = os.Environ()
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout

	err := (term.TTY{In: os.Stdin, TryDev: true}).Safe(com.Run)
	return err
}

// Push pushes a docker image and updates the latest ref
// Uses docker cli to avoid dealing with auth
func Push(image string) error {
	com := exec.Command("/usr/bin/docker", []string{
		"push",
		image,
	}...)
	com.Env = os.Environ()
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout

	err := (term.TTY{In: os.Stdin, TryDev: true}).Safe(com.Run)
	if err != nil {
		return err
	}

	subs := strings.Split(image, ":'")
	latest := fmt.Sprintf("%s:latest", subs[0])

	Tag(image, latest)

	log.Printf("updating latest tag")

	com = exec.Command("/usr/bin/docker", []string{
		"push",
		latest,
	}...)
	com.Env = os.Environ()
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout

	err = (term.TTY{In: os.Stdin, TryDev: true}).Safe(com.Run)
	return err
}

// Logs retrieves container logs
func Logs(id string, args []string) error {
	// joins logs w/ args then id
	dockerArgs := append(append([]string{"logs"}, args...), id)
	log.Tracef("exec /usr/bin/docker %v", dockerArgs)

	com := exec.Command("/usr/bin/docker", dockerArgs...)
	com.Env = os.Environ()
	com.Stderr = os.Stderr
	com.Stdin = os.Stdin
	com.Stdout = os.Stdout

	err := (term.TTY{In: os.Stdin, TryDev: true}).Safe(com.Run)
	return err
}
