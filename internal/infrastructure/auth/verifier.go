package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationauth "github.com/l4l4dev/fj/internal/application/auth"
	"github.com/l4l4dev/fj/internal/application/config"
	"github.com/l4l4dev/fj/internal/infrastructure/forgejo"
)

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

// ForgejoVerifier checks a token by requesting the authenticated user.
type ForgejoVerifier struct {
	version    string
	httpClient doer
}

func NewForgejoVerifier(version string) ForgejoVerifier {
	return ForgejoVerifier{version: version}
}

func (verifier ForgejoVerifier) Verify(ctx context.Context, instance config.Instance, token string) (string, error) {
	client := forgejo.NewClient(instance, applicationauth.NewCredential(token), verifier.version, verifier.httpClient)
	response, err := client.Do(ctx, http.MethodGet, "/api/v1/user", nil)
	if err != nil {
		return "", translateVerifyError(err)
	}
	defer response.Body.Close()

	var decoded struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil || decoded.Login == "" {
		return "", apperror.New(apperror.Remote, "auth login", "could not reach the instance")
	}
	return decoded.Login, nil
}

func translateVerifyError(err error) error {
	var remote forgejo.RemoteError
	if errors.As(err, &remote) {
		switch remote.StatusCode() {
		case http.StatusUnauthorized, http.StatusForbidden:
			return apperror.New(apperror.Authentication, "auth login", "token was rejected by the instance")
		}
	}
	return apperror.New(apperror.Remote, "auth login", "could not reach the instance")
}

var _ applicationauth.TokenVerifier = ForgejoVerifier{}
