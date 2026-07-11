package issue

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type AssignUseCase struct{ assigner Assigner }

func NewAssignUseCase(assigner Assigner) AssignUseCase { return AssignUseCase{assigner: assigner} }

func (u AssignUseCase) Execute(ctx context.Context, request AssignRequest) (Assignment, error) {
	if err := validateAssignmentTarget(request.Owner, request.Name, request.Number, "assign issue"); err != nil {
		return Assignment{}, err
	}
	if strings.TrimSpace(request.Username) == "" || strings.EqualFold(strings.TrimSpace(request.Username), "none") {
		return Assignment{}, apperror.NewValidation("assign issue", "username is required")
	}
	if u.assigner == nil {
		return Assignment{}, apperror.New(apperror.Internal, "assign issue", "")
	}
	return u.assigner.Assign(ctx, request)
}

type UnassignUseCase struct{ unassigner Unassigner }

func NewUnassignUseCase(unassigner Unassigner) UnassignUseCase {
	return UnassignUseCase{unassigner: unassigner}
}

func (u UnassignUseCase) Execute(ctx context.Context, request UnassignRequest) error {
	if err := validateAssignmentTarget(request.Owner, request.Name, request.Number, "unassign issue"); err != nil {
		return err
	}
	if u.unassigner == nil {
		return apperror.New(apperror.Internal, "unassign issue", "")
	}
	return u.unassigner.Unassign(ctx, request)
}

func validateAssignmentTarget(owner, name string, number int, operation string) error {
	if strings.TrimSpace(owner) == "" || strings.TrimSpace(name) == "" {
		return apperror.NewValidation(operation, "OWNER/NAME owner and name are required")
	}
	if number < 1 {
		return apperror.NewValidation(operation, "issue number must be a positive integer")
	}
	return nil
}
