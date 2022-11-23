package main

import (
	"github_actions_usage_calculator/cmd"
	"os"
)

func main() {
	repo := os.Args[1]
	startDate := os.Args[2]
	endDate := os.Args[3]
	token := os.Getenv("GITHUB_TOKEN")

	usage, err := cmd.Run(repo, startDate, endDate, token)
	if err != nil {
		panic(err)
	}

	printer := cmd.SwitchPrinter()
	err = printer.Print(repo, startDate, endDate, usage)
	if err != nil {
		panic(err)
	}
}
