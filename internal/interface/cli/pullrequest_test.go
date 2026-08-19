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

type pullRequestReviewSubmitterStub struct {
	submission applicationpullrequest.ReviewSubmission
}

func (stub *pullRequestReviewSubmitterStub) SubmitReview(_ context.Context, submission applicationpullrequest.ReviewSubmission) (applicationpullrequest.SubmittedReview, error) {
	stub.submission = submission
	return applicationpullrequest.SubmittedReview{State: "APPROVED"}, nil
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

func TestPullRequestReviewOutputAndRequest(t *testing.T) {
	submitter := &pullRequestReviewSubmitterStub{}
	command := newPullRequestReviewCommand(submitter)
	command.SetArgs([]string{"alice/project", "12", "--outcome", "approve", "--body", "Looks good"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if output.String() != "Review submitted: APPROVED\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if submitter.submission.Owner != "alice" || submitter.submission.Name != "project" || submitter.submission.Number != 12 || submitter.submission.Event != applicationpullrequest.ReviewEventApprove || submitter.submission.Body != "Looks good" {
		t.Fatalf("unexpected submission: %+v", submitter.submission)
	}
}

func TestPullRequestReviewRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{"alice/project"},
		{"invalid", "12", "--outcome", "approve"},
		{"alice/project", "0", "--outcome", "approve"},
		{"alice/project", "not-a-number", "--outcome", "approve"},
		{"alice/project", "12"},
		{"alice/project", "12", "--outcome", "approved"},
		{"alice/project", "12", "--outcome", "comment"},
		{"alice/project", "12", "--outcome", "request-changes", "--body", "  "},
	}
	for _, args := range tests {
		command := newPullRequestReviewCommand(&pullRequestReviewSubmitterStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}

type pullRequestUpdaterStub struct {
	request applicationpullrequest.UpdateRequest
	detail  applicationpullrequest.PullRequestDetail
}

func (stub *pullRequestUpdaterStub) Update(_ context.Context, request applicationpullrequest.UpdateRequest) (applicationpullrequest.PullRequestDetail, error) {
	stub.request = request
	if stub.detail.Number == 0 {
		stub.detail.Number = request.Number
	}
	return stub.detail, nil
}

func TestPullRequestUpdateSendsOnlyChangedFields(t *testing.T) {
	updater := &pullRequestUpdaterStub{}
	command := newPullRequestUpdateCommand(updater)
	command.SetArgs([]string{"alice/project", "12", "--title", "New title"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if updater.request.Title == nil || *updater.request.Title != "New title" || updater.request.Body != nil {
		t.Fatalf("unexpected request: %+v", updater.request)
	}
	if output.String() != "Pull request updated: #12\nChanged fields: title\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestPullRequestUpdateRejectsInvalidInput(t *testing.T) {
	tests := [][]string{
		{"alice/project"},
		{"invalid", "12", "--title", "New"},
		{"alice/project", "0", "--title", "New"},
		{"alice/project", "12"},
	}
	for _, args := range tests {
		command := newPullRequestUpdateCommand(&pullRequestUpdaterStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for %v", args)
		}
	}
}

type pullRequestCommentViewerStub struct {
	request  applicationpullrequest.ListCommentsRequest
	comments []applicationpullrequest.Comment
}

func (stub *pullRequestCommentViewerStub) ListComments(_ context.Context, request applicationpullrequest.ListCommentsRequest) ([]applicationpullrequest.Comment, error) {
	stub.request = request
	return stub.comments, nil
}

type pullRequestCommentCreatorStub struct {
	request applicationpullrequest.AddCommentRequest
}

func (stub *pullRequestCommentCreatorStub) AddComment(_ context.Context, request applicationpullrequest.AddCommentRequest) (applicationpullrequest.Comment, error) {
	stub.request = request
	return applicationpullrequest.Comment{ID: 9, Body: request.Body}, nil
}

func TestPullRequestCommentListOutput(t *testing.T) {
	viewer := &pullRequestCommentViewerStub{comments: []applicationpullrequest.Comment{{ID: 5, Body: "First"}}}
	command := newPullRequestCommentListCommand(viewer)
	command.SetArgs([]string{"alice/project", "12"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if viewer.request.Number != 12 {
		t.Fatalf("unexpected request: %+v", viewer.request)
	}
	if output.String() != "Comments:\n- #5 First\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestPullRequestCommentAddOutputAndRequest(t *testing.T) {
	creator := &pullRequestCommentCreatorStub{}
	command := newPullRequestCommentAddCommand(creator)
	command.SetArgs([]string{"alice/project", "12", "--body", "Looks good"})
	var output strings.Builder
	command.SetOut(&output)
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if creator.request.Number != 12 || creator.request.Body != "Looks good" {
		t.Fatalf("unexpected request: %+v", creator.request)
	}
	if output.String() != "Comment:\n#9 Looks good\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestPullRequestCommentRejectsInvalidInput(t *testing.T) {
	listTests := [][]string{
		{"alice/project"},
		{"invalid", "12"},
		{"alice/project", "0"},
	}
	for _, args := range listTests {
		command := newPullRequestCommentListCommand(&pullRequestCommentViewerStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for list %v", args)
		}
	}
	addTests := [][]string{
		{"alice/project", "12"},
		{"alice/project", "12", "--body", "  "},
	}
	for _, args := range addTests {
		command := newPullRequestCommentAddCommand(&pullRequestCommentCreatorStub{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("expected validation error for add %v", args)
		}
	}
}
