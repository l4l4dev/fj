package issue

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type State string

const (
	StateOpen   State = "open"
	StateClosed State = "closed"
)

type Issue struct {
	Number int
	Title  string
	State  State
}

type ListRequest struct {
	Owner string
	Name  string
	Page  int
	Limit int
	State State
}

type Page struct {
	Issues    []Issue
	Page      int
	Limit     int
	MorePages bool
}

type Lister interface {
	List(context.Context, ListRequest) (Page, error)
}

type ListUseCase struct{ lister Lister }

func NewListUseCase(lister Lister) ListUseCase { return ListUseCase{lister: lister} }

func (useCase ListUseCase) Execute(ctx context.Context, request ListRequest) (Page, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return Page{}, apperror.NewValidation("list issues", "OWNER/NAME owner and name are required")
	}
	if request.Page < 1 {
		return Page{}, apperror.NewValidation("list issues", "page must be at least 1")
	}
	if request.Limit < 1 {
		return Page{}, apperror.NewValidation("list issues", "limit must be at least 1")
	}
	if request.State != StateOpen && request.State != StateClosed && request.State != "all" {
		return Page{}, apperror.NewValidation("list issues", "state must be open, closed, or all")
	}
	if useCase.lister == nil {
		return Page{}, apperror.New(apperror.Internal, "list issues", "")
	}
	return useCase.lister.List(ctx, request)
}
