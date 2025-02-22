package main

import (
	"context"
	"fmt"

	"github.com/cyperpunksamurai/go-build-it/pkg/types"
	"github.com/cyperpunksamurai/go-build-it/pkg/utils"
)

func main() {
	// Init Env
	utils.Initenv()

	// Print Env
	fmt.Println("AMQP URL:", utils.Getenv("AMQP_URL"))

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
	for _, cmd := range runTask.Dockerfile.Lines {
		fmt.Printf("- command: %s\n", cmd)
	}
}
