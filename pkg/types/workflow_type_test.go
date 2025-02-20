package types

import (
	"testing"

	"github.com/cyperpunksamurai/go-build-it/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// TestUnmarshalWorkflowType tests the UnmarshalWorkflowType function with various test cases.
func TestUnmarshalWorkflowType(t *testing.T) {
	// Here we define a few test yaml strings, and their expected results
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
		check   func(*testing.T, *WorkflowType)
	}{
		{
			name: "valid workflow with all fields",
			yaml: `
version: "1.0"
name: Build and Test
jobs:
  test:
    name: Run Tests
    image: golang:latest
    requires:
      - build
    steps:
      - use: checkout
      - name: run tests
        run: go test ./...
  build:
    image: golang:latest
    steps:
      - use: checkout
      - run: go build -o main .
`,
			wantErr: false,
			check: func(t *testing.T, w *WorkflowType) {
				assert.Equal(t, "1.0", *w.Version)
				assert.Equal(t, "Build and Test", *w.Name)
				assert.Len(t, w.Jobs, 2)

				// Check test job
				testJob := w.Jobs["test"]
				assert.Equal(t, "Run Tests", *testJob.Name)
				assert.Equal(t, "golang:latest", testJob.Image)
				assert.Equal(t, []string{"build"}, testJob.Requires)
				assert.Len(t, testJob.Steps, 2)
				assert.Equal(t, "checkout", *testJob.Steps[0].Uses)
				assert.Equal(t, "run tests", *testJob.Steps[1].Name)
				assert.Equal(t, "go test ./...", *testJob.Steps[1].Run)

				// Check build job
				buildJob := w.Jobs["build"]
				assert.Equal(t, "build", *buildJob.Name) // default name from key
				assert.Equal(t, "golang:latest", buildJob.Image)
				assert.Len(t, buildJob.Steps, 2)
			},
		},
		{
			name: "invalid step with both use and run",
			yaml: `
jobs:
  test:
    image: golang:latest
    steps:
      - use: checkout
        run: echo "invalid"
`,
			wantErr: true,
		},
		{
			name: "invalid step with neither use nor run",
			yaml: `
jobs:
  test:
    image: golang:latest
    steps:
      - name: invalid step
`,
			wantErr: true,
		},
		{
			name: "minimal valid workflow",
			yaml: `
jobs:
  test:
    image: golang:latest
    steps:
      - run: echo "hello"
`,
			wantErr: false,
			check: func(t *testing.T, w *WorkflowType) {
				assert.Nil(t, w.Version)
				assert.Nil(t, w.Name)
				assert.Len(t, w.Jobs, 1)
				assert.Equal(t, "test", *w.Jobs["test"].Name)
				assert.Equal(t, StepKindRun, w.Jobs["test"].Steps[0].GetKind())
			},
		},
	}

	// Test Each Case with Expected Results
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load Test Yaml
			workflow, err := UnmarshalWorkflowType(tt.yaml)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			// Assert There are no errors
			assert.NoError(t, err)
			if tt.check != nil {
				tt.check(t, workflow)
			}
		})
	}
}

// TestStepType_GetName tests the GetName function of the StepType struct
func TestStepType_GetName(t *testing.T) {
	tests := []struct {
		name     string
		step     StepType
		expected string
	}{
		{
			name: "with explicit name",
			step: StepType{
				Name: utils.StrPtr("explicit name"),
				Run:  utils.StrPtr("echo hello"),
			},
			expected: "explicit name",
		},
		{
			name: "with uses only",
			step: StepType{
				Uses: utils.StrPtr("checkout"),
			},
			expected: "uses: checkout",
		},
		{
			name: "with run only",
			step: StepType{
				Run: utils.StrPtr("echo hello"),
			},
			expected: "run: echo hello",
		},
		{
			name:     "empty step",
			step:     StepType{},
			expected: "",
		},
	}

	// Test Each Case with Expected Results
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assert the Cases Match
			assert.Equal(t, tt.expected, tt.step.GetName())
		})
	}
}
