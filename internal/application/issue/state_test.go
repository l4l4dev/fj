package issue

import (
	"context"
	"testing"
)

type stateChangerStub struct{ request ChangeStateRequest }

func (stub *stateChangerStub) ChangeState(_ context.Context, request ChangeStateRequest) (IssueDetail, error) {
	stub.request = request
	return IssueDetail{Number: request.Number, State: request.State}, nil
}

func TestChangeStateUseCaseDelegates(t *testing.T) {
	stub := &stateChangerStub{}
	result, err := NewChangeStateUseCase(stub).Execute(context.Background(), ChangeStateRequest{Owner: "alice", Name: "project", Number: 12, State: StateClosed})
	if err != nil || result.State != StateClosed || stub.request.State != StateClosed {
		t.Fatalf("unexpected result: %+v request=%+v err=%v", result, stub.request, err)
	}
}

func TestChangeStateUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []ChangeStateRequest{{Name: "project", Number: 1, State: StateOpen}, {Owner: "alice", Number: 1, State: StateOpen}, {Owner: "alice", Name: "project", Number: 0, State: StateOpen}, {Owner: "alice", Name: "project", Number: 1, State: State("all")}} {
		if _, err := NewChangeStateUseCase(&stateChangerStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}
