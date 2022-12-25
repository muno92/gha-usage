package cmd

import (
	"fmt"
	"ghausage/config"
	"ghausage/github"
	"io"
	"log"
)

type CountCommand struct {
	Logger *log.Logger
}

func (c CountCommand) Run(stdout io.Writer, repo string, startDate string, endDate string, token string) error {
	r, err := github.NewRange(startDate, endDate)
	if err != nil {
		return err
	}

	client := github.Client{Token: token, Logger: c.Logger}

	w, err := github.FetchWorkflowRuns(repo, client, r, config.PerPage, 1)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(stdout, "%s workflow run count (from %s to %s): %d\n", repo, startDate, endDate, w.TotalCount)
	if err != nil {
		return err
	}

	return nil
}
