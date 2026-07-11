package repository

import (
	"context"
	"github.com/l4l4dev/fj/internal/application/apperror"
	"strings"
)

type CreateUseCase struct{ creator Creator }

func NewCreateUseCase(creator Creator) CreateUseCase { return CreateUseCase{creator: creator} }

func (useCase CreateUseCase) Execute(ctx context.Context, request CreateRequest) (Repository, error) {
	if strings.TrimSpace(request.Name) == "" {
		return Repository{}, apperror.NewValidation("create repository", "repository name is required")
	}
	return useCase.creator.Create(ctx, request)
}
