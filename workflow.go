package github_actions_usage_calculator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WorkflowRun struct {
	TotalCount int `json:"total_count"`
}

func FetchWorkflowRun(repo string, token string, targetRange Range) (*WorkflowRun, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/actions/runs?created=%s..%s",
		repo,
		targetRange.Start.Format("2006-01-02"),
		targetRange.End.Format("2006-01-02"),
	)

	client := &http.Client{}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

	w := WorkflowRun{}

	if err := json.Unmarshal(body, &w); err != nil {
		return nil, err
	}

	return &w, nil
}
