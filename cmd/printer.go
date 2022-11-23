package cmd

import (
	"fmt"
	"github_actions_usage_calculator/github"
	"os"
)

type Printer interface {
	Print(repo string, startDate string, endDate string, usage github.Usage) error
}

func SwitchPrinter() Printer {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return GitHubActionsPrinter{}
	}
	return CommandLinePrinter{}
}

type CommandLinePrinter struct {
}

func (p CommandLinePrinter) Print(repo string, startDate string, endDate string, usage github.Usage) error {
	h, err := usage.HumanReadable()
	if err != nil {
		return err
	}

	fmt.Printf("%s (%s ~ %s) usage\n", repo, startDate, endDate)
	fmt.Printf("Linux: %s (%ds)\n", h.Linux, usage.Linux)
	fmt.Printf("Windows: %s (%ds)\n", h.Windows, usage.Windows)
	fmt.Printf("Mac: %s (%ds)\n", h.Mac, usage.Mac)

	return nil
}

type GitHubActionsPrinter struct {
}

func (p GitHubActionsPrinter) Print(repo string, startDate string, endDate string, usage github.Usage) error {
	h, err := usage.HumanReadable()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(os.Getenv("GITHUB_STEP_SUMMARY"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("### %s GitHub Action Usage (%s ~ %s)  \n", repo, startDate, endDate)
	message += "| Linux | Windows | Mac |\n"
	message += "| --- | --- | --- |\n"
	message += fmt.Sprintf("| %s (%ds) | %s (%ds) | %s (%ds) |\n", h.Linux, usage.Linux, h.Windows, usage.Windows, h.Mac, usage.Mac)

	_, err = fmt.Fprint(f, message)
	if err != nil {
		return err
	}
	return nil
}
