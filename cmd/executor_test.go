package cmd

import (
	"errors"
	"ghausage/github"
	"log"
	"os"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")

	tests := []struct {
		name                string
		repo                string
		startDate           string
		endDate             string
		expectedErrorExists bool
		expectedUsage       github.Usage
	}{
		{
			name:                "workflow run count is less than 100",
			repo:                "muno92/resharper_inspectcode",
			startDate:           "2022-01-01",
			endDate:             "2022-01-31",
			expectedErrorExists: false,
			expectedUsage: github.Usage{
				Linux:   7369,
				Windows: 14430,
				Mac:     8211,
			},
		},
		{
			name:                "workflow run count is between 100 and 1000",
			repo:                "muno92/life_log",
			startDate:           "2022-03-01",
			endDate:             "2022-03-05",
			expectedErrorExists: false,
			expectedUsage: github.Usage{
				Linux:   3262,
				Windows: 0,
				Mac:     0,
			},
		},
		{
			name:                "workflow run count is more than 1000",
			repo:                "muno92/life_log",
			startDate:           "2022-03-01",
			endDate:             "2022-03-15",
			expectedErrorExists: true,
			expectedUsage:       github.Usage{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage, err := Run(tt.repo, tt.startDate, tt.endDate, token, log.Default())

			errorExists := err != nil
			if tt.expectedErrorExists != errorExists {
				t.Errorf("expected error exists is %v, got %v\n%v", tt.expectedErrorExists, errorExists, err)
			}

			if usage != tt.expectedUsage {
				t.Errorf("Expected usage is %v, got %v", tt.expectedUsage, usage)
			}
		})
	}
}

func TestIsRunnable(t *testing.T) {
	tests := []struct {
		name               string
		remainingRateLimit int
		totalWorkflowRuns  int
		runnable           bool
		expectedError      error
	}{
		{
			name:               "total workflow runs is equal to rate limit",
			remainingRateLimit: 910,
			totalWorkflowRuns:  900,
			runnable:           true,
			expectedError:      nil,
		},
		{
			name:               "total workflow runs is equal to rate limit and just 1000",
			remainingRateLimit: 1010,
			totalWorkflowRuns:  1000,
			runnable:           true,
			expectedError:      nil,
		},
		{
			name:               "total workflow runs is equal to rate limit, but over 1000",
			remainingRateLimit: 1012,
			totalWorkflowRuns:  1001,
			runnable:           false,
			expectedError:      errors.New("count of workflow run (1001) must be less than or equal to 1000 (because GitHub API will not return over 1000 records even with pagination)"),
		},
		{
			name:               "total workflow runs is less than rate limit",
			remainingRateLimit: 910,
			totalWorkflowRuns:  901,
			runnable:           false,
			expectedError:      errors.New("please try again at 2022-01-23T01:23:45Z, because rate limit remaining (910) is less than expected fetch count (911)"),
		},
		{
			name:               "total workflow runs is greater than rate limit",
			remainingRateLimit: 913,
			totalWorkflowRuns:  901,
			runnable:           true,
			expectedError:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resetTime, _ := time.Parse(time.RFC3339, "2022-01-23T01:23:45+00:00")
			rateLimits := github.RateLimits{
				Resources: github.Resource{
					Core: github.RateLimit{
						Remaining: tt.remainingRateLimit,
						Reset:     resetTime.UnixMilli(),
					},
				},
			}
			workflowRuns := github.WorkflowRuns{TotalCount: tt.totalWorkflowRuns}

			runnable, err := IsRunnable(rateLimits, workflowRuns)

			if runnable != tt.runnable {
				t.Errorf("Expected %v, got %v", tt.runnable, runnable)
			}

			if tt.runnable {
				return
			}

			if err.Error() != tt.expectedError.Error() {
				t.Errorf("\nExpected error is \n\t%v,\ngot \n\t%v", tt.expectedError, err)
			}
		})
	}
}
