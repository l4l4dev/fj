package pullrequest

import (
	"context"
	"testing"
)

type commentViewerStub struct {
	request  ListCommentsRequest
	comments []Comment
}

func (stub *commentViewerStub) ListComments(_ context.Context, request ListCommentsRequest) ([]Comment, error) {
	stub.request = request
	return stub.comments, nil
}

type commentCreatorStub struct {
	request AddCommentRequest
	comment Comment
}

func (stub *commentCreatorStub) AddComment(_ context.Context, request AddCommentRequest) (Comment, error) {
	stub.request = request
	return stub.comment, nil
}

func TestListCommentsUseCaseForwardsTarget(t *testing.T) {
	stub := &commentViewerStub{comments: []Comment{{ID: 5, Body: "First"}}}
	comments, err := NewListCommentsUseCase(stub).Execute(context.Background(), ListCommentsRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || len(comments) != 1 || comments[0].ID != 5 {
		t.Fatalf("unexpected result: %+v err=%v", comments, err)
	}
	if stub.request.Number != 12 || stub.request.Owner != "alice" || stub.request.Name != "project" {
		t.Fatalf("unexpected forwarded target: %+v", stub.request)
	}
}

func TestListCommentsUseCaseValidatesTarget(t *testing.T) {
	tests := []ListCommentsRequest{
		{Owner: "", Name: "project", Number: 1},
		{Owner: "alice", Name: "", Number: 1},
		{Owner: "alice", Name: "project", Number: 0},
	}
	for _, request := range tests {
		if _, err := NewListCommentsUseCase(&commentViewerStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}

func TestAddCommentUseCaseRequiresBody(t *testing.T) {
	if _, err := NewAddCommentUseCase(&commentCreatorStub{}).Execute(context.Background(), AddCommentRequest{Owner: "alice", Name: "project", Number: 12, Body: "  "}); err == nil {
		t.Fatal("expected validation error for empty body")
	}
}

func TestAddCommentUseCaseForwardsRequest(t *testing.T) {
	stub := &commentCreatorStub{comment: Comment{ID: 9, Body: "Looks good"}}
	comment, err := NewAddCommentUseCase(stub).Execute(context.Background(), AddCommentRequest{Owner: "alice", Name: "project", Number: 12, Body: "Looks good"})
	if err != nil || comment.ID != 9 {
		t.Fatalf("unexpected result: %+v err=%v", comment, err)
	}
	if stub.request.Number != 12 || stub.request.Body != "Looks good" {
		t.Fatalf("unexpected forwarded request: %+v", stub.request)
	}
}

func TestCommentUseCasesRequireDependencies(t *testing.T) {
	if _, err := NewListCommentsUseCase(nil).Execute(context.Background(), ListCommentsRequest{Owner: "alice", Name: "project", Number: 1}); err == nil {
		t.Fatal("expected error for nil viewer")
	}
	if _, err := NewAddCommentUseCase(nil).Execute(context.Background(), AddCommentRequest{Owner: "alice", Name: "project", Number: 1, Body: "Hi"}); err == nil {
		t.Fatal("expected error for nil creator")
	}
}
