package repository

import (
	"context"
	"fmt"
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
		return Repository{}, fmt.Errorf("repository owner is required")
	}
	if strings.TrimSpace(request.Name) == "" {
		return Repository{}, fmt.Errorf("repository name is required")
	}
	return useCase.getter.Get(ctx, request)
}
