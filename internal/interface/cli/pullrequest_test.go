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

type pullRequestCreatorStub struct {
	request applicationpullrequest.CreateRequest
}

func (stub *pullRequestCreatorStub) Create(_ context.Context, request applicationpullrequest.CreateRequest) (applicationpullrequest.PullRequestDetail, error) {
	stub.request = request
	return applicationpullrequest.PullRequestDetail{Number: 7, Title: request.Title, HeadBranch: request.HeadBranch, BaseBranch: request.BaseBranch}, nil
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

func TestPullRequestCreateOutputAndRequest(t *testing.T) {
	creator := &pullRequestCreatorStub{}
	command := newPullRequestCreateCommand(creator)
	command.SetArgs([]string{"alice/project", "--title", "Improve flow", "--head", "feature", "--base", "main"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	want := "Pull request created: #7\nTitle: Improve flow\nHead branch: feature\nBase branch: main\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if creator.request.Owner != "alice" || creator.request.Name != "project" || creator.request.HeadBranch != "feature" || creator.request.BaseBranch != "main" {
		t.Fatalf("unexpected request: %+v", creator.request)
	}
}

func TestPullRequestCreateRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{"invalid", "--title", "Improve flow", "--head", "feature", "--base", "main"},
		{"alice/project", "--head", "feature", "--base", "main"},
		{"alice/project", "--title", "Improve flow", "--base", "main"},
		{"alice/project", "--title", "Improve flow", "--head", "feature"},
		{"alice/project", "--title", "Improve flow", "--head", "main", "--base", "main"},
		{"alice/project", "--title", "Improve flow", "--head", "alice:feature", "--base", "main"},
	}
	for _, args := range tests {
		command := newPullRequestCreateCommand(&pullRequestCreatorStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}
