package environment

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type (
	// The Docker type implements Environment and uses the local environment's
	// docker set up to execute a command.
	Docker struct {
		image     string
		volumeDir string
		cli       *client.Client
	}
)

// NewDockerEnvironment creates a new Environment implementation that will use the local
// docker client to execute commands. All files within the 'volumeDir' will be mounted to
// the docker image under the /module directory.
func NewDockerEnvironment(image, volumeDir string) (Environment, error) {
	cli, err := client.NewEnvClient()

	if err != nil {
		return nil, err
	}

	return &Docker{
		image:     image,
		volumeDir: volumeDir,
		cli:       cli,
	}, nil
}

// Execute runs the provided command within the configured docker image
func (d *Docker) Execute(ctx context.Context, command string) error {
	log.Debugf("Pulling image %s", d.image)

	if _, err := d.cli.ImagePull(ctx, d.image, types.ImagePullOptions{}); err != nil {
		return err
	}

	log.Debugf("Creating container for %s", d.image)
	body, err := d.cli.ContainerCreate(ctx,
		&container.Config{
			Image:      d.image,
			WorkingDir: "/module",
			Cmd:        strings.Fields(command),
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: d.volumeDir,
					Target: "/module",
				},
			},
		}, nil, "")

	if err != nil {
		return err
	}

	log.Debugf("Starting container %s", body.ID)
	if err := d.cli.ContainerStart(ctx, body.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	log.Debugf("Creating exec configuration for container %s", body.ID)
	resp, err := d.cli.ContainerExecCreate(ctx, body.ID, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          strings.Fields(command),
	})

	if err != nil {
		return err
	}

	log.Debugf("Starting execution of command in container %s", body.ID)
	if err := d.cli.ContainerExecStart(ctx, resp.ID, types.ExecStartCheck{}); err != nil {
		return err
	}

	log.Debugf("Getting log stream for container %s", body.ID)
	logs, err := d.cli.ContainerLogs(ctx, body.ID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
	})

	if err != nil {
		return err
	}

	log.Debugf("Waiting for container %s to exit", body.ID)

	go func() {
		defer logs.Close()
		scanner := bufio.NewScanner(logs)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	if _, err := d.cli.ContainerWait(ctx, body.ID); err != nil {
		return err
	}

	log.Debugf("Inspecting exec %s in container %s", resp.ID, body.ID)
	inspect, err := d.cli.ContainerExecInspect(ctx, resp.ID)

	if err != nil {
		return err
	}

	log.Debugf("Exec %s exited with code %d in container %s", resp.ID, inspect.ExitCode, body.ID)

	return nil
}
