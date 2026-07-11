package cli

import (
	"fmt"
	"io"

	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
)

type pullRequestPresenter struct{}

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
