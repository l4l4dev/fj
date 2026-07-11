package cli

import (
	"fmt"
	"strconv"
	"strings"

	applicationissue "github.com/l4l4dev/fj/internal/application/issue"
	"github.com/spf13/cobra"
)

func newIssueCommand(lister applicationissue.Lister, inspector applicationissue.Inspector, creator applicationissue.Creator, updater applicationissue.Updater) *cobra.Command {
	command := &cobra.Command{Use: "issue", Short: "Manage issues"}
	command.AddCommand(newIssueListCommand(lister))
	command.AddCommand(newIssueInspectCommand(inspector))
	command.AddCommand(newIssueCreateCommand(creator))
	command.AddCommand(newIssueUpdateCommand(updater))
	return command
}

func newIssueUpdateCommand(updater applicationissue.Updater) *cobra.Command {
	var instance, title, body string
	var titleSet, bodySet bool
	command := &cobra.Command{Use: "update OWNER/NAME NUMBER", Short: "Update an issue", Args: func(command *cobra.Command, args []string) error {
		titleSet = command.Flags().Changed("title")
		bodySet = command.Flags().Changed("body")
		if len(args) != 2 {
			return newCommandError(categoryValidation, "update issue", fmt.Errorf("OWNER/NAME and issue number are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "update issue", err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, "update issue", fmt.Errorf("issue number must be a positive integer"))
		}
		if !titleSet && !bodySet {
			return newCommandError(categoryValidation, "update issue", fmt.Errorf("at least one issue field is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if updater == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			updater = dependencies.IssueUpdater
			if updater == nil {
				return newCommandError(categoryInternal, "update issue", fmt.Errorf("issue updater unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		number, _ := strconv.Atoi(args[1])
		request := applicationissue.UpdateRequest{Owner: parts[0], Name: parts[1], Number: number}
		if titleSet {
			request.Title = &title
		}
		if bodySet {
			request.Body = &body
		}
		result, err := applicationissue.NewUpdateUseCase(updater).Execute(command.Context(), request)
		if err != nil {
			return mapApplicationError(err, "update issue")
		}
		fields := make([]string, 0, 2)
		if titleSet {
			fields = append(fields, "title")
		}
		if bodySet {
			fields = append(fields, "body")
		}
		return (issuePresenter{}).PresentUpdated(command.OutOrStdout(), result, fields)
	}}
	command.Flags().StringVar(&title, "title", "", "issue title")
	command.Flags().StringVar(&body, "body", "", "issue body")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newIssueCreateCommand(creator applicationissue.Creator) *cobra.Command {
	var instance, title, body string
	command := &cobra.Command{Use: "create OWNER/NAME", Short: "Create an issue", Args: func(command *cobra.Command, args []string) error {
		if len(args) != 1 {
			return newCommandError(categoryValidation, "create issue", fmt.Errorf("OWNER/NAME is required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "create issue", err)
		}
		if strings.TrimSpace(title) == "" {
			return newCommandError(categoryValidation, "create issue", fmt.Errorf("issue title is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if creator == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			creator = dependencies.IssueCreator
			if creator == nil {
				return newCommandError(categoryInternal, "create issue", fmt.Errorf("issue creator unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		result, err := applicationissue.NewCreateUseCase(creator).Execute(command.Context(), applicationissue.CreateRequest{Owner: parts[0], Name: parts[1], Title: title, Body: body})
		if err != nil {
			return mapApplicationError(err, "create issue")
		}
		return (issuePresenter{}).PresentInspect(command.OutOrStdout(), result)
	}}
	command.Flags().StringVar(&title, "title", "", "issue title")
	command.Flags().StringVar(&body, "body", "", "issue body")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newIssueInspectCommand(inspector applicationissue.Inspector) *cobra.Command {
	var instance string
	command := &cobra.Command{Use: "inspect OWNER/NAME NUMBER", Short: "Inspect an issue", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return newCommandError(categoryValidation, "inspect issue", fmt.Errorf("OWNER/NAME and issue number are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "inspect issue", err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, "inspect issue", fmt.Errorf("issue number must be a positive integer"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if inspector == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			inspector = dependencies.IssueInspector
			if inspector == nil {
				return newCommandError(categoryInternal, "inspect issue", fmt.Errorf("issue inspector unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		number, _ := strconv.Atoi(args[1])
		result, err := applicationissue.NewInspectUseCase(inspector).Execute(command.Context(), applicationissue.InspectRequest{Owner: parts[0], Name: parts[1], Number: number})
		if err != nil {
			return mapApplicationError(err, "inspect issue")
		}
		return (issuePresenter{}).PresentInspect(command.OutOrStdout(), result)
	}}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
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
