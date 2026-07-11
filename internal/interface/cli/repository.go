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
	var creator applicationrepository.Creator
	if candidate, ok := service.(applicationrepository.Creator); ok {
		creator = candidate
	}
	command.AddCommand(newRepositoryCreateCommand(creator))
	var updater applicationrepository.Updater
	if candidate, ok := service.(applicationrepository.Updater); ok {
		updater = candidate
	}
	command.AddCommand(newRepositoryUpdateCommand(updater))
	var archiver applicationrepository.Archiver
	if candidate, ok := service.(applicationrepository.Archiver); ok {
		archiver = candidate
	}
	command.AddCommand(newRepositoryArchiveCommand(archiver, true))
	command.AddCommand(newRepositoryArchiveCommand(archiver, false))
	return command
}

func newRepositoryArchiveCommand(archiver applicationrepository.Archiver, archived bool) *cobra.Command {
	name := "restore"
	operation := "restore repository"
	if archived {
		name = "archive"
		operation = "archive repository"
	}
	var instance string
	command := &cobra.Command{Use: name + " OWNER/NAME", Short: name + " a repository", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return newCommandError(categoryValidation, operation, fmt.Errorf("OWNER/NAME is required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, operation, err)
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if archiver == nil {
			service, err := composeRepositoryService(command.Context(), instance)
			if err != nil {
				return err
			}
			var ok bool
			archiver, ok = service.(applicationrepository.Archiver)
			if !ok {
				return newCommandError(categoryInternal, operation, fmt.Errorf("repository service does not support archive operations"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		result, err := applicationrepository.NewArchiveUseCase(archiver).Execute(command.Context(), applicationrepository.ArchiveRequest{Owner: parts[0], Name: parts[1], Archived: archived})
		if err != nil {
			return mapArchiveRepositoryError(err, operation)
		}
		fmt.Fprintf(command.OutOrStdout(), "Repository: %s/%s\nArchived: %t\n", result.Owner, result.Name, result.Archived)
		return nil
	}}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newRepositoryUpdateCommand(updater applicationrepository.Updater) *cobra.Command {
	var description, visibility, instance string
	var descriptionSet, visibilitySet bool
	command := &cobra.Command{
		Use: "update OWNER/NAME", Short: "Update repository metadata",
		Args: func(command *cobra.Command, args []string) error {
			descriptionSet = command.Flags().Changed("description")
			visibilitySet = command.Flags().Changed("visibility")
			if len(args) != 1 {
				return newCommandError(categoryValidation, "update repository", fmt.Errorf("OWNER/NAME is required"))
			}
			if err := validateRepositoryTarget(args[0]); err != nil {
				return newCommandError(categoryValidation, "update repository", err)
			}
			if !descriptionSet && !visibilitySet {
				return newCommandError(categoryValidation, "update repository", fmt.Errorf("at least one repository field is required"))
			}
			if visibilitySet && visibility != "public" && visibility != "private" {
				return newCommandError(categoryValidation, "update repository", fmt.Errorf("visibility must be public or private"))
			}
			return nil
		},
		RunE: func(command *cobra.Command, args []string) error {
			if updater == nil {
				service, err := composeRepositoryService(command.Context(), instance)
				if err != nil {
					return err
				}
				var ok bool
				updater, ok = service.(applicationrepository.Updater)
				if !ok {
					return newCommandError(categoryInternal, "update repository", fmt.Errorf("repository service does not support updates"))
				}
			}
			parts := strings.SplitN(args[0], "/", 2)
			request := applicationrepository.UpdateRequest{Owner: parts[0], Name: parts[1]}
			if descriptionSet {
				request.Description = &description
			}
			if visibilitySet {
				private := visibility == "private"
				request.Private = &private
			}
			result, err := applicationrepository.NewUpdateUseCase(updater).Execute(command.Context(), request)
			if err != nil {
				return mapUpdateRepositoryError(err)
			}
			fields := make([]string, 0, 2)
			if descriptionSet {
				fields = append(fields, "description")
			}
			if visibilitySet {
				fields = append(fields, "visibility")
			}
			printUpdatedRepository(command, result, fields)
			return nil
		},
	}
	command.Flags().StringVar(&description, "description", "", "repository description")
	command.Flags().StringVar(&visibility, "visibility", "", "repository visibility (public or private)")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newRepositoryCreateCommand(creator applicationrepository.Creator) *cobra.Command {
	var description, visibility, instance string
	command := &cobra.Command{
		Use: "create NAME", Short: "Create a repository",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 || strings.TrimSpace(args[0]) == "" {
				return newCommandError(categoryValidation, "create repository", fmt.Errorf("repository name is required"))
			}
			if visibility != "public" && visibility != "private" {
				return newCommandError(categoryValidation, "create repository", fmt.Errorf("visibility must be public or private"))
			}
			return nil
		},
		RunE: func(command *cobra.Command, args []string) error {
			if creator == nil {
				service, err := composeRepositoryService(command.Context(), instance)
				if err != nil {
					return err
				}
				var ok bool
				creator, ok = service.(applicationrepository.Creator)
				if !ok {
					return newCommandError(categoryInternal, "create repository", fmt.Errorf("repository service does not support creation"))
				}
			}
			result, err := applicationrepository.NewCreateUseCase(creator).Execute(command.Context(), applicationrepository.CreateRequest{Name: args[0], Description: description, Private: visibility == "private"})
			if err != nil {
				return mapCreateRepositoryError(err)
			}
			printRepository(command, result)
			return nil
		},
	}
	command.Flags().StringVar(&description, "description", "", "repository description")
	command.Flags().StringVar(&visibility, "visibility", "private", "repository visibility (public or private)")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
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
	defaultBranch := repository.DefaultBranch
	if defaultBranch == "" {
		defaultBranch = "-"
	}
	fmt.Fprintf(command.OutOrStdout(), "Repository: %s/%s\nDescription: %s\nPrivate: %t\nArchived: %t\nDefault branch: %s\n", repository.Owner, repository.Name, description, repository.Private, repository.Archived, defaultBranch)
}

func printUpdatedRepository(command *cobra.Command, repository applicationrepository.Repository, fields []string) {
	description := repository.Description
	if description == "" {
		description = "-"
	}
	branch := repository.DefaultBranch
	if branch == "" {
		branch = "-"
	}
	fmt.Fprintf(command.OutOrStdout(), "Repository: %s/%s\nChanged fields: %s\nDescription: %s\nPrivate: %t\nArchived: %t\nDefault branch: %s\n", repository.Owner, repository.Name, strings.Join(fields, ", "), description, repository.Private, repository.Archived, branch)
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
