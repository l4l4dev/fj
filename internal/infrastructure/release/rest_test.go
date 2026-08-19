package release

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationrelease "github.com/l4l4dev/fj/internal/application/release"
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

type statusError int

func (err statusError) Error() string   { return "remote failure" }
func (err statusError) StatusCode() int { return int(err) }

type errorTransport struct{ err error }

func (transport errorTransport) Do(context.Context, string, string, url.Values) (*http.Response, error) {
	return nil, transport.err
}

func (transport errorTransport) DoJSON(context.Context, string, string, url.Values, []byte) (*http.Response, error) {
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
		response = `{"id":1,"tag_name":"v1.0.0","name":"First release","draft":true,"prerelease":false}`
	}
	return &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(strings.NewReader(response))}, nil
}

func TestRESTAdapterList(t *testing.T) {
	transport := &stubTransport{body: `[{"id":1,"tag_name":"v1.0.0","name":"First release","draft":false,"prerelease":false},{"id":2,"tag_name":"v1.1.0-rc1","name":"Release candidate","draft":false,"prerelease":true},{"id":3,"tag_name":"v0.1.0","name":"Draft","draft":true,"prerelease":false}]`}
	result, err := NewRESTAdapter(transport).List(context.Background(), applicationrelease.ListRequest{Owner: "alice", Name: "project", Page: 2, Limit: 20})
	if err != nil || len(result) != 3 {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if result[0].TagName != "v1.0.0" || result[0].Title != "First release" || result[0].Draft || result[0].Prerelease {
		t.Fatalf("unexpected release 0: %+v", result[0])
	}
	if result[1].Prerelease != true {
		t.Fatalf("unexpected release 1: %+v", result[1])
	}
	if result[2].Draft != true {
		t.Fatalf("unexpected release 2: %+v", result[2])
	}
	if transport.path != "/api/v1/repos/alice/project/releases" || transport.query.Get("page") != "2" || transport.query.Get("limit") != "20" {
		t.Fatalf("unexpected request: path=%s query=%v", transport.path, transport.query)
	}
}

func TestRESTAdapterListEmpty(t *testing.T) {
	result, err := NewRESTAdapter(&stubTransport{body: `[]`}).List(context.Background(), applicationrelease.ListRequest{Owner: "alice", Name: "project", Page: 1, Limit: 20})
	if err != nil || result == nil || len(result) != 0 {
		t.Fatalf("unexpected empty result: %#v err=%v", result, err)
	}
}

func TestRESTAdapterListMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{http.StatusUnauthorized, apperror.Authentication},
		{http.StatusForbidden, apperror.Authentication},
		{http.StatusNotFound, apperror.NotFound},
		{http.StatusInternalServerError, apperror.Remote},
	}
	for _, test := range tests {
		_, err := NewRESTAdapter(errorTransport{err: statusError(test.status)}).List(context.Background(), applicationrelease.ListRequest{Owner: "alice", Name: "project", Page: 1, Limit: 20})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %#v", test.status, err)
		}
	}
}

func TestRESTAdapterListMapsNotFoundSafely(t *testing.T) {
	_, err := NewRESTAdapter(errorTransport{err: statusError(http.StatusNotFound)}).List(context.Background(), applicationrelease.ListRequest{Owner: "example-owner", Name: "example-repository", Page: 1, Limit: 20})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "repository not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
	}
}

func TestRESTAdapterInspect(t *testing.T) {
	transport := &stubTransport{body: `{"id":1,"tag_name":"v1.0.0","name":"First release","body":"Release notes","draft":false,"prerelease":false,"assets":[{"id":11,"name":"fj_darwin_arm64.tar.gz","size":1024}]}`}
	result, err := NewRESTAdapter(transport).Inspect(context.Background(), applicationrelease.InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
	if err != nil || result.TagName != "v1.0.0" || result.Title != "First release" || result.Notes != "Release notes" || result.Draft || result.Prerelease {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if len(result.Assets) != 1 || result.Assets[0].ID != 11 || result.Assets[0].Name != "fj_darwin_arm64.tar.gz" || result.Assets[0].Size != 1024 {
		t.Fatalf("unexpected assets: %+v", result.Assets)
	}
	if transport.path != "/api/v1/repos/alice/project/releases/tags/v1.0.0" {
		t.Fatalf("unexpected path: %s", transport.path)
	}
}

func TestRESTAdapterInspectEscapesTag(t *testing.T) {
	transport := &stubTransport{body: `{}`}
	_, err := NewRESTAdapter(transport).Inspect(context.Background(), applicationrelease.InspectRequest{Owner: "alice", Name: "project", Tag: "v1/rc1"})
	if err != nil {
		t.Fatal(err)
	}
	if transport.path != "/api/v1/repos/alice/project/releases/tags/v1%2Frc1" {
		t.Fatalf("unexpected path: %s", transport.path)
	}
}

func TestRESTAdapterInspectMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{http.StatusUnauthorized, apperror.Authentication},
		{http.StatusForbidden, apperror.Authentication},
		{http.StatusNotFound, apperror.NotFound},
		{http.StatusInternalServerError, apperror.Remote},
	}
	for _, test := range tests {
		_, err := NewRESTAdapter(errorTransport{err: statusError(test.status)}).Inspect(context.Background(), applicationrelease.InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %#v", test.status, err)
		}
	}
}

func TestRESTAdapterInspectMapsNotFoundSafely(t *testing.T) {
	_, err := NewRESTAdapter(errorTransport{err: statusError(http.StatusNotFound)}).Inspect(context.Background(), applicationrelease.InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "release not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
	}
}

func TestRESTAdapterCreateSendsDraftPayload(t *testing.T) {
	transport := &jsonStubTransport{}
	request := applicationrelease.CreateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: "First release", Notes: "Notes", Prerelease: true}
	result, err := NewRESTAdapter(transport).Create(context.Background(), request)
	if err != nil || result.TagName != "v1.0.0" || result.Title != "First release" {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if transport.method != http.MethodPost || transport.path != "/api/v1/repos/alice/project/releases" {
		t.Fatalf("unexpected request: method=%s path=%s", transport.method, transport.path)
	}
	if string(transport.body) != `{"tag_name":"v1.0.0","name":"First release","body":"Notes","draft":true,"prerelease":true}` {
		t.Fatalf("unexpected body: %s", transport.body)
	}
}

func TestRESTAdapterCreateMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{http.StatusUnauthorized, apperror.Authentication},
		{http.StatusForbidden, apperror.Authentication},
		{http.StatusNotFound, apperror.NotFound},
		{http.StatusConflict, apperror.Conflict},
		{http.StatusUnprocessableEntity, apperror.Validation},
		{http.StatusInternalServerError, apperror.Remote},
	}
	for _, test := range tests {
		_, err := NewRESTAdapter(&jsonStubTransport{err: statusError(test.status)}).Create(context.Background(), applicationrelease.CreateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: "First release"})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %#v", test.status, err)
		}
	}
}

func TestRESTAdapterCreateMapsNotFoundSafely(t *testing.T) {
	_, err := NewRESTAdapter(&jsonStubTransport{err: statusError(http.StatusNotFound)}).Create(context.Background(), applicationrelease.CreateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: "First release"})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "repository not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
	}
}
