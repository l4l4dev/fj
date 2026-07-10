package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fj",
		Short: "AI-first CLI for Forgejo",
		RunE: func(command *cobra.Command, _ []string) error {
			return command.Help()
		},
	}
}
