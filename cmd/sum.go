package cmd

import (
	"fmt"
	"ghausage/config"
	"ghausage/github"
	"io"
	"log"
	"math"
	"time"
)

type WorkflowRunResult struct {
	WorkflowRuns github.WorkflowRuns
	Error        error
}

type SumCommand struct {
	Logger *log.Logger
}

func (s SumCommand) Run(stdout io.Writer, repo string, startDate string, endDate string, token string) error {
	usage, err := SumUsage(repo, startDate, endDate, token, s.Logger)
	if err != nil {
		return err
	}

	printer, err := NewPrinter(stdout)
	if err != nil {
		return err
	}

	err = printer.Print(repo, startDate, endDate, usage)
	if err != nil {
		return err
	}

	return nil
}

func SumUsage(repo string, startDate string, endDate string, token string, logger *log.Logger) (github.Usage, error) {
	targetRange, err := github.NewRange(startDate, endDate)
	if err != nil {
		return github.Usage{}, err
	}

	client := github.NewClient(token, logger)

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

	logger.Printf("Workflow run count: %d\n", workflowRuns.TotalCount)

	totalPage := TotalPage(workflowRuns)
	allWorkflowRuns := make([]github.WorkflowRun, workflowRuns.TotalCount)
	allWorkflowRuns = workflowRuns.WorkflowRuns

	logger.Printf("Complete fetch workflow run with pagination (1/%d)\n", totalPage)

	// total_count is over 100
	if totalPage > 1 {
		wc := make(chan WorkflowRunResult)
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
			logger.Printf("Complete fetch workflow run with pagination (%d/%d)\n", j, totalPage)
		}
	}

	uc := make(chan github.UsageResult)
	for _, w := range allWorkflowRuns {
		go func(w github.WorkflowRun) {
			u, err := w.Usage(client)
			if err != nil {
				uc <- github.UsageResult{Usage: github.Usage{}, Error: err}
				return
			}
			uc <- github.UsageResult{Usage: u, Error: nil}
		}(w)
	}

	usage := github.Usage{}
	for k := 0; k < workflowRuns.TotalCount; k++ {
		u := <-uc
		if u.Error != nil {
			return github.Usage{}, u.Error
		}

		usage = usage.Plus(u.Usage)

		logger.Printf("Complete fetch job (%d/%d)\n", k+1, workflowRuns.TotalCount)
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
		"please try again at %s, because rate limit remaining (%d) is less than expected fetch count (%d)",
		time.Unix(limits.Resources.Core.Reset, 0).In(time.UTC).Format(time.RFC3339),
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
