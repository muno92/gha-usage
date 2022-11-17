package github

import (
	"encoding/json"
	"time"
)

type JobRuns struct {
	TotalCount int `json:"total_count"`
	Jobs       []Job
}

type Job struct {
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	Labels      []string
}

// Usage has each OS GitHub Actions runner execution time in seconds
type Usage struct {
	// Order by cost
	// https://docs.github.com/ja/billing/managing-billing-for-github-actions/about-billing-for-github-actions
	Linux   int64
	Windows int64
	Mac     int64
}

func FetchJobRuns(client Client, jobsUrl string) (JobRuns, error) {
	body, err := client.Get(jobsUrl)
	if err != nil {
		return JobRuns{}, err
	}

	j := JobRuns{}
	if err := json.Unmarshal(body, &j); err != nil {
		return JobRuns{}, err
	}

	return j, nil
}

func FetchUsage(client Client, jobsUrl string) (Usage, error) {
	jobRuns, err := FetchJobRuns(client, jobsUrl)
	if err != nil {
		return Usage{}, err
	}

	u := Usage{}
	for _, job := range jobRuns.Jobs {
		u.Linux += int64(job.CompletedAt.Sub(job.StartedAt).Seconds())
	}

	return u, nil
}
