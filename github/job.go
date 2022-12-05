package github

import (
	"fmt"
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
	Steps       []Step
}

type Step struct {
	Name        string
	Status      string
	Conclusion  string
	number      int
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
}

//go:generate stringer -type=RunnerType
type RunnerType int

const (
	Linux RunnerType = iota + 1
	Windows
	Mac
	SelfHosted
)

// Usage has each OS GitHub Actions runner execution time in seconds
type Usage struct {
	// Order by cost
	// https://docs.github.com/ja/billing/managing-billing-for-github-actions/about-billing-for-github-actions
	Linux      int64
	Windows    int64
	Mac        int64
	SelfHosted int64
}

func (u Usage) HumanReadable() (HumanReadableUsage, error) {
	linux, err := ToString(u.Linux)
	if err != nil {
		return HumanReadableUsage{}, err
	}
	windows, err := ToString(u.Windows)
	if err != nil {
		return HumanReadableUsage{}, err
	}
	mac, err := ToString(u.Mac)
	if err != nil {
		return HumanReadableUsage{}, err
	}
	selfHosted, err := ToString(u.SelfHosted)
	if err != nil {
		return HumanReadableUsage{}, err
	}

	return HumanReadableUsage{
		Linux:      linux,
		Windows:    windows,
		Mac:        mac,
		SelfHosted: selfHosted,
	}, nil
}

func ToString(seconds int64) (string, error) {
	s, err := time.ParseDuration(fmt.Sprintf("%ds", seconds))
	if err != nil {
		return "", err
	}
	return s.String(), nil
}

type HumanReadableUsage struct {
	Linux      string
	Windows    string
	Mac        string
	SelfHosted string
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
		if job.RunnerType() == Linux {
			u.Linux += job.Usage()
			continue
		}
		u.SelfHosted += job.Usage()
	}

	return u
}

func (j Job) Usage() int64 {
	if j.CompletedAt.IsZero() {
		stepUsage := j.StepUsage()
		if stepUsage != 0 {
			return stepUsage
		}
		return 0
	}
	return int64(j.CompletedAt.Sub(j.StartedAt).Seconds())
}

func (j Job) StepUsage() int64 {
	var u int64
	for _, step := range j.Steps {
		u += step.Usage()
	}
	return u
}

func (s Step) Usage() int64 {
	if s.CompletedAt.IsZero() {
		return 0
	}
	return int64(s.CompletedAt.Sub(s.StartedAt).Seconds())
}

func (j Job) RunnerType() RunnerType {
	for _, l := range j.Labels {
		label := strings.ToLower(l)
		if IsLinuxRunner(label) {
			return Linux
		}
		if IsWindowsRunner(label) {
			return Windows
		}
		if IsMacRunner(label) {
			return Mac
		}
	}
	return SelfHosted
}

func IsLinuxRunner(label string) bool {
	if label == "ubuntu-latest" {
		return true
	}
	if label == "ubuntu-22.04" {
		return true
	}
	if label == "ubuntu-20.04" {
		return true
	}
	if label == "ubuntu-18.04" {
		return true
	}
	return false
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
