package release

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type UploadAssetUseCase struct{ uploader AssetUploader }

func NewUploadAssetUseCase(uploader AssetUploader) UploadAssetUseCase {
	return UploadAssetUseCase{uploader: uploader}
}

func (u UploadAssetUseCase) Execute(ctx context.Context, request UploadAssetRequest) (Asset, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return Asset{}, apperror.NewValidation("upload release asset", "OWNER/NAME owner and name are required")
	}
	if strings.TrimSpace(request.Tag) == "" {
		return Asset{}, apperror.NewValidation("upload release asset", "release tag is required")
	}
	if strings.TrimSpace(request.AssetName) == "" {
		return Asset{}, apperror.NewValidation("upload release asset", "asset name is required")
	}
	if len(request.Content) == 0 {
		return Asset{}, apperror.NewValidation("upload release asset", "asset file is empty")
	}
	if u.uploader == nil {
		return Asset{}, apperror.New(apperror.Internal, "upload release asset", "")
	}
	return u.uploader.UploadAsset(ctx, request)
}

type DeleteAssetUseCase struct{ deleter AssetDeleter }

func NewDeleteAssetUseCase(deleter AssetDeleter) DeleteAssetUseCase {
	return DeleteAssetUseCase{deleter: deleter}
}

func (u DeleteAssetUseCase) Execute(ctx context.Context, request DeleteAssetRequest) error {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return apperror.NewValidation("delete release asset", "OWNER/NAME owner and name are required")
	}
	if strings.TrimSpace(request.Tag) == "" {
		return apperror.NewValidation("delete release asset", "release tag is required")
	}
	if strings.TrimSpace(request.AssetName) == "" {
		return apperror.NewValidation("delete release asset", "asset name is required")
	}
	if u.deleter == nil {
		return apperror.New(apperror.Internal, "delete release asset", "")
	}
	return u.deleter.DeleteAsset(ctx, request)
}
