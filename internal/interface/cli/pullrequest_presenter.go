package cli

import (
	"fmt"
	"io"
	"strings"

	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
)

type pullRequestPresenter struct{}

func (pullRequestPresenter) PresentSubmittedReview(w io.Writer, review applicationpullrequest.SubmittedReview) error {
	_, err := fmt.Fprintf(w, "Review submitted: %s\n", review.State)
	return err
}

func (pullRequestPresenter) PresentUpdated(w io.Writer, detail applicationpullrequest.PullRequestDetail, fields []string) error {
	_, err := fmt.Fprintf(w, "Pull request updated: #%d\nChanged fields: %s\n", detail.Number, strings.Join(fields, ", "))
	return err
}

func (pullRequestPresenter) PresentComments(w io.Writer, comments []applicationpullrequest.Comment) error {
	if _, err := fmt.Fprintln(w, "Comments:"); err != nil {
		return err
	}
	if len(comments) == 0 {
		_, err := fmt.Fprintln(w, "No comments found.")
		return err
	}
	for _, comment := range comments {
		if _, err := fmt.Fprintf(w, "- #%d %s\n", comment.ID, comment.Body); err != nil {
			return err
		}
	}
	return nil
}

func (pullRequestPresenter) PresentComment(w io.Writer, comment applicationpullrequest.Comment) error {
	_, err := fmt.Fprintf(w, "Comment:\n#%d %s\n", comment.ID, comment.Body)
	return err
}

func (pullRequestPresenter) PresentMerged(w io.Writer, number int, method string) error {
	_, err := fmt.Fprintf(w, "Pull request merged: #%d (%s)\n", number, method)
	return err
}

func (pullRequestPresenter) PresentClosed(w io.Writer, detail applicationpullrequest.PullRequestDetail) error {
	_, err := fmt.Fprintf(w, "Pull request closed: #%d\nState: %s\n", detail.Number, detail.State)
	return err
}

func (pullRequestPresenter) PresentStatus(w io.Writer, status applicationpullrequest.PullRequestStatus) error {
	_, err := fmt.Fprintf(w, "Pull request: #%d\nReview: %s\nChecks: %s\nMergeable: %s\n", status.Number, status.Review, status.Check, status.Mergeable)
	return err
}

func (pullRequestPresenter) PresentCreated(w io.Writer, detail applicationpullrequest.PullRequestDetail) error {
	_, err := fmt.Fprintf(w, "Pull request created: #%d\nTitle: %s\nHead branch: %s\nBase branch: %s\n", detail.Number, detail.Title, detail.HeadBranch, detail.BaseBranch)
	return err
}

func (pullRequestPresenter) PresentInspect(w io.Writer, detail applicationpullrequest.PullRequestDetail) error {
	body := detail.Body
	if body == "" {
		body = "-"
	}
	head := detail.HeadBranch
	if head == "" {
		head = "-"
	}
	base := detail.BaseBranch
	if base == "" {
		base = "-"
	}
	_, err := fmt.Fprintf(w, "Pull request: #%d\nTitle: %s\nState: %s\nHead branch: %s\nBase branch: %s\nBody: %s\n", detail.Number, detail.Title, detail.State, head, base, body)
	return err
}

func (pullRequestPresenter) Present(w io.Writer, pullRequests []applicationpullrequest.PullRequest) error {
	if _, err := fmt.Fprintln(w, "Pull requests:"); err != nil {
		return err
	}
	if len(pullRequests) == 0 {
		_, err := fmt.Fprintln(w, "No pull requests found.")
		return err
	}
	for _, pullRequest := range pullRequests {
		if _, err := fmt.Fprintf(w, "#%d %s [%s]\n  %s -> %s\n", pullRequest.Number, pullRequest.Title, pullRequest.State, pullRequest.HeadBranch, pullRequest.BaseBranch); err != nil {
			return err
		}
	}
	return nil
}
