package cli

import (
	"fmt"
	"strconv"
	"strings"

	applicationissue "github.com/l4l4dev/fj/internal/application/issue"
	"github.com/spf13/cobra"
)

func newIssueCommand(lister applicationissue.Lister, inspector applicationissue.Inspector, creator applicationissue.Creator, updater applicationissue.Updater, stateChanger applicationissue.StateChanger, commentViewer applicationissue.CommentViewer, commentCreator applicationissue.CommentCreator, labelAdder applicationissue.LabelAdder, labelRemover applicationissue.LabelRemover, milestoneSetter applicationissue.MilestoneSetter, milestoneClearer applicationissue.MilestoneClearer) *cobra.Command {
	command := &cobra.Command{Use: "issue", Short: "Manage issues"}
	command.AddCommand(newIssueListCommand(lister))
	command.AddCommand(newIssueInspectCommand(inspector))
	command.AddCommand(newIssueCreateCommand(creator))
	command.AddCommand(newIssueUpdateCommand(updater))
	command.AddCommand(newIssueStateCommand(stateChanger))
	command.AddCommand(newIssueCommentCommand(commentViewer, commentCreator))
	command.AddCommand(newIssueLabelCommand(labelAdder, labelRemover))
	command.AddCommand(newIssueMilestoneCommand(milestoneSetter, milestoneClearer))
	return command
}

func newIssueMilestoneCommand(setter applicationissue.MilestoneSetter, clearer applicationissue.MilestoneClearer) *cobra.Command {
	command := &cobra.Command{Use: "milestone", Short: "Manage issue milestones"}
	command.AddCommand(newIssueMilestoneSetCommand(setter))
	command.AddCommand(newIssueMilestoneClearCommand(clearer))
	return command
}

func newIssueMilestoneSetCommand(setter applicationissue.MilestoneSetter) *cobra.Command {
	var instance string
	command := &cobra.Command{Use: "set OWNER/NAME NUMBER MILESTONE", Short: "Set issue milestone", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 3 {
			return newCommandError(categoryValidation, "set issue milestone", fmt.Errorf("OWNER/NAME, issue number, and milestone are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "set issue milestone", err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, "set issue milestone", fmt.Errorf("issue number must be a positive integer"))
		}
		if strings.TrimSpace(args[2]) == "" {
			return newCommandError(categoryValidation, "set issue milestone", fmt.Errorf("milestone is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if setter == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			setter = dependencies.MilestoneSetter
			if setter == nil {
				return newCommandError(categoryInternal, "set issue milestone", fmt.Errorf("milestone setter unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		number, _ := strconv.Atoi(args[1])
		milestone, err := applicationissue.NewSetMilestoneUseCase(setter).Execute(command.Context(), applicationissue.SetMilestoneRequest{Owner: parts[0], Name: parts[1], Number: number, Milestone: args[2]})
		if err != nil {
			return mapApplicationError(err, "set issue milestone")
		}
		return (issuePresenter{}).PresentMilestoneSet(command.OutOrStdout(), number, milestone)
	}}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newIssueMilestoneClearCommand(clearer applicationissue.MilestoneClearer) *cobra.Command {
	var instance string
	command := &cobra.Command{Use: "clear OWNER/NAME NUMBER", Short: "Clear issue milestone", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return newCommandError(categoryValidation, "clear issue milestone", fmt.Errorf("OWNER/NAME and issue number are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "clear issue milestone", err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, "clear issue milestone", fmt.Errorf("issue number must be a positive integer"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if clearer == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			clearer = dependencies.MilestoneClearer
			if clearer == nil {
				return newCommandError(categoryInternal, "clear issue milestone", fmt.Errorf("milestone clearer unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		number, _ := strconv.Atoi(args[1])
		if err := applicationissue.NewClearMilestoneUseCase(clearer).Execute(command.Context(), applicationissue.ClearMilestoneRequest{Owner: parts[0], Name: parts[1], Number: number}); err != nil {
			return mapApplicationError(err, "clear issue milestone")
		}
		return (issuePresenter{}).PresentMilestoneCleared(command.OutOrStdout(), number)
	}}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newIssueLabelCommand(adder applicationissue.LabelAdder, remover applicationissue.LabelRemover) *cobra.Command {
	command := &cobra.Command{Use: "label", Short: "Manage issue labels"}
	command.AddCommand(newIssueLabelAddCommand(adder))
	command.AddCommand(newIssueLabelRemoveCommand(remover))
	return command
}

func newIssueLabelAddCommand(adder applicationissue.LabelAdder) *cobra.Command {
	return newIssueLabelMutationCommand("add", "add issue label", adder, nil)
}

func newIssueLabelRemoveCommand(remover applicationissue.LabelRemover) *cobra.Command {
	return newIssueLabelMutationCommand("remove", "remove issue label", nil, remover)
}

func newIssueLabelMutationCommand(name, operation string, adder applicationissue.LabelAdder, remover applicationissue.LabelRemover) *cobra.Command {
	var instance string
	command := &cobra.Command{Use: name + " OWNER/NAME NUMBER LABEL", Short: operation, Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 3 {
			return newCommandError(categoryValidation, operation, fmt.Errorf("OWNER/NAME, issue number, and label are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, operation, err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, operation, fmt.Errorf("issue number must be a positive integer"))
		}
		if strings.TrimSpace(args[2]) == "" {
			return newCommandError(categoryValidation, operation, fmt.Errorf("label is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		number, _ := strconv.Atoi(args[1])
		if name == "add" {
			if adder == nil {
				dependencies, err := composeRepositoryDependencies(command.Context(), instance)
				if err != nil {
					return err
				}
				adder = dependencies.LabelAdder
				if adder == nil {
					return newCommandError(categoryInternal, operation, fmt.Errorf("label adder unavailable"))
				}
			}
			parts := strings.SplitN(args[0], "/", 2)
			label, err := applicationissue.NewAddLabelUseCase(adder).Execute(command.Context(), applicationissue.AddLabelRequest{Owner: parts[0], Name: parts[1], Number: number, Label: args[2]})
			if err != nil {
				return mapApplicationError(err, operation)
			}
			return (issuePresenter{}).PresentLabelAdded(command.OutOrStdout(), number, label)
		}
		if remover == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			remover = dependencies.LabelRemover
			if remover == nil {
				return newCommandError(categoryInternal, operation, fmt.Errorf("label remover unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		label, err := applicationissue.NewRemoveLabelUseCase(remover).Execute(command.Context(), applicationissue.RemoveLabelRequest{Owner: parts[0], Name: parts[1], Number: number, Label: args[2]})
		if err != nil {
			return mapApplicationError(err, operation)
		}
		return (issuePresenter{}).PresentLabelRemoved(command.OutOrStdout(), number, label)
	}}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newIssueCommentCommand(viewer applicationissue.CommentViewer, creator applicationissue.CommentCreator) *cobra.Command {
	command := &cobra.Command{Use: "comment", Short: "Manage issue comments"}
	command.AddCommand(newIssueCommentListCommand(viewer))
	command.AddCommand(newIssueCommentAddCommand(creator))
	return command
}

func newIssueCommentListCommand(viewer applicationissue.CommentViewer) *cobra.Command {
	var instance string
	command := &cobra.Command{Use: "list OWNER/NAME NUMBER", Short: "List issue comments", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return newCommandError(categoryValidation, "list issue comments", fmt.Errorf("OWNER/NAME and issue number are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "list issue comments", err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, "list issue comments", fmt.Errorf("issue number must be a positive integer"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if viewer == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			viewer = dependencies.CommentViewer
			if viewer == nil {
				return newCommandError(categoryInternal, "list issue comments", fmt.Errorf("comment viewer unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		number, _ := strconv.Atoi(args[1])
		comments, err := applicationissue.NewListCommentsUseCase(viewer).Execute(command.Context(), applicationissue.ListCommentsRequest{Owner: parts[0], Name: parts[1], Number: number})
		if err != nil {
			return mapApplicationError(err, "list issue comments")
		}
		return (issuePresenter{}).PresentComments(command.OutOrStdout(), comments)
	}}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newIssueCommentAddCommand(creator applicationissue.CommentCreator) *cobra.Command {
	var instance, body string
	command := &cobra.Command{Use: "add OWNER/NAME NUMBER", Short: "Add an issue comment", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return newCommandError(categoryValidation, "add issue comment", fmt.Errorf("OWNER/NAME and issue number are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "add issue comment", err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, "add issue comment", fmt.Errorf("issue number must be a positive integer"))
		}
		if strings.TrimSpace(body) == "" {
			return newCommandError(categoryValidation, "add issue comment", fmt.Errorf("comment body is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if creator == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			creator = dependencies.CommentCreator
			if creator == nil {
				return newCommandError(categoryInternal, "add issue comment", fmt.Errorf("comment creator unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		number, _ := strconv.Atoi(args[1])
		comment, err := applicationissue.NewAddCommentUseCase(creator).Execute(command.Context(), applicationissue.AddCommentRequest{Owner: parts[0], Name: parts[1], Number: number, Body: body})
		if err != nil {
			return mapApplicationError(err, "add issue comment")
		}
		return (issuePresenter{}).PresentComment(command.OutOrStdout(), comment)
	}}
	command.Flags().StringVar(&body, "body", "", "comment body")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newIssueStateCommand(changer applicationissue.StateChanger) *cobra.Command {
	var instance, state string
	command := &cobra.Command{Use: "state OWNER/NAME NUMBER", Short: "Change issue state", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return newCommandError(categoryValidation, "change issue state", fmt.Errorf("OWNER/NAME and issue number are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "change issue state", err)
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return newCommandError(categoryValidation, "change issue state", fmt.Errorf("issue number must be a positive integer"))
		}
		if state != string(applicationissue.StateOpen) && state != string(applicationissue.StateClosed) {
			return newCommandError(categoryValidation, "change issue state", fmt.Errorf("state must be open or closed"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if changer == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			changer = dependencies.IssueStateChanger
			if changer == nil {
				return newCommandError(categoryInternal, "change issue state", fmt.Errorf("issue state changer unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		number, _ := strconv.Atoi(args[1])
		result, err := applicationissue.NewChangeStateUseCase(changer).Execute(command.Context(), applicationissue.ChangeStateRequest{Owner: parts[0], Name: parts[1], Number: number, State: applicationissue.State(state)})
		if err != nil {
			return mapApplicationError(err, "change issue state")
		}
		return (issuePresenter{}).PresentState(command.OutOrStdout(), result)
	}}
	command.Flags().StringVar(&state, "state", "", "issue state (open or closed)")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
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
