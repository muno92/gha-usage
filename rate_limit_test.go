package github_actions_usage_calculator

import (
	"os"
	"testing"
)

func TestFetchRateLimit(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")

	rateLimit, err := FetchRateLimit(token)

	if err != nil {
		panic(err)
	}

	if rateLimit.Resources.Core.Limit != 5000 {
		t.Errorf("Expected Rate Limit is 5000 (because token is Personal Access Token), but got %v", rateLimit.Resources.Core.Limit)
	}
}
