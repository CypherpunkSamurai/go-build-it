package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cyperpunksamurai/go-build-it/pkg/docker"
	"github.com/cyperpunksamurai/go-build-it/pkg/types"
	"github.com/cyperpunksamurai/go-build-it/pkg/utils"
)

func main() {
	// Init Env
	utils.Initenv()

	// create a context for main
	ctx := context.Background()
	defer ctx.Done()

	// Read Workflow File
	workflow, err := types.UnmarshalWorkflowTypeFromFile("example.yml")
	if err != nil {
		panic(err)
	}

	// FromWorkflowType returns a pointer
	runTask, err := types.FromWorkflowType(workflow, "test")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded workflow task: %+v\n", runTask.Name)
	fmt.Printf("- image: %s\n", runTask.Image)
	for _, line := range runTask.DockerFileLines {
		fmt.Printf("- command: %s\n", line)
	}

	// Run Task
	ctx, cancel := context.WithTimeout(ctx, time.Duration(5*time.Minute))
	defer cancel()
	err = docker.RunTaskWithDocker(ctx, runTask)
	if err != nil {
		panic(err)
	}

	// // Listen for Tasks
	// rabbitmq := queue.NewRabbitMQClient("go-build-worker")
	// // connect
	// err = rabbitmq.Connect(utils.Getenv("AMQP_URL"))
	// if err != nil {
	// 	panic(err)
	// }

	// // Listen on ci_tasks queue
	// err = rabbitmq.ConsumeQueue("ci_tasks", func(msg []byte) error {
	// 	fmt.Printf("[main] Received a message: %s\n", string(msg))
	// 	return nil
	// })
	// if err != nil {
	// 	panic(err)
	// }
}
