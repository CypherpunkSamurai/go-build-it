package docker

import (
	"context"
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
