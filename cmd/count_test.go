package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestCountCommand(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")

	tests := []struct {
		name          string
		repo          string
		startDate     string
		endDate       string
		expectedCount int
	}{
		{
			name:          "count per month",
			repo:          "muno92/resharper_inspectcode",
			startDate:     "2022-11-01",
			endDate:       "2022-11-30",
			expectedCount: 21,
		},
		{
			name:          "count per week",
			repo:          "muno92/life_log",
			startDate:     "2022-11-06",
			endDate:       "2022-11-12",
			expectedCount: 386,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := new(bytes.Buffer)

			err := CountCommand{Logger: log.Default()}.Run(stdout, tt.repo, tt.startDate, tt.endDate, token)
			if err != nil {
				t.Error(err)
			}

			expected := fmt.Sprintf("%s workflow run count (from %s to %s): %d\n", tt.repo, tt.startDate, tt.endDate, tt.expectedCount)
			actual := stdout.String()

			if actual != expected {
				t.Errorf("Expected message is\n\t%v,\ngot\n\t%v", expected, actual)
			}
		})
	}
}
