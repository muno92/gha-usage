package github

import (
	"encoding/json"
	"fmt"
)

type WorkflowRuns struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

type WorkflowRun struct {
	JobsUrl string `json:"jobs_url"`
}

func (w WorkflowRun) JobRuns(client Client) (JobRuns, error) {
	body, err := client.Get(w.JobsUrl)
	if err != nil {
		return JobRuns{}, err
	}

	j := JobRuns{}
	if err := json.Unmarshal(body, &j); err != nil {
		return JobRuns{}, err
	}

	return j, nil
}

func (w WorkflowRun) Usage(client Client) (Usage, error) {
	jobRuns, err := w.JobRuns(client)
	if err != nil {
		return Usage{}, err
	}

	return jobRuns.Usage(), nil
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
