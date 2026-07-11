package repository

import (
	"context"
	"testing"
)

type updaterStub struct{ request UpdateRequest }

func (stub *updaterStub) Update(_ context.Context, request UpdateRequest) (Repository, error) {
	stub.request = request
	return Repository{Owner: request.Owner, Name: request.Name}, nil
}

func TestUpdateUseCaseDelegatesExplicitFields(t *testing.T) {
	description := ""
	private := false
	stub := &updaterStub{}
	if _, err := NewUpdateUseCase(stub).Execute(context.Background(), UpdateRequest{Owner: "alice", Name: "project", Description: &description, Private: &private}); err != nil {
		t.Fatal(err)
	}
	if stub.request.Description == nil || *stub.request.Description != "" || stub.request.Private == nil || *stub.request.Private {
		t.Fatalf("request = %+v", stub.request)
	}
}

func TestUpdateUseCaseRejectsInvalidRequests(t *testing.T) {
	for _, request := range []UpdateRequest{{Name: "project", Description: stringPtr("x")}, {Owner: "alice", Description: nil}} {
		if _, err := NewUpdateUseCase(&updaterStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("request %+v accepted", request)
		}
	}
}

func stringPtr(value string) *string { return &value }
