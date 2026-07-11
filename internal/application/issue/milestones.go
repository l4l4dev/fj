package issue

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type SetMilestoneUseCase struct{ setter MilestoneSetter }

func NewSetMilestoneUseCase(setter MilestoneSetter) SetMilestoneUseCase {
	return SetMilestoneUseCase{setter: setter}
}

func (useCase SetMilestoneUseCase) Execute(ctx context.Context, request SetMilestoneRequest) (Milestone, error) {
	if err := validateMilestoneTarget(request.Owner, request.Name, request.Number, "set issue milestone"); err != nil {
		return Milestone{}, err
	}
	if strings.TrimSpace(request.Milestone) == "" {
		return Milestone{}, apperror.NewValidation("set issue milestone", "milestone is required")
	}
	if useCase.setter == nil {
		return Milestone{}, apperror.New(apperror.Internal, "set issue milestone", "")
	}
	return useCase.setter.SetMilestone(ctx, request)
}

type ClearMilestoneUseCase struct{ clearer MilestoneClearer }

func NewClearMilestoneUseCase(clearer MilestoneClearer) ClearMilestoneUseCase {
	return ClearMilestoneUseCase{clearer: clearer}
}

func (useCase ClearMilestoneUseCase) Execute(ctx context.Context, request ClearMilestoneRequest) error {
	if err := validateMilestoneTarget(request.Owner, request.Name, request.Number, "clear issue milestone"); err != nil {
		return err
	}
	if useCase.clearer == nil {
		return apperror.New(apperror.Internal, "clear issue milestone", "")
	}
	return useCase.clearer.ClearMilestone(ctx, request)
}

func validateMilestoneTarget(owner, name string, number int, operation string) error {
	if strings.TrimSpace(owner) == "" || strings.TrimSpace(name) == "" {
		return apperror.NewValidation(operation, "OWNER/NAME owner and name are required")
	}
	if number < 1 {
		return apperror.NewValidation(operation, "issue number must be a positive integer")
	}
	return nil
}
