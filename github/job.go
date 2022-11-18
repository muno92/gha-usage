package github

import (
	"encoding/json"
	"strings"
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

//go:generate stringer -type=RunnerType
type RunnerType int

const (
	Linux RunnerType = iota + 1
	Windows
	Mac
)

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

	return jobRuns.Usage(), nil
}

func (j JobRuns) Usage() Usage {
	u := Usage{}
	for _, job := range j.Jobs {
		if job.RunnerType() == Windows {
			u.Windows += job.Usage()
			continue
		}
		if job.RunnerType() == Mac {
			u.Mac += job.Usage()
			continue
		}
		u.Linux += job.Usage()
	}

	return u
}

func (j Job) Usage() int64 {
	return int64(j.CompletedAt.Sub(j.StartedAt).Seconds())
}

func (j Job) RunnerType() RunnerType {
	for _, l := range j.Labels {
		label := strings.ToLower(l)
		if IsWindowsRunner(label) {
			return Windows
		}
		if IsMacRunner(label) {
			return Mac
		}
	}
	return Linux
}

func IsWindowsRunner(label string) bool {
	if label == "windows-latest" {
		return true
	}
	if label == "windows-2022" {
		return true
	}
	if label == "windows-2019" {
		return true
	}
	return false
}

func IsMacRunner(label string) bool {
	if label == "macos-latest" {
		return true
	}
	if label == "macos-12" {
		return true
	}
	if label == "macos-11" {
		return true
	}
	// Deprecated
	if label == "macos-10.15" {
		return true
	}
	return false
}
