package release

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type CreateUseCase struct{ creator Creator }

func NewCreateUseCase(creator Creator) CreateUseCase { return CreateUseCase{creator: creator} }

func (u CreateUseCase) Execute(ctx context.Context, request CreateRequest) (ReleaseDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return ReleaseDetail{}, apperror.NewValidation("create release", "OWNER/NAME owner and name are required")
	}
	if strings.TrimSpace(request.Tag) == "" {
		return ReleaseDetail{}, apperror.NewValidation("create release", "release tag is required")
	}
	if strings.ContainsAny(request.Tag, " \t\n") {
		return ReleaseDetail{}, apperror.NewValidation("create release", "release tag must not contain whitespace")
	}
	if strings.TrimSpace(request.Title) == "" {
		return ReleaseDetail{}, apperror.NewValidation("create release", "release title is required")
	}
	if u.creator == nil {
		return ReleaseDetail{}, apperror.New(apperror.Internal, "create release", "")
	}
	return u.creator.Create(ctx, request)
}
