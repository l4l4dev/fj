package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/l4l4dev/fj/internal/application/config"
)

type providerFunc func(context.Context, config.CredentialReference) (Credential, error)

func (provider providerFunc) Credential(ctx context.Context, reference config.CredentialReference) (Credential, error) {
	return provider(ctx, reference)
}

func TestResolverSuppliesCredential(t *testing.T) {
	const reference config.CredentialReference = "work-token"
	const secret = "secret-token"
	var receivedReference config.CredentialReference
	resolver := NewResolver(providerFunc(func(_ context.Context, received config.CredentialReference) (Credential, error) {
		receivedReference = received
		return NewCredential(secret), nil
	}))

	credential, err := resolver.Resolve(context.Background(), reference)
	if err != nil {
		t.Fatal(err)
	}
	if receivedReference != reference {
		t.Errorf("provider reference = %q, want %q", receivedReference, reference)
	}
	if credential.Value() != secret {
		t.Errorf("credential value = %q, want supplied value", credential.Value())
	}
	for _, diagnostic := range []string{fmt.Sprint(credential), fmt.Sprintf("%#v", credential)} {
		if strings.Contains(diagnostic, secret) || diagnostic != "[REDACTED]" {
			t.Errorf("credential diagnostic = %q", diagnostic)
		}
	}
}

func TestResolverRedactsProviderError(t *testing.T) {
	resolver := NewResolver(providerFunc(func(context.Context, config.CredentialReference) (Credential, error) {
		return Credential{}, errors.New("secret-token is unavailable")
	}))

	_, err := resolver.Resolve(context.Background(), "work-token")
	if !errors.Is(err, ErrCredentialUnavailable) {
		t.Fatalf("Resolve() error = %v", err)
	}
	if got := err.Error(); got != "resolve credential: credential unavailable" {
		t.Errorf("Resolve() error = %q", got)
	}
	if strings.Contains(err.Error(), "secret-token") {
		t.Errorf("Resolve() error exposes credential material: %q", err)
	}
}
