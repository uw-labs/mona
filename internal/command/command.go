// Package command contains methods called by the CLI to manage
// a mona project.
package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/apex/log"
	"github.com/davidsbond/mona/internal/files"
	"github.com/davidsbond/mona/pkg/hashdir"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/hashicorp/go-multierror"
)

type (
	changeType int
	rangeFn    func(*files.ModuleFile) error
)

const (
	changeTypeBuild changeType = 0
	changeTypeTest  changeType = 1
	changeTypeLint  changeType = 2
)

func getChangedModules(pj *files.ProjectFile, change changeType) ([]*files.ModuleFile, error) {
	modules, err := files.FindModules(pj.Location)

	if err != nil {
		return nil, err
	}

	lock, err := files.LoadLockFile(pj.Location)

	if err != nil {
		return nil, err
	}

	var out []*files.ModuleFile
	for _, modInfo := range modules {
		lockInfo, ok := lock.Modules[modInfo.Name]

		if !ok {
			// If the module isn't present in the lock file, it's probably
			// new and needs adding to the lock file.
			out = append(out, modInfo)
			continue
		}

		// Generate a new hash for the module directory
		exclude := append(pj.Exclude, modInfo.Exclude...)
		newHash, err := hashdir.Generate(modInfo.Location, exclude...)

		if err != nil {
			return nil, err
		}

		// Check if we need to build/lint/test
		diff := false
		switch change {
		case changeTypeBuild:
			diff = lockInfo.BuildHash != newHash
		case changeTypeTest:
			diff = lockInfo.TestHash != newHash
		case changeTypeLint:
			diff = lockInfo.LintHash != newHash
		}

		if diff {
			out = append(out, modInfo)
		}
	}

	return out, nil
}

func streamOutputs(outputs ...io.ReadCloser) {
	for _, output := range outputs {
		go func(o io.ReadCloser) {
			defer o.Close()

			scanner := bufio.NewScanner(o)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
			}
		}(output)
	}
}

func streamCommand(command, wd string) error {
	parts := strings.Split(command, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = wd

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()

	if err != nil {
		return err
	}

	streamOutputs(stdout, stderr)

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func rangeChangedModules(pj *files.ProjectFile, fn rangeFn, ct changeType) error {
	changed, err := getChangedModules(pj, ct)

	if err != nil || len(changed) == 0 {
		return err
	}

	lock, err := files.LoadLockFile(pj.Location)

	if err != nil {
		return err
	}

	var errs []error
	for _, module := range changed {
		if err := fn(module); err != nil {
			errs = append(errs, fmt.Errorf("module %s: %s", module.Name, err.Error()))
			continue
		}

		exclude := append(pj.Exclude, module.Exclude...)
		newHash, err := hashdir.Generate(module.Location, exclude...)

		if err != nil {
			return err
		}

		lockInfo, modInLock := lock.Modules[module.Name]

		if !modInLock {
			log.Debugf("Detected new module %s at %s, adding to lock file", module.Name, module.Location)

			if err := files.AddModule(lock, pj.Location, module.Name); err != nil {
				return err
			}

			lockInfo = lock.Modules[module.Name]
		}

		switch ct {
		case changeTypeBuild:
			lockInfo.BuildHash = newHash
		case changeTypeTest:
			lockInfo.TestHash = newHash
		case changeTypeLint:
			lockInfo.LintHash = newHash
		}

		if err := files.UpdateLockFile(pj.Location, lock); err != nil {
			return err
		}
	}

	return multierror.Append(nil, errs...).ErrorOrNil()
}

func runInImage(image, cmd, localDir string) error {
	cli, err := client.NewEnvClient()

	if err != nil {
		return err
	}

	ctx := context.Background()

	log.Debugf("Pulling image %s", image)
	if _, err := cli.ImagePull(ctx, image, types.ImagePullOptions{}); err != nil {
		return err
	}

	log.Debugf("Creating container for %s", image)
	body, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:      image,
			WorkingDir: "/module",
			Cmd:        strings.Fields(cmd),
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: localDir,
					Target: "/module",
				},
			},
		}, nil, "")

	if err != nil {
		return err
	}

	log.Debugf("Starting container %s", body.ID)
	if err := cli.ContainerStart(ctx, body.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	log.Debugf("Creating exec configuration for container %s", body.ID)
	resp, err := cli.ContainerExecCreate(ctx, body.ID, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          strings.Fields(cmd),
	})

	if err != nil {
		return err
	}

	log.Debugf("Starting execution of command in container %s", body.ID)
	if err := cli.ContainerExecStart(ctx, resp.ID, types.ExecStartCheck{}); err != nil {
		return err
	}

	log.Debugf("Getting log stream for container %s", body.ID)
	logs, err := cli.ContainerLogs(ctx, body.ID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
	})

	if err != nil {
		return err
	}

	log.Debugf("Waiting for container %s to exit", body.ID)
	streamOutputs(logs)

	if _, err := cli.ContainerWait(ctx, body.ID); err != nil {
		return err
	}

	log.Debugf("Inspecting exec %s in container %s", resp.ID, body.ID)
	inspect, err := cli.ContainerExecInspect(ctx, resp.ID)

	if err != nil {
		return err
	}

	log.Debugf("Exec %s exited with code %d in container %s", resp.ID, inspect.ExitCode, body.ID)

	return nil
}
