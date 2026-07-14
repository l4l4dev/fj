package cli

import (
	"fmt"
	"io"

	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
)

type pullRequestPresenter struct{}

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
