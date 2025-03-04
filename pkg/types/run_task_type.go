package types

import (
	"fmt"
)

// RunTaskType represents a task to be run with Docker API
type RunTaskType struct {
	Name                string              // the name of this task
	Image               string              // the image to use
	DockerShellCommands []DockerCommandType // our docker commands
	GitUrl              string
}

// DockerCommandType represents a single command to be run in the dockerimage
type DockerCommandType struct {
	cmd string            // command to run in the container
	env map[string]string // env key value pairs
	cwd string            // current working directory (nill to use default directory)
}

// ToString returns the docker command as a string
func (dc *DockerCommandType) ToString() string {
	return dc.cmd
}

// Validate ensures the run task is properly configured
func (rt *RunTaskType) Validate() error {
	if rt.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	// validate that no command is null
	for _, dockercommand := range rt.DockerShellCommands {
		if dockercommand.ToString() == "" {
			return fmt.Errorf("docker command cannot be empty")
		}
	}
	return nil
}

// FromWorkflowType creates a RunTaskType from a WorkflowType
func FromWorkflowType(wf *WorkflowType, jobID string) (*RunTaskType, error) {
	job, exists := wf.Jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job %s does not exist", jobID)
	}

	// Create a new RunTaskType
	t := &RunTaskType{
		Name:                *job.Name,
		Image:               job.Image,
		DockerShellCommands: []DockerCommandType{},
	}

	// check if git url is set
	if wf.GitUrl != nil {
		t.GitUrl = *wf.GitUrl
	}

	// Add steps
	for _, step := range job.Steps {
		switch step.GetKind() {
		case StepKindUses:
			if *step.Uses == "checkout" && t.GitUrl != "" {
				t.DockerShellCommands = append(t.DockerShellCommands, DockerCommandType{
					cmd: fmt.Sprintf("git clone %s .", t.GitUrl),
					cwd: "",
				})
			}
		case StepKindRun:
			t.DockerShellCommands = append(t.DockerShellCommands, DockerCommandType{
				cmd: *step.Run,
				cwd: "",
			})
		}
	}

	return t, nil
}
