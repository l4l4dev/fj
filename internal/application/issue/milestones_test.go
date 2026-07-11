package issue

import (
	"context"
	"testing"
)

type milestoneSetterStub struct{ request SetMilestoneRequest }

func (stub *milestoneSetterStub) SetMilestone(_ context.Context, request SetMilestoneRequest) (Milestone, error) {
	stub.request = request
	return Milestone{ID: 1, Title: request.Milestone}, nil
}

type milestoneClearerStub struct{ request ClearMilestoneRequest }

func (stub *milestoneClearerStub) ClearMilestone(_ context.Context, request ClearMilestoneRequest) error {
	stub.request = request
	return nil
}

func TestMilestoneUseCasesDelegate(t *testing.T) {
	setter := &milestoneSetterStub{}
	if _, err := NewSetMilestoneUseCase(setter).Execute(context.Background(), SetMilestoneRequest{Owner: "alice", Name: "project", Number: 1, Milestone: "v1"}); err != nil || setter.request.Milestone != "v1" {
		t.Fatalf("unexpected set delegation: %+v err=%v", setter.request, err)
	}
	clearer := &milestoneClearerStub{}
	if err := NewClearMilestoneUseCase(clearer).Execute(context.Background(), ClearMilestoneRequest{Owner: "alice", Name: "project", Number: 1}); err != nil || clearer.request.Number != 1 {
		t.Fatalf("unexpected clear delegation: %+v err=%v", clearer.request, err)
	}
}

func TestSetMilestoneRejectsInvalidInput(t *testing.T) {
	for _, title := range []string{"", "  "} {
		if _, err := NewSetMilestoneUseCase(&milestoneSetterStub{}).Execute(context.Background(), SetMilestoneRequest{Owner: "alice", Name: "project", Number: 1, Milestone: title}); err == nil {
			t.Fatalf("expected validation error for %q", title)
		}
	}
}
