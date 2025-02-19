package main

import (
	"context"

	"github.com/cyperpunksamurai/go-build-it/pkg/docker"
	"github.com/cyperpunksamurai/go-build-it/pkg/utils"
)

func main() {
	// create a context for main
	ctx := context.Background()
	defer ctx.Done()

	// test pulling alpine:latest container
	reader, err := docker.PullImage(ctx, "alpine:latest")
	if err != nil {
		panic(err)
	}
	// close the reader
	defer reader.Close()

	// read from reader and print
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			break
		}
		utils.Logger().Println(string(buf[:n]))
	}
}
