package repository

import (
	"context"
	"fmt"
)

type ListUseCase struct {
	service Service
}

func NewListUseCase(service Service) ListUseCase {
	return ListUseCase{service: service}
}

func (useCase ListUseCase) Execute(ctx context.Context, request ListRequest) ([]Repository, error) {
	if request.Page < 1 {
		return nil, fmt.Errorf("page must be at least 1")
	}
	if request.Limit < 1 {
		return nil, fmt.Errorf("limit must be at least 1")
	}
	return useCase.service.List(ctx, request)
}
