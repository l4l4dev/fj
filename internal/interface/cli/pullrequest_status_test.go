package cli

import (
	"context"
	"strings"
	"testing"

	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
)

type pullRequestStatusViewerStub struct {
	request applicationpullrequest.StatusRequest
}

func (stub *pullRequestStatusViewerStub) ViewStatus(_ context.Context, request applicationpullrequest.StatusRequest) (applicationpullrequest.StatusData, error) {
	stub.request = request
	return applicationpullrequest.StatusData{Number: request.Number, ReviewsAvailable: true, Reviews: []applicationpullrequest.Review{{ID: 1, ReviewerID: 1, State: "APPROVED"}}, ChecksAvailable: true, Checks: []string{"pending"}, Mergeable: applicationpullrequest.MergeableYes}, nil
}

func TestPullRequestStatusOutput(t *testing.T) {
	viewer := &pullRequestStatusViewerStub{}
	command := newPullRequestStatusCommand(viewer)
	command.SetArgs([]string{"alice/project", "12"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	want := "Pull request: #12\nReview: success\nChecks: pending\nMergeable: yes\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if viewer.request.Owner != "alice" || viewer.request.Name != "project" || viewer.request.Number != 12 {
		t.Fatalf("unexpected request: %+v", viewer.request)
	}
}

func TestPullRequestStatusRejectsInvalidInput(t *testing.T) {
	for _, args := range [][]string{{"invalid", "1"}, {"alice/project"}, {"alice/project", "0"}, {"alice/project", "invalid"}} {
		command := newPullRequestStatusCommand(&pullRequestStatusViewerStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}
