package repository

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	applicationrepository "github.com/l4l4dev/fj/internal/application/repository"
)

type transportFunc func(context.Context, string, string, url.Values) (*http.Response, error)

func (transport transportFunc) Do(ctx context.Context, method, path string, query url.Values) (*http.Response, error) {
	return transport(ctx, method, path, query)
}

type jsonTransportFunc func(context.Context, string, string, url.Values, []byte) (*http.Response, error)

func (transport jsonTransportFunc) Do(ctx context.Context, method, path string, query url.Values) (*http.Response, error) {
	return transport(ctx, method, path, query, nil)
}

func (transport jsonTransportFunc) DoJSON(ctx context.Context, method, path string, query url.Values, body []byte) (*http.Response, error) {
	return transport(ctx, method, path, query, body)
}

type transportError struct {
	operation string
	status    int
	secret    string
}

func (err transportError) Error() string {
	return err.secret
}

func (err transportError) Operation() string {
	return err.operation
}

func (err transportError) StatusCode() int {
	return err.status
}

func TestRESTAdapterListsRepositoriesWithOnlyPageAndLimit(t *testing.T) {
	ctx := context.Background()
	var receivedMethod, receivedPath string
	var receivedQuery url.Values
	adapter := NewRESTAdapter(transportFunc(func(received context.Context, method, path string, query url.Values) (*http.Response, error) {
		if received != ctx {
			t.Errorf("context was not passed through")
		}
		receivedMethod, receivedPath, receivedQuery = method, path, query
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[{"name":"project","owner":{"login":"octo"}}]`)),
		}, nil
	}))

	result, err := adapter.List(ctx, applicationrepository.ListRequest{Page: 2, Limit: 25})
	if err != nil {
		t.Fatal(err)
	}
	if receivedMethod != http.MethodGet || receivedPath != "/api/v1/user/repos" {
		t.Errorf("request = %s %s", receivedMethod, receivedPath)
	}
	if len(receivedQuery) != 2 || receivedQuery.Get("page") != "2" || receivedQuery.Get("limit") != "25" {
		t.Errorf("query = %#v", receivedQuery)
	}
	if len(result) != 1 || result[0] != (applicationrepository.Repository{Owner: "octo", Name: "project"}) {
		t.Errorf("result = %#v", result)
	}
}

func TestRESTAdapterReturnsNonNilEmptySlice(t *testing.T) {
	adapter := NewRESTAdapter(transportFunc(func(context.Context, string, string, url.Values) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("[]"))}, nil
	}))

	result, err := adapter.List(context.Background(), applicationrepository.ListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil || len(result) != 0 {
		t.Errorf("result = %#v, want non-nil empty slice", result)
	}
}

func TestRESTAdapterGetsRepositoryWithEncodedPath(t *testing.T) {
	var receivedMethod, receivedPath string
	adapter := NewRESTAdapter(transportFunc(func(_ context.Context, method, path string, query url.Values) (*http.Response, error) {
		receivedMethod, receivedPath = method, path
		if query != nil {
			t.Errorf("query = %#v, want nil", query)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"name":"project","owner":{"login":"alice"},"description":"demo","private":true,"archived":false,"default_branch":"main"}`))}, nil
	}))

	result, err := adapter.Get(context.Background(), applicationrepository.GetRequest{Owner: "alice/team", Name: "project"})
	if err != nil {
		t.Fatal(err)
	}
	if receivedMethod != http.MethodGet || receivedPath != "/api/v1/repos/alice%2Fteam/project" {
		t.Fatalf("request = %s %s", receivedMethod, receivedPath)
	}
	want := applicationrepository.Repository{Owner: "alice", Name: "project", Description: "demo", Private: true, DefaultBranch: "main"}
	if result != want {
		t.Fatalf("result = %+v, want %+v", result, want)
	}
}

func TestRESTAdapterCreatesRepositoryWithJSONBody(t *testing.T) {
	var method, path string
	var requestBody string
	adapter := NewRESTAdapter(jsonTransportFunc(func(_ context.Context, receivedMethod, receivedPath string, query url.Values, body []byte) (*http.Response, error) {
		method, path = receivedMethod, receivedPath
		requestBody = string(body)
		if query != nil {
			t.Errorf("query = %#v", query)
		}
		return &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(strings.NewReader(`{"name":"project","owner":{"login":"alice"},"private":true}`))}, nil
	}))
	result, err := adapter.Create(context.Background(), applicationrepository.CreateRequest{Name: "project", Description: "demo", Private: true})
	if err != nil {
		t.Fatal(err)
	}
	if method != http.MethodPost || path != "/api/v1/user/repos" {
		t.Fatalf("request = %s %s", method, path)
	}
	if !strings.Contains(requestBody, `"name":"project"`) || !strings.Contains(requestBody, `"description":"demo"`) || !strings.Contains(requestBody, `"private":true`) {
		t.Fatalf("body = %s", requestBody)
	}
	if result.Owner != "alice" || result.Name != "project" || !result.Private {
		t.Fatalf("result = %+v", result)
	}
}

func TestRESTAdapterTranslatesFailuresSafely(t *testing.T) {
	const secret = "secret-token"
	tests := []struct {
		name string
		err  error
		body string
	}{
		{name: "typed transport error", err: transportError{operation: "request", status: http.StatusUnauthorized, secret: secret}},
		{name: "raw transport error", err: errors.New(secret)},
		{name: "malformed JSON", body: secret + " not json"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			adapter := NewRESTAdapter(transportFunc(func(context.Context, string, string, url.Values) (*http.Response, error) {
				if test.err != nil {
					return nil, test.err
				}
				return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(test.body))}, nil
			}))

			_, err := adapter.List(context.Background(), applicationrepository.ListRequest{})
			var remoteError applicationrepository.RemoteError
			if !errors.As(err, &remoteError) {
				t.Fatalf("List() error = %v, want Application RemoteError", err)
			}
			if strings.Contains(err.Error(), secret) || strings.Contains(remoteError.Error(), secret) {
				t.Errorf("translated error exposes sensitive data: %q", err)
			}
		})
	}
}
