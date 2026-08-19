package pullrequest

import (
	"context"
	"errors"
	"testing"
)

type mergerStub struct {
	request MergeRequest
	err     error
}

func (stub *mergerStub) Merge(_ context.Context, request MergeRequest) error {
	stub.request = request
	return stub.err
}

type closerStub struct {
	request CloseRequest
	detail  PullRequestDetail
}

func (stub *closerStub) Close(_ context.Context, request CloseRequest) (PullRequestDetail, error) {
	stub.request = request
	return stub.detail, nil
}

func TestMergeUseCaseForwardsExplicitMethod(t *testing.T) {
	stub := &mergerStub{}
	err := NewMergeUseCase(stub).Execute(context.Background(), MergeRequest{Owner: "alice", Name: "project", Number: 12, Method: MergeMethodSquash})
	if err != nil || stub.request.Method != MergeMethodSquash || stub.request.Number != 12 {
		t.Fatalf("unexpected request: %+v err=%v", stub.request, err)
	}
}

func TestMergeUseCaseRejectsInvalidInput(t *testing.T) {
	tests := []MergeRequest{
		{Owner: "", Name: "project", Number: 1, Method: MergeMethodMerge},
		{Owner: "alice", Name: "project", Number: 0, Method: MergeMethodMerge},
		{Owner: "alice", Name: "project", Number: 1},
		{Owner: "alice", Name: "project", Number: 1, Method: "fast-forward"},
	}
	for _, request := range tests {
		if err := NewMergeUseCase(&mergerStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}

func TestMergeUseCasePropagatesRemoteFailure(t *testing.T) {
	remoteFailure := errors.New("not mergeable")
	err := NewMergeUseCase(&mergerStub{err: remoteFailure}).Execute(context.Background(), MergeRequest{Owner: "alice", Name: "project", Number: 12, Method: MergeMethodMerge})
	if !errors.Is(err, remoteFailure) {
		t.Fatalf("expected remote failure to propagate, got %v", err)
	}
}

func TestCloseUseCaseForwardsTarget(t *testing.T) {
	stub := &closerStub{detail: PullRequestDetail{Number: 12, State: StateClosed}}
	detail, err := NewCloseUseCase(stub).Execute(context.Background(), CloseRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || detail.State != StateClosed || stub.request.Number != 12 {
		t.Fatalf("unexpected result: %+v err=%v", detail, err)
	}
}

func TestCompletionUseCasesRequireDependencies(t *testing.T) {
	if err := NewMergeUseCase(nil).Execute(context.Background(), MergeRequest{Owner: "alice", Name: "project", Number: 1, Method: MergeMethodMerge}); err == nil {
		t.Fatal("expected error for nil merger")
	}
	if _, err := NewCloseUseCase(nil).Execute(context.Background(), CloseRequest{Owner: "alice", Name: "project", Number: 1}); err == nil {
		t.Fatal("expected error for nil closer")
	}
}
