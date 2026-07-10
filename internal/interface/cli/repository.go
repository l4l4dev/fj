package cli

import (
	"fmt"

	applicationrepository "github.com/l4l4dev/fj/internal/application/repository"
	"github.com/spf13/cobra"
)

func newRepositoryCommand(service applicationrepository.Service) *cobra.Command {
	command := &cobra.Command{Use: "repo", Short: "Manage repositories"}
	command.AddCommand(newRepositoryListCommand(service))
	return command
}

func newRepositoryListCommand(service applicationrepository.Service) *cobra.Command {
	var instance string
	var page, limit int
	command := &cobra.Command{
		Use:   "list",
		Short: "List repositories",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if page < 1 {
				return newCommandError(categoryValidation, "list repositories", fmt.Errorf("page must be at least 1"))
			}
			if limit < 1 {
				return newCommandError(categoryValidation, "list repositories", fmt.Errorf("limit must be at least 1"))
			}
			if service == nil {
				var err error
				service, err = composeRepositoryService(command.Context(), instance)
				if err != nil {
					return err
				}
			}
			result, err := applicationrepository.NewListUseCase(service).Execute(command.Context(), applicationrepository.ListRequest{Page: page, Limit: limit})
			if err != nil {
				return mapRepositoryError(err)
			}
			if len(result) == 0 {
				fmt.Fprintln(command.OutOrStdout(), "No repositories found.")
				return nil
			}
			for _, repository := range result {
				fmt.Fprintf(command.OutOrStdout(), "%s/%s\n", repository.Owner, repository.Name)
			}
			return nil
		},
	}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	command.Flags().IntVar(&page, "page", 1, "page number")
	command.Flags().IntVar(&limit, "limit", 30, "number of repositories per page")
	return command
}
