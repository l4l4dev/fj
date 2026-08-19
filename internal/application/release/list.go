package release

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type ListUseCase struct{ lister Lister }

func NewListUseCase(lister Lister) ListUseCase { return ListUseCase{lister: lister} }

func (u ListUseCase) Execute(ctx context.Context, request ListRequest) ([]Release, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return nil, apperror.NewValidation("list releases", "OWNER/NAME owner and name are required")
	}
	if request.Page < 1 {
		return nil, apperror.NewValidation("list releases", "page must be at least 1")
	}
	if request.Limit < 1 || request.Limit > 100 {
		return nil, apperror.NewValidation("list releases", "limit must be between 1 and 100")
	}
	if u.lister == nil {
		return nil, apperror.New(apperror.Internal, "list releases", "")
	}
	return u.lister.List(ctx, request)
}
