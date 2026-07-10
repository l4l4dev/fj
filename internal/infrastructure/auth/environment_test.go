package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	applicationauth "github.com/l4l4dev/fj/internal/application/auth"
	"github.com/l4l4dev/fj/internal/application/config"
)

func TestEnvironmentProviderSuppliesCredentialThroughApplicationResolver(t *testing.T) {
	const secret = "secret-token"
	provider := EnvironmentProvider{lookupEnv: func(name string) (string, bool) {
		if name == "FORGEJO_TOKEN" {
			return secret, true
		}
		return "", false
	}}

	credential, err := applicationauth.NewResolver(provider).Resolve(context.Background(), "FORGEJO_TOKEN")
	if err != nil {
		t.Fatal(err)
	}
	if credential.Value() != secret {
		t.Errorf("credential value = %q, want supplied value", credential.Value())
	}
	for _, diagnostic := range []string{fmt.Sprint(credential), fmt.Sprintf("%#v", credential)} {
		if diagnostic != "[REDACTED]" || strings.Contains(diagnostic, secret) {
			t.Errorf("credential diagnostic = %q", diagnostic)
		}
	}
}

func TestEnvironmentProviderRejectsUnavailableCredentialsSafely(t *testing.T) {
	tests := []struct {
		name      string
		reference config.CredentialReference
		lookup    func(string) (string, bool)
	}{
		{
			name:      "empty reference",
			reference: "",
			lookup:    func(string) (string, bool) { return "secret-token", true },
		},
		{
			name:      "missing variable",
			reference: "FORGEJO_TOKEN",
			lookup:    func(string) (string, bool) { return "", false },
		},
		{
			name:      "empty variable",
			reference: "FORGEJO_TOKEN",
			lookup:    func(string) (string, bool) { return "", true },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider := EnvironmentProvider{lookupEnv: test.lookup}
			_, err := applicationauth.NewResolver(provider).Resolve(context.Background(), test.reference)
			if !errors.Is(err, applicationauth.ErrCredentialUnavailable) {
				t.Fatalf("Resolve() error = %v, want ErrCredentialUnavailable", err)
			}
			if strings.Contains(err.Error(), "secret-token") || (test.reference != "" && strings.Contains(err.Error(), string(test.reference))) {
				t.Errorf("Resolve() error exposes sensitive data: %q", err)
			}
		})
	}
}
