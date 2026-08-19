package pullrequest

import (
	"context"
	"errors"
	"testing"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type updaterStub struct {
	request UpdateRequest
	detail  PullRequestDetail
	err     error
}

func (stub *updaterStub) Update(_ context.Context, request UpdateRequest) (PullRequestDetail, error) {
	stub.request = request
	return stub.detail, stub.err
}

func TestUpdateUseCaseForwardsOnlySuppliedFields(t *testing.T) {
	title := "New title"
	stub := &updaterStub{detail: PullRequestDetail{Number: 12, Title: "New title"}}
	result, err := NewUpdateUseCase(stub).Execute(context.Background(), UpdateRequest{Owner: "alice", Name: "project", Number: 12, Title: &title})
	if err != nil || result.Number != 12 {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if stub.request.Title == nil || *stub.request.Title != "New title" || stub.request.Body != nil {
		t.Fatalf("unexpected forwarded request: %+v", stub.request)
	}
}

func TestUpdateUseCaseValidatesInput(t *testing.T) {
	title := "Title"
	empty := " "
	tests := []UpdateRequest{
		{Owner: "", Name: "project", Number: 1, Title: &title},
		{Owner: "alice", Name: "", Number: 1, Title: &title},
		{Owner: "alice", Name: "project", Number: 0, Title: &title},
		{Owner: "alice", Name: "project", Number: 1},
		{Owner: "alice", Name: "project", Number: 1, Title: &empty},
	}
	for _, request := range tests {
		if _, err := NewUpdateUseCase(&updaterStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		} else {
			var validation apperror.ValidationError
			if !errors.As(err, &validation) {
				t.Fatalf("expected validation error, got %v", err)
			}
		}
	}
}

func TestUpdateUseCaseRequiresUpdater(t *testing.T) {
	title := "Title"
	if _, err := NewUpdateUseCase(nil).Execute(context.Background(), UpdateRequest{Owner: "alice", Name: "project", Number: 1, Title: &title}); err == nil {
		t.Fatal("expected error for nil updater")
	}
}
