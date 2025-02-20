package types

import (
	"testing"

	"github.com/cyperpunksamurai/go-build-it/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// TestFromWorkflowType - Tests FromWorkflowType function so it returns valid RunTaskType
func TestFromWorkflowType(t *testing.T) {
	// define a list of test cases (we defined only 2)
	tests := []struct {
		name     string
		workflow *WorkflowType
		jobName  string
		want     *RunTaskType
		wantErr  bool
	}{
		{
			name: "valid job",
			workflow: &WorkflowType{
				Jobs: map[string]JobType{
					"test": {
						Image: "golang:latest",
						Steps: []StepType{
							{Uses: utils.StrPtr("checkout")},
							{Run: utils.StrPtr("go test ./...")},
						},
					},
				},
			},
			jobName: "test",
			want: &RunTaskType{
				ImageName: "golang:latest",
				Commands: RunCommandType{
					ShellCommands: []string{
						"echo 'cloning git... :P'",
						"go test ./...",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "non-existent job",
			workflow: &WorkflowType{
				Jobs: map[string]JobType{},
			},
			jobName: "non-existent",
			want:    nil,
			wantErr: true,
		},
	}

	// test each test case with expectations
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// check each case
			got, err := RunTaskType{}.FromWorkflowType(tt.workflow, tt.jobName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}

}
