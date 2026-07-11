package issue

import (
	"context"
	"testing"
)

type inspectorStub struct{ request InspectRequest }

func (stub *inspectorStub) Inspect(_ context.Context, request InspectRequest) (IssueDetail, error) {
	stub.request = request
	return IssueDetail{Number: request.Number, Title: "Example", State: StateOpen}, nil
}

func TestInspectUseCaseValidatesAndDelegates(t *testing.T) {
	stub := &inspectorStub{}
	result, err := NewInspectUseCase(stub).Execute(context.Background(), InspectRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || result.Number != 12 || stub.request.Owner != "alice" {
		t.Fatalf("unexpected result: %+v request=%+v err=%v", result, stub.request, err)
	}
}

func TestInspectUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []InspectRequest{{Name: "project", Number: 1}, {Owner: "alice", Number: 1}, {Owner: "alice", Name: "project", Number: 0}} {
		if _, err := NewInspectUseCase(&inspectorStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}
