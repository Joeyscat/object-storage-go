package main

import (
	"fmt"
	"os"

	"github.com/joeyscat/object-storage-go/oshell"
)

func main() {
	if err := oshell.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
