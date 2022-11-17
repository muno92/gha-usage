package github_actions_usage_calculator

import (
	"github_actions_usage_calculator/github"
)

// Usage has each OS GitHub Actions runner execution time in seconds
type Usage struct {
	// Order by cost
	// https://docs.github.com/ja/billing/managing-billing-for-github-actions/about-billing-for-github-actions
	Linux   int64
	Windows int64
	Mac     int64
}

func FetchUsage(client github.Client, jobsUrl string) (Usage, error) {
	jobRuns, err := github.FetchJobRuns(client, jobsUrl)
	if err != nil {
		return Usage{}, err
	}

	u := Usage{}
	for _, job := range jobRuns.Jobs {
		u.Linux += int64(job.CompletedAt.Sub(job.StartedAt).Seconds())
	}

	return u, nil
}
