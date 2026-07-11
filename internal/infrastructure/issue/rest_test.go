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

type jsonStubTransport struct {
	stubTransport
	method       string
	body         []byte
	responseBody string
}

func (stub *jsonStubTransport) Do(_ context.Context, method, path string, _ url.Values) (*http.Response, error) {
	stub.method, stub.path = method, path
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(stub.stubTransport.body))}, nil
}

func (stub *jsonStubTransport) DoJSON(_ context.Context, method, path string, _ url.Values, body []byte) (*http.Response, error) {
	stub.method, stub.path, stub.body = method, path, body
	response := stub.responseBody
	if response == "" {
		response = `{"number":13,"title":"Created","state":"open","body":"Body"}`
	}
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(response))}, nil
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

func TestRESTAdapterCreate(t *testing.T) {
	transport := &jsonStubTransport{}
	result, err := NewRESTAdapter(transport).Create(context.Background(), applicationissue.CreateRequest{Owner: "alice", Name: "project", Title: "Created", Body: "Body"})
	if err != nil || result.Number != 13 || transport.method != http.MethodPost || transport.path != "/api/v1/repos/alice/project/issues" {
		t.Fatalf("unexpected create result: %+v method=%s path=%s err=%v", result, transport.method, transport.path, err)
	}
	if string(transport.body) != `{"title":"Created","body":"Body"}` {
		t.Fatalf("unexpected request body: %s", transport.body)
	}
}

func TestRESTAdapterUpdateSendsOnlySpecifiedFields(t *testing.T) {
	transport := &jsonStubTransport{}
	title := "Updated"
	_, err := NewRESTAdapter(transport).Update(context.Background(), applicationissue.UpdateRequest{Owner: "alice", Name: "project", Number: 12, Title: &title})
	if err != nil || transport.method != http.MethodPatch || transport.path != "/api/v1/repos/alice/project/issues/12" {
		t.Fatalf("unexpected update request: method=%s path=%s err=%v", transport.method, transport.path, err)
	}
	if string(transport.body) != `{"title":"Updated"}` {
		t.Fatalf("unexpected update body: %s", transport.body)
	}
	body := ""
	_, err = NewRESTAdapter(transport).Update(context.Background(), applicationissue.UpdateRequest{Owner: "alice", Name: "project", Number: 12, Body: &body})
	if err != nil || string(transport.body) != `{"body":""}` {
		t.Fatalf("unexpected empty body update: %s err=%v", transport.body, err)
	}
}

func TestRESTAdapterChangeState(t *testing.T) {
	transport := &jsonStubTransport{}
	_, err := NewRESTAdapter(transport).ChangeState(context.Background(), applicationissue.ChangeStateRequest{Owner: "alice", Name: "project", Number: 12, State: applicationissue.StateClosed})
	if err != nil || transport.method != http.MethodPatch || transport.path != "/api/v1/repos/alice/project/issues/12" {
		t.Fatalf("unexpected state request: method=%s path=%s err=%v", transport.method, transport.path, err)
	}
	if string(transport.body) != `{"state":"closed"}` {
		t.Fatalf("unexpected state body: %s", transport.body)
	}
}

func TestRESTAdapterComments(t *testing.T) {
	transport := &stubTransport{body: `[{"id":1,"body":"hello"}]`}
	comments, err := NewRESTAdapter(transport).ListComments(context.Background(), applicationissue.ListCommentsRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || len(comments) != 1 || comments[0].Body != "hello" || transport.path != "/api/v1/repos/alice/project/issues/12/comments" {
		t.Fatalf("unexpected comments: %+v path=%s err=%v", comments, transport.path, err)
	}
	jsonTransport := &jsonStubTransport{}
	comment, err := NewRESTAdapter(jsonTransport).AddComment(context.Background(), applicationissue.AddCommentRequest{Owner: "alice", Name: "project", Number: 12, Body: "hello"})
	if err != nil || comment.Body != "Body" || jsonTransport.method != http.MethodPost || jsonTransport.path != "/api/v1/repos/alice/project/issues/12/comments" {
		t.Fatalf("unexpected add comment: %+v method=%s path=%s err=%v", comment, jsonTransport.method, jsonTransport.path, err)
	}
	if string(jsonTransport.body) != `{"body":"hello"}` {
		t.Fatalf("unexpected comment body: %s", jsonTransport.body)
	}
}

func TestRESTAdapterLabels(t *testing.T) {
	transport := &jsonStubTransport{}
	transport.stubTransport.body = `[]`
	transport.responseBody = `[{"id":42,"name":"bug"}]`
	label, err := NewRESTAdapter(transport).AddLabel(context.Background(), applicationissue.AddLabelRequest{Owner: "alice", Name: "project", Number: 12, Label: "bug"})
	if err != nil || label.ID != 42 || label.Name != "bug" || transport.method != http.MethodPost || transport.path != "/api/v1/repos/alice/project/issues/12/labels" || string(transport.body) != `{"labels":["bug"]}` {
		t.Fatalf("unexpected add label: %+v method=%s path=%s body=%s err=%v", label, transport.method, transport.path, transport.body, err)
	}
	transport.stubTransport.body = `[{"id":42,"name":"bug"}]`
	transport.method, transport.path, transport.body = "", "", nil
	label, err = NewRESTAdapter(transport).RemoveLabel(context.Background(), applicationissue.RemoveLabelRequest{Owner: "alice", Name: "project", Number: 12, Label: "bug"})
	if err != nil || label.ID != 42 || transport.method != http.MethodDelete || transport.path != "/api/v1/repos/alice/project/issues/12/labels/42" {
		t.Fatalf("unexpected remove label: %+v method=%s path=%s err=%v", label, transport.method, transport.path, err)
	}
}

type milestoneTransport struct {
	responses []string
	methods   []string
	paths     []string
	bodies    [][]byte
}

func (stub *milestoneTransport) next(method, path string, body []byte) (*http.Response, error) {
	stub.methods = append(stub.methods, method)
	stub.paths = append(stub.paths, path)
	stub.bodies = append(stub.bodies, body)
	response := stub.responses[0]
	stub.responses = stub.responses[1:]
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(response))}, nil
}

func (stub *milestoneTransport) Do(_ context.Context, method, path string, _ url.Values) (*http.Response, error) {
	return stub.next(method, path, nil)
}

func (stub *milestoneTransport) DoJSON(_ context.Context, method, path string, _ url.Values, body []byte) (*http.Response, error) {
	return stub.next(method, path, body)
}

func TestRESTAdapterMilestoneIdempotencyAndConversion(t *testing.T) {
	transport := &milestoneTransport{responses: []string{`[{"id":7,"title":"v1"}]`, `{"number":12,"milestone":null}`, `{"number":12,"milestone":{"id":7,"title":"v1"}}`}}
	milestone, err := NewRESTAdapter(transport).SetMilestone(context.Background(), applicationissue.SetMilestoneRequest{Owner: "alice", Name: "project", Number: 12, Milestone: "v1"})
	if err != nil || milestone.ID != 7 || milestone.Title != "v1" || len(transport.methods) != 3 || transport.methods[2] != http.MethodPatch || string(transport.bodies[2]) != `{"milestone":7}` {
		t.Fatalf("unexpected milestone set: %+v methods=%v paths=%v bodies=%v err=%v", milestone, transport.methods, transport.paths, transport.bodies, err)
	}
	transport = &milestoneTransport{responses: []string{`{"number":12,"milestone":null}`}}
	if err := NewRESTAdapter(transport).ClearMilestone(context.Background(), applicationissue.ClearMilestoneRequest{Owner: "alice", Name: "project", Number: 12}); err != nil || len(transport.methods) != 1 {
		t.Fatalf("unexpected clear idempotency: methods=%v err=%v", transport.methods, err)
	}
}
