package issue

import (
	"context"
	"testing"
)

type labelAdderStub struct{ request AddLabelRequest }

func (stub *labelAdderStub) AddLabel(_ context.Context, request AddLabelRequest) (Label, error) {
	stub.request = request
	return Label{ID: 1, Name: request.Label}, nil
}

type labelRemoverStub struct{ request RemoveLabelRequest }

func (stub *labelRemoverStub) RemoveLabel(_ context.Context, request RemoveLabelRequest) (Label, error) {
	stub.request = request
	return Label{ID: 1, Name: request.Label}, nil
}

func TestLabelUseCasesDelegate(t *testing.T) {
	adder := &labelAdderStub{}
	if _, err := NewAddLabelUseCase(adder).Execute(context.Background(), AddLabelRequest{Owner: "alice", Name: "project", Number: 1, Label: "bug"}); err != nil || adder.request.Label != "bug" {
		t.Fatalf("unexpected add delegation: request=%+v err=%v", adder.request, err)
	}
	remover := &labelRemoverStub{}
	if _, err := NewRemoveLabelUseCase(remover).Execute(context.Background(), RemoveLabelRequest{Owner: "alice", Name: "project", Number: 1, Label: "bug"}); err != nil || remover.request.Label != "bug" {
		t.Fatalf("unexpected remove delegation: request=%+v err=%v", remover.request, err)
	}
}

func TestLabelUseCasesRejectInvalidInput(t *testing.T) {
	for _, label := range []string{"", "  "} {
		if _, err := NewAddLabelUseCase(&labelAdderStub{}).Execute(context.Background(), AddLabelRequest{Owner: "alice", Name: "project", Number: 1, Label: label}); err == nil {
			t.Fatalf("expected validation error for %q", label)
		}
	}
}
