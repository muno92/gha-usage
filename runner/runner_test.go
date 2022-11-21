package runner

import (
	"github_actions_usage_calculator/github"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")

	tests := []struct {
		name          string
		repo          string
		targetMonth   string
		expectedUsage github.Usage
	}{
		{
			name:        "2022-01",
			repo:        "muno92/resharper_inspectcode",
			targetMonth: "2022-01",
			expectedUsage: github.Usage{
				Linux:   7369,
				Windows: 14430,
				Mac:     8211,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage, err := Run(tt.repo, tt.targetMonth, token)
			if err != nil {
				panic(err)
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
	}{
		{
			name:               "total workflow runs equal to rate limit",
			remainingRateLimit: 1010,
			totalWorkflowRuns:  1000,
			runnable:           true,
		},
		{
			name:               "total workflow runs less than rate limit",
			remainingRateLimit: 1011,
			totalWorkflowRuns:  1001,
			runnable:           false,
		},
		{
			name:               "total workflow runs greater than rate limit",
			remainingRateLimit: 1013,
			totalWorkflowRuns:  1001,
			runnable:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rateLimits := github.RateLimits{Resources: github.Resource{Core: github.RateLimit{Remaining: tt.remainingRateLimit}}}
			workflowRuns := github.WorkflowRuns{TotalCount: tt.totalWorkflowRuns}

			runnable := IsRunnable(rateLimits, workflowRuns)

			if runnable != tt.runnable {
				t.Errorf("Expected %v, got %v", tt.runnable, runnable)
			}
		})
	}
}
