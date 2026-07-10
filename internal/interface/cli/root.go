package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:           "fj",
		Short:         "AI-first CLI for Forgejo",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(command *cobra.Command, _ []string) error {
			return command.Help()
		},
	}
	command.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return newCommandError(categoryValidation, "execute command", err)
	})
	return command
}
