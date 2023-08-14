package main

import (
	"fmt"
	"os"
)

func main() {
	AnalysisExifProcess()
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
