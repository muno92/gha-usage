package github

import (
	"os"
	"testing"
	"time"
)

func TestFetchWorkflowRuns(t *testing.T) {
	client := Client{Token: os.Getenv("GITHUB_TOKEN")}

	tests := []struct {
		name                     string
		repo                     string
		targetRange              Range
		perPage                  int
		page                     int
		expectedTotalCount       int
		expectedWorkflowRunCount int
	}{
		{
			name: "Public repo, first page",
			repo: "muno92/resharper_inspectcode",
			targetRange: Range{
				Start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2022, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
			perPage:                  30,
			page:                     1,
			expectedTotalCount:       60,
			expectedWorkflowRunCount: 30,
		},
		{
			name: "Public repo, last page",
			repo: "muno92/resharper_inspectcode",
			targetRange: Range{
				Start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2022, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
			perPage:                  25,
			page:                     3,
			expectedTotalCount:       60,
			expectedWorkflowRunCount: 10,
		},
		{
			name: "Private repo",
			repo: "muno92/dummy_private_repo",
			targetRange: Range{
				Start: time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2022, time.January, 31, 0, 0, 0, 0, time.UTC),
			},
			perPage:                  30,
			page:                     1,
			expectedTotalCount:       0,
			expectedWorkflowRunCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workflowRuns, err := FetchWorkflowRuns(tt.repo, client, tt.targetRange, tt.perPage, tt.page)

			if err != nil {
				panic(err)
			}

			if workflowRuns.TotalCount != tt.expectedTotalCount {
				t.Errorf("Expected total count is %v, but got %v.", tt.expectedTotalCount, workflowRuns.TotalCount)
			}
			if len(workflowRuns.WorkflowRuns) != tt.expectedWorkflowRunCount {
				t.Errorf("Expected job count is %v, but got %v.", tt.expectedWorkflowRunCount, len(workflowRuns.WorkflowRuns))
			}
		})
	}

}
