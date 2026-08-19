package release

import (
	"context"
	"testing"
)

type inspectorStub struct{ request InspectRequest }

func (s *inspectorStub) Inspect(_ context.Context, request InspectRequest) (ReleaseDetail, error) {
	s.request = request
	return ReleaseDetail{ID: 1, TagName: request.Tag}, nil
}

func TestInspectUseCaseDelegates(t *testing.T) {
	stub := &inspectorStub{}
	result, err := NewInspectUseCase(stub).Execute(context.Background(), InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
	if err != nil || result.TagName != "v1.0.0" || stub.request.Owner != "alice" || stub.request.Name != "project" || stub.request.Tag != "v1.0.0" {
		t.Fatalf("unexpected result: %+v request=%+v err=%v", result, stub.request, err)
	}
}

func TestInspectUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []InspectRequest{
		{Name: "project", Tag: "v1.0.0"},
		{Owner: "alice", Tag: "v1.0.0"},
		{Owner: "alice", Name: "project"},
		{Owner: "alice", Name: "project", Tag: "  "},
	} {
		if _, err := NewInspectUseCase(&inspectorStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}

func TestInspectUseCaseRejectsNilInspector(t *testing.T) {
	if _, err := NewInspectUseCase(nil).Execute(context.Background(), InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"}); err == nil {
		t.Fatal("expected internal error for nil inspector")
	}
}
