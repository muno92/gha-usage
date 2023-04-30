package github

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestFetchWorkflowRuns(t *testing.T) {
	client := NewClient(os.Getenv("GITHUB_TOKEN"), log.Default())

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
				t.Error(err)
			}

			if workflowRuns.TotalCount != tt.expectedTotalCount {
				t.Errorf("Expected total count is %v, but got %v.", tt.expectedTotalCount, workflowRuns.TotalCount)
			}
			if len(workflowRuns.WorkflowRuns) != tt.expectedWorkflowRunCount {
				t.Errorf("Expected workflow run count is %v, but got %v.", tt.expectedWorkflowRunCount, len(workflowRuns.WorkflowRuns))
			}
		})
	}

}

func TestWorkflowRunUsage(t *testing.T) {
	client := NewClient(os.Getenv("GITHUB_TOKEN"), log.Default())

	tests := []struct {
		name          string
		jobsUrl       string
		expectedUsage Usage
	}{
		{
			name:    "single job",
			jobsUrl: "https://api.github.com/repos/muno92/resharper_inspectcode/actions/runs/3370110776/jobs",
			expectedUsage: Usage{
				Linux:      15,
				Windows:    0,
				Mac:        0,
				SelfHosted: 0,
			},
		},
		{
			name:    "multi job",
			jobsUrl: "https://api.github.com/repos/muno92/resharper_inspectcode/actions/runs/3370110771/jobs",
			expectedUsage: Usage{
				Linux:      287,
				Windows:    622,
				Mac:        284,
				SelfHosted: 0,
			},
		},
		{
			name:    "job per workflow is over 100",
			jobsUrl: "https://api.github.com/repos/phpstan/phpstan-src/actions/runs/3281861062/jobs",
			expectedUsage: Usage{
				Linux:      7084,
				Windows:    61,
				Mac:        0,
				SelfHosted: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WorkflowRun{JobsUrl: tt.jobsUrl}

			usage, err := w.Usage(client)
			if err != nil {
				t.Error(err)
			}

			if usage != tt.expectedUsage {
				t.Errorf("Expected usage is %v, got %v", tt.expectedUsage, usage)
			}
		})
	}
}
