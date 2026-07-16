package pullrequest

import (
	"context"
	"testing"
)

type statusViewerStub struct{ request StatusRequest }

func (stub *statusViewerStub) ViewStatus(_ context.Context, request StatusRequest) (StatusData, error) {
	stub.request = request
	return StatusData{Number: request.Number, ReviewsAvailable: true, Reviews: []Review{{ID: 1, ReviewerID: 1, State: "APPROVED"}}, ChecksAvailable: true, Checks: []string{"pending"}, Mergeable: MergeableYes}, nil
}

func TestAggregateReviewsUsesLatestEffectiveReviewPerReviewer(t *testing.T) {
	tests := []struct {
		name    string
		reviews []Review
		want    AggregateState
	}{
		{name: "request changes then approve", reviews: []Review{{ID: 1, ReviewerID: 10, State: "REQUEST_CHANGES"}, {ID: 2, ReviewerID: 10, State: "APPROVED"}}, want: AggregateSuccess},
		{name: "approve then request changes", reviews: []Review{{ID: 1, ReviewerID: 10, State: "APPROVED"}, {ID: 2, ReviewerID: 10, State: "REQUEST_CHANGES"}}, want: AggregateFailed},
		{name: "different reviewers remain effective", reviews: []Review{{ID: 2, ReviewerID: 10, State: "APPROVED"}, {ID: 3, ReviewerID: 20, State: "REQUEST_CHANGES"}}, want: AggregateFailed},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := aggregateReviews(test.reviews, true, 0); got != test.want {
				t.Fatalf("aggregate = %s, want %s", got, test.want)
			}
		})
	}
}

func TestAggregateReviewsUnavailable(t *testing.T) {
	if got := aggregateReviews(nil, true, 0); got != AggregateUnavailable {
		t.Fatalf("empty aggregate = %s", got)
	}
	if got := aggregateReviews([]Review{{ID: 1, ReviewerID: 1, State: "APPROVED"}}, false, 0); got != AggregateUnavailable {
		t.Fatalf("unavailable aggregate = %s", got)
	}
}

func TestAggregateReviewsPreservesPendingAndConservativeStates(t *testing.T) {
	tests := []struct {
		name      string
		reviews   []Review
		requested int
		want      AggregateState
	}{
		{name: "requested reviewer pending", reviews: []Review{{ID: 1, ReviewerID: 1, State: "APPROVED"}}, requested: 1, want: AggregatePending},
		{name: "unknown never succeeds", reviews: []Review{{ID: 1, ReviewerID: 1, State: "APPROVED"}, {ID: 2, ReviewerID: 2, State: "UNKNOWN"}}, want: AggregateUnavailable},
		{name: "dismissed ignored", reviews: []Review{{ID: 1, ReviewerID: 1, State: "REQUEST_CHANGES", Dismissed: true}}, want: AggregateUnavailable},
		{name: "stale ignored", reviews: []Review{{ID: 1, ReviewerID: 1, State: "APPROVED", Stale: true}}, want: AggregateUnavailable},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := aggregateReviews(test.reviews, true, test.requested); got != test.want {
				t.Fatalf("aggregate = %s, want %s", got, test.want)
			}
		})
	}
}

func TestAggregateChecks(t *testing.T) {
	tests := []struct {
		name   string
		checks []string
		want   AggregateState
	}{
		{name: "success", checks: []string{"success", "success"}, want: AggregateSuccess},
		{name: "pending", checks: []string{"success", "pending"}, want: AggregatePending},
		{name: "failed", checks: []string{"pending", "failure"}, want: AggregateFailed},
		{name: "warning fails", checks: []string{"warning"}, want: AggregateFailed},
		{name: "unknown never succeeds", checks: []string{"success", "unknown"}, want: AggregateUnavailable},
		{name: "empty", want: AggregateUnavailable},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := aggregateChecks(test.checks, true); got != test.want {
				t.Fatalf("aggregate = %s, want %s", got, test.want)
			}
		})
	}
}

func TestStatusUseCase(t *testing.T) {
	viewer := &statusViewerStub{}
	request := StatusRequest{Owner: "alice", Name: "project", Number: 12}
	result, err := NewStatusUseCase(viewer).Execute(context.Background(), request)
	if err != nil || result.Number != 12 || viewer.request != request {
		t.Fatalf("unexpected result: %+v request=%+v err=%v", result, viewer.request, err)
	}
}

func TestStatusUseCaseRejectsInvalidInput(t *testing.T) {
	tests := []StatusRequest{{Name: "project", Number: 1}, {Owner: "alice", Number: 1}, {Owner: "alice", Name: "project"}}
	for _, request := range tests {
		if _, err := NewStatusUseCase(&statusViewerStub{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("expected validation error for %+v", request)
		}
	}
}
