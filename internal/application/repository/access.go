package repository

import (
	"context"
	"fmt"
	"github.com/l4l4dev/fj/internal/application/apperror"
	"strings"
)

type AccessUseCase struct{ viewer AccessViewer }

func NewAccessUseCase(viewer AccessViewer) AccessUseCase { return AccessUseCase{viewer: viewer} }
func (useCase AccessUseCase) Execute(ctx context.Context, request AccessRequest) (RepositoryAccess, error) {
	if strings.TrimSpace(request.Owner) == "" {
		return RepositoryAccess{}, apperror.NewValidation("view repository access", "repository owner is required")
	}
	if strings.TrimSpace(request.Name) == "" {
		return RepositoryAccess{}, apperror.NewValidation("view repository access", "repository name is required")
	}
	if useCase.viewer == nil {
		return RepositoryAccess{}, fmt.Errorf("access viewer is unavailable")
	}
	return useCase.viewer.ViewAccess(ctx, request)
}
