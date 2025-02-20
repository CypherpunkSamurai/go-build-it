package types

import "fmt"

// RunCommandType represents a steps in shell commands
type RunCommandType struct {
	ShellCommands []string
}

// RunTaskType represents a job in a workflow
type RunTaskType struct {
	Name      string
	ImageName string
	Commands  RunCommandType
}

// Validate ensures the task is properly configured
func (t *RunTaskType) Validate() error {
	return nil
}

// FromWorkflowType creates a RunTaskType from a WorkflowType
func (RunTaskType) FromWorkflowType(wf *WorkflowType, jobId string) (*RunTaskType, error) {
	// create a new RunTaskType
	var t RunTaskType

	// check job name exists
	job, exists := wf.Jobs[jobId]
	if !exists {
		return nil, fmt.Errorf("job %s does not exist", jobId)
	}

	// get job name | or set id
	t.Name = *job.Name

	// get docker image name
	t.ImageName = job.Image

	// iterate over steps
	for _, step := range job.Steps {
		// check step uses
		if step.Uses != nil {
			// check if uses "checkout"
			if *step.Uses == "checkout" {
				// TODO: Implement Real Git Cloning
				t.Commands.ShellCommands = append(t.Commands.ShellCommands, "echo 'cloning git... :P'")
			}
		}

		// check step run
		if step.Run != nil {
			t.Commands.ShellCommands = append(t.Commands.ShellCommands, *step.Run)
		}
	}

	return &t, nil
}
