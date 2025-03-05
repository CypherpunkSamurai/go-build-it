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
				GitUrl: utils.StrPtr("https://github.com/cyperpunksamurai/go-build-it.git"),
				Jobs: map[string]JobType{
					"test": {
						Name:  utils.StrPtr("test"),
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
				Name:  "test",
				Image: "golang:latest",
				DockerFileLines: []string{
					// remember to add the initial commands in test
					"RUN mkdir -p /workdir",
					"WORKDIR /workdir",
					"RUN cd /workdir && git clone https://github.com/cyperpunksamurai/go-build-it.git .",
					"RUN go test ./...",
				},
				GitUrl: "https://github.com/cyperpunksamurai/go-build-it.git",
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
		//
		t.Run(tt.name, func(t *testing.T) {
			// Load Test Yaml
			got, err := FromWorkflowType(tt.workflow, tt.jobName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			// Assert There are no errors
			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
