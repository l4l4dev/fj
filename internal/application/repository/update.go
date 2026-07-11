package repository

import (
	"context"
	"github.com/l4l4dev/fj/internal/application/apperror"
	"strings"
)

type UpdateUseCase struct{ updater Updater }

func newValidationError(message string) error {
	return apperror.NewValidation("update repository", message)
}

func NewUpdateUseCase(updater Updater) UpdateUseCase { return UpdateUseCase{updater: updater} }

func (useCase UpdateUseCase) Execute(ctx context.Context, request UpdateRequest) (Repository, error) {
	if strings.TrimSpace(request.Owner) == "" {
		return Repository{}, newValidationError("repository owner is required")
	}
	if strings.TrimSpace(request.Name) == "" {
		return Repository{}, newValidationError("repository name is required")
	}
	if request.Description == nil && request.Private == nil {
		return Repository{}, newValidationError("at least one repository field is required")
	}
	return useCase.updater.Update(ctx, request)
}
