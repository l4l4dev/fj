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
	creator   applicationrelease.Creator
	updater   applicationrelease.Updater
	publisher applicationrelease.Publisher
}

func newReleaseCommand(dependencies releaseDependencies) *cobra.Command {
	command := &cobra.Command{Use: "release", Short: "Manage releases"}
	command.AddCommand(newReleaseListCommand(dependencies.lister))
	command.AddCommand(newReleaseInspectCommand(dependencies.inspector))
	command.AddCommand(newReleaseCreateCommand(dependencies.creator))
	command.AddCommand(newReleaseUpdateCommand(dependencies.updater))
	command.AddCommand(newReleasePublishCommand(dependencies.inspector, dependencies.publisher))
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

func newReleasePublishCommand(inspector applicationrelease.Inspector, publisher applicationrelease.Publisher) *cobra.Command {
	var instance string
	command := &cobra.Command{Use: "publish OWNER/NAME TAG", Short: "Publish a draft release", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return newCommandError(categoryValidation, "publish release", fmt.Errorf("OWNER/NAME and release tag are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "publish release", err)
		}
		if strings.TrimSpace(args[1]) == "" {
			return newCommandError(categoryValidation, "publish release", fmt.Errorf("release tag is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if inspector == nil || publisher == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			if inspector == nil {
				inspector = dependencies.ReleaseInspector
			}
			if publisher == nil {
				publisher = dependencies.ReleasePublisher
			}
			if inspector == nil || publisher == nil {
				return newCommandError(categoryInternal, "publish release", fmt.Errorf("release publisher unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		result, err := applicationrelease.NewPublishUseCase(inspector, publisher).Execute(command.Context(), applicationrelease.InspectRequest{Owner: parts[0], Name: parts[1], Tag: args[1]})
		if err != nil {
			return mapApplicationError(err, "publish release")
		}
		return (releasePresenter{}).PresentPublished(command.OutOrStdout(), result)
	}}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newReleaseUpdateCommand(updater applicationrelease.Updater) *cobra.Command {
	var instance, title, notes string
	var prerelease bool
	var titleSet, notesSet, prereleaseSet bool
	command := &cobra.Command{Use: "update OWNER/NAME TAG", Short: "Update release metadata", Args: func(command *cobra.Command, args []string) error {
		titleSet = command.Flags().Changed("title")
		notesSet = command.Flags().Changed("notes")
		prereleaseSet = command.Flags().Changed("prerelease")
		if len(args) != 2 {
			return newCommandError(categoryValidation, "update release", fmt.Errorf("OWNER/NAME and release tag are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "update release", err)
		}
		if strings.TrimSpace(args[1]) == "" {
			return newCommandError(categoryValidation, "update release", fmt.Errorf("release tag is required"))
		}
		if !titleSet && !notesSet && !prereleaseSet {
			return newCommandError(categoryValidation, "update release", fmt.Errorf("at least one release field is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if updater == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			updater = dependencies.ReleaseUpdater
			if updater == nil {
				return newCommandError(categoryInternal, "update release", fmt.Errorf("release updater unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		request := applicationrelease.UpdateRequest{Owner: parts[0], Name: parts[1], Tag: args[1]}
		fields := make([]string, 0, 3)
		if titleSet {
			request.Title = &title
			fields = append(fields, "title")
		}
		if notesSet {
			request.Notes = &notes
			fields = append(fields, "notes")
		}
		if prereleaseSet {
			request.Prerelease = &prerelease
			fields = append(fields, "prerelease")
		}
		result, err := applicationrelease.NewUpdateUseCase(updater).Execute(command.Context(), request)
		if err != nil {
			return mapApplicationError(err, "update release")
		}
		return (releasePresenter{}).PresentUpdated(command.OutOrStdout(), result, fields)
	}}
	command.Flags().StringVar(&title, "title", "", "release title")
	command.Flags().StringVar(&notes, "notes", "", "release notes")
	command.Flags().BoolVar(&prerelease, "prerelease", false, "mark the release as a prerelease")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newReleaseCreateCommand(creator applicationrelease.Creator) *cobra.Command {
	var instance, tag, title, notes string
	var prerelease bool
	command := &cobra.Command{Use: "create OWNER/NAME", Short: "Create a draft release", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return newCommandError(categoryValidation, "create release", fmt.Errorf("OWNER/NAME is required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "create release", err)
		}
		if strings.TrimSpace(tag) == "" {
			return newCommandError(categoryValidation, "create release", fmt.Errorf("release tag is required"))
		}
		if strings.ContainsAny(tag, " \t\n") {
			return newCommandError(categoryValidation, "create release", fmt.Errorf("release tag must not contain whitespace"))
		}
		if strings.TrimSpace(title) == "" {
			return newCommandError(categoryValidation, "create release", fmt.Errorf("release title is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if creator == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			creator = dependencies.ReleaseCreator
			if creator == nil {
				return newCommandError(categoryInternal, "create release", fmt.Errorf("release creator unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		result, err := applicationrelease.NewCreateUseCase(creator).Execute(command.Context(), applicationrelease.CreateRequest{Owner: parts[0], Name: parts[1], Tag: tag, Title: title, Notes: notes, Prerelease: prerelease})
		if err != nil {
			return mapApplicationError(err, "create release")
		}
		return (releasePresenter{}).PresentCreated(command.OutOrStdout(), result)
	}}
	command.Flags().StringVar(&tag, "tag", "", "release tag")
	command.Flags().StringVar(&title, "title", "", "release title")
	command.Flags().StringVar(&notes, "notes", "", "release notes")
	command.Flags().BoolVar(&prerelease, "prerelease", false, "mark the release as a prerelease")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}
