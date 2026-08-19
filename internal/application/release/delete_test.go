package release

import (
	"context"
	"errors"
	"testing"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type deleteInspectorStub struct {
	request InspectRequest
	detail  ReleaseDetail
	err     error
}

func (s *deleteInspectorStub) Inspect(_ context.Context, request InspectRequest) (ReleaseDetail, error) {
	s.request = request
	return s.detail, s.err
}

type deleterStub struct {
	called  bool
	request DeleteRequest
	err     error
}

func (s *deleterStub) Delete(_ context.Context, request DeleteRequest) error {
	s.called = true
	s.request = request
	return s.err
}

func TestDeleteUseCaseDeletesResolvedRelease(t *testing.T) {
	inspector := &deleteInspectorStub{detail: ReleaseDetail{ID: 7, TagName: "v1.0.0", Title: "First release"}}
	deleter := &deleterStub{}
	result, err := NewDeleteUseCase(inspector, deleter).Execute(context.Background(), InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
	if err != nil || result.TagName != "v1.0.0" || result.Title != "First release" {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if !deleter.called || deleter.request.ID != 7 || deleter.request.Owner != "alice" || deleter.request.Name != "project" {
		t.Fatalf("unexpected delete request: %+v", deleter.request)
	}
}

func TestDeleteUseCasePropagatesInspectError(t *testing.T) {
	inspector := &deleteInspectorStub{err: apperror.New(apperror.NotFound, "inspect release", "release not found")}
	deleter := &deleterStub{}
	_, err := NewDeleteUseCase(inspector, deleter).Execute(context.Background(), InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound {
		t.Fatalf("unexpected error: %#v", err)
	}
	if deleter.called {
		t.Fatal("deleter must not be called when inspect fails")
	}
}

func TestDeleteUseCasePropagatesDeleteError(t *testing.T) {
	inspector := &deleteInspectorStub{detail: ReleaseDetail{ID: 7, TagName: "v1.0.0"}}
	deleter := &deleterStub{err: apperror.New(apperror.Conflict, "delete release", "release could not be deleted in its current state")}
	_, err := NewDeleteUseCase(inspector, deleter).Execute(context.Background(), InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.Conflict {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestDeleteUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []InspectRequest{
		{Name: "project", Tag: "v1.0.0"},
		{Owner: "alice", Tag: "v1.0.0"},
		{Owner: "alice", Name: "project"},
		{Owner: "alice", Name: "project", Tag: "  "},
	} {
		if _, err := NewDeleteUseCase(&deleteInspectorStub{}, &deleterStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}

func TestDeleteUseCaseRejectsNilDependencies(t *testing.T) {
	request := InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"}
	if _, err := NewDeleteUseCase(nil, &deleterStub{}).Execute(context.Background(), request); err == nil {
		t.Fatal("expected internal error for nil inspector")
	}
	if _, err := NewDeleteUseCase(&deleteInspectorStub{}, nil).Execute(context.Background(), request); err == nil {
		t.Fatal("expected internal error for nil deleter")
	}
}
