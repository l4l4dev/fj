package repository

import (
	"context"
	"strings"
)

type ArchiveUseCase struct{ archiver Archiver }

func NewArchiveUseCase(archiver Archiver) ArchiveUseCase { return ArchiveUseCase{archiver: archiver} }

func (useCase ArchiveUseCase) Execute(ctx context.Context, request ArchiveRequest) (Repository, error) {
	if strings.TrimSpace(request.Owner) == "" {
		return Repository{}, ValidationError{message: "repository owner is required"}
	}
	if strings.TrimSpace(request.Name) == "" {
		return Repository{}, ValidationError{message: "repository name is required"}
	}
	return useCase.archiver.SetArchived(ctx, request)
}
