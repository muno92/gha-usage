package github_actions_usage_calculator

import (
	"fmt"
	"github_actions_usage_calculator/config"
	"github_actions_usage_calculator/github"
	"math"
)

func Run(repo string, targetMonth string, token string) (github.Usage, error) {
	targetRange, err := github.DecideRange(targetMonth)
	if err != nil {
		return github.Usage{}, err
	}

	client := github.Client{Token: token}

	workflowRuns, err := github.FetchWorkflowRuns(repo, client, targetRange, config.PerPage, 1)
	if err != nil {
		return github.Usage{}, err
	}

	rateLimit, err := github.FetchRateLimit(client)
	if err != nil {
		return github.Usage{}, err
	}

	if !IsRunnable(rateLimit, workflowRuns) {
		return github.Usage{}, fmt.Errorf(
			"rate limit remaing (%d) is less than expected fetch count (%d)",
			rateLimit.Resources.Core.Remaining,
			ExpectedFetchCount(workflowRuns),
		)
	}

	totalPage := TotalPage(workflowRuns)
	allWorkflowRuns := make([]github.WorkflowRun, workflowRuns.TotalCount)
	copy(allWorkflowRuns, workflowRuns.WorkflowRuns)

	// total_count is over 100
	if totalPage > 1 {
		for i := 2; i <= totalPage; i++ {
			w, err := github.FetchWorkflowRuns(repo, client, targetRange, config.PerPage, i)
			if err != nil {
				return github.Usage{}, err
			}

			allWorkflowRuns = append(allWorkflowRuns, w.WorkflowRuns...)
		}
	}

	usage := github.Usage{}
	for _, w := range allWorkflowRuns {
		u, err := w.Usage(client)
		if err != nil {
			return github.Usage{}, err
		}

		usage.Linux += u.Linux
		usage.Windows += u.Windows
		usage.Mac += u.Mac
	}

	return usage, nil
}

func IsRunnable(limits github.RateLimits, runs github.WorkflowRuns) bool {
	return ExpectedFetchCount(runs) <= limits.Resources.Core.Remaining
}

func ExpectedFetchCount(runs github.WorkflowRuns) int {
	return TotalPage(runs) + runs.TotalCount
}

func TotalPage(runs github.WorkflowRuns) int {
	return int(math.Ceil(float64(runs.TotalCount) / float64(config.PerPage)))
}
