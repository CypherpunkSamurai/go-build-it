package types

import (
	"fmt"
	"strings"
)

// RunTaskType represents a task to be run with Docker API
type RunTaskType struct {
	Name       string
	Image      string
	Dockerfile DockerfileType
}

// DockerfileType represents the dockerfile in lines
type DockerfileType struct {
	Lines []string
}

// ToString returns the dockerfile as a string
func (df *DockerfileType) ToString() string {
	return strings.Join(df.Lines, "\n")
}

// Validate ensures the run task is properly configured
func (rt *RunTaskType) Validate() error {
	if rt.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if rt.Dockerfile.ToString() == "" {
		return fmt.Errorf("dockerfile cannot be empty")
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
		Dockerfile: DockerfileType{
			Lines: []string{
				"FROM " + job.Image,
				"RUN mkdir -p /workdir",
				"RUN chmod -R +x /workdir",
				"RUN chown -R $USER /workdir",
				"WORKDIR /workdir",
			},
		},
	}

	// Add steps
	for _, step := range job.Steps {
		switch step.GetKind() {
		case StepKindUses:
			if *step.Uses == "checkout" {
				t.Dockerfile.Lines = append(t.Dockerfile.Lines, "ADD . .")
			}
		case StepKindRun:
			t.Dockerfile.Lines = append(t.Dockerfile.Lines, fmt.Sprintf("RUN %s", *step.Run))
		}
	}

	return t, nil
}
