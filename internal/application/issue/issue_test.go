package issue

import (
	"context"
	"testing"
)

type stubLister struct {
	request ListRequest
}

func (stub *stubLister) List(_ context.Context, request ListRequest) (Page, error) {
	stub.request = request
	return Page{Page: request.Page, Limit: request.Limit}, nil
}

func TestListUseCaseValidatesAndDelegates(t *testing.T) {
	stub := &stubLister{}
	_, err := NewListUseCase(stub).Execute(context.Background(), ListRequest{Owner: "alice", Name: "project", Page: 2, Limit: 10, State: StateOpen})
	if err != nil || stub.request.Page != 2 || stub.request.State != StateOpen {
		t.Fatalf("unexpected delegation: request=%+v err=%v", stub.request, err)
	}
}

func TestListUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []ListRequest{{Owner: "", Name: "repo", Page: 1, Limit: 1, State: StateOpen}, {Owner: "o", Name: "r", Page: 0, Limit: 1, State: StateOpen}, {Owner: "o", Name: "r", Page: 1, Limit: 1, State: State("bad")}} {
		if _, err := NewListUseCase(&stubLister{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}
