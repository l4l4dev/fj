package issue

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type ListCommentsUseCase struct{ viewer CommentViewer }

func NewListCommentsUseCase(viewer CommentViewer) ListCommentsUseCase {
	return ListCommentsUseCase{viewer: viewer}
}

func (useCase ListCommentsUseCase) Execute(ctx context.Context, request ListCommentsRequest) ([]Comment, error) {
	if err := validateCommentTarget(request.Owner, request.Name, request.Number, "list issue comments"); err != nil {
		return nil, err
	}
	if useCase.viewer == nil {
		return nil, apperror.New(apperror.Internal, "list issue comments", "")
	}
	return useCase.viewer.ListComments(ctx, request)
}

type AddCommentUseCase struct{ creator CommentCreator }

func NewAddCommentUseCase(creator CommentCreator) AddCommentUseCase {
	return AddCommentUseCase{creator: creator}
}

func (useCase AddCommentUseCase) Execute(ctx context.Context, request AddCommentRequest) (Comment, error) {
	if err := validateCommentTarget(request.Owner, request.Name, request.Number, "add issue comment"); err != nil {
		return Comment{}, err
	}
	if strings.TrimSpace(request.Body) == "" {
		return Comment{}, apperror.NewValidation("add issue comment", "comment body is required")
	}
	if useCase.creator == nil {
		return Comment{}, apperror.New(apperror.Internal, "add issue comment", "")
	}
	return useCase.creator.AddComment(ctx, request)
}

func validateCommentTarget(owner, name string, number int, operation string) error {
	if strings.TrimSpace(owner) == "" || strings.TrimSpace(name) == "" {
		return apperror.NewValidation(operation, "OWNER/NAME owner and name are required")
	}
	if number < 1 {
		return apperror.NewValidation(operation, "issue number must be a positive integer")
	}
	return nil
}
