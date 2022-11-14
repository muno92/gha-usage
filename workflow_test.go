package github_actions_usage_calculator

import (
	"os"
	"testing"
	"time"
)

func TestFetchWorkflowRuns(t *testing.T) {
	client := Client{Token: os.Getenv("GITHUB_TOKEN")}

	tests := []struct {
		name               string
		repo               string
		targetRange        Range
		expectedTotalCount int
	}{
		{
			name: "Public repo",
			repo: "muno92/resharper_inspectcode",
			targetRange: Range{
				Start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2022, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
			expectedTotalCount: 60,
		},
		{
			name: "Private repo",
			repo: "muno92/dummy_private_repo",
			targetRange: Range{
				Start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2022, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
			expectedTotalCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workflowRun, err := FetchWorkflowRun(tt.repo, client, tt.targetRange)

			if err != nil {
				panic(err)
			}

			if workflowRun.TotalCount != tt.expectedTotalCount {
				t.Errorf("Expected total count is %v, but got %v.", tt.expectedTotalCount, workflowRun.TotalCount)
			}
		})
	}

}
