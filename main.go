package main

import (
	"fmt"
	"os"

	"github.com/bantl23/yabba/cmd"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	return cmd.Execute()
}
