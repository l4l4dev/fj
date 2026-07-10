package repository

import (
	"context"
	"testing"
)

type getterStub struct {
	request GetRequest
}

func (stub *getterStub) Get(_ context.Context, request GetRequest) (Repository, error) {
	stub.request = request
	return Repository{Owner: request.Owner, Name: request.Name, DefaultBranch: "main"}, nil
}

func TestInspectUseCaseValidatesAndDelegates(t *testing.T) {
	stub := &getterStub{}
	repository, err := NewInspectUseCase(stub).Execute(context.Background(), GetRequest{Owner: "alice", Name: "project"})
	if err != nil {
		t.Fatal(err)
	}
	if stub.request != (GetRequest{Owner: "alice", Name: "project"}) {
		t.Fatalf("request = %+v", stub.request)
	}
	if repository.Owner != "alice" || repository.Name != "project" {
		t.Fatalf("repository = %+v", repository)
	}
}

func TestInspectUseCaseRejectsMissingTargetParts(t *testing.T) {
	for _, request := range []GetRequest{{Name: "project"}, {Owner: "alice"}} {
		if _, err := NewInspectUseCase(&getterStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("request %+v was accepted", request)
		}
	}
}
