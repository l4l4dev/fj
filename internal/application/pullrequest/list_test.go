package pullrequest

import (
	"context"
	"testing"
)

type listerStub struct{ request ListRequest }

func (s *listerStub) List(_ context.Context, request ListRequest) ([]PullRequest, error) {
	s.request = request
	return []PullRequest{{Number: 1}}, nil
}

func TestListUseCaseDelegates(t *testing.T) {
	stub := &listerStub{}
	result, err := NewListUseCase(stub).Execute(context.Background(), ListRequest{Owner: "alice", Name: "project", Page: 1, Limit: 20, State: StateOpen})
	if err != nil || len(result) != 1 || stub.request.State != StateOpen {
		t.Fatalf("unexpected result: %+v request=%+v err=%v", result, stub.request, err)
	}
}

func TestListUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []ListRequest{
		{Name: "project", Page: 1, Limit: 20, State: StateOpen},
		{Owner: "alice", Name: "project", Page: 0, Limit: 20, State: StateOpen},
		{Owner: "alice", Name: "project", Page: 1, Limit: 0, State: StateOpen},
		{Owner: "alice", Name: "project", Page: 1, Limit: 101, State: StateOpen},
		{Owner: "alice", Name: "project", Page: 1, Limit: 20, State: "draft"},
	} {
		if _, err := NewListUseCase(&listerStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}
