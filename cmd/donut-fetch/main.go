package main

import (
	"fmt"
	"os"

	"github.com/neflyte/donut-fetch/cmd/donut-fetch/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("*  program error: %s", err.Error())
		os.Exit(1)
	}
}
