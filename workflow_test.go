package github_actions_usage_calculator

import (
	"testing"
	"time"
)

func TestFetchWorkflowRuns(t *testing.T) {
	tests := []struct {
		name               string
		repo               string
		token              string
		targetRange        Range
		expectedTotalCount int
	}{
		{
			name:  "Public repo with no token",
			repo:  "muno92/resharper_inspectcode",
			token: "",
			targetRange: Range{
				Start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2022, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
			expectedTotalCount: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workflowRun, err := FetchWorkflowRun(tt.repo, tt.token, tt.targetRange)

			if err != nil {
				panic(err)
			}

			if workflowRun.TotalCount != tt.expectedTotalCount {
				t.Errorf("Expected total count is %v, but got %v.", tt.expectedTotalCount, workflowRun.TotalCount)
			}
		})
	}

}
