package main

import (
	"os"

	"github.com/l4l4dev/fj/internal/interface/cli"
	applicationversion "github.com/l4l4dev/fj/internal/version"
)

func main() {
	if code := cli.Execute(cli.NewRootCommandWithVersion(cli.RepositoryDependencies{}, applicationversion.Current()), os.Args[1:]); code != 0 {
		os.Exit(code)
	}
}
