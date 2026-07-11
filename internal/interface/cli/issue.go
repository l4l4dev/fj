package cli

import (
	"fmt"
	"strings"

	applicationissue "github.com/l4l4dev/fj/internal/application/issue"
	"github.com/spf13/cobra"
)

func newIssueCommand(lister applicationissue.Lister) *cobra.Command {
	command := &cobra.Command{Use: "issue", Short: "Manage issues"}
	command.AddCommand(newIssueListCommand(lister))
	return command
}

func newIssueListCommand(lister applicationissue.Lister) *cobra.Command {
	var instance, state string
	var assignees, labels []string
	var page, limit int
	command := &cobra.Command{Use: "list OWNER/NAME", Short: "List issues", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return newCommandError(categoryValidation, "list issues", fmt.Errorf("OWNER/NAME is required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "list issues", err)
		}
		if state != "open" && state != "closed" && state != "all" {
			return newCommandError(categoryValidation, "list issues", fmt.Errorf("state must be open, closed, or all"))
		}
		if len(assignees) > 1 || len(labels) > 1 {
			return newCommandError(categoryValidation, "list issues", fmt.Errorf("repeated issue filters are not supported"))
		}
		if len(assignees) == 1 && strings.TrimSpace(assignees[0]) == "" {
			return newCommandError(categoryValidation, "list issues", fmt.Errorf("assignee must not be empty"))
		}
		if len(labels) == 1 && strings.TrimSpace(labels[0]) == "" {
			return newCommandError(categoryValidation, "list issues", fmt.Errorf("label must not be empty"))
		}
		if len(assignees) == 1 && len(labels) == 1 {
			return newCommandError(categoryValidation, "list issues", fmt.Errorf("only one issue filter may be specified"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if page < 1 {
			return newCommandError(categoryValidation, "list issues", fmt.Errorf("page must be at least 1"))
		}
		if limit < 1 {
			return newCommandError(categoryValidation, "list issues", fmt.Errorf("limit must be at least 1"))
		}
		if lister == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			lister = dependencies.Issues
			if lister == nil {
				return newCommandError(categoryInternal, "list issues", fmt.Errorf("issue lister unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		filter := applicationissue.IssueFilter{}
		if len(assignees) == 1 {
			filter.Assignee = assignees[0]
		}
		if len(labels) == 1 {
			filter.Label = labels[0]
		}
		result, err := applicationissue.NewListUseCase(lister).Execute(command.Context(), applicationissue.ListRequest{Owner: parts[0], Name: parts[1], Page: page, Limit: limit, State: applicationissue.State(state), Filter: filter})
		if err != nil {
			return mapApplicationError(err, "list issues")
		}
		return (issuePresenter{}).PresentList(command.OutOrStdout(), result)
	}}
	command.Flags().IntVar(&page, "page", 1, "page number")
	command.Flags().IntVar(&limit, "limit", 30, "page size")
	command.Flags().StringVar(&state, "state", "open", "issue state (open, closed, or all)")
	command.Flags().StringArrayVar(&assignees, "assignee", nil, "filter by assignee")
	command.Flags().StringArrayVar(&labels, "label", nil, "filter by label")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}
