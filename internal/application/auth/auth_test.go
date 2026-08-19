package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/l4l4dev/fj/internal/application/config"
)

type providerFunc func(context.Context, config.Instance) (Credential, error)

func (provider providerFunc) Credential(ctx context.Context, instance config.Instance) (Credential, error) {
	return provider(ctx, instance)
}

func TestResolverSuppliesCredential(t *testing.T) {
	instance := config.Instance{Name: "work", Endpoint: "https://forgejo.example", Credential: "work-token"}
	const secret = "secret-token"
	var receivedInstance config.Instance
	resolver := NewResolver(providerFunc(func(_ context.Context, received config.Instance) (Credential, error) {
		receivedInstance = received
		return NewCredential(secret), nil
	}))

	credential, err := resolver.Resolve(context.Background(), instance)
	if err != nil {
		t.Fatal(err)
	}
	if receivedInstance != instance {
		t.Errorf("provider instance = %#v, want %#v", receivedInstance, instance)
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

func TestResolverTriesProvidersInOrder(t *testing.T) {
	var calls []string
	resolver := NewResolver(
		providerFunc(func(context.Context, config.Instance) (Credential, error) {
			calls = append(calls, "first")
			return Credential{}, fmt.Errorf("environment: %w", ErrCredentialUnavailable)
		}),
		providerFunc(func(context.Context, config.Instance) (Credential, error) {
			calls = append(calls, "second")
			return NewCredential("stored-value"), nil
		}),
		providerFunc(func(context.Context, config.Instance) (Credential, error) {
			calls = append(calls, "third")
			return NewCredential("never-reached"), nil
		}),
	)

	credential, err := resolver.Resolve(context.Background(), config.Instance{Name: "work"})
	if err != nil {
		t.Fatal(err)
	}
	if credential.Value() != "stored-value" {
		t.Errorf("credential value = %q, want the second provider value", credential.Value())
	}
	if strings.Join(calls, ",") != "first,second" {
		t.Errorf("provider calls = %v, want first then second only", calls)
	}
}

func TestResolverReportsUnavailableWhenNoProviderHasCredential(t *testing.T) {
	resolver := NewResolver(providerFunc(func(context.Context, config.Instance) (Credential, error) {
		return Credential{}, ErrCredentialUnavailable
	}))

	_, err := resolver.Resolve(context.Background(), config.Instance{Name: "work"})
	if !errors.Is(err, ErrCredentialUnavailable) {
		t.Fatalf("Resolve() error = %v", err)
	}
	if got := err.Error(); got != "resolve credential: credential unavailable" {
		t.Errorf("Resolve() error = %q", got)
	}
}

func TestResolverAbortsOnRealProviderError(t *testing.T) {
	failure := errors.New("credentials file is malformed")
	var secondCalled bool
	resolver := NewResolver(
		providerFunc(func(context.Context, config.Instance) (Credential, error) {
			return Credential{}, failure
		}),
		providerFunc(func(context.Context, config.Instance) (Credential, error) {
			secondCalled = true
			return NewCredential("stored-value"), nil
		}),
	)

	_, err := resolver.Resolve(context.Background(), config.Instance{Name: "work"})
	if !errors.Is(err, failure) {
		t.Fatalf("Resolve() error = %v, want the provider error", err)
	}
	if secondCalled {
		t.Error("Resolve() continued after a real provider error")
	}
}
