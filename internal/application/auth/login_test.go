package auth

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/l4l4dev/fj/internal/application/config"
)

type verifierStub struct {
	instance config.Instance
	token    string
	username string
	err      error
	calls    int
}

func (stub *verifierStub) Verify(_ context.Context, instance config.Instance, token string) (string, error) {
	stub.calls++
	stub.instance = instance
	stub.token = token
	return stub.username, stub.err
}

type storeStub struct {
	instanceName string
	token        string
	path         string
	err          error
	calls        int
}

func (stub *storeStub) Save(instanceName, token string) (string, error) {
	stub.calls++
	stub.instanceName = instanceName
	stub.token = token
	return stub.path, stub.err
}

func TestLoginStoresVerifiedToken(t *testing.T) {
	instance := config.Instance{Name: "work", Endpoint: "https://forgejo.example"}
	verifier := &verifierStub{username: "example-owner"}
	store := &storeStub{path: "/tmp/config/fj/credentials.toml"}

	result, err := NewLoginUseCase(verifier, store).Execute(context.Background(), LoginRequest{Instance: instance, Token: "  supplied-value \n"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Username != "example-owner" || result.Path != store.path {
		t.Fatalf("Execute() result = %#v", result)
	}
	if verifier.instance != instance {
		t.Errorf("verifier instance = %#v, want %#v", verifier.instance, instance)
	}
	if verifier.token != "supplied-value" || store.token != "supplied-value" {
		t.Error("Execute() did not trim the token before use")
	}
	if store.instanceName != instance.Name {
		t.Errorf("store instance name = %q, want %q", store.instanceName, instance.Name)
	}
}

func TestLoginDoesNotStoreRejectedToken(t *testing.T) {
	rejection := errors.New("token was rejected")
	verifier := &verifierStub{err: rejection}
	store := &storeStub{}

	_, err := NewLoginUseCase(verifier, store).Execute(context.Background(), LoginRequest{Instance: config.Instance{Name: "work"}, Token: "supplied-value"})
	if !errors.Is(err, rejection) {
		t.Fatalf("Execute() error = %v, want the verifier error", err)
	}
	if store.calls != 0 {
		t.Error("Execute() stored a token that failed verification")
	}
}

func TestLoginRejectsBlankToken(t *testing.T) {
	verifier := &verifierStub{}
	store := &storeStub{}

	_, err := NewLoginUseCase(verifier, store).Execute(context.Background(), LoginRequest{Instance: config.Instance{Name: "work"}, Token: "   \t "})
	if err == nil || err.Error() != "token is required" {
		t.Fatalf("Execute() error = %v, want a validation error", err)
	}
	if verifier.calls != 0 || store.calls != 0 {
		t.Error("Execute() used its dependencies for a blank token")
	}
}

func TestLoginReportsMissingDependencies(t *testing.T) {
	for name, useCase := range map[string]LoginUseCase{
		"no verifier": NewLoginUseCase(nil, &storeStub{}),
		"no store":    NewLoginUseCase(&verifierStub{}, nil),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := useCase.Execute(context.Background(), LoginRequest{Instance: config.Instance{Name: "work"}, Token: "supplied-value"})
			if err == nil || !strings.Contains(err.Error(), "internal error") {
				t.Fatalf("Execute() error = %v, want an internal error", err)
			}
			if strings.Contains(err.Error(), "supplied-value") {
				t.Errorf("Execute() error exposes the token: %q", err)
			}
		})
	}
}
