package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cyperpunksamurai/go-build-it/pkg/types"
	"github.com/docker/docker/api/types/container"
)

// RunTaskWithDocker - Runs a RunTask with Docker API
func RunTaskWithDocker(ctx context.Context, runTask *types.RunTaskType) error {
	client := GetDockerClient()
	if client == nil {
		return fmt.Errorf("docker client is not initialized")
	}

	// Pull Image (just in case)
	r, err := PullImage(ctx, runTask.Image)
	if err != nil {
		return err
	}

	// write to discard
	// io.Copy(io.Discard, r)

	// stream r to stdout
	io.Copy(os.Stdout, r)

	// create a container with a shell
	resp, err := GetDockerClient().ContainerCreate(ctx, &container.Config{
		Image:        runTask.Image,
		Cmd:          []string{"/bin/sh"},
		Tty:          true,
		OpenStdin:    true,
		StdinOnce:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}

	// start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	// attach io
	waiter, err := cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		return err
	}
	defer waiter.Close()

	go io.Copy(os.Stdout, waiter.Reader)
	go io.Copy(os.Stderr, waiter.Reader)

	// run shell commands
	for _, line := range runTask.DockerFileLines {
		linebytes := []byte(line + "\n")
		fmt.Println("Running command: ", string(linebytes))
		_, err = waiter.Conn.Write(linebytes)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}

	// wait for container to exit
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	return nil
}
