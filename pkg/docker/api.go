package docker

import (
	"context"
	"io"
	"sync"

	"github.com/cyperpunksamurai/go-build-it/pkg/utils"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// Docker Client Instance
var (
	cli  *client.Client
	once sync.Once
)

// GetDockerClient returns the Docker client instance
func GetDockerClient() *client.Client {
	// thread-safe implementation (classic go :P)
	once.Do(func() {
		var err error
		cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			utils.Logger().Printf("Error creating docker client: %v", err)
			cli = nil
		}
	})
	return cli
}

// PullImage pulls a Docker image from a registry
// Parameters:
//   - ctx: The context for the operation
//   - imageName: The name of the image to pull
//   - options: Optional parameters for pulling the image
//
// Returns:
//   - io.ReadCloser: A reader for the pull operation output
//   - error: An error if the operation fails
func PullImage(ctx context.Context, imageName string) (io.ReadCloser, error) {
	return GetDockerClient().ImagePull(ctx, imageName, image.PullOptions{})
}
