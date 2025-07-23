package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
