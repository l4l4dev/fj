package release

import (
	"context"
	"errors"
	"testing"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type publishInspectorStub struct {
	request InspectRequest
	detail  ReleaseDetail
	err     error
}

func (s *publishInspectorStub) Inspect(_ context.Context, request InspectRequest) (ReleaseDetail, error) {
	s.request = request
	return s.detail, s.err
}

type publisherStub struct {
	called  bool
	request PublishRequest
	detail  ReleaseDetail
}

func (s *publisherStub) Publish(_ context.Context, request PublishRequest) (ReleaseDetail, error) {
	s.called = true
	s.request = request
	return s.detail, nil
}

func TestPublishUseCasePublishesDraft(t *testing.T) {
	inspector := &publishInspectorStub{detail: ReleaseDetail{ID: 7, TagName: "v1.0.0", Draft: true}}
	publisher := &publisherStub{detail: ReleaseDetail{ID: 7, TagName: "v1.0.0", Draft: false}}
	result, err := NewPublishUseCase(inspector, publisher).Execute(context.Background(), InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
	if err != nil || result.Draft {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if !publisher.called || publisher.request.ID != 7 || publisher.request.Owner != "alice" || publisher.request.Name != "project" {
		t.Fatalf("unexpected publish request: %+v", publisher.request)
	}
}

func TestPublishUseCaseRejectsAlreadyPublished(t *testing.T) {
	inspector := &publishInspectorStub{detail: ReleaseDetail{ID: 7, TagName: "v1.0.0", Draft: false}}
	publisher := &publisherStub{}
	_, err := NewPublishUseCase(inspector, publisher).Execute(context.Background(), InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.Conflict {
		t.Fatalf("unexpected error: %#v", err)
	}
	if publisher.called {
		t.Fatal("publisher must not be called when the release is already published")
	}
}

func TestPublishUseCasePropagatesInspectError(t *testing.T) {
	inspector := &publishInspectorStub{err: apperror.New(apperror.NotFound, "inspect release", "release not found")}
	publisher := &publisherStub{}
	_, err := NewPublishUseCase(inspector, publisher).Execute(context.Background(), InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound {
		t.Fatalf("unexpected error: %#v", err)
	}
	if publisher.called {
		t.Fatal("publisher must not be called when inspect fails")
	}
}

func TestPublishUseCaseRejectsInvalidInput(t *testing.T) {
	for _, request := range []InspectRequest{
		{Name: "project", Tag: "v1.0.0"},
		{Owner: "alice", Tag: "v1.0.0"},
		{Owner: "alice", Name: "project"},
		{Owner: "alice", Name: "project", Tag: "  "},
	} {
		if _, err := NewPublishUseCase(&publishInspectorStub{}, &publisherStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}

func TestPublishUseCaseRejectsNilDependencies(t *testing.T) {
	request := InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"}
	if _, err := NewPublishUseCase(nil, &publisherStub{}).Execute(context.Background(), request); err == nil {
		t.Fatal("expected internal error for nil inspector")
	}
	if _, err := NewPublishUseCase(&publishInspectorStub{}, nil).Execute(context.Background(), request); err == nil {
		t.Fatal("expected internal error for nil publisher")
	}
}
