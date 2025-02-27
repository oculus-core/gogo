package main

import (
	"fmt"
	"os"

	"github.com/oculus-core/gogo/cmd/gogo"
)

func main() {
	// Execute the root command
	if err := gogo.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
