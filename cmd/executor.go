package cmd

import (
	"fmt"
	"github_actions_usage_calculator/config"
	"github_actions_usage_calculator/github"
	"math"
)

type WorkflowRunResult struct {
	WorkflowRuns github.WorkflowRuns
	Error        error
}

type UsageResult struct {
	Usage github.Usage
	Error error
}

func Run(repo string, startDate string, endDate string, token string) (github.Usage, error) {
	targetRange, err := github.NewRange(startDate, endDate)
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

	isRunnable, err := IsRunnable(rateLimit, workflowRuns)
	if !isRunnable {
		return github.Usage{}, err
	}

	fmt.Printf("Workflow run count: %d\n", workflowRuns.TotalCount)

	totalPage := TotalPage(workflowRuns)
	allWorkflowRuns := make([]github.WorkflowRun, workflowRuns.TotalCount)
	allWorkflowRuns = workflowRuns.WorkflowRuns

	fmt.Printf("Complete fetch workflow run with pagination (1/%d)\n", totalPage)

	// total_count is over 100
	if totalPage > 1 {
		wc := make(chan WorkflowRunResult)
		defer close(wc)
		for i := 2; i <= totalPage; i++ {
			go func(page int) {
				w, err := github.FetchWorkflowRuns(repo, client, targetRange, config.PerPage, page)
				if err != nil {
					wc <- WorkflowRunResult{WorkflowRuns: github.WorkflowRuns{}, Error: err}
					return
				}
				wc <- WorkflowRunResult{WorkflowRuns: w, Error: nil}
			}(i)
		}

		for j := 2; j <= totalPage; j++ {
			w := <-wc
			if w.Error != nil {
				return github.Usage{}, w.Error
			}

			allWorkflowRuns = append(allWorkflowRuns, w.WorkflowRuns.WorkflowRuns...)
			fmt.Printf("Complete fetch workflow run with pagination (%d/%d)\n", j, totalPage)
		}
	}

	uc := make(chan UsageResult)

	usage := github.Usage{}
	defer close(uc)
	for _, w := range allWorkflowRuns {
		go func(w github.WorkflowRun) {
			u, err := w.Usage(client)
			if err != nil {
				uc <- UsageResult{Usage: github.Usage{}, Error: err}
				return
			}
			uc <- UsageResult{Usage: u, Error: nil}
		}(w)
	}

	for k := 0; k < workflowRuns.TotalCount; k++ {
		u := <-uc
		if u.Error != nil {
			return github.Usage{}, u.Error
		}

		usage.Linux += u.Usage.Linux
		usage.Windows += u.Usage.Windows
		usage.Mac += u.Usage.Mac

		fmt.Printf("Complete fetch job (%d/%d)\n", k+1, workflowRuns.TotalCount)
	}

	return usage, nil
}

func IsRunnable(limits github.RateLimits, runs github.WorkflowRuns) (bool, error) {
	if !jobIsAcquirable(runs) {
		return false, fmt.Errorf(
			"count of workflow run (%d) must be less than or equal to 1000 (because GitHub API will not return over 1000 records even with pagination)",
			runs.TotalCount,
		)
	}

	return RateLimitIsEnough(limits, runs), fmt.Errorf(
		"rate limit remaining (%d) is less than expected fetch count (%d)",
		limits.Resources.Core.Remaining,
		ExpectedFetchCount(runs),
	)
}

func jobIsAcquirable(runs github.WorkflowRuns) bool {
	return runs.TotalCount <= 1000
}

func RateLimitIsEnough(limits github.RateLimits, runs github.WorkflowRuns) bool {
	return ExpectedFetchCount(runs) <= limits.Resources.Core.Remaining
}

func ExpectedFetchCount(runs github.WorkflowRuns) int {
	return TotalPage(runs) + runs.TotalCount
}

func TotalPage(runs github.WorkflowRuns) int {
	return int(math.Ceil(float64(runs.TotalCount) / float64(config.PerPage)))
}
