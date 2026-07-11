package issue

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type ChangeStateUseCase struct{ changer StateChanger }

func NewChangeStateUseCase(changer StateChanger) ChangeStateUseCase {
	return ChangeStateUseCase{changer: changer}
}

func (useCase ChangeStateUseCase) Execute(ctx context.Context, request ChangeStateRequest) (IssueDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return IssueDetail{}, apperror.NewValidation("change issue state", "OWNER/NAME owner and name are required")
	}
	if request.Number < 1 {
		return IssueDetail{}, apperror.NewValidation("change issue state", "issue number must be a positive integer")
	}
	if request.State != StateOpen && request.State != StateClosed {
		return IssueDetail{}, apperror.NewValidation("change issue state", "state must be open or closed")
	}
	if useCase.changer == nil {
		return IssueDetail{}, apperror.New(apperror.Internal, "change issue state", "")
	}
	return useCase.changer.ChangeState(ctx, request)
}
