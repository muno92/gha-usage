package github_actions_usage_calculator

import (
	"github_actions_usage_calculator/config"
	"github_actions_usage_calculator/github"
	"math"
)

func IsRunnable(limits github.RateLimits, runs github.WorkflowRuns) bool {
	workflowFetchCount := int(math.Ceil(float64(runs.TotalCount) / float64(config.PerPage)))

	return workflowFetchCount+runs.TotalCount <= limits.Resources.Core.Remaining
}
