package main

import (
	"fmt"
	"os"

	"codeberg.org/jhot/db-check/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Success")
}
