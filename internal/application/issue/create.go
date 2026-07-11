package issue

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type CreateUseCase struct{ creator Creator }

func NewCreateUseCase(creator Creator) CreateUseCase { return CreateUseCase{creator: creator} }

func (useCase CreateUseCase) Execute(ctx context.Context, request CreateRequest) (IssueDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return IssueDetail{}, apperror.NewValidation("create issue", "OWNER/NAME owner and name are required")
	}
	if strings.TrimSpace(request.Title) == "" {
		return IssueDetail{}, apperror.NewValidation("create issue", "issue title is required")
	}
	if useCase.creator == nil {
		return IssueDetail{}, apperror.New(apperror.Internal, "create issue", "")
	}
	return useCase.creator.Create(ctx, request)
}
