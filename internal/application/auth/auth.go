package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/l4l4dev/fj/internal/application/config"
)

var ErrCredentialUnavailable = errors.New("credential unavailable")

type Provider interface {
	Credential(context.Context, config.CredentialReference) (Credential, error)
}

type Resolver struct {
	provider Provider
}

func NewResolver(provider Provider) Resolver {
	return Resolver{provider: provider}
}

func (resolver Resolver) Resolve(ctx context.Context, reference config.CredentialReference) (Credential, error) {
	credential, err := resolver.provider.Credential(ctx, reference)
	if err != nil {
		return Credential{}, fmt.Errorf("resolve credential: %w", ErrCredentialUnavailable)
	}
	return credential, nil
}

type Credential struct {
	value string
}

func NewCredential(value string) Credential {
	return Credential{value: value}
}

func (credential Credential) Value() string {
	return credential.value
}

func (Credential) String() string {
	return "[REDACTED]"
}

func (Credential) GoString() string {
	return "[REDACTED]"
}
