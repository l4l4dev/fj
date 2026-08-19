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
	method   string
	path     string
	body     []byte
	err      error
	response string
}

func (stub *jsonStubTransport) Do(context.Context, string, string, url.Values) (*http.Response, error) {
	return nil, errors.New("unexpected Do call")
}

func (stub *jsonStubTransport) DoJSON(_ context.Context, method, path string, _ url.Values, body []byte) (*http.Response, error) {
	stub.method, stub.path, stub.body = method, path, body
	if stub.err != nil {
		return nil, stub.err
	}
	response := stub.response
	if response == "" {
		response = `{"number":7,"title":"Improve flow","state":"open","head":{"ref":"feature"},"base":{"ref":"main"}}`
	}
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

func TestRESTAdapterSubmitReviewPayload(t *testing.T) {
	transport := &jsonStubTransport{response: `{"id":42,"state":"APPROVED","body":"not printed"}`}
	result, err := NewRESTAdapter(transport).SubmitReview(context.Background(), applicationpullrequest.ReviewSubmission{Owner: "alice", Name: "project", Number: 12, Event: applicationpullrequest.ReviewEventApprove})
	if err != nil || result.State != "APPROVED" {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if transport.method != http.MethodPost || transport.path != "/api/v1/repos/alice/project/pulls/12/reviews" {
		t.Fatalf("unexpected request: method=%s path=%s", transport.method, transport.path)
	}
	if string(transport.body) != `{"event":"APPROVE"}` {
		t.Fatalf("unexpected body: %s", transport.body)
	}

	transport = &jsonStubTransport{response: `{"state":"REQUEST_CHANGES"}`}
	_, err = NewRESTAdapter(transport).SubmitReview(context.Background(), applicationpullrequest.ReviewSubmission{Owner: "alice", Name: "project", Number: 12, Event: applicationpullrequest.ReviewEventRequestChanges, Body: "Please add a test"})
	if err != nil || string(transport.body) != `{"event":"REQUEST_CHANGES","body":"Please add a test"}` {
		t.Fatalf("unexpected body submission: body=%s err=%v", transport.body, err)
	}
}

func TestRESTAdapterSubmitReviewMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{status: http.StatusUnauthorized, category: apperror.Authentication},
		{status: http.StatusForbidden, category: apperror.Authentication},
		{status: http.StatusNotFound, category: apperror.NotFound},
		{status: http.StatusConflict, category: apperror.Conflict},
		{status: http.StatusUnprocessableEntity, category: apperror.Validation},
		{status: http.StatusInternalServerError, category: apperror.Remote},
	}
	for _, test := range tests {
		_, err := NewRESTAdapter(&jsonStubTransport{err: statusError(test.status)}).SubmitReview(context.Background(), applicationpullrequest.ReviewSubmission{})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: error=%#v, want category %v", test.status, err, test.category)
		}
	}
}

func TestRESTAdapterUpdateSendsOnlySuppliedFields(t *testing.T) {
	title := "New title"
	stub := &jsonStubTransport{response: `{"number":12,"title":"New title","state":"open","body":"Same","head":{"ref":"feature"},"base":{"ref":"main"}}`}
	result, err := NewRESTAdapter(stub).Update(context.Background(), applicationpullrequest.UpdateRequest{Owner: "alice", Name: "project", Number: 12, Title: &title})
	if err != nil || result.Number != 12 || result.Title != "New title" {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if stub.method != http.MethodPatch || stub.path != "/api/v1/repos/alice/project/pulls/12" {
		t.Fatalf("unexpected request: %s %s", stub.method, stub.path)
	}
	if string(stub.body) != `{"title":"New title"}` {
		t.Fatalf("unexpected payload: %s", stub.body)
	}
}

func TestRESTAdapterUpdateMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{401, apperror.Authentication},
		{403, apperror.Authentication},
		{404, apperror.NotFound},
		{409, apperror.Validation},
		{422, apperror.Validation},
		{500, apperror.Remote},
	}
	title := "New title"
	for _, test := range tests {
		stub := &jsonStubTransport{err: statusError(test.status)}
		_, err := NewRESTAdapter(stub).Update(context.Background(), applicationpullrequest.UpdateRequest{Owner: "alice", Name: "project", Number: 12, Title: &title})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %v", test.status, err)
		}
	}
}

type commentStubTransport struct {
	paths      []string
	verifyErr  error
	listBody   string
	jsonMethod string
	jsonPath   string
	jsonBody   []byte
	jsonErr    error
}

func (stub *commentStubTransport) Do(_ context.Context, _ string, path string, _ url.Values) (*http.Response, error) {
	stub.paths = append(stub.paths, path)
	if strings.Contains(path, "/pulls/") {
		if stub.verifyErr != nil {
			return nil, stub.verifyErr
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"number":12}`))}, nil
	}
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(stub.listBody))}, nil
}

func (stub *commentStubTransport) DoJSON(_ context.Context, method, path string, _ url.Values, body []byte) (*http.Response, error) {
	stub.jsonMethod, stub.jsonPath, stub.jsonBody = method, path, body
	if stub.jsonErr != nil {
		return nil, stub.jsonErr
	}
	return &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(strings.NewReader(`{"id":9,"body":"Looks good"}`))}, nil
}

func TestRESTAdapterListCommentsVerifiesPullRequest(t *testing.T) {
	stub := &commentStubTransport{listBody: `[{"id":5,"body":"First"}]`}
	comments, err := NewRESTAdapter(stub).ListComments(context.Background(), applicationpullrequest.ListCommentsRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || len(comments) != 1 || comments[0].ID != 5 || comments[0].Body != "First" {
		t.Fatalf("unexpected result: %+v err=%v", comments, err)
	}
	if len(stub.paths) != 2 || stub.paths[0] != "/api/v1/repos/alice/project/pulls/12" || stub.paths[1] != "/api/v1/repos/alice/project/issues/12/comments" {
		t.Fatalf("unexpected request paths: %v", stub.paths)
	}
}

func TestRESTAdapterListCommentsRejectsMissingPullRequest(t *testing.T) {
	stub := &commentStubTransport{verifyErr: statusError(404)}
	_, err := NewRESTAdapter(stub).ListComments(context.Background(), applicationpullrequest.ListCommentsRequest{Owner: "alice", Name: "project", Number: 12})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "pull request not found" {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stub.paths) != 1 {
		t.Fatalf("comment endpoint must not be called when the pull request is missing: %v", stub.paths)
	}
}

func TestRESTAdapterAddCommentVerifiesAndPosts(t *testing.T) {
	stub := &commentStubTransport{}
	comment, err := NewRESTAdapter(stub).AddComment(context.Background(), applicationpullrequest.AddCommentRequest{Owner: "alice", Name: "project", Number: 12, Body: "Looks good"})
	if err != nil || comment.ID != 9 {
		t.Fatalf("unexpected result: %+v err=%v", comment, err)
	}
	if len(stub.paths) != 1 || stub.paths[0] != "/api/v1/repos/alice/project/pulls/12" {
		t.Fatalf("unexpected verification path: %v", stub.paths)
	}
	if stub.jsonMethod != http.MethodPost || stub.jsonPath != "/api/v1/repos/alice/project/issues/12/comments" || string(stub.jsonBody) != `{"body":"Looks good"}` {
		t.Fatalf("unexpected request: %s %s %s", stub.jsonMethod, stub.jsonPath, stub.jsonBody)
	}
}

func TestRESTAdapterAddCommentMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{401, apperror.Authentication},
		{403, apperror.Authentication},
		{422, apperror.Validation},
		{500, apperror.Remote},
	}
	for _, test := range tests {
		stub := &commentStubTransport{jsonErr: statusError(test.status)}
		_, err := NewRESTAdapter(stub).AddComment(context.Background(), applicationpullrequest.AddCommentRequest{Owner: "alice", Name: "project", Number: 12, Body: "Hi"})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %v", test.status, err)
		}
	}
}
