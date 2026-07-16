package pullrequest

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type StatusUseCase struct{ viewer StatusViewer }

func NewStatusUseCase(viewer StatusViewer) StatusUseCase {
	return StatusUseCase{viewer: viewer}
}

func (useCase StatusUseCase) Execute(ctx context.Context, request StatusRequest) (PullRequestStatus, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return PullRequestStatus{}, apperror.NewValidation("view pull request status", "OWNER/NAME owner and name are required")
	}
	if request.Number < 1 {
		return PullRequestStatus{}, apperror.NewValidation("view pull request status", "pull request number must be a positive integer")
	}
	if useCase.viewer == nil {
		return PullRequestStatus{}, apperror.New(apperror.Internal, "view pull request status", "")
	}
	data, err := useCase.viewer.ViewStatus(ctx, request)
	if err != nil {
		return PullRequestStatus{}, err
	}
	return PullRequestStatus{
		Number:    data.Number,
		Review:    aggregateReviews(data.Reviews, data.ReviewsAvailable, data.RequestedReviewers),
		Check:     aggregateChecks(data.Checks, data.ChecksAvailable),
		Mergeable: data.Mergeable,
	}, nil
}

func aggregateReviews(reviews []Review, available bool, requestedReviewers int) AggregateState {
	if !available {
		return AggregateUnavailable
	}

	latest := make(map[int64]Review)
	for _, review := range reviews {
		if review.Dismissed || review.Stale {
			continue
		}
		if review.ID < 1 || review.ReviewerID < 1 {
			return AggregateUnavailable
		}
		current, exists := latest[review.ReviewerID]
		if !exists || review.ID > current.ID {
			latest[review.ReviewerID] = review
		}
	}

	approved, failed, pending, unknown := false, false, requestedReviewers > 0, false
	for _, review := range latest {
		switch strings.ToUpper(strings.TrimSpace(review.State)) {
		case "APPROVED":
			approved = true
		case "REQUEST_CHANGES":
			failed = true
		case "PENDING", "REQUEST_REVIEW":
			pending = true
		case "COMMENT":
		default:
			unknown = true
		}
	}
	if failed {
		return AggregateFailed
	}
	if pending {
		return AggregatePending
	}
	if unknown {
		return AggregateUnavailable
	}
	if approved {
		return AggregateSuccess
	}
	return AggregateUnavailable
}

func aggregateChecks(checks []string, available bool) AggregateState {
	if !available || len(checks) == 0 {
		return AggregateUnavailable
	}
	pending, failed, unknown := false, false, false
	for _, check := range checks {
		switch strings.ToLower(strings.TrimSpace(check)) {
		case "success":
		case "pending":
			pending = true
		case "failure", "error", "warning":
			failed = true
		default:
			unknown = true
		}
	}
	if failed {
		return AggregateFailed
	}
	if pending {
		return AggregatePending
	}
	if unknown {
		return AggregateUnavailable
	}
	return AggregateSuccess
}
