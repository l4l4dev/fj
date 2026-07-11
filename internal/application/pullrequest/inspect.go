package pullrequest

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type InspectUseCase struct{ inspector PullRequestInspector }

func NewInspectUseCase(inspector PullRequestInspector) InspectUseCase {
	return InspectUseCase{inspector: inspector}
}

func (u InspectUseCase) Execute(ctx context.Context, request InspectRequest) (PullRequestDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return PullRequestDetail{}, apperror.NewValidation("inspect pull request", "OWNER/NAME owner and name are required")
	}
	if request.Number < 1 {
		return PullRequestDetail{}, apperror.NewValidation("inspect pull request", "pull request number must be a positive integer")
	}
	if u.inspector == nil {
		return PullRequestDetail{}, apperror.New(apperror.Internal, "inspect pull request", "")
	}
	return u.inspector.Inspect(ctx, request)
}
