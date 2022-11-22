package main

import (
	"fmt"
	"github_actions_usage_calculator/cmd"
	"os"
)

func main() {
	repo := os.Args[1]
	targetMonth := os.Args[2]
	token := os.Getenv("GITHUB_TOKEN")

	usage, err := cmd.Run(repo, targetMonth, token)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s %s usage\n", repo, targetMonth)
	fmt.Printf("Linux: %ds\n", usage.Linux)
	fmt.Printf("Windows: %ds\n", usage.Windows)
	fmt.Printf("Mac: %ds\n", usage.Mac)
}
