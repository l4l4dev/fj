package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationauth "github.com/l4l4dev/fj/internal/application/auth"
	"github.com/l4l4dev/fj/internal/application/config"
	infrastructureconfig "github.com/l4l4dev/fj/internal/infrastructure/config"
)

type credentialsFile struct {
	Instances map[string]credentialsEntry `toml:"instances"`
}

type credentialsEntry struct {
	Token string `toml:"token"`
}

// FileProvider reads tokens stored by `fj auth login`.
type FileProvider struct {
	path func() (string, error)
}

func NewFileProvider() FileProvider {
	return FileProvider{path: infrastructureconfig.CredentialsPath}
}

func (provider FileProvider) Credential(_ context.Context, instance config.Instance) (applicationauth.Credential, error) {
	path, err := provider.path()
	if err != nil {
		return applicationauth.Credential{}, err
	}
	credentials, err := readCredentials(path)
	if err != nil {
		return applicationauth.Credential{}, err
	}
	entry, ok := credentials.Instances[instance.Name]
	if !ok || entry.Token == "" {
		return applicationauth.Credential{}, applicationauth.ErrCredentialUnavailable
	}
	return applicationauth.NewCredential(entry.Token), nil
}

// FileStore writes tokens to the credentials file with owner-only permissions.
type FileStore struct {
	path func() (string, error)
}

func NewFileStore() FileStore {
	return FileStore{path: infrastructureconfig.CredentialsPath}
}

func (store FileStore) Save(instanceName, token string) (string, error) {
	path, err := store.path()
	if err != nil {
		return "", err
	}
	credentials, err := readCredentials(path)
	if err != nil && !errors.Is(err, applicationauth.ErrCredentialUnavailable) {
		return "", apperror.New(apperror.Internal, "auth login", err.Error())
	}
	if credentials.Instances == nil {
		credentials.Instances = map[string]credentialsEntry{}
	}
	credentials.Instances[instanceName] = credentialsEntry{Token: token}

	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return "", apperror.New(apperror.Internal, "auth login", "could not create the credentials directory at "+directory)
	}

	temporaryPath := path + fmt.Sprintf(".tmp%d", os.Getpid())
	writeError := apperror.New(apperror.Internal, "auth login", "could not write the credentials file at "+path)
	file, err := os.OpenFile(temporaryPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return "", writeError
	}
	if err := toml.NewEncoder(file).Encode(credentials); err != nil {
		file.Close()
		os.Remove(temporaryPath)
		return "", writeError
	}
	if err := file.Close(); err != nil {
		os.Remove(temporaryPath)
		return "", writeError
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		os.Remove(temporaryPath)
		return "", writeError
	}
	return path, nil
}

// readCredentials reports ErrCredentialUnavailable when the file is absent, and
// a path-only error for any other failure so token material is never exposed.
func readCredentials(path string) (credentialsFile, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return credentialsFile{}, applicationauth.ErrCredentialUnavailable
		}
		return credentialsFile{}, fmt.Errorf("read credentials file %s", path)
	}
	var decoded credentialsFile
	if _, err := toml.Decode(string(contents), &decoded); err != nil {
		return credentialsFile{}, fmt.Errorf("credentials file contains malformed TOML: %s", path)
	}
	return decoded, nil
}

var (
	_ applicationauth.Provider        = FileProvider{}
	_ applicationauth.CredentialStore = FileStore{}
)
