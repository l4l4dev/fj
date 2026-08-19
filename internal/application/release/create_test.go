package release

import (
	"context"
	"testing"
)

type creatorStub struct{ request CreateRequest }

func (s *creatorStub) Create(_ context.Context, request CreateRequest) (ReleaseDetail, error) {
	s.request = request
	return ReleaseDetail{TagName: request.Tag, Title: request.Title, Draft: true}, nil
}

func TestCreateUseCaseDelegates(t *testing.T) {
	stub := &creatorStub{}
	result, err := NewCreateUseCase(stub).Execute(context.Background(), CreateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: "First release", Notes: "Notes", Prerelease: true})
	if err != nil || result.TagName != "v1.0.0" || stub.request.Owner != "alice" || stub.request.Name != "project" || stub.request.Prerelease != true {
		t.Fatalf("unexpected result: %+v request=%+v err=%v", result, stub.request, err)
	}
}

func TestCreateUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []CreateRequest{
		{Name: "project", Tag: "v1.0.0", Title: "First release"},
		{Owner: "alice", Tag: "v1.0.0", Title: "First release"},
		{Owner: "alice", Name: "project", Title: "First release"},
		{Owner: "alice", Name: "project", Tag: "  ", Title: "First release"},
		{Owner: "alice", Name: "project", Tag: "v 1", Title: "First release"},
		{Owner: "alice", Name: "project", Tag: "v1.0.0"},
		{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: "  "},
	} {
		if _, err := NewCreateUseCase(&creatorStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}

func TestCreateUseCaseRejectsNilCreator(t *testing.T) {
	if _, err := NewCreateUseCase(nil).Execute(context.Background(), CreateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: "First release"}); err == nil {
		t.Fatal("expected internal error for nil creator")
	}
}
