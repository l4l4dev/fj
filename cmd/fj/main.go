package main

import (
	"os"

	"github.com/l4l4dev/fj/internal/interface/cli"
)

func main() {
	if code := cli.Execute(cli.NewRootCommand(), os.Args[1:]); code != 0 {
		os.Exit(code)
	}
}
