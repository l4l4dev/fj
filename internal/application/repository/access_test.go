package repository

import (
	"context"
	"testing"
)

type accessStub struct{}

func (accessStub) ViewAccess(_ context.Context, request AccessRequest) (RepositoryAccess, error) {
	return RepositoryAccess{Owner: request.Owner, Name: request.Name, Collaborators: []Collaborator{{Username: "bob", Permission: PermissionAdmin}}}, nil
}
func TestAccessUseCase(t *testing.T) {
	result, err := NewAccessUseCase(accessStub{}).Execute(context.Background(), AccessRequest{Owner: "alice", Name: "project"})
	if err != nil || len(result.Collaborators) != 1 || result.Collaborators[0].Permission != PermissionAdmin {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}
