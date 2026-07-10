package repository

import (
	"context"
	"fmt"
	"strings"
)

type CreateUseCase struct{ creator Creator }

func NewCreateUseCase(creator Creator) CreateUseCase { return CreateUseCase{creator: creator} }

func (useCase CreateUseCase) Execute(ctx context.Context, request CreateRequest) (Repository, error) {
	if strings.TrimSpace(request.Name) == "" {
		return Repository{}, fmt.Errorf("repository name is required")
	}
	return useCase.creator.Create(ctx, request)
}
