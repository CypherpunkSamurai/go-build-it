package main

import (
	"context"
	"fmt"

	"github.com/cyperpunksamurai/go-build-it/pkg/types"
)

func main() {
	// create a context for main
	ctx := context.Background()
	defer ctx.Done()

	// Read Workflow File
	workflow, err := types.UnmarshalWorkflowTypeFromFile("example.yml")
	if err != nil {
		panic(err)
	}

	// Create an empty RunTaskType to call the method
	var rt types.RunTaskType

	// FromWorkflowType returns a pointer
	runTask, err := rt.FromWorkflowType(workflow, "test")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded workflow task: %+v\n", runTask.Name)
	fmt.Printf("- image: %s\n", runTask.ImageName)
	for _, cmd := range runTask.Commands.ShellCommands {
		fmt.Printf("- command: %s\n", cmd)
	}
}
