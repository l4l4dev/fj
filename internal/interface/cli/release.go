package cli

import (
	"fmt"
	"strings"

	applicationrelease "github.com/l4l4dev/fj/internal/application/release"
	"github.com/spf13/cobra"
)

type releaseDependencies struct {
	lister    applicationrelease.Lister
	inspector applicationrelease.Inspector
}

func newReleaseCommand(dependencies releaseDependencies) *cobra.Command {
	command := &cobra.Command{Use: "release", Short: "Manage releases"}
	command.AddCommand(newReleaseListCommand(dependencies.lister))
	command.AddCommand(newReleaseInspectCommand(dependencies.inspector))
	return command
}

func newReleaseListCommand(lister applicationrelease.Lister) *cobra.Command {
	var instance string
	var page, limit int
	command := &cobra.Command{Use: "list OWNER/NAME", Short: "List releases", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return newCommandError(categoryValidation, "list releases", fmt.Errorf("OWNER/NAME is required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "list releases", err)
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if page < 1 {
			return newCommandError(categoryValidation, "list releases", fmt.Errorf("page must be at least 1"))
		}
		if limit < 1 || limit > 100 {
			return newCommandError(categoryValidation, "list releases", fmt.Errorf("limit must be between 1 and 100"))
		}
		if lister == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			lister = dependencies.Releases
			if lister == nil {
				return newCommandError(categoryInternal, "list releases", fmt.Errorf("release lister unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		result, err := applicationrelease.NewListUseCase(lister).Execute(command.Context(), applicationrelease.ListRequest{Owner: parts[0], Name: parts[1], Page: page, Limit: limit})
		if err != nil {
			return mapApplicationError(err, "list releases")
		}
		return (releasePresenter{}).Present(command.OutOrStdout(), result)
	}}
	command.Flags().IntVar(&page, "page", 1, "page number")
	command.Flags().IntVar(&limit, "limit", 20, "page size")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newReleaseInspectCommand(inspector applicationrelease.Inspector) *cobra.Command {
	var instance string
	command := &cobra.Command{Use: "inspect OWNER/NAME TAG", Short: "Inspect a release", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return newCommandError(categoryValidation, "inspect release", fmt.Errorf("OWNER/NAME and release tag are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "inspect release", err)
		}
		if strings.TrimSpace(args[1]) == "" {
			return newCommandError(categoryValidation, "inspect release", fmt.Errorf("release tag is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if inspector == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			inspector = dependencies.ReleaseInspector
			if inspector == nil {
				return newCommandError(categoryInternal, "inspect release", fmt.Errorf("release inspector unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		result, err := applicationrelease.NewInspectUseCase(inspector).Execute(command.Context(), applicationrelease.InspectRequest{Owner: parts[0], Name: parts[1], Tag: args[1]})
		if err != nil {
			return mapApplicationError(err, "inspect release")
		}
		return (releasePresenter{}).PresentInspect(command.OutOrStdout(), result)
	}}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}
