package cmd

import (
	"fmt"
	"ghausage/config"
)

func PrintVersion() {
	fmt.Printf("%s version %s\n", config.Name, config.Version)
}
