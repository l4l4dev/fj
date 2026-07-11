package issue

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type UpdateUseCase struct{ updater Updater }

func NewUpdateUseCase(updater Updater) UpdateUseCase { return UpdateUseCase{updater: updater} }

func (useCase UpdateUseCase) Execute(ctx context.Context, request UpdateRequest) (IssueDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return IssueDetail{}, apperror.NewValidation("update issue", "OWNER/NAME owner and name are required")
	}
	if request.Number < 1 {
		return IssueDetail{}, apperror.NewValidation("update issue", "issue number must be a positive integer")
	}
	if request.Title == nil && request.Body == nil {
		return IssueDetail{}, apperror.NewValidation("update issue", "at least one issue field is required")
	}
	if useCase.updater == nil {
		return IssueDetail{}, apperror.New(apperror.Internal, "update issue", "")
	}
	return useCase.updater.Update(ctx, request)
}
