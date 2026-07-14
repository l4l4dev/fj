package pullrequest

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type CreateUseCase struct{ creator PullRequestCreator }

func NewCreateUseCase(creator PullRequestCreator) CreateUseCase {
	return CreateUseCase{creator: creator}
}

func (useCase CreateUseCase) Execute(ctx context.Context, request CreateRequest) (PullRequestDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return PullRequestDetail{}, apperror.NewValidation("create pull request", "OWNER/NAME owner and name are required")
	}
	if strings.TrimSpace(request.Title) == "" {
		return PullRequestDetail{}, apperror.NewValidation("create pull request", "pull request title is required")
	}
	if strings.TrimSpace(request.HeadBranch) == "" {
		return PullRequestDetail{}, apperror.NewValidation("create pull request", "head branch is required")
	}
	if strings.Contains(request.HeadBranch, ":") {
		return PullRequestDetail{}, apperror.NewValidation("create pull request", "cross-fork head branches are not supported")
	}
	if strings.TrimSpace(request.BaseBranch) == "" {
		return PullRequestDetail{}, apperror.NewValidation("create pull request", "base branch is required")
	}
	if strings.Contains(request.BaseBranch, ":") {
		return PullRequestDetail{}, apperror.NewValidation("create pull request", "base branch must belong to the target repository")
	}
	if request.HeadBranch == request.BaseBranch {
		return PullRequestDetail{}, apperror.NewValidation("create pull request", "head and base branches must differ")
	}
	if useCase.creator == nil {
		return PullRequestDetail{}, apperror.New(apperror.Internal, "create pull request", "")
	}
	return useCase.creator.Create(ctx, request)
}
