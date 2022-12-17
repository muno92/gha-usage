package github

import (
	"encoding/json"
	"fmt"
	"ghausage/config"
)

type WorkflowRuns struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

type WorkflowRun struct {
	JobsUrl string `json:"jobs_url"`
}

func (w WorkflowRun) JobRuns(client Client, page int) (JobRuns, error) {
	body, err := client.Get(fmt.Sprintf(
		"%s?per_page=%d&page=%d",
		w.JobsUrl,
		// Maybe, job per workflow is under 100
		config.PerPage,
		page,
	))
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
	jobRuns, err := w.JobRuns(client, 1)
	if err != nil {
		return Usage{}, err
	}

	u := jobRuns.Usage()

	totalPage := jobRuns.TotalPage()

	if totalPage > 1 {
		uc := make(chan UsageResult)
		for i := 2; i <= totalPage; i++ {
			go func(page int) {
				jobRuns, err := w.JobRuns(client, page)
				if err != nil {
					uc <- UsageResult{Usage: Usage{}, Error: err}
					return
				}
				uc <- UsageResult{Usage: jobRuns.Usage(), Error: nil}
			}(i)
		}

		for j := 2; j <= totalPage; j++ {
			result := <-uc
			if result.Error != nil {
				return Usage{}, result.Error
			}

			u = u.Plus(result.Usage)
		}
	}

	return u, nil
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
