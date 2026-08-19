package release

import (
	"context"
	"testing"
)

type updaterStub struct{ request UpdateRequest }

func (s *updaterStub) Update(_ context.Context, request UpdateRequest) (ReleaseDetail, error) {
	s.request = request
	return ReleaseDetail{TagName: request.Tag}, nil
}

func TestUpdateUseCaseForwardsOnlySuppliedFields(t *testing.T) {
	stub := &updaterStub{}
	title := "New title"
	result, err := NewUpdateUseCase(stub).Execute(context.Background(), UpdateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: &title})
	if err != nil || result.TagName != "v1.0.0" {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if stub.request.Title == nil || *stub.request.Title != "New title" {
		t.Fatalf("unexpected title: %+v", stub.request)
	}
	if stub.request.Notes != nil || stub.request.Prerelease != nil {
		t.Fatalf("unexpected non-nil fields: %+v", stub.request)
	}
}

func TestUpdateUseCaseRejectsInvalidInput(t *testing.T) {
	title := "New title"
	blank := "  "
	for _, request := range []UpdateRequest{
		{Name: "project", Tag: "v1.0.0", Title: &title},
		{Owner: "alice", Tag: "v1.0.0", Title: &title},
		{Owner: "alice", Name: "project", Title: &title},
		{Owner: "alice", Name: "project", Tag: "  ", Title: &title},
		{Owner: "alice", Name: "project", Tag: "v1.0.0"},
		{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: &blank},
	} {
		if _, err := NewUpdateUseCase(&updaterStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}

func TestUpdateUseCaseAcceptsPrereleaseOnly(t *testing.T) {
	stub := &updaterStub{}
	prerelease := false
	if _, err := NewUpdateUseCase(stub).Execute(context.Background(), UpdateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Prerelease: &prerelease}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.request.Prerelease == nil || *stub.request.Prerelease != false {
		t.Fatalf("unexpected prerelease: %+v", stub.request)
	}
}

func TestUpdateUseCaseRejectsNilUpdater(t *testing.T) {
	title := "New title"
	if _, err := NewUpdateUseCase(nil).Execute(context.Background(), UpdateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: &title}); err == nil {
		t.Fatal("expected internal error for nil updater")
	}
}
