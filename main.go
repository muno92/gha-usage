package main

import (
	"flag"
	"ghausage/cmd"
	"os"
)

var v bool

func init() {
	flag.BoolVar(&v, "version", false, "Show Version")
}

func main() {
	flag.Parse()
	if v {
		cmd.PrintVersion()
		return
	}

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
