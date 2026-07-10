package repository

import (
	"context"
	"testing"
)

type creatorStub struct{ request CreateRequest }

func (stub *creatorStub) Create(_ context.Context, request CreateRequest) (Repository, error) {
	stub.request = request
	return Repository{Owner: "alice", Name: request.Name, Private: request.Private}, nil
}

func TestCreateUseCaseValidatesAndDelegates(t *testing.T) {
	stub := &creatorStub{}
	result, err := NewCreateUseCase(stub).Execute(context.Background(), CreateRequest{Name: "project", Description: "demo", Private: true})
	if err != nil {
		t.Fatal(err)
	}
	if stub.request.Name != "project" || !stub.request.Private || stub.request.Description != "demo" {
		t.Fatalf("request = %+v", stub.request)
	}
	if result.Owner != "alice" || result.Name != "project" {
		t.Fatalf("result = %+v", result)
	}
}

func TestCreateUseCaseRejectsBlankName(t *testing.T) {
	for _, name := range []string{"", "   "} {
		if _, err := NewCreateUseCase(&creatorStub{}).Execute(context.Background(), CreateRequest{Name: name}); err == nil {
			t.Fatalf("name %q accepted", name)
		}
	}
}
