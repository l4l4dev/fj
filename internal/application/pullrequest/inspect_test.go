package pullrequest

import (
	"context"
	"testing"
)

type inspectorStub struct{ request InspectRequest }

func (s *inspectorStub) Inspect(_ context.Context, request InspectRequest) (PullRequestDetail, error) {
	s.request = request
	return PullRequestDetail{Number: request.Number}, nil
}

func TestInspectUseCaseDelegates(t *testing.T) {
	stub := &inspectorStub{}
	result, err := NewInspectUseCase(stub).Execute(context.Background(), InspectRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || result.Number != 12 || stub.request.Number != 12 {
		t.Fatalf("unexpected result: %+v request=%+v err=%v", result, stub.request, err)
	}
}

func TestInspectUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []InspectRequest{{Name: "project", Number: 1}, {Owner: "alice", Name: "project", Number: 0}} {
		if _, err := NewInspectUseCase(&inspectorStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}
