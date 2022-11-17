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
