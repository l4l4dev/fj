package issue

import (
	"context"
	"testing"
)

type updaterStub struct{ request UpdateRequest }

func (stub *updaterStub) Update(_ context.Context, request UpdateRequest) (IssueDetail, error) {
	stub.request = request
	return IssueDetail{Number: request.Number, Title: valueOrEmpty(request.Title), Body: valueOrEmpty(request.Body)}, nil
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func TestUpdateUseCaseDelegatesExplicitFields(t *testing.T) {
	stub := &updaterStub{}
	title := "Updated"
	result, err := NewUpdateUseCase(stub).Execute(context.Background(), UpdateRequest{Owner: "alice", Name: "project", Number: 12, Title: &title})
	if err != nil || result.Title != title || stub.request.Body != nil {
		t.Fatalf("unexpected result: %+v request=%+v err=%v", result, stub.request, err)
	}
}

func TestUpdateUseCaseAllowsEmptyBody(t *testing.T) {
	stub := &updaterStub{}
	body := ""
	if _, err := NewUpdateUseCase(stub).Execute(context.Background(), UpdateRequest{Owner: "alice", Name: "project", Number: 12, Body: &body}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateUseCaseRejectsInvalidInput(t *testing.T) {
	title := "title"
	for _, request := range []UpdateRequest{{Name: "project", Number: 1, Title: &title}, {Owner: "alice", Number: 1, Title: &title}, {Owner: "alice", Name: "project", Number: 0, Title: &title}, {Owner: "alice", Name: "project", Number: 1}} {
		if _, err := NewUpdateUseCase(&updaterStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}
