package types

import (
	"fmt"

	"github.com/goccy/go-yaml"
)

// StepType represents the type of step (uses or run)
type StepKind string

const (
	StepKindUses StepKind = "uses"
	StepKindRun  StepKind = "run"
)

// WorkflowType represents a workflow
type WorkflowType struct {
	Version *string            `yaml:"version,omitempty"` // version of the workflow (optional)
	Name    *string            `yaml:"name,omitempty"`    // name of the workflow (optional)
	Jobs    map[string]JobType `yaml:"jobs"`              // map of job name to job definition
}

// JobType represents a job in a workflow
type JobType struct {
	Name     *string    `yaml:"name,omitempty"`     // name of the job (optional)
	Requires []string   `yaml:"requires,omitempty"` // list of jobs that this job depends on
	Image    string     `yaml:"image"`              // docker image to run the job
	Steps    []StepType `yaml:"steps"`              // list of steps to execute
}

// StepType represents a step in a job
type StepType struct {
	Uses *string `yaml:"use,omitempty"`  // use the step defined in the workflow
	Name *string `yaml:"name,omitempty"` // name of the step (optional)
	Run  *string `yaml:"run,omitempty"`  // run command to execute
}

// GetName returns the name of the step, if set, otherwise returns the step identifier.
// We use this instead of setting .Name to preserve the original yaml formatting
func (s *StepType) GetName() string {
	// if the step has a name
	if s.Name != nil {
		return *s.Name
	}
	// if the step is a 'uses' type
	if s.Uses != nil {
		return fmt.Sprintf("uses: %s", *s.Uses)
	}
	// if the step is a 'run' type
	if s.Run != nil {
		return fmt.Sprintf("run: %s", *s.Run)
	}
	return ""
}

// Validate ensures the step is properly configured
func (s *StepType) Validate() error {
	if s.Uses != nil && s.Run != nil {
		return fmt.Errorf("step cannot have both 'use' and 'run' fields")
	}
	if s.Uses == nil && s.Run == nil {
		return fmt.Errorf("step must have either 'use' or 'run' field")
	}
	return nil
}

// GetKind returns the type of step (uses or run)
func (s *StepType) GetKind() StepKind {
	if s.Uses != nil {
		return StepKindUses
	}
	return StepKindRun
}

// UnmarshalWorkflowType unmarshalls a workflow from yaml string
func UnmarshalWorkflowType(yml string) (*WorkflowType, error) {
	v := &WorkflowType{}
	if err := yaml.Unmarshal([]byte(yml), v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	// Validate all steps
	for jobName, job := range v.Jobs {
		if job.Name == nil {
			// If job name is not provided, use the key as the name
			name := jobName
			job.Name = &name
			v.Jobs[jobName] = job
		}
		for _, step := range job.Steps {
			if err := step.Validate(); err != nil {
				return nil, fmt.Errorf("invalid step in job %s: %w", jobName, err)
			}
		}
	}

	return v, nil
}
