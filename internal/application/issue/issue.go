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

type IssueDetail struct {
	Number int
	Title  string
	State  State
	Body   string
}

type IssueFilter struct {
	Assignee string
	Label    string
}

type ListRequest struct {
	Owner  string
	Name   string
	Page   int
	Limit  int
	State  State
	Filter IssueFilter
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

type InspectRequest struct {
	Owner  string
	Name   string
	Number int
}

type Inspector interface {
	Inspect(context.Context, InspectRequest) (IssueDetail, error)
}

type CreateRequest struct {
	Owner string
	Name  string
	Title string
	Body  string
}

type Creator interface {
	Create(context.Context, CreateRequest) (IssueDetail, error)
}

type UpdateRequest struct {
	Owner  string
	Name   string
	Number int
	Title  *string
	Body   *string
}

type Updater interface {
	Update(context.Context, UpdateRequest) (IssueDetail, error)
}

type ChangeStateRequest struct {
	Owner  string
	Name   string
	Number int
	State  State
}

type StateChanger interface {
	ChangeState(context.Context, ChangeStateRequest) (IssueDetail, error)
}

type Comment struct {
	ID   int64
	Body string
}

type ListCommentsRequest struct {
	Owner  string
	Name   string
	Number int
}

type AddCommentRequest struct {
	Owner  string
	Name   string
	Number int
	Body   string
}

type CommentViewer interface {
	ListComments(context.Context, ListCommentsRequest) ([]Comment, error)
}

type CommentCreator interface {
	AddComment(context.Context, AddCommentRequest) (Comment, error)
}

type Label struct {
	ID   int64
	Name string
}

type AddLabelRequest struct {
	Owner  string
	Name   string
	Number int
	Label  string
}

type RemoveLabelRequest struct {
	Owner  string
	Name   string
	Number int
	Label  string
}

type LabelAdder interface {
	AddLabel(context.Context, AddLabelRequest) (Label, error)
}

type LabelRemover interface {
	RemoveLabel(context.Context, RemoveLabelRequest) (Label, error)
}

type Milestone struct {
	ID    int64
	Title string
}

type SetMilestoneRequest struct {
	Owner     string
	Name      string
	Number    int
	Milestone string
}

type ClearMilestoneRequest struct {
	Owner  string
	Name   string
	Number int
}

type MilestoneSetter interface {
	SetMilestone(context.Context, SetMilestoneRequest) (Milestone, error)
}

type MilestoneClearer interface {
	ClearMilestone(context.Context, ClearMilestoneRequest) error
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
