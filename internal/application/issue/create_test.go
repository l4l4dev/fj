package issue

import (
	"context"
	"testing"
)

type creatorStub struct{ request CreateRequest }

func (stub *creatorStub) Create(_ context.Context, request CreateRequest) (IssueDetail, error) {
	stub.request = request
	return IssueDetail{Number: 13, Title: request.Title, Body: request.Body}, nil
}

func TestCreateUseCaseValidatesAndDelegates(t *testing.T) {
	stub := &creatorStub{}
	result, err := NewCreateUseCase(stub).Execute(context.Background(), CreateRequest{Owner: "alice", Name: "project", Title: "Fix it"})
	if err != nil || result.Title != "Fix it" || stub.request.Body != "" {
		t.Fatalf("unexpected result: %+v request=%+v err=%v", result, stub.request, err)
	}
}

func TestCreateUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []CreateRequest{{Name: "project", Title: "title"}, {Owner: "alice", Title: "title"}, {Owner: "alice", Name: "project", Title: "  "}} {
		if _, err := NewCreateUseCase(&creatorStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}
