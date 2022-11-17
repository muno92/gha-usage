package github

import (
	"os"
	"testing"
)

func TestFetchJobUsage(t *testing.T) {
	client := Client{Token: os.Getenv("GITHUB_TOKEN")}

	tests := []struct {
		name          string
		jobsUrl       string
		expectedUsage Usage
	}{
		{
			name:    "single job",
			jobsUrl: "https://api.github.com/repos/muno92/resharper_inspectcode/actions/runs/3370110776/jobs",
			expectedUsage: Usage{
				Linux:   15,
				Windows: 0,
				Mac:     0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage, err := FetchUsage(client, tt.jobsUrl)
			if err != nil {
				panic(err)
			}

			if usage != tt.expectedUsage {
				t.Errorf("Expected usage is %v, got %v", tt.expectedUsage, usage)
			}
		})
	}
}

func TestJobRunnerType(t *testing.T) {
	tests := []struct {
		name               string
		job                Job
		expectedRunnerType RunnerType
	}{
		{
			name:               "ubuntu-latest",
			job:                Job{Labels: []string{"ubuntu-latest"}},
			expectedRunnerType: Linux,
		},
		{
			name:               "ubuntu-latest on second label",
			job:                Job{Labels: []string{"test", "ubuntu-latest"}},
			expectedRunnerType: Linux,
		},
		{
			name:               "ubuntu-22.04",
			job:                Job{Labels: []string{"ubuntu-22.04"}},
			expectedRunnerType: Linux,
		},
		{
			name:               "ubuntu-20.04",
			job:                Job{Labels: []string{"ubuntu-20.04"}},
			expectedRunnerType: Linux,
		},
		{
			name:               "ubuntu-18.04",
			job:                Job{Labels: []string{"ubuntu-18.04"}},
			expectedRunnerType: Linux,
		},
		{
			name:               "windows-latest",
			job:                Job{Labels: []string{"windows-latest"}},
			expectedRunnerType: Windows,
		},
		{
			name:               "windows-2022",
			job:                Job{Labels: []string{"windows-2022"}},
			expectedRunnerType: Windows,
		},
		{
			name:               "windows-2019",
			job:                Job{Labels: []string{"windows-2019"}},
			expectedRunnerType: Windows,
		},
		{
			name:               "macos-latest",
			job:                Job{Labels: []string{"macos-latest"}},
			expectedRunnerType: Mac,
		},
		{
			name:               "macos-12",
			job:                Job{Labels: []string{"macos-12"}},
			expectedRunnerType: Mac,
		},
		{
			name:               "macos-11",
			job:                Job{Labels: []string{"macos-11"}},
			expectedRunnerType: Mac,
		},
		{
			name:               "macos-10.15",
			job:                Job{Labels: []string{"macos-10.15"}},
			expectedRunnerType: Mac,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.job.RunnerType()

			if r != tt.expectedRunnerType {
				t.Errorf("Expected %v, got %v", tt.expectedRunnerType, r)
			}
		})
	}
}
