package main

import (
	"os"

	"github.com/l4l4dev/fj/internal/interface/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
