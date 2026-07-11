package cli

import (
	"fmt"
	"io"

	applicationversion "github.com/l4l4dev/fj/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCommand(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the fj version",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return printVersion(command.OutOrStdout(), version)
		},
	}
}

func printVersion(w io.Writer, version string) error {
	_, err := fmt.Fprintln(w, version)
	return err
}

func defaultVersion() string { return applicationversion.Current() }
