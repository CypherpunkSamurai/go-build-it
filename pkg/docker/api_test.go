package docker

import (
	"context"
	"io"
	"testing"
	"time"
)

// isDockerAvailable checks if Docker daemon is accessible
// useful in ci environments
func isDockerAvailable(_ *testing.T) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client := GetDockerClient()
	if client == nil {
		return false
	}

	_, err := client.Ping(ctx)
	return err == nil
}

// TestGetDockerClient tests the GetDockerClient function
// tests the docker client for functionality
func TestGetDockerClient(t *testing.T) {
	// test client
	t.Run("Non-nil client", func(t *testing.T) {
		client := GetDockerClient()
		if client == nil {
			t.Error("Expected a non-nil Docker client, but got nil")
		}
	})

	// test singleton behavior
	t.Run("Singleton behavior", func(t *testing.T) {
		client := GetDockerClient()
		secondClient := GetDockerClient()
		if client != secondClient {
			t.Error("Expected the same client instance on second call, but got different instances")
		}
	})

	// Note: Tests After This Require Docker Daemon. Thus We Check Docker Daemon Connectivity
	if !isDockerAvailable(t) {
		t.Skip("Docker daemon is not available - skipping connectivity tests")
	}

	// check ping
	t.Run("Client functionality", func(t *testing.T) {
		client := GetDockerClient()
		ctx := context.Background()
		_, err := client.Ping(ctx)
		if err != nil {
			t.Errorf("Docker client ping failed: %v", err)
		}
	})

	// check server version
	t.Run("Server version", func(t *testing.T) {
		client := GetDockerClient()
		ctx := context.Background()
		version, err := client.ServerVersion(ctx)
		if err != nil {
			t.Errorf("Docker server version failed: %v", err)
		}
		t.Logf("Docker server version: %s (API: %s)", version.Version, version.APIVersion)
	})
}

// TestPullImage tests the PullImage function
// we will use hello-world image
func TestPullImage(t *testing.T) {
	// skip test if docker daemon is not available
	if !isDockerAvailable(t) {
		t.Skip("Docker daemon is not available - skipping pull image test")
	}

	// test pull hello-world image
	t.Run("Pull hello-world image", func(t *testing.T) {
		ctx := context.Background()
		reader, err := PullImage(ctx, "hello-world:latest")
		if err != nil {
			t.Fatalf("Failed to pull image: %v", err)
		}
		defer reader.Close()

		// Verify that we can read from the reader
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			t.Errorf("Failed to read from pull reader: %v", err)
		}
		if n == 0 {
			t.Error("Expected to read some data from pull reader, but got none")
		}
	})

	// test pull logs
	t.Run("Pull logs", func(t *testing.T) {
		ctx := context.Background()
		reader, err := PullImage(ctx, "hello-world:latest")
		if err != nil {
			t.Fatalf("Failed to pull image: %v", err)
		}
		defer reader.Close()

		// Verify that we can read from the reader
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			t.Errorf("Failed to read from pull reader: %v", err)
		}
		if n == 0 {
			t.Error("Expected to read some data from pull reader, but got none")
		}
	})

	// test pull non-existent image
	t.Run("Pull non-existent image", func(t *testing.T) {
		ctx := context.Background()
		_, err := PullImage(ctx, "non-existent-image:latest")
		if err == nil {
			t.Error("Expected an error when pulling non-existent image, but got nil")
		}
	})
}
