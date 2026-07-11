package cli

import (
	"context"
	"errors"

	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationauth "github.com/l4l4dev/fj/internal/application/auth"
	applicationissue "github.com/l4l4dev/fj/internal/application/issue"
	applicationrepository "github.com/l4l4dev/fj/internal/application/repository"
	infrastructureauth "github.com/l4l4dev/fj/internal/infrastructure/auth"
	infrastructureconfig "github.com/l4l4dev/fj/internal/infrastructure/config"
	"github.com/l4l4dev/fj/internal/infrastructure/forgejo"
	infrastructureissue "github.com/l4l4dev/fj/internal/infrastructure/issue"
	infrastructurerepository "github.com/l4l4dev/fj/internal/infrastructure/repository"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command { return NewRootCommandWithDependencies(RepositoryDependencies{}) }

type RepositoryDependencies struct {
	List              applicationrepository.Service
	Inspect           applicationrepository.Getter
	Create            applicationrepository.Creator
	Update            applicationrepository.Updater
	Archive           applicationrepository.Archiver
	Access            applicationrepository.AccessViewer
	Issues            applicationissue.Lister
	IssueInspector    applicationissue.Inspector
	IssueCreator      applicationissue.Creator
	IssueUpdater      applicationissue.Updater
	IssueStateChanger applicationissue.StateChanger
	CommentViewer     applicationissue.CommentViewer
	CommentCreator    applicationissue.CommentCreator
}

func NewRootCommandWithDependencies(dependencies RepositoryDependencies) *cobra.Command {
	command := &cobra.Command{Use: "fj", Short: "AI-first CLI for Forgejo", Args: cobra.NoArgs, SilenceErrors: true, SilenceUsage: true, RunE: func(command *cobra.Command, _ []string) error { return command.Help() }}
	command.AddCommand(newRepositoryCommand(dependencies))
	command.AddCommand(newIssueCommand(dependencies.Issues, dependencies.IssueInspector, dependencies.IssueCreator, dependencies.IssueUpdater, dependencies.IssueStateChanger, dependencies.CommentViewer, dependencies.CommentCreator))
	command.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return newCommandError(categoryValidation, "execute command", err)
	})
	return command
}

func NewRootCommandWithRepositoryService(service applicationrepository.Service) *cobra.Command {
	return NewRootCommandWithDependencies(legacyRepositoryDependencies(service))
}

func legacyRepositoryDependencies(service applicationrepository.Service) RepositoryDependencies {
	dependencies := RepositoryDependencies{List: service}
	if value, ok := service.(applicationrepository.Getter); ok {
		dependencies.Inspect = value
	}
	if value, ok := service.(applicationrepository.Creator); ok {
		dependencies.Create = value
	}
	if value, ok := service.(applicationrepository.Updater); ok {
		dependencies.Update = value
	}
	if value, ok := service.(applicationrepository.Archiver); ok {
		dependencies.Archive = value
	}
	if value, ok := service.(applicationrepository.AccessViewer); ok {
		dependencies.Access = value
	}
	return dependencies
}

func composeRepositoryDependencies(ctx context.Context, instanceName string) (RepositoryDependencies, error) {
	configuration, err := infrastructureconfig.Load()
	if err != nil {
		return RepositoryDependencies{}, newCommandError(categoryValidation, "load configuration", err)
	}
	instance, err := configuration.SelectInstance(instanceName)
	if err != nil {
		return RepositoryDependencies{}, newCommandError(categoryValidation, "select instance", err)
	}
	credential, err := applicationauth.NewResolver(infrastructureauth.NewEnvironmentProvider()).Resolve(ctx, instance.Credential)
	if err != nil {
		return RepositoryDependencies{}, newCommandError(categoryAuthentication, "resolve credential", err)
	}
	adapter := infrastructurerepository.NewRESTAdapter(forgejo.NewClient(instance, credential, "dev", nil))
	issueAdapter := infrastructureissue.NewRESTAdapter(forgejo.NewClient(instance, credential, "dev", nil))
	return RepositoryDependencies{List: adapter, Inspect: adapter, Create: adapter, Update: adapter, Archive: adapter, Access: adapter, Issues: issueAdapter, IssueInspector: issueAdapter, IssueCreator: issueAdapter, IssueUpdater: issueAdapter, IssueStateChanger: issueAdapter, CommentViewer: issueAdapter, CommentCreator: issueAdapter}, nil
}

func mapApplicationError(err error, operation string) error {
	if err == nil {
		return nil
	}
	var appErr apperror.Error
	if errors.As(err, &appErr) {
		category := categoryInternal
		switch appErr.Category {
		case apperror.Validation:
			category = categoryValidation
		case apperror.Authentication:
			category = categoryAuthentication
		case apperror.NotFound, apperror.Conflict, apperror.Remote:
			category = categoryRemote
		}
		message := appErr.Message
		if message == "" {
			return newCommandError(category, operation, err)
		}
		return newCommandErrorWithMessage(category, operation, message, err)
	}
	var remote applicationrepository.RemoteError
	if errors.As(err, &remote) {
		category := categoryRemote
		message := ""
		switch remote.Category() {
		case apperror.Authentication:
			category = categoryAuthentication
		case apperror.NotFound:
			category = categoryRemote
			message = "repository not found"
		case apperror.Conflict:
			category = categoryRemote
			if operation == "create repository" {
				message = "repository already exists"
			} else if operation == "update repository" {
				message = "repository update conflict"
			}
		}
		if message != "" {
			return newCommandErrorWithMessage(category, operation, message, err)
		}
		return newCommandError(category, operation, err)
	}
	var validation apperror.ValidationError
	if errors.As(err, &validation) {
		return newCommandError(categoryValidation, operation, err)
	}
	return newCommandError(categoryInternal, operation, err)
}
