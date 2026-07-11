package cli

import (
	"fmt"
	"strings"

	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
	"github.com/spf13/cobra"
)

func newPullRequestCommand(lister applicationpullrequest.PullRequestLister) *cobra.Command {
	command := &cobra.Command{Use: "pr", Short: "Manage pull requests"}
	command.AddCommand(newPullRequestListCommand(lister))
	return command
}

func newPullRequestListCommand(lister applicationpullrequest.PullRequestLister) *cobra.Command {
	var instance, state string
	var page, limit int
	command := &cobra.Command{Use: "list OWNER/NAME", Short: "List pull requests", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return newCommandError(categoryValidation, "list pull requests", fmt.Errorf("OWNER/NAME is required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "list pull requests", err)
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if page < 1 {
			return newCommandError(categoryValidation, "list pull requests", fmt.Errorf("page must be at least 1"))
		}
		if limit < 1 || limit > 100 {
			return newCommandError(categoryValidation, "list pull requests", fmt.Errorf("limit must be between 1 and 100"))
		}
		if state != "open" && state != "closed" && state != "all" {
			return newCommandError(categoryValidation, "list pull requests", fmt.Errorf("state must be open, closed, or all"))
		}
		if lister == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			lister = dependencies.PullRequests
			if lister == nil {
				return newCommandError(categoryInternal, "list pull requests", fmt.Errorf("pull request lister unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		result, err := applicationpullrequest.NewListUseCase(lister).Execute(command.Context(), applicationpullrequest.ListRequest{Owner: parts[0], Name: parts[1], Page: page, Limit: limit, State: applicationpullrequest.State(state)})
		if err != nil {
			return mapApplicationError(err, "list pull requests")
		}
		return (pullRequestPresenter{}).Present(command.OutOrStdout(), result)
	}}
	command.Flags().IntVar(&page, "page", 1, "page number")
	command.Flags().IntVar(&limit, "limit", 20, "page size")
	command.Flags().StringVar(&state, "state", "open", "pull request state (open, closed, or all)")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}
