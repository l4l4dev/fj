package cli

import (
	"context"
	"strings"
	"testing"

	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
)

type pullRequestListerStub struct {
	result []applicationpullrequest.PullRequest
}

func (s pullRequestListerStub) List(context.Context, applicationpullrequest.ListRequest) ([]applicationpullrequest.PullRequest, error) {
	return s.result, nil
}

func TestPullRequestListOutput(t *testing.T) {
	command := newPullRequestListCommand(pullRequestListerStub{result: []applicationpullrequest.PullRequest{{Number: 12, Title: "Improve flow", State: applicationpullrequest.StateOpen, HeadBranch: "feature", BaseBranch: "main"}}})
	command.SetArgs([]string{"alice/project"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	want := "Pull requests:\n#12 Improve flow [open]\n  feature -> main\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestPullRequestListEmptyOutput(t *testing.T) {
	command := newPullRequestListCommand(pullRequestListerStub{})
	command.SetArgs([]string{"alice/project"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if output.String() != "Pull requests:\nNo pull requests found.\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestPullRequestListRejectsInvalidFlags(t *testing.T) {
	for _, args := range [][]string{{"alice/project", "--page", "0"}, {"alice/project", "--limit", "101"}, {"alice/project", "--state", "merged"}} {
		command := newPullRequestListCommand(pullRequestListerStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}
