package github_actions_usage_calculator

import (
	"encoding/json"
	"fmt"
)

type WorkflowRun struct {
	TotalCount int `json:"total_count"`
}

func FetchWorkflowRun(repo string, client Client, targetRange Range) (*WorkflowRun, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/actions/runs?created=%s..%s",
		repo,
		targetRange.Start.Format("2006-01-02"),
		targetRange.End.Format("2006-01-02"),
	)

	body, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	w := WorkflowRun{}
	if err := json.Unmarshal(body, &w); err != nil {
		return nil, err
	}

	return &w, nil
}
