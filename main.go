package main

import (
	"flag"
	"ghausage/cmd"
	"log"
	"os"
)

var v bool
var countCmd *flag.FlagSet
var repo string
var startDate string
var endDate string

func init() {
	countCmd = flag.NewFlagSet("count", flag.ExitOnError)
	countCmd.StringVar(&repo, "repo", "", "target repository (format: OWNER/REPOSITORY)")
	countCmd.StringVar(&startDate, "start", "", "start date of count range (format: yyyy-mm-dd)")
	countCmd.StringVar(&endDate, "end", "", "end date of count range (format: yyyy-mm-dd)")

	flag.BoolVar(&v, "version", false, "Show Version")
}

func main() {
	flag.Parse()
	if v {
		cmd.PrintVersion()
		return
	}

	token := os.Getenv("GITHUB_TOKEN")

	switch os.Args[1] {
	case "count":
		err := countCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if repo == "" || startDate == "" || endDate == "" {
			println("Usage of count:")
			countCmd.PrintDefaults()
			os.Exit(1)
		}

		err = cmd.CountCommand{Logger: log.Default()}.Run(os.Stdout, repo, startDate, endDate, token)
		if err != nil {
			panic(err)
		}
		return
	default:
	}

	repo := os.Args[1]
	startDate := os.Args[2]
	endDate := os.Args[3]

	usage, err := cmd.Run(repo, startDate, endDate, token, log.Default())
	if err != nil {
		panic(err)
	}

	printer := cmd.SwitchPrinter()
	err = printer.Print(repo, startDate, endDate, usage)
	if err != nil {
		panic(err)
	}
}
