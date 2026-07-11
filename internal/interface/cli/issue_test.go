package cli

import (
	"context"
	"strings"
	"testing"

	applicationissue "github.com/l4l4dev/fj/internal/application/issue"
)

type issueListerStub struct{ page applicationissue.Page }

func (stub issueListerStub) List(context.Context, applicationissue.ListRequest) (applicationissue.Page, error) {
	return stub.page, nil
}

func TestIssueListPresenterOutput(t *testing.T) {
	command := newIssueListCommand(issueListerStub{page: applicationissue.Page{Issues: []applicationissue.Issue{{Number: 12, Title: "Fix authentication flow", State: applicationissue.StateOpen}}, Page: 1, Limit: 30}})
	var output strings.Builder
	command.SetOut(&output)
	command.SetArgs([]string{"alice/project"})
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	want := "Issues:\n- #12 Fix authentication flow [open]\n\nPage: 1\nLimit: 30\nMore pages: false\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestIssueListRejectsInvalidState(t *testing.T) {
	command := newIssueListCommand(issueListerStub{})
	command.SetArgs([]string{"alice/project", "--state", "invalid"})
	if err := command.Execute(); err == nil {
		t.Fatal("expected validation error")
	}
}
