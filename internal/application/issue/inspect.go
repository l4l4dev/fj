package issue

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type InspectUseCase struct{ inspector Inspector }

func NewInspectUseCase(inspector Inspector) InspectUseCase {
	return InspectUseCase{inspector: inspector}
}

func (useCase InspectUseCase) Execute(ctx context.Context, request InspectRequest) (IssueDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return IssueDetail{}, apperror.NewValidation("inspect issue", "OWNER/NAME owner and name are required")
	}
	if request.Number < 1 {
		return IssueDetail{}, apperror.NewValidation("inspect issue", "issue number must be a positive integer")
	}
	if useCase.inspector == nil {
		return IssueDetail{}, apperror.New(apperror.Internal, "inspect issue", "")
	}
	return useCase.inspector.Inspect(ctx, request)
}
