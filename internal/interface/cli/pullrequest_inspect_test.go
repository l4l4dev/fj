package cli

import (
	"context"
	"strings"
	"testing"

	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
)

type pullRequestInspectorStub struct{}

func (pullRequestInspectorStub) Inspect(_ context.Context, request applicationpullrequest.InspectRequest) (applicationpullrequest.PullRequestDetail, error) {
	return applicationpullrequest.PullRequestDetail{Number: request.Number, Title: "Improve flow", State: applicationpullrequest.StateOpen, HeadBranch: "feature", BaseBranch: "main", Body: "Details"}, nil
}

func TestPullRequestInspectOutput(t *testing.T) {
	command := newPullRequestInspectCommand(pullRequestInspectorStub{})
	command.SetArgs([]string{"alice/project", "12"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	want := "Pull request: #12\nTitle: Improve flow\nState: open\nHead branch: feature\nBase branch: main\nBody: Details\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestPullRequestInspectRejectsInvalidNumber(t *testing.T) {
	command := newPullRequestInspectCommand(pullRequestInspectorStub{})
	command.SetArgs([]string{"alice/project", "0"})
	if err := command.Execute(); err == nil {
		t.Fatal("expected validation error")
	}
}
