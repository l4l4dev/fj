package cli

import (
	"fmt"
	"strings"

	applicationrepository "github.com/l4l4dev/fj/internal/application/repository"
	"github.com/spf13/cobra"
)

func newRepositoryCommand(service applicationrepository.Service) *cobra.Command {
	command := &cobra.Command{Use: "repo", Short: "Manage repositories"}
	command.AddCommand(newRepositoryListCommand(service))
	var getter applicationrepository.Getter
	if candidate, ok := service.(applicationrepository.Getter); ok {
		getter = candidate
	}
	command.AddCommand(newRepositoryInspectCommand(getter))
	return command
}

func newRepositoryInspectCommand(getter applicationrepository.Getter) *cobra.Command {
	var instance string
	command := &cobra.Command{
		Use:   "inspect OWNER/NAME",
		Short: "Inspect a repository",
		Args: func(command *cobra.Command, args []string) error {
			if len(args) != 1 {
				return newCommandError(categoryValidation, "inspect repository", fmt.Errorf("OWNER/NAME is required"))
			}
			if err := validateRepositoryTarget(args[0]); err != nil {
				return newCommandError(categoryValidation, "inspect repository", err)
			}
			return nil
		},
		RunE: func(command *cobra.Command, args []string) error {
			if getter == nil {
				service, err := composeRepositoryService(command.Context(), instance)
				if err != nil {
					return err
				}
				var ok bool
				getter, ok = service.(applicationrepository.Getter)
				if !ok {
					return newCommandError(categoryInternal, "inspect repository", fmt.Errorf("repository service does not support inspection"))
				}
			}
			parts := strings.Split(args[0], "/")
			owner, name := parts[0], parts[1]
			result, err := applicationrepository.NewInspectUseCase(getter).Execute(command.Context(), applicationrepository.GetRequest{Owner: owner, Name: name})
			if err != nil {
				return mapInspectRepositoryError(err)
			}
			printRepository(command, result)
			return nil
		},
	}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func validateRepositoryTarget(target string) error {
	if strings.Count(target, "/") != 1 {
		return fmt.Errorf("OWNER/NAME must contain exactly one slash")
	}
	parts := strings.SplitN(target, "/", 2)
	if parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("OWNER/NAME owner and name are required")
	}
	return nil
}

func printRepository(command *cobra.Command, repository applicationrepository.Repository) {
	description := repository.Description
	if description == "" {
		description = "-"
	}
	fmt.Fprintf(command.OutOrStdout(), "Repository: %s/%s\nDescription: %s\nPrivate: %t\nArchived: %t\nDefault branch: %s\n", repository.Owner, repository.Name, description, repository.Private, repository.Archived, repository.DefaultBranch)
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
