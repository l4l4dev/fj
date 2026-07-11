package pullrequest

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
)

type stubTransport struct {
	path  string
	query url.Values
	body  string
}

func (s *stubTransport) Do(_ context.Context, _ string, path string, query url.Values) (*http.Response, error) {
	s.path, s.query = path, query
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(s.body))}, nil
}

func TestRESTAdapterList(t *testing.T) {
	transport := &stubTransport{body: `[{"number":12,"title":"Improve flow","state":"open","head":{"ref":"feature"},"base":{"ref":"main"}}]`}
	result, err := NewRESTAdapter(transport).List(context.Background(), applicationpullrequest.ListRequest{Owner: "alice", Name: "project", Page: 2, Limit: 20, State: applicationpullrequest.StateOpen})
	if err != nil || len(result) != 1 || result[0].Number != 12 || result[0].HeadBranch != "feature" || result[0].BaseBranch != "main" {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if transport.path != "/api/v1/repos/alice/project/pulls" || transport.query.Get("page") != "2" || transport.query.Get("limit") != "20" || transport.query.Get("state") != "open" {
		t.Fatalf("unexpected request: path=%s query=%v", transport.path, transport.query)
	}
}

func TestRESTAdapterListEmpty(t *testing.T) {
	result, err := NewRESTAdapter(&stubTransport{body: `[]`}).List(context.Background(), applicationpullrequest.ListRequest{Owner: "alice", Name: "project", Page: 1, Limit: 20, State: applicationpullrequest.StateOpen})
	if err != nil || result == nil || len(result) != 0 {
		t.Fatalf("unexpected empty result: %#v err=%v", result, err)
	}
}

func TestRESTAdapterInspect(t *testing.T) {
	transport := &stubTransport{body: `{"number":12,"title":"Improve flow","state":"open","body":"Details","head":{"ref":"feature"},"base":{"ref":"main"}}`}
	result, err := NewRESTAdapter(transport).Inspect(context.Background(), applicationpullrequest.InspectRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || result.Number != 12 || result.Body != "Details" || result.HeadBranch != "feature" || result.BaseBranch != "main" {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if transport.path != "/api/v1/repos/alice/project/pulls/12" {
		t.Fatalf("unexpected path: %s", transport.path)
	}
}
