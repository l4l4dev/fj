package pullrequest

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type UpdateUseCase struct{ updater Updater }

func NewUpdateUseCase(updater Updater) UpdateUseCase { return UpdateUseCase{updater: updater} }

func (useCase UpdateUseCase) Execute(ctx context.Context, request UpdateRequest) (PullRequestDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return PullRequestDetail{}, apperror.NewValidation("update pull request", "OWNER/NAME owner and name are required")
	}
	if request.Number < 1 {
		return PullRequestDetail{}, apperror.NewValidation("update pull request", "pull request number must be a positive integer")
	}
	if request.Title == nil && request.Body == nil {
		return PullRequestDetail{}, apperror.NewValidation("update pull request", "at least one pull request field is required")
	}
	if request.Title != nil && strings.TrimSpace(*request.Title) == "" {
		return PullRequestDetail{}, apperror.NewValidation("update pull request", "pull request title must not be empty")
	}
	if useCase.updater == nil {
		return PullRequestDetail{}, apperror.New(apperror.Internal, "update pull request", "")
	}
	return useCase.updater.Update(ctx, request)
}
