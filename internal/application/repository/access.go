package repository

import (
	"context"
	"fmt"
	"strings"
)

type AccessUseCase struct{ viewer AccessViewer }

func NewAccessUseCase(viewer AccessViewer) AccessUseCase { return AccessUseCase{viewer: viewer} }
func (useCase AccessUseCase) Execute(ctx context.Context, request AccessRequest) (RepositoryAccess, error) {
	if strings.TrimSpace(request.Owner) == "" {
		return RepositoryAccess{}, ValidationError{message: "repository owner is required"}
	}
	if strings.TrimSpace(request.Name) == "" {
		return RepositoryAccess{}, ValidationError{message: "repository name is required"}
	}
	if useCase.viewer == nil {
		return RepositoryAccess{}, fmt.Errorf("access viewer is unavailable")
	}
	return useCase.viewer.ViewAccess(ctx, request)
}
