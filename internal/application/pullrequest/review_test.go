package pullrequest

import (
	"context"
	"testing"
)

type reviewSubmitterStub struct{ submission ReviewSubmission }

func (stub *reviewSubmitterStub) SubmitReview(_ context.Context, submission ReviewSubmission) (SubmittedReview, error) {
	stub.submission = submission
	return SubmittedReview{State: string(submission.Event)}, nil
}

func TestSubmitReviewMapsOutcomes(t *testing.T) {
	tests := []struct {
		outcome ReviewOutcome
		body    string
		want    ReviewEvent
	}{
		{outcome: ReviewOutcomeComment, body: "Looks good overall", want: ReviewEventComment},
		{outcome: ReviewOutcomeApprove, want: ReviewEventApprove},
		{outcome: ReviewOutcomeRequestChanges, body: "Please add a test", want: ReviewEventRequestChanges},
	}
	for _, test := range tests {
		stub := &reviewSubmitterStub{}
		result, err := NewSubmitReviewUseCase(stub).Execute(context.Background(), SubmitReviewRequest{Owner: "alice", Name: "project", Number: 12, Outcome: test.outcome, Body: test.body})
		if err != nil || stub.submission.Event != test.want || result.State != string(test.want) {
			t.Fatalf("outcome %q: result=%+v submission=%+v err=%v", test.outcome, result, stub.submission, err)
		}
	}
}

func TestSubmitReviewRejectsUnsupportedOutcome(t *testing.T) {
	for _, outcome := range []ReviewOutcome{"", "approved", "APPROVE"} {
		_, err := NewSubmitReviewUseCase(&reviewSubmitterStub{}).Execute(context.Background(), SubmitReviewRequest{Owner: "alice", Name: "project", Number: 12, Outcome: outcome, Body: "body"})
		if err == nil {
			t.Fatalf("expected validation error for outcome %q", outcome)
		}
	}
}

func TestSubmitReviewValidatesBodyByOutcome(t *testing.T) {
	for _, outcome := range []ReviewOutcome{ReviewOutcomeComment, ReviewOutcomeRequestChanges} {
		for _, body := range []string{"", " \t\n"} {
			_, err := NewSubmitReviewUseCase(&reviewSubmitterStub{}).Execute(context.Background(), SubmitReviewRequest{Owner: "alice", Name: "project", Number: 12, Outcome: outcome, Body: body})
			if err == nil {
				t.Fatalf("expected validation error for outcome %q and body %q", outcome, body)
			}
		}
	}
	if _, err := NewSubmitReviewUseCase(&reviewSubmitterStub{}).Execute(context.Background(), SubmitReviewRequest{Owner: "alice", Name: "project", Number: 12, Outcome: ReviewOutcomeApprove}); err != nil {
		t.Fatalf("approve without body returned error: %v", err)
	}
}

func TestSubmitReviewValidatesTarget(t *testing.T) {
	requests := []SubmitReviewRequest{
		{Name: "project", Number: 12, Outcome: ReviewOutcomeApprove},
		{Owner: "alice", Number: 12, Outcome: ReviewOutcomeApprove},
		{Owner: "alice", Name: "project", Number: 0, Outcome: ReviewOutcomeApprove},
	}
	for _, request := range requests {
		if _, err := NewSubmitReviewUseCase(&reviewSubmitterStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for request %+v", request)
		}
	}
}
