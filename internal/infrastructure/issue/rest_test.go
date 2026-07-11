package issue

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	applicationissue "github.com/l4l4dev/fj/internal/application/issue"
)

type stubTransport struct {
	path  string
	query url.Values
	body  string
}

func (stub *stubTransport) Do(_ context.Context, _ string, path string, query url.Values) (*http.Response, error) {
	stub.path, stub.query = path, query
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(stub.body))}, nil
}

func TestRESTAdapterList(t *testing.T) {
	transport := &stubTransport{body: `[{"number":12,"title":"Fix it","state":"open"}]`}
	page, err := NewRESTAdapter(transport).List(context.Background(), applicationissue.ListRequest{Owner: "alice", Name: "project", Page: 1, Limit: 30, State: applicationissue.StateOpen})
	if err != nil || len(page.Issues) != 1 || page.Issues[0].State != applicationissue.StateOpen {
		t.Fatalf("unexpected result: %+v err=%v", page, err)
	}
	if transport.path != "/api/v1/repos/alice/project/issues" || transport.query.Get("state") != "open" {
		t.Fatalf("unexpected request: path=%s query=%v", transport.path, transport.query)
	}
}

func TestRESTAdapterMapsIssueFilter(t *testing.T) {
	transport := &stubTransport{body: `[]`}
	_, err := NewRESTAdapter(transport).List(context.Background(), applicationissue.ListRequest{Owner: "alice", Name: "project", Page: 1, Limit: 30, State: applicationissue.StateOpen, Filter: applicationissue.IssueFilter{Assignee: "bob"}})
	if err != nil || transport.query.Get("assignee") != "bob" || transport.query.Get("labels") != "" {
		t.Fatalf("unexpected filter mapping: query=%v err=%v", transport.query, err)
	}
	transport.body = `[]`
	_, err = NewRESTAdapter(transport).List(context.Background(), applicationissue.ListRequest{Owner: "alice", Name: "project", Page: 1, Limit: 30, State: applicationissue.StateOpen, Filter: applicationissue.IssueFilter{Label: "bug"}})
	if err != nil || transport.query.Get("labels") != "bug" || transport.query.Get("assignee") != "" {
		t.Fatalf("unexpected label mapping: query=%v err=%v", transport.query, err)
	}
}

func TestRESTAdapterInspect(t *testing.T) {
	transport := &stubTransport{body: `{"number":12,"title":"Fix it","state":"open","body":"Details"}`}
	result, err := NewRESTAdapter(transport).Inspect(context.Background(), applicationissue.InspectRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || result.Number != 12 || result.Body != "Details" {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if transport.path != "/api/v1/repos/alice/project/issues/12" {
		t.Fatalf("unexpected path: %s", transport.path)
	}
}
