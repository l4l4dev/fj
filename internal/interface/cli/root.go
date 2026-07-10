package cli

import (
	"context"
	"errors"

	applicationauth "github.com/l4l4dev/fj/internal/application/auth"
	applicationrepository "github.com/l4l4dev/fj/internal/application/repository"
	infrastructureauth "github.com/l4l4dev/fj/internal/infrastructure/auth"
	infrastructureconfig "github.com/l4l4dev/fj/internal/infrastructure/config"
	"github.com/l4l4dev/fj/internal/infrastructure/forgejo"
	infrastructurerepository "github.com/l4l4dev/fj/internal/infrastructure/repository"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	return NewRootCommandWithRepositoryService(nil)
}

func NewRootCommandWithRepositoryService(service applicationrepository.Service) *cobra.Command {
	command := &cobra.Command{
		Use:           "fj",
		Short:         "AI-first CLI for Forgejo",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(command *cobra.Command, _ []string) error {
			return command.Help()
		},
	}
	command.AddCommand(newRepositoryCommand(service))
	command.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return newCommandError(categoryValidation, "execute command", err)
	})
	return command
}

func composeRepositoryService(ctx context.Context, instanceName string) (applicationrepository.Service, error) {
	configuration, err := infrastructureconfig.Load()
	if err != nil {
		return nil, newCommandError(categoryValidation, "load configuration", err)
	}
	instance, err := configuration.SelectInstance(instanceName)
	if err != nil {
		return nil, newCommandError(categoryValidation, "select instance", err)
	}
	credential, err := applicationauth.NewResolver(infrastructureauth.NewEnvironmentProvider()).Resolve(ctx, instance.Credential)
	if err != nil {
		return nil, newCommandError(categoryAuthentication, "resolve credential", err)
	}
	client := forgejo.NewClient(instance, credential, "dev", nil)
	return infrastructurerepository.NewRESTAdapter(client), nil
}

func mapRepositoryError(err error) error {
	if err == nil {
		return nil
	}
	var remoteError applicationrepository.RemoteError
	if errors.As(err, &remoteError) {
		if remoteError.StatusCode() == 401 || remoteError.StatusCode() == 403 {
			return newCommandError(categoryAuthentication, "list repositories", err)
		}
		return newCommandError(categoryRemote, "list repositories", err)
	}
	return newCommandError(categoryValidation, "list repositories", err)
}

func mapInspectRepositoryError(err error) error {
	if err == nil {
		return nil
	}
	var remoteError applicationrepository.RemoteError
	if errors.As(err, &remoteError) {
		if remoteError.StatusCode() == 401 || remoteError.StatusCode() == 403 {
			return newCommandError(categoryAuthentication, "inspect repository", err)
		}
		if remoteError.StatusCode() == 404 {
			return newCommandErrorWithMessage(categoryRemote, "inspect repository", "repository not found", err)
		}
		return newCommandError(categoryRemote, "inspect repository", err)
	}
	return newCommandError(categoryValidation, "inspect repository", err)
}
