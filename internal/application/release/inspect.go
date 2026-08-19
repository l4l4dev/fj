package release

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type InspectUseCase struct{ inspector Inspector }

func NewInspectUseCase(inspector Inspector) InspectUseCase {
	return InspectUseCase{inspector: inspector}
}

func (u InspectUseCase) Execute(ctx context.Context, request InspectRequest) (ReleaseDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return ReleaseDetail{}, apperror.NewValidation("inspect release", "OWNER/NAME owner and name are required")
	}
	if strings.TrimSpace(request.Tag) == "" {
		return ReleaseDetail{}, apperror.NewValidation("inspect release", "release tag is required")
	}
	if u.inspector == nil {
		return ReleaseDetail{}, apperror.New(apperror.Internal, "inspect release", "")
	}
	return u.inspector.Inspect(ctx, request)
}
