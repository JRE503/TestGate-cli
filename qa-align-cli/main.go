package main

import (
	"fmt"
	"os"
	"qa-align-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Engine system halt: %v\n", err)
		os.Exit(1)
	}
}
