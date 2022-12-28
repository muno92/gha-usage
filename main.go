package main

import (
	"flag"
	"ghausage/cmd"
	"log"
	"os"
)

var v bool
var countCmd *flag.FlagSet
var sumCmd *flag.FlagSet
var repo string
var startDate string
var endDate string

func init() {
	countCmd = flag.NewFlagSet("count", flag.ExitOnError)
	SetFlagOption(countCmd)
	countCmd.Usage = PrintCountHelp

	sumCmd = flag.NewFlagSet("sum", flag.ExitOnError)
	SetFlagOption(sumCmd)
	sumCmd.Usage = PrintSumHelp

	flag.BoolVar(&v, "v", false, "Show Version")
	flag.BoolVar(&v, "version", false, "Show Version")
}

func SetFlagOption(flag *flag.FlagSet) {
	flag.StringVar(&repo, "repo", "", "target repository (format: OWNER/REPO)")
	flag.StringVar(&startDate, "start", "", "start date of count range (format: yyyy-mm-dd)")
	flag.StringVar(&endDate, "end", "", "end date of count range (format: yyyy-mm-dd)")
}

func main() {
	flag.Usage = PrintHelp
	flag.Parse()
	if v {
		cmd.PrintVersion()
		return
	}

	token := os.Getenv("GITHUB_TOKEN")

	if len(os.Args) == 1 {
		PrintHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "count":
		err := countCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if repo == "" || startDate == "" || endDate == "" {
			PrintCountHelp()
			os.Exit(1)
		}

		err = cmd.CountCommand{Logger: log.Default()}.Run(os.Stdout, repo, startDate, endDate, token)
		if err != nil {
			panic(err)
		}
		return
	case "sum":
		err := sumCmd.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		if repo == "" || startDate == "" || endDate == "" {
			PrintSumHelp()
			os.Exit(1)
		}

		err = cmd.SumCommand{Logger: log.Default()}.Run(os.Stdout, repo, startDate, endDate, token)
		if err != nil {
			panic(err)
		}
		return
	default:
	}
}

func PrintCountHelp() {
	println("Usage of count: ./ghausage count --repo OWNER/REPO --start YYYY-MM-DD --end YYYY-MM-DD")
	countCmd.PrintDefaults()
}

func PrintSumHelp() {
	println("Usage of sum: ./ghausage sum --repo OWNER/REPO --start YYYY-MM-DD --end YYYY-MM-DD")
	sumCmd.PrintDefaults()
}

func PrintHelp() {
	println("Usage of ./ghausage")
	println("./ghausage [-v | --version] [-h | --help]")
	println("./ghausage count    Count GitHub Actions workflow runs of specified repository")
	println("./ghausage sum      Sum GitHub Actions usage time of specified repository")
}
