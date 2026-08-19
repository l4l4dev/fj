package release

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type DeleteUseCase struct {
	inspector Inspector
	deleter   Deleter
}

func NewDeleteUseCase(inspector Inspector, deleter Deleter) DeleteUseCase {
	return DeleteUseCase{inspector: inspector, deleter: deleter}
}

func (u DeleteUseCase) Execute(ctx context.Context, request InspectRequest) (ReleaseDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return ReleaseDetail{}, apperror.NewValidation("delete release", "OWNER/NAME owner and name are required")
	}
	if strings.TrimSpace(request.Tag) == "" {
		return ReleaseDetail{}, apperror.NewValidation("delete release", "release tag is required")
	}
	if u.inspector == nil || u.deleter == nil {
		return ReleaseDetail{}, apperror.New(apperror.Internal, "delete release", "")
	}
	detail, err := u.inspector.Inspect(ctx, request)
	if err != nil {
		return ReleaseDetail{}, err
	}
	if err := u.deleter.Delete(ctx, DeleteRequest{Owner: request.Owner, Name: request.Name, ID: detail.ID}); err != nil {
		return ReleaseDetail{}, err
	}
	return detail, nil
}
