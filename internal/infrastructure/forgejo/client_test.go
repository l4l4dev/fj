package forgejo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	applicationauth "github.com/l4l4dev/fj/internal/application/auth"
	"github.com/l4l4dev/fj/internal/application/config"
)

type doerFunc func(*http.Request) (*http.Response, error)

func (doer doerFunc) Do(request *http.Request) (*http.Response, error) {
	return doer(request)
}

func TestClientBuildsAuthenticatedRequest(t *testing.T) {
	const secret = "secret-token"
	var received *http.Request
	client := NewClient(config.Instance{Endpoint: "https://forgejo.example/"}, applicationauth.NewCredential(secret), "0.1.0", doerFunc(func(request *http.Request) (*http.Response, error) {
		received = request
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("[]"))}, nil
	}))

	response, err := client.Do(context.Background(), http.MethodGet, "/api/v1/user/repos", url.Values{"page": {"2"}, "limit": {"10"}})
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if received.URL.String() != "https://forgejo.example/api/v1/user/repos?limit=10&page=2" {
		t.Errorf("request URL = %q", received.URL.String())
	}
	if received.Header.Get("Authorization") != "token "+secret {
		t.Errorf("authorization header = %q", received.Header.Get("Authorization"))
	}
	if received.Header.Get("User-Agent") != "fj/0.1.0" {
		t.Errorf("user-agent = %q", received.Header.Get("User-Agent"))
	}
	if strings.Contains(received.URL.String(), secret) {
		t.Errorf("request URL exposes credential: %q", received.URL.String())
	}
}

func TestClientReturnsSafeRemoteErrors(t *testing.T) {
	const secret = "secret-token"
	tests := []struct {
		name       string
		statusCode int
		cause      error
	}{
		{name: "non success response", statusCode: http.StatusUnauthorized},
		{name: "transport failure", cause: errors.New(secret)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient(config.Instance{Endpoint: "https://forgejo.example"}, applicationauth.NewCredential(secret), "0.1.0", doerFunc(func(*http.Request) (*http.Response, error) {
				if test.cause != nil {
					return nil, test.cause
				}
				return &http.Response{StatusCode: test.statusCode, Body: io.NopCloser(strings.NewReader(secret))}, nil
			}))

			_, err := client.Do(context.Background(), http.MethodGet, "api/v1/user/repos", nil)
			var remoteError RemoteError
			if !errors.As(err, &remoteError) {
				t.Fatalf("Do() error = %v, want RemoteError", err)
			}
			if remoteError.Operation() != "request" {
				t.Errorf("RemoteError operation = %q, want request", remoteError.Operation())
			}
			if remoteError.StatusCode() != test.statusCode {
				t.Errorf("RemoteError status = %d, want %d", remoteError.StatusCode(), test.statusCode)
			}
			if strings.Contains(err.Error(), secret) {
				t.Errorf("remote error exposes sensitive value: %q", err)
			}
			if errors.Is(err, test.cause) {
				t.Errorf("RemoteError exposes raw transport cause")
			}
			if strings.Contains(fmt.Sprintf("%#v", remoteError), secret) {
				t.Errorf("RemoteError diagnostic exposes sensitive value")
			}
		})
	}
}

func TestDefaultHTTPClientHasTimeoutAndCrossHostRedirectProtection(t *testing.T) {
	httpClient := newHTTPClient()
	if httpClient.Timeout != defaultHTTPTimeout {
		t.Errorf("timeout = %s, want %s", httpClient.Timeout, defaultHTTPTimeout)
	}
	previous, _ := http.NewRequest(http.MethodGet, "https://forgejo.example/api", nil)
	redirect, _ := http.NewRequest(http.MethodGet, "https://other.example/api", nil)
	if err := httpClient.CheckRedirect(redirect, []*http.Request{previous}); !errors.Is(err, http.ErrUseLastResponse) {
		t.Errorf("CheckRedirect() error = %v, want http.ErrUseLastResponse", err)
	}
}

func TestClientAcceptsInjectedHTTPClient(t *testing.T) {
	called := false
	client := NewClient(config.Instance{Endpoint: "https://forgejo.example"}, applicationauth.NewCredential("token"), "test", doerFunc(func(*http.Request) (*http.Response, error) {
		called = true
		return &http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(strings.NewReader(""))}, nil
	}))
	response, err := client.Do(context.Background(), http.MethodGet, "api/v1/user/repos", nil)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if !called {
		t.Error("injected HTTP client was not called")
	}
}
