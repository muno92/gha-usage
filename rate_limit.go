package github_actions_usage_calculator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RateLimit struct {
	Limit     int
	Remaining int
	Reset     int
	Used      int
}

type Resource struct {
	// Other object (e.g. search) is not needed
	// https://docs.github.com/en/rest/rate-limit
	Core RateLimit
}

type RateLimits struct {
	Resources Resource
	// Deprecated
	Rate RateLimit
}

func FetchRateLimit(token string) (*RateLimits, error) {
	client := &http.Client{}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/rate_limit", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r := RateLimits{}
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	return &r, nil
}
