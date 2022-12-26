package cmd

import (
	"fmt"
	"ghausage/github"
	"io"
	"os"
)

type Printer interface {
	Print(stdout io.Writer, repo string, startDate string, endDate string, usage github.Usage) error
}

func SwitchPrinter() Printer {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return GitHubActionsPrinter{}
	}
	return CommandLinePrinter{}
}

type CommandLinePrinter struct {
}

func (p CommandLinePrinter) Print(stdout io.Writer, repo string, startDate string, endDate string, usage github.Usage) error {
	h, err := usage.HumanReadable()
	if err != nil {
		return err
	}

	message := fmt.Sprintf("%s (%s ~ %s) usage\n", repo, startDate, endDate)
	message += fmt.Sprintf("Linux: %s (%ds)\n", h.Linux, usage.Linux)
	message += fmt.Sprintf("Windows: %s (%ds)\n", h.Windows, usage.Windows)
	message += fmt.Sprintf("Mac: %s (%ds)\n", h.Mac, usage.Mac)
	message += fmt.Sprintf("self-hosted runner: %s (%ds)\n", h.SelfHosted, usage.SelfHosted)

	_, err = fmt.Fprint(stdout, message)
	if err != nil {
		return err
	}

	return nil
}

type GitHubActionsPrinter struct {
}

func (p GitHubActionsPrinter) Print(stdout io.Writer, repo string, startDate string, endDate string, usage github.Usage) error {
	h, err := usage.HumanReadable()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(os.Getenv("GITHUB_STEP_SUMMARY"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	message := "### GitHub Action Usage  \n"
	message += "#### Repository  \n"
	message += fmt.Sprintf("%s (from %s to %s)  \n", repo, startDate, endDate)
	message += "#### Usage time  \n"
	message += "| Linux | Windows | Mac | self-hosted runner |\n"
	message += "| --- | --- | --- | --- |\n"
	message += fmt.Sprintf(
		"| %s (%ds) | %s (%ds) | %s (%ds) | %s (%ds) |\n",
		h.Linux,
		usage.Linux,
		h.Windows,
		usage.Windows,
		h.Mac,
		usage.Mac,
		h.SelfHosted,
		usage.SelfHosted,
	)

	_, err = fmt.Fprint(f, message)
	if err != nil {
		return err
	}
	return nil
}
