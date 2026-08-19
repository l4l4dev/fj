package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/l4l4dev/fj/internal/application/apperror"
	"github.com/l4l4dev/fj/internal/application/config"
)

type doerFunc func(*http.Request) (*http.Response, error)

func (do doerFunc) Do(request *http.Request) (*http.Response, error) { return do(request) }

func response(statusCode int, body string) *http.Response {
	return &http.Response{StatusCode: statusCode, Body: io.NopCloser(strings.NewReader(body)), Header: http.Header{}}
}

func TestForgejoVerifierReturnsUserName(t *testing.T) {
	var request *http.Request
	verifier := ForgejoVerifier{version: "0.0.0-test", httpClient: doerFunc(func(received *http.Request) (*http.Response, error) {
		request = received
		return response(http.StatusOK, `{"login":"example-owner"}`), nil
	})}

	username, err := verifier.Verify(context.Background(), config.Instance{Name: "work", Endpoint: "https://forgejo.example"}, "supplied-value")
	if err != nil {
		t.Fatal(err)
	}
	if username != "example-owner" {
		t.Errorf("Verify() = %q, want the login name", username)
	}
	if request.URL.String() != "https://forgejo.example/api/v1/user" {
		t.Errorf("request URL = %q", request.URL)
	}
	if request.Header.Get("Authorization") != "token supplied-value" {
		t.Error("request did not carry the supplied token")
	}
}

func TestForgejoVerifierReportsRejectedToken(t *testing.T) {
	for _, statusCode := range []int{http.StatusUnauthorized, http.StatusForbidden} {
		verifier := ForgejoVerifier{version: "0.0.0-test", httpClient: doerFunc(func(*http.Request) (*http.Response, error) {
			return response(statusCode, `{"message":"unauthorized"}`), nil
		})}

		_, err := verifier.Verify(context.Background(), config.Instance{Name: "work", Endpoint: "https://forgejo.example"}, "supplied-value")
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != apperror.Authentication {
			t.Fatalf("Verify() error = %v, want an authentication error", err)
		}
		if err.Error() != "auth login: token was rejected by the instance" {
			t.Errorf("Verify() error = %q", err)
		}
		if strings.Contains(err.Error(), "supplied-value") {
			t.Errorf("Verify() error exposes the token: %q", err)
		}
	}
}

func TestForgejoVerifierReportsUnreachableInstance(t *testing.T) {
	verifier := ForgejoVerifier{version: "0.0.0-test", httpClient: doerFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusInternalServerError, ""), nil
	})}

	_, err := verifier.Verify(context.Background(), config.Instance{Name: "work", Endpoint: "https://forgejo.example"}, "supplied-value")
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.Remote {
		t.Fatalf("Verify() error = %v, want a remote error", err)
	}
}
