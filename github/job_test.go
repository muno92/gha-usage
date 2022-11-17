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
