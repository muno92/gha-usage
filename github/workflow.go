package github

import (
	"encoding/json"
	"fmt"
)

type WorkflowRun struct {
	JobsUrl string `json:"jobs_url"`
}

type WorkflowRuns struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

func FetchWorkflowRuns(repo string, client Client, targetRange Range, perPage int, page int) (WorkflowRuns, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/actions/runs?created=%s..%s&per_page=%d&page=%d",
		repo,
		targetRange.Start.Format("2006-01-02"),
		targetRange.End.Format("2006-01-02"),
		perPage,
		page,
	)

	body, err := client.Get(url)
	if err != nil {
		return WorkflowRuns{}, err
	}

	w := WorkflowRuns{}
	if err := json.Unmarshal(body, &w); err != nil {
		return WorkflowRuns{}, err
	}

	return w, nil
}
