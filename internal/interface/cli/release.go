package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	applicationrelease "github.com/l4l4dev/fj/internal/application/release"
	"github.com/spf13/cobra"
)

type releaseDependencies struct {
	lister        applicationrelease.Lister
	inspector     applicationrelease.Inspector
	creator       applicationrelease.Creator
	updater       applicationrelease.Updater
	publisher     applicationrelease.Publisher
	assetUploader applicationrelease.AssetUploader
	assetDeleter  applicationrelease.AssetDeleter
}

func newReleaseCommand(dependencies releaseDependencies) *cobra.Command {
	command := &cobra.Command{Use: "release", Short: "Manage releases"}
	command.AddCommand(newReleaseListCommand(dependencies.lister))
	command.AddCommand(newReleaseInspectCommand(dependencies.inspector))
	command.AddCommand(newReleaseCreateCommand(dependencies.creator))
	command.AddCommand(newReleaseUpdateCommand(dependencies.updater))
	command.AddCommand(newReleasePublishCommand(dependencies.inspector, dependencies.publisher))
	command.AddCommand(newReleaseAssetCommand(dependencies.inspector, dependencies.assetUploader, dependencies.assetDeleter))
	return command
}

func newReleaseAssetCommand(inspector applicationrelease.Inspector, uploader applicationrelease.AssetUploader, deleter applicationrelease.AssetDeleter) *cobra.Command {
	command := &cobra.Command{Use: "asset", Short: "Manage release assets"}
	command.AddCommand(newReleaseAssetListCommand(inspector))
	command.AddCommand(newReleaseAssetUploadCommand(uploader))
	command.AddCommand(newReleaseAssetDeleteCommand(deleter))
	return command
}

func newReleaseAssetListCommand(inspector applicationrelease.Inspector) *cobra.Command {
	var instance string
	command := &cobra.Command{Use: "list OWNER/NAME TAG", Short: "List release assets", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return newCommandError(categoryValidation, "list release assets", fmt.Errorf("OWNER/NAME and release tag are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "list release assets", err)
		}
		if strings.TrimSpace(args[1]) == "" {
			return newCommandError(categoryValidation, "list release assets", fmt.Errorf("release tag is required"))
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
				return newCommandError(categoryInternal, "list release assets", fmt.Errorf("release inspector unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		result, err := applicationrelease.NewInspectUseCase(inspector).Execute(command.Context(), applicationrelease.InspectRequest{Owner: parts[0], Name: parts[1], Tag: args[1]})
		if err != nil {
			return mapApplicationError(err, "list release assets")
		}
		return (releasePresenter{}).PresentAssets(command.OutOrStdout(), result)
	}}
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newReleaseAssetUploadCommand(uploader applicationrelease.AssetUploader) *cobra.Command {
	var instance, name string
	command := &cobra.Command{Use: "upload OWNER/NAME TAG FILE", Short: "Upload a release asset", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 3 {
			return newCommandError(categoryValidation, "upload release asset", fmt.Errorf("OWNER/NAME, release tag, and asset file are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "upload release asset", err)
		}
		if strings.TrimSpace(args[1]) == "" {
			return newCommandError(categoryValidation, "upload release asset", fmt.Errorf("release tag is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		content, err := os.ReadFile(args[2])
		if err != nil {
			return newCommandError(categoryValidation, "upload release asset", err)
		}
		assetName := name
		if strings.TrimSpace(assetName) == "" {
			assetName = filepath.Base(args[2])
		}
		if uploader == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			uploader = dependencies.ReleaseAssetUploader
			if uploader == nil {
				return newCommandError(categoryInternal, "upload release asset", fmt.Errorf("release asset uploader unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		result, err := applicationrelease.NewUploadAssetUseCase(uploader).Execute(command.Context(), applicationrelease.UploadAssetRequest{Owner: parts[0], Name: parts[1], Tag: args[1], AssetName: assetName, Content: content})
		if err != nil {
			return mapApplicationError(err, "upload release asset")
		}
		return (releasePresenter{}).PresentAssetUploaded(command.OutOrStdout(), result, args[1])
	}}
	command.Flags().StringVar(&name, "name", "", "asset name (defaults to the file name)")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
	return command
}

func newReleaseAssetDeleteCommand(deleter applicationrelease.AssetDeleter) *cobra.Command {
	var instance string
	var confirm bool
	command := &cobra.Command{Use: "delete OWNER/NAME TAG ASSET_NAME", Short: "Delete a release asset", Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 3 {
			return newCommandError(categoryValidation, "delete release asset", fmt.Errorf("OWNER/NAME, release tag, and asset name are required"))
		}
		if err := validateRepositoryTarget(args[0]); err != nil {
			return newCommandError(categoryValidation, "delete release asset", err)
		}
		if strings.TrimSpace(args[1]) == "" {
			return newCommandError(categoryValidation, "delete release asset", fmt.Errorf("release tag is required"))
		}
		if strings.TrimSpace(args[2]) == "" {
			return newCommandError(categoryValidation, "delete release asset", fmt.Errorf("asset name is required"))
		}
		return nil
	}, RunE: func(command *cobra.Command, args []string) error {
		if !confirm {
			return newCommandError(categoryValidation, "delete release asset", fmt.Errorf("deleting an asset requires --confirm"))
		}
		if deleter == nil {
			dependencies, err := composeRepositoryDependencies(command.Context(), instance)
			if err != nil {
				return err
			}
			deleter = dependencies.ReleaseAssetDeleter
			if deleter == nil {
				return newCommandError(categoryInternal, "delete release asset", fmt.Errorf("release asset deleter unavailable"))
			}
		}
		parts := strings.SplitN(args[0], "/", 2)
		if err := applicationrelease.NewDeleteAssetUseCase(deleter).Execute(command.Context(), applicationrelease.DeleteAssetRequest{Owner: parts[0], Name: parts[1], Tag: args[1], AssetName: args[2]}); err != nil {
			return mapApplicationError(err, "delete release asset")
		}
		return (releasePresenter{}).PresentAssetDeleted(command.OutOrStdout(), args[1], args[2])
	}}
	command.Flags().BoolVar(&confirm, "confirm", false, "confirm deleting the asset")
	command.Flags().StringVar(&instance, "instance", "", "configured Forgejo instance profile")
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
