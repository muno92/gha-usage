package github_actions_usage_calculator

import (
	"github_actions_usage_calculator/config"
	"github_actions_usage_calculator/github"
	"math"
)

func IsRunnable(limits github.RateLimits, runs github.WorkflowRuns) bool {
	c := len(runs.WorkflowRuns)
	workflowFetchCount := int(math.Ceil(float64(c) / float64(config.PerPage)))

	return workflowFetchCount+c <= limits.Resources.Core.Remaining
}
