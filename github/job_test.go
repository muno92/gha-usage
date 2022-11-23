package github

import (
	"testing"
)

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
			name:               "Windows-latest",
			job:                Job{Labels: []string{"Windows-latest"}},
			expectedRunnerType: Windows,
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
			name:               "macOS-latest",
			job:                Job{Labels: []string{"macOS-latest"}},
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

func TestHumanReadableUsage(t *testing.T) {
	usage := Usage{
		Linux:   59,
		Windows: 60,
		Mac:     3600,
	}
	expected := HumanReadableUsage{
		Linux:   "59s",
		Windows: "1m0s",
		Mac:     "1h0m0s",
	}

	actual, err := usage.HumanReadable()
	if err != nil {
		panic(err)
	}
	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
