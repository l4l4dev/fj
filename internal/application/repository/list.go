package repository

import (
	"context"
	"github.com/l4l4dev/fj/internal/application/apperror"
)

type ListUseCase struct {
	service Service
}

func NewListUseCase(service Service) ListUseCase {
	return ListUseCase{service: service}
}

func (useCase ListUseCase) Execute(ctx context.Context, request ListRequest) ([]Repository, error) {
	if request.Page < 1 {
		return nil, apperror.NewValidation("list repositories", "page must be at least 1")
	}
	if request.Limit < 1 {
		return nil, apperror.NewValidation("list repositories", "limit must be at least 1")
	}
	return useCase.service.List(ctx, request)
}
