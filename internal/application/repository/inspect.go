package repository

import (
	"context"
	"github.com/l4l4dev/fj/internal/application/apperror"
	"strings"
)

type InspectUseCase struct {
	getter Getter
}

func NewInspectUseCase(getter Getter) InspectUseCase {
	return InspectUseCase{getter: getter}
}

func (useCase InspectUseCase) Execute(ctx context.Context, request GetRequest) (Repository, error) {
	if strings.TrimSpace(request.Owner) == "" {
		return Repository{}, apperror.NewValidation("inspect repository", "repository owner is required")
	}
	if strings.TrimSpace(request.Name) == "" {
		return Repository{}, apperror.NewValidation("inspect repository", "repository name is required")
	}
	return useCase.getter.Get(ctx, request)
}
