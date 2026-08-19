package auth

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
	"github.com/l4l4dev/fj/internal/application/config"
)

// TokenVerifier checks a token against an instance and returns the user name it
// authenticates as.
type TokenVerifier interface {
	Verify(ctx context.Context, instance config.Instance, token string) (string, error)
}

// CredentialStore persists a token for an instance profile and returns the
// storage path so the caller can report where it was written.
type CredentialStore interface {
	Save(instanceName, token string) (string, error)
}

type LoginRequest struct {
	Instance config.Instance
	Token    string
}

type LoginResult struct {
	Username string
	Path     string
}

type LoginUseCase struct {
	verifier TokenVerifier
	store    CredentialStore
}

func NewLoginUseCase(verifier TokenVerifier, store CredentialStore) LoginUseCase {
	return LoginUseCase{verifier: verifier, store: store}
}

func (useCase LoginUseCase) Execute(ctx context.Context, request LoginRequest) (LoginResult, error) {
	token := strings.TrimSpace(request.Token)
	if token == "" {
		return LoginResult{}, apperror.NewValidation("auth login", "token is required")
	}
	if useCase.verifier == nil || useCase.store == nil {
		return LoginResult{}, apperror.New(apperror.Internal, "auth login", "")
	}

	username, err := useCase.verifier.Verify(ctx, request.Instance, token)
	if err != nil {
		return LoginResult{}, err
	}

	path, err := useCase.store.Save(request.Instance.Name, token)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{Username: username, Path: path}, nil
}
