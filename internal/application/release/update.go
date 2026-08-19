package release

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type UpdateUseCase struct{ updater Updater }

func NewUpdateUseCase(updater Updater) UpdateUseCase { return UpdateUseCase{updater: updater} }

func (u UpdateUseCase) Execute(ctx context.Context, request UpdateRequest) (ReleaseDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return ReleaseDetail{}, apperror.NewValidation("update release", "OWNER/NAME owner and name are required")
	}
	if strings.TrimSpace(request.Tag) == "" {
		return ReleaseDetail{}, apperror.NewValidation("update release", "release tag is required")
	}
	if request.Title == nil && request.Notes == nil && request.Prerelease == nil {
		return ReleaseDetail{}, apperror.NewValidation("update release", "at least one release field is required")
	}
	if request.Title != nil && strings.TrimSpace(*request.Title) == "" {
		return ReleaseDetail{}, apperror.NewValidation("update release", "release title must not be empty")
	}
	if u.updater == nil {
		return ReleaseDetail{}, apperror.New(apperror.Internal, "update release", "")
	}
	return u.updater.Update(ctx, request)
}
