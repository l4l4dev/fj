package repository

import (
	"context"
	"testing"
)

type archiverStub struct{ request ArchiveRequest }

func (stub *archiverStub) SetArchived(_ context.Context, request ArchiveRequest) (Repository, error) {
	stub.request = request
	return Repository{Owner: request.Owner, Name: request.Name, Archived: request.Archived}, nil
}

func TestArchiveUseCaseDelegatesArchiveState(t *testing.T) {
	stub := &archiverStub{}
	result, err := NewArchiveUseCase(stub).Execute(context.Background(), ArchiveRequest{Owner: "alice", Name: "project", Archived: true})
	if err != nil {
		t.Fatal(err)
	}
	if !stub.request.Archived || !result.Archived {
		t.Fatalf("request/result = %+v/%+v", stub.request, result)
	}
}
