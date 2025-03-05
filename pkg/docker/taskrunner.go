package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cyperpunksamurai/go-build-it/pkg/types"
	"github.com/cyperpunksamurai/go-build-it/pkg/utils"
	dockerTypes "github.com/docker/docker/api/types"
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

	// Create a Dockerfile Tar Ball
	fmt.Println(strings.Join(runTask.DockerFileLines, "\n"))
	dockerfileTar, err := utils.TarFile(ctx, "dockerfile.tar", "Dockerfile", []byte(strings.Join(runTask.DockerFileLines, "\n")))
	if err != nil {
		return err
	}

	// Build
	response, err := client.ImageBuild(ctx, dockerfileTar, dockerTypes.ImageBuildOptions{
		NoCache: true,
		// Tags:    []string{runTask.Image},
		Memory: 1024 * 1024 * 1024, // 1GB
		// SuppressOutput: true,
	})
	if err != nil {
		return err
	}

	// read response
	io.Copy(os.Stdout, response.Body)

	return nil
}
