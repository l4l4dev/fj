package cli

import (
	"fmt"
	"strconv"
	"strings"

	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
	"github.com/spf13/cobra"
)

func newPullRequestCommand(lister applicationpullrequest.PullRequestLister) *cobra.Command {
	return newPullRequestCommandWithInspector(lister, nil)
}

func newPullRequestCommandWithInspector(lister applicationpullrequest.PullRequestLister, inspector applicationpullrequest.PullRequestInspector) *cobra.Command {
	return newPullRequestCommandWithDependencies(lister, inspector, nil, nil)
}

func newPullRequestCommandWithDependencies(lister applicationpullrequest.PullRequestLister, inspector applicationpullrequest.PullRequestInspector, creator applicationpullrequest.PullRequestCreator, statusViewer applicationpullrequest.StatusViewer) *cobra.Command {
	command := &cobra.Command{Use: "pr", Short: "Manage pull requests"}
	command.AddCommand(newPullRequestListCommand(lister))
	command.AddCommand(newPullRequestInspectCommand(inspector))
	command.AddCommand(newPullRequestCreateCommand(creator))
	command.AddCommand(newPullRequestStatusCommand(statusViewer))
	return command
}

func newPullRequestStatusCommand(viewer applicationpullrequest.StatusViewer) *cobra.Command {
	var instance string
	command := &cobra.Command{Use: "status OWNER/NAME NUMBER", Short: "View pull request review and check status", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return newCommandError(categoryValidation, "view pull request status", fmt.Errorf("OWNER/NAME and pull request number are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "view pull request status", err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, "view pull request status", fmt.Errorf("pull request number must be a positive integer"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if viewer == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			viewer = dependencies.PullRequestStatusViewer
			if viewer == nil {
				return newCommandError(categoryInternal, "view pull request status", fmt.Errorf("pull request status viewer unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		number, _ := strconv.Atoi(args[1])
		result, err := applicationpullrequest.NewStatusUseCase(viewer).Execute(command.Context(), applicationpullrequest.StatusRequest{Owner: parts[0], Name: parts[1], Number: number})
		if err != nil {
			return mapApplicationError(err, "view pull request status")
		}
		return (pullRequestPresenter{}).PresentStatus(command.OutOrStdout(), result)
	}}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newPullRequestCreateCommand(creator applicationpullrequest.PullRequestCreator) *cobra.Command {
	var instance, title, head, base string
	command := &cobra.Command{Use: "create OWNER/NAME", Short: "Create a pull request", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return newCommandError(categoryValidation, "create pull request", fmt.Errorf("OWNER/NAME is required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "create pull request", err)
		}
		if strings.TrimSpace(title) == "" {
			return newCommandError(categoryValidation, "create pull request", fmt.Errorf("pull request title is required"))
		}
		if strings.TrimSpace(head) == "" {
			return newCommandError(categoryValidation, "create pull request", fmt.Errorf("head branch is required"))
		}
		if strings.Contains(head, ":") {
			return newCommandError(categoryValidation, "create pull request", fmt.Errorf("cross-fork head branches are not supported"))
		}
		if strings.TrimSpace(base) == "" {
			return newCommandError(categoryValidation, "create pull request", fmt.Errorf("base branch is required"))
		}
		if strings.Contains(base, ":") {
			return newCommandError(categoryValidation, "create pull request", fmt.Errorf("base branch must belong to the target repository"))
		}
		if head == base {
			return newCommandError(categoryValidation, "create pull request", fmt.Errorf("head and base branches must differ"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if creator == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			creator = dependencies.PullRequestCreator
			if creator == nil {
				return newCommandError(categoryInternal, "create pull request", fmt.Errorf("pull request creator unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		result, err := applicationpullrequest.NewCreateUseCase(creator).Execute(command.Context(), applicationpullrequest.CreateRequest{Owner: parts[0], Name: parts[1], Title: title, HeadBranch: head, BaseBranch: base})
		if err != nil {
			return mapApplicationError(err, "create pull request")
		}
		return (pullRequestPresenter{}).PresentCreated(command.OutOrStdout(), result)
	}}
	command.Flags().StringVar(&title, "title", "", "pull request title")
	command.Flags().StringVar(&head, "head", "", "source branch")
	command.Flags().StringVar(&base, "base", "", "target branch")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newPullRequestInspectCommand(inspector applicationpullrequest.PullRequestInspector) *cobra.Command {
	var instance string
	command := &cobra.Command{Use: "inspect OWNER/NAME NUMBER", Short: "Inspect a pull request", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return newCommandError(categoryValidation, "inspect pull request", fmt.Errorf("OWNER/NAME and pull request number are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "inspect pull request", err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, "inspect pull request", fmt.Errorf("pull request number must be a positive integer"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if inspector == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			inspector = dependencies.PullRequestInspector
			if inspector == nil {
				return newCommandError(categoryInternal, "inspect pull request", fmt.Errorf("pull request inspector unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		number, _ := strconv.Atoi(args[1])
		result, err := applicationpullrequest.NewInspectUseCase(inspector).Execute(command.Context(), applicationpullrequest.InspectRequest{Owner: parts[0], Name: parts[1], Number: number})
		if err != nil {
			return mapApplicationError(err, "inspect pull request")
		}
		return (pullRequestPresenter{}).PresentInspect(command.OutOrStdout(), result)
	}}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
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
