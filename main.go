package main

import (
	"fmt"
	"os"

	"github.com/jhot21/db-check/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
