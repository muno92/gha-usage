package github_actions_usage_calculator

import (
	"github_actions_usage_calculator/github"
	"math"
)

const (
	PerPage int = 100
)

func IsRunnable(limits github.RateLimits, runs github.WorkflowRuns) bool {
	c := len(runs.WorkflowRuns)
	workflowFetchCount := int(math.Ceil(float64(c) / float64(PerPage)))

	return workflowFetchCount+c <= limits.Resources.Core.Remaining
}
