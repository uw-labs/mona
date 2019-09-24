package environment

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
// the docker image under the /app directory.
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
	log.Debugf("Executing command '%s' in image %s", command, d.image)
	log.Debugf("Pulling image %s", d.image)

	if _, err := d.cli.ImagePull(ctx, d.image, types.ImagePullOptions{}); err != nil {
		return err
	}

	log.Debugf("Creating container for %s", d.image)
	body, err := d.cli.ContainerCreate(ctx,
		&container.Config{
			Image:      d.image,
			WorkingDir: "/app",
			Cmd:        strings.Split(command, " "),
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: d.volumeDir,
					Target: "/app",
				},
			},
		}, nil, "")

	if err != nil {
		return err
	}

	for _, warning := range body.Warnings {
		log.Warn(warning)
	}

	log.Debugf("Starting container %s", body.ID)
	if err := d.cli.ContainerStart(ctx, body.ID, types.ContainerStartOptions{}); err != nil {
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

	go func(o io.ReadCloser) {
		defer o.Close()

		scanner := bufio.NewScanner(o)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}
	}(logs)

	if _, err := d.cli.ContainerWait(ctx, body.ID); err != nil {
		return err
	}

	return nil
}
