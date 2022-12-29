package github

import (
	"log"
	"os"
	"testing"
)

func TestFetchRateLimit(t *testing.T) {
	client := NewClient(os.Getenv("GITHUB_TOKEN"), log.Default())

	rateLimit, err := FetchRateLimit(client)

	if err != nil {
		t.Error(err)
	}

	if rateLimit.Resources.Core.Limit != 5000 {
		t.Errorf("Expected Rate Limit is 5000 (because token is Personal Access Token), but got %v", rateLimit.Resources.Core.Limit)
	}
}
