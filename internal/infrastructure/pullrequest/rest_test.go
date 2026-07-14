package pullrequest

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
)

type stubTransport struct {
	path  string
	query url.Values
	body  string
}

type statusError int

func (err statusError) Error() string   { return "remote failure" }
func (err statusError) StatusCode() int { return int(err) }

type errorTransport struct{ err error }

func (transport errorTransport) Do(context.Context, string, string, url.Values) (*http.Response, error) {
	return nil, transport.err
}

type jsonStubTransport struct {
	method string
	path   string
	body   []byte
	err    error
}

func (stub *jsonStubTransport) Do(context.Context, string, string, url.Values) (*http.Response, error) {
	return nil, errors.New("unexpected Do call")
}

func (stub *jsonStubTransport) DoJSON(_ context.Context, method, path string, _ url.Values, body []byte) (*http.Response, error) {
	stub.method, stub.path, stub.body = method, path, body
	if stub.err != nil {
		return nil, stub.err
	}
	response := `{"number":7,"title":"Improve flow","state":"open","head":{"ref":"feature"},"base":{"ref":"main"}}`
	return &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(strings.NewReader(response))}, nil
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

func TestRESTAdapterListMapsNotFoundSafely(t *testing.T) {
	_, err := NewRESTAdapter(errorTransport{err: statusError(http.StatusNotFound)}).List(context.Background(), applicationpullrequest.ListRequest{Owner: "example-owner", Name: "example-repository", Page: 1, Limit: 20, State: applicationpullrequest.StateOpen})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "repository not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
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

func TestRESTAdapterCreate(t *testing.T) {
	transport := &jsonStubTransport{}
	request := applicationpullrequest.CreateRequest{Owner: "alice", Name: "project", Title: "Improve flow", HeadBranch: "feature", BaseBranch: "main"}
	result, err := NewRESTAdapter(transport).Create(context.Background(), request)
	if err != nil || result.Number != 7 || result.HeadBranch != "feature" || result.BaseBranch != "main" {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if transport.method != http.MethodPost || transport.path != "/api/v1/repos/alice/project/pulls" {
		t.Fatalf("unexpected request: method=%s path=%s", transport.method, transport.path)
	}
	if string(transport.body) != `{"title":"Improve flow","head":"feature","base":"main"}` {
		t.Fatalf("unexpected body: %s", transport.body)
	}
}

func TestRESTAdapterCreateMapsConflictWithoutRemoteDetails(t *testing.T) {
	transport := &jsonStubTransport{err: statusError(http.StatusUnprocessableEntity)}
	_, err := NewRESTAdapter(transport).Create(context.Background(), applicationpullrequest.CreateRequest{})
	if err == nil || err.Error() != "create pull request: pull request branches are invalid or conflict with an existing pull request" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRESTAdapterCreateMapsNotFoundWithoutMisdiagnosingTarget(t *testing.T) {
	transport := &jsonStubTransport{err: statusError(http.StatusNotFound)}
	_, err := NewRESTAdapter(transport).Create(context.Background(), applicationpullrequest.CreateRequest{})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "repository or branch not found" {
		t.Fatalf("unexpected error: %#v", err)
	}
}
