package github_actions_usage_calculator

import (
	"github_actions_usage_calculator/github"
	"testing"
)

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
