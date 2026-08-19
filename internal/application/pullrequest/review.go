package pullrequest

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type SubmitReviewUseCase struct{ submitter ReviewSubmitter }

func NewSubmitReviewUseCase(submitter ReviewSubmitter) SubmitReviewUseCase {
	return SubmitReviewUseCase{submitter: submitter}
}

func (useCase SubmitReviewUseCase) Execute(ctx context.Context, request SubmitReviewRequest) (SubmittedReview, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return SubmittedReview{}, apperror.NewValidation("submit pull request review", "OWNER/NAME owner and name are required")
	}
	if request.Number < 1 {
		return SubmittedReview{}, apperror.NewValidation("submit pull request review", "pull request number must be a positive integer")
	}
	event, ok := reviewEvent(request.Outcome)
	if !ok {
		return SubmittedReview{}, apperror.NewValidation("submit pull request review", "outcome must be comment, approve, or request-changes")
	}
	if (request.Outcome == ReviewOutcomeComment || request.Outcome == ReviewOutcomeRequestChanges) && strings.TrimSpace(request.Body) == "" {
		return SubmittedReview{}, apperror.NewValidation("submit pull request review", "review body is required for comment and request-changes outcomes")
	}
	if useCase.submitter == nil {
		return SubmittedReview{}, apperror.New(apperror.Internal, "submit pull request review", "")
	}
	return useCase.submitter.SubmitReview(ctx, ReviewSubmission{
		Owner: request.Owner, Name: request.Name, Number: request.Number, Event: event, Body: request.Body,
	})
}

func reviewEvent(outcome ReviewOutcome) (ReviewEvent, bool) {
	switch outcome {
	case ReviewOutcomeComment:
		return ReviewEventComment, true
	case ReviewOutcomeApprove:
		return ReviewEventApprove, true
	case ReviewOutcomeRequestChanges:
		return ReviewEventRequestChanges, true
	default:
		return "", false
	}
}
