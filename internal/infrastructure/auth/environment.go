package auth

import (
	"context"
	"os"

	applicationauth "github.com/l4l4dev/fj/internal/application/auth"
	"github.com/l4l4dev/fj/internal/application/config"
)

type EnvironmentProvider struct {
	lookupEnv func(string) (string, bool)
}

func NewEnvironmentProvider() EnvironmentProvider {
	return EnvironmentProvider{lookupEnv: os.LookupEnv}
}

func (provider EnvironmentProvider) Credential(_ context.Context, reference config.CredentialReference) (applicationauth.Credential, error) {
	if reference == "" || provider.lookupEnv == nil {
		return applicationauth.Credential{}, applicationauth.ErrCredentialUnavailable
	}

	value, ok := provider.lookupEnv(string(reference))
	if !ok || value == "" {
		return applicationauth.Credential{}, applicationauth.ErrCredentialUnavailable
	}
	return applicationauth.NewCredential(value), nil
}

var _ applicationauth.Provider = EnvironmentProvider{}
