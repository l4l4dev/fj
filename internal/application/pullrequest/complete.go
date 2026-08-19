package pullrequest

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type MergeUseCase struct{ merger Merger }

func NewMergeUseCase(merger Merger) MergeUseCase { return MergeUseCase{merger: merger} }

func (useCase MergeUseCase) Execute(ctx context.Context, request MergeRequest) error {
	if err := validateCompletionTarget(request.Owner, request.Name, request.Number, "merge pull request"); err != nil {
		return err
	}
	switch request.Method {
	case MergeMethodMerge, MergeMethodRebase, MergeMethodSquash:
	default:
		return apperror.NewValidation("merge pull request", "merge method must be merge, rebase, or squash")
	}
	if useCase.merger == nil {
		return apperror.New(apperror.Internal, "merge pull request", "")
	}
	return useCase.merger.Merge(ctx, request)
}

type CloseUseCase struct{ closer Closer }

func NewCloseUseCase(closer Closer) CloseUseCase { return CloseUseCase{closer: closer} }

func (useCase CloseUseCase) Execute(ctx context.Context, request CloseRequest) (PullRequestDetail, error) {
	if err := validateCompletionTarget(request.Owner, request.Name, request.Number, "close pull request"); err != nil {
		return PullRequestDetail{}, err
	}
	if useCase.closer == nil {
		return PullRequestDetail{}, apperror.New(apperror.Internal, "close pull request", "")
	}
	return useCase.closer.Close(ctx, request)
}

func validateCompletionTarget(owner, name string, number int, operation string) error {
	if strings.TrimSpace(owner) == "" || strings.TrimSpace(name) == "" {
		return apperror.NewValidation(operation, "OWNER/NAME owner and name are required")
	}
	if number < 1 {
		return apperror.NewValidation(operation, "pull request number must be a positive integer")
	}
	return nil
}
