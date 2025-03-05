package types

import (
	"fmt"
)

// RunTaskType represents a task to be run with Docker API
type RunTaskType struct {
	Name            string   // the name of this task
	Image           string   // the image to use
	DockerFileLines []string // our docker commands
	GitUrl          string
}

// Validate ensures the run task is properly configured
func (rt *RunTaskType) Validate() error {
	if rt.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	// validate that no command is null
	for _, dockerFileLine := range rt.DockerFileLines {
		if dockerFileLine == "" {
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
		Name:  *job.Name,
		Image: job.Image,
		DockerFileLines: []string{
			"FROM " + job.Image,
			"RUN mkdir -p /workdir",
			"WORKDIR /workdir",
		},
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
				t.DockerFileLines = append(t.DockerFileLines,
					fmt.Sprintf("RUN cd /workdir && git clone %s .", t.GitUrl),
				)
			}
		case StepKindRun:
			t.DockerFileLines = append(t.DockerFileLines, fmt.Sprintf("RUN %s", *step.Run))
		default:
			return nil, fmt.Errorf("unknown step kind %s", step.GetKind())
		}
	}

	return t, nil
}
