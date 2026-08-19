package release

import (
	"context"
	"testing"
)

type listerStub struct{ request ListRequest }

func (s *listerStub) List(_ context.Context, request ListRequest) ([]Release, error) {
	s.request = request
	return []Release{{ID: 1}}, nil
}

func TestListUseCaseDelegates(t *testing.T) {
	stub := &listerStub{}
	result, err := NewListUseCase(stub).Execute(context.Background(), ListRequest{Owner: "alice", Name: "project", Page: 1, Limit: 20})
	if err != nil || len(result) != 1 || stub.request.Owner != "alice" || stub.request.Name != "project" {
		t.Fatalf("unexpected result: %+v request=%+v err=%v", result, stub.request, err)
	}
}

func TestListUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []ListRequest{
		{Name: "project", Page: 1, Limit: 20},
		{Owner: "alice", Page: 1, Limit: 20},
		{Owner: "alice", Name: "project", Page: 0, Limit: 20},
		{Owner: "alice", Name: "project", Page: 1, Limit: 0},
		{Owner: "alice", Name: "project", Page: 1, Limit: 101},
	} {
		if _, err := NewListUseCase(&listerStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}

func TestListUseCaseRejectsNilLister(t *testing.T) {
	if _, err := NewListUseCase(nil).Execute(context.Background(), ListRequest{Owner: "alice", Name: "project", Page: 1, Limit: 20}); err == nil {
		t.Fatal("expected internal error for nil lister")
	}
}
