package pullrequest

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type ListUseCase struct{ lister PullRequestLister }

func NewListUseCase(lister PullRequestLister) ListUseCase { return ListUseCase{lister: lister} }

func (u ListUseCase) Execute(ctx context.Context, request ListRequest) ([]PullRequest, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return nil, apperror.NewValidation("list pull requests", "OWNER/NAME owner and name are required")
	}
	if request.Page < 1 {
		return nil, apperror.NewValidation("list pull requests", "page must be at least 1")
	}
	if request.Limit < 1 || request.Limit > 100 {
		return nil, apperror.NewValidation("list pull requests", "limit must be between 1 and 100")
	}
	if request.State != StateOpen && request.State != StateClosed && request.State != StateAll {
		return nil, apperror.NewValidation("list pull requests", "state must be open, closed, or all")
	}
	if u.lister == nil {
		return nil, apperror.New(apperror.Internal, "list pull requests", "")
	}
	return u.lister.List(ctx, request)
}
