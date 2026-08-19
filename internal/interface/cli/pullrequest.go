package cli

import (
	"fmt"
	"strconv"
	"strings"

	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
	"github.com/spf13/cobra"
)

type pullRequestDependencies struct {
	lister          applicationpullrequest.PullRequestLister
	inspector       applicationpullrequest.PullRequestInspector
	creator         applicationpullrequest.PullRequestCreator
	updater         applicationpullrequest.Updater
	statusViewer    applicationpullrequest.StatusViewer
	reviewSubmitter applicationpullrequest.ReviewSubmitter
}

func newPullRequestCommand(dependencies pullRequestDependencies) *cobra.Command {
	command := &cobra.Command{Use: "pr", Short: "Manage pull requests"}
	command.AddCommand(newPullRequestListCommand(dependencies.lister))
	command.AddCommand(newPullRequestInspectCommand(dependencies.inspector))
	command.AddCommand(newPullRequestCreateCommand(dependencies.creator))
	command.AddCommand(newPullRequestUpdateCommand(dependencies.updater))
	command.AddCommand(newPullRequestStatusCommand(dependencies.statusViewer))
	command.AddCommand(newPullRequestReviewCommand(dependencies.reviewSubmitter))
	return command
}

func newPullRequestUpdateCommand(updater applicationpullrequest.Updater) *cobra.Command {
	var instance, title, body string
	var titleSet, bodySet bool
	command := &cobra.Command{Use: "update OWNER/NAME NUMBER", Short: "Update a pull request", Args: func(command *cobra.Command, args []string) error {
		titleSet = command.Flags().Changed("title")
		bodySet = command.Flags().Changed("body")
		if len(args) != 2 {
			return newCommandError(categoryValidation, "update pull request", fmt.Errorf("OWNER/NAME and pull request number are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "update pull request", err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, "update pull request", fmt.Errorf("pull request number must be a positive integer"))
		}
		if !titleSet && !bodySet {
			return newCommandError(categoryValidation, "update pull request", fmt.Errorf("at least one pull request field is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if updater == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			updater = dependencies.PullRequestUpdater
			if updater == nil {
				return newCommandError(categoryInternal, "update pull request", fmt.Errorf("pull request updater unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		number, _ := strconv.Atoi(args[1])
		request := applicationpullrequest.UpdateRequest{Owner: parts[0], Name: parts[1], Number: number}
		fields := make([]string, 0, 2)
		if titleSet {
			request.Title = &title
			fields = append(fields, "title")
		}
		if bodySet {
			request.Body = &body
			fields = append(fields, "body")
		}
		result, err := applicationpullrequest.NewUpdateUseCase(updater).Execute(command.Context(), request)
		if err != nil {
			return mapApplicationError(err, "update pull request")
		}
		return (pullRequestPresenter{}).PresentUpdated(command.OutOrStdout(), result, fields)
	}}
	command.Flags().StringVar(&title, "title", "", "pull request title")
	command.Flags().StringVar(&body, "body", "", "pull request body")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newPullRequestReviewCommand(submitter applicationpullrequest.ReviewSubmitter) *cobra.Command {
	var instance, outcome, body string
	command := &cobra.Command{Use: "review OWNER/NAME NUMBER", Short: "Submit a pull request review", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return newCommandError(categoryValidation, "submit pull request review", fmt.Errorf("OWNER/NAME and pull request number are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "submit pull request review", err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, "submit pull request review", fmt.Errorf("pull request number must be a positive integer"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if submitter == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			submitter = dependencies.PullRequestReviewSubmitter
			if submitter == nil {
				return newCommandError(categoryInternal, "submit pull request review", fmt.Errorf("pull request review submitter unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		number, _ := strconv.Atoi(args[1])
		result, err := applicationpullrequest.NewSubmitReviewUseCase(submitter).Execute(command.Context(), applicationpullrequest.SubmitReviewRequest{Owner: parts[0], Name: parts[1], Number: number, Outcome: applicationpullrequest.ReviewOutcome(outcome), Body: body})
		if err != nil {
			return mapApplicationError(err, "submit pull request review")
		}
		return (pullRequestPresenter{}).PresentSubmittedReview(command.OutOrStdout(), result)
	}}
	command.Flags().StringVar(&outcome, "outcome", "", "review outcome (comment, approve, or request-changes)")
	command.Flags().StringVar(&body, "body", "", "review body")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
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
