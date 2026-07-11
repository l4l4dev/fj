package repository

import (
	"context"
	"github.com/l4l4dev/fj/internal/application/apperror"
	"strings"
)

type ArchiveUseCase struct{ archiver Archiver }

func NewArchiveUseCase(archiver Archiver) ArchiveUseCase { return ArchiveUseCase{archiver: archiver} }

func (useCase ArchiveUseCase) Execute(ctx context.Context, request ArchiveRequest) (Repository, error) {
	if strings.TrimSpace(request.Owner) == "" {
		return Repository{}, apperror.NewValidation("archive repository", "repository owner is required")
	}
	if strings.TrimSpace(request.Name) == "" {
		return Repository{}, apperror.NewValidation("archive repository", "repository name is required")
	}
	return useCase.archiver.SetArchived(ctx, request)
}
