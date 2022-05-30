package main

import (
	"fmt"
	"os"
)

func ExitWithError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
