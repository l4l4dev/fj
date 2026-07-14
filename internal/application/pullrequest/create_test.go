package pullrequest

import (
	"context"
	"testing"
)

type creatorStub struct{ request CreateRequest }

func (stub *creatorStub) Create(_ context.Context, request CreateRequest) (PullRequestDetail, error) {
	stub.request = request
	return PullRequestDetail{Number: 7, Title: request.Title, HeadBranch: request.HeadBranch, BaseBranch: request.BaseBranch}, nil
}

func TestCreateUseCase(t *testing.T) {
	creator := &creatorStub{}
	request := CreateRequest{Owner: "alice", Name: "project", Title: "Improve flow", HeadBranch: "feature", BaseBranch: "main"}
	result, err := NewCreateUseCase(creator).Execute(context.Background(), request)
	if err != nil || result.Number != 7 || creator.request != request {
		t.Fatalf("unexpected result: %+v request=%+v err=%v", result, creator.request, err)
	}
}

func TestCreateUseCaseRejectsInvalidInput(t *testing.T) {
	valid := CreateRequest{Owner: "alice", Name: "project", Title: "Improve flow", HeadBranch: "feature", BaseBranch: "main"}
	tests := []CreateRequest{
		{Name: valid.Name, Title: valid.Title, HeadBranch: valid.HeadBranch, BaseBranch: valid.BaseBranch},
		{Owner: valid.Owner, Name: valid.Name, HeadBranch: valid.HeadBranch, BaseBranch: valid.BaseBranch},
		{Owner: valid.Owner, Name: valid.Name, Title: valid.Title, BaseBranch: valid.BaseBranch},
		{Owner: valid.Owner, Name: valid.Name, Title: valid.Title, HeadBranch: valid.HeadBranch},
		{Owner: valid.Owner, Name: valid.Name, Title: valid.Title, HeadBranch: "main", BaseBranch: "main"},
		{Owner: valid.Owner, Name: valid.Name, Title: valid.Title, HeadBranch: "alice:feature", BaseBranch: valid.BaseBranch},
	}
	for _, request := range tests {
		if _, err := NewCreateUseCase(&creatorStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}
