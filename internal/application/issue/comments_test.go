package issue

import (
	"context"
	"testing"
)

type commentViewerStub struct{ request ListCommentsRequest }

func (stub *commentViewerStub) ListComments(_ context.Context, request ListCommentsRequest) ([]Comment, error) {
	stub.request = request
	return []Comment{{ID: 1, Body: "hello"}}, nil
}

type commentCreatorStub struct{ request AddCommentRequest }

func (stub *commentCreatorStub) AddComment(_ context.Context, request AddCommentRequest) (Comment, error) {
	stub.request = request
	return Comment{ID: 2, Body: request.Body}, nil
}

func TestCommentUseCasesDelegate(t *testing.T) {
	viewer := &commentViewerStub{}
	comments, err := NewListCommentsUseCase(viewer).Execute(context.Background(), ListCommentsRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || len(comments) != 1 || viewer.request.Number != 12 {
		t.Fatalf("unexpected list result: %+v request=%+v err=%v", comments, viewer.request, err)
	}
	creator := &commentCreatorStub{}
	comment, err := NewAddCommentUseCase(creator).Execute(context.Background(), AddCommentRequest{Owner: "alice", Name: "project", Number: 12, Body: "hello"})
	if err != nil || comment.Body != "hello" || creator.request.Body != "hello" {
		t.Fatalf("unexpected add result: %+v request=%+v err=%v", comment, creator.request, err)
	}
}

func TestAddCommentRejectsEmptyBody(t *testing.T) {
	for _, body := range []string{"", "  "} {
		if _, err := NewAddCommentUseCase(&commentCreatorStub{}).Execute(context.Background(), AddCommentRequest{Owner: "alice", Name: "project", Number: 1, Body: body}); err == nil {
			t.Fatalf("expected validation error for body %q", body)
		}
	}
}
