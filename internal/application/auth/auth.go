package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/l4l4dev/fj/internal/application/config"
)

var ErrCredentialUnavailable = errors.New("credential unavailable")

type Provider interface {
	Credential(context.Context, config.Instance) (Credential, error)
}

type Resolver struct {
	providers []Provider
}

func NewResolver(providers ...Provider) Resolver {
	return Resolver{providers: providers}
}

// Resolve asks each provider in order and returns the first credential found.
// A provider reporting ErrCredentialUnavailable is skipped; any other failure
// aborts resolution.
func (resolver Resolver) Resolve(ctx context.Context, instance config.Instance) (Credential, error) {
	for _, provider := range resolver.providers {
		credential, err := provider.Credential(ctx, instance)
		if err == nil {
			return credential, nil
		}
		if !errors.Is(err, ErrCredentialUnavailable) {
			return Credential{}, err
		}
	}
	return Credential{}, fmt.Errorf("resolve credential: %w", ErrCredentialUnavailable)
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
