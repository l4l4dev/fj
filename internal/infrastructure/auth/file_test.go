package auth

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	applicationauth "github.com/l4l4dev/fj/internal/application/auth"
	"github.com/l4l4dev/fj/internal/application/config"
)

func credentialsPathIn(t *testing.T) func() (string, error) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fj", "credentials.toml")
	return func() (string, error) { return path, nil }
}

func TestFileStoreWritesOwnerOnlyCredentials(t *testing.T) {
	path := credentialsPathIn(t)
	store := FileStore{path: path}

	written, err := store.Save("work", "stored-value")
	if err != nil {
		t.Fatal(err)
	}
	expected, _ := path()
	if written != expected {
		t.Fatalf("Save() path = %q, want %q", written, expected)
	}
	info, err := os.Stat(written)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("credentials file mode = %v, want -rw-------", info.Mode().Perm())
	}
	directory, err := os.Stat(filepath.Dir(written))
	if err != nil {
		t.Fatal(err)
	}
	if directory.Mode().Perm() != 0o700 {
		t.Errorf("credentials directory mode = %v, want drwx------", directory.Mode().Perm())
	}
}

func TestFileStoreRoundTripsThroughFileProvider(t *testing.T) {
	path := credentialsPathIn(t)
	if _, err := (FileStore{path: path}).Save("work", "stored-value"); err != nil {
		t.Fatal(err)
	}

	credential, err := (FileProvider{path: path}).Credential(context.Background(), config.Instance{Name: "work"})
	if err != nil {
		t.Fatal(err)
	}
	if credential.Value() != "stored-value" {
		t.Errorf("credential value = %q, want the stored value", credential.Value())
	}
}

func TestFileStorePreservesOtherInstances(t *testing.T) {
	path := credentialsPathIn(t)
	store := FileStore{path: path}
	if _, err := store.Save("personal", "personal-value"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Save("work", "work-value"); err != nil {
		t.Fatal(err)
	}

	provider := FileProvider{path: path}
	for name, want := range map[string]string{"personal": "personal-value", "work": "work-value"} {
		credential, err := provider.Credential(context.Background(), config.Instance{Name: name})
		if err != nil {
			t.Fatalf("Credential(%q) error = %v", name, err)
		}
		if credential.Value() != want {
			t.Errorf("Credential(%q) = %q, want %q", name, credential.Value(), want)
		}
	}
}

func TestFileStoreReplacesExistingToken(t *testing.T) {
	path := credentialsPathIn(t)
	store := FileStore{path: path}
	if _, err := store.Save("work", "old-value"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Save("work", "new-value"); err != nil {
		t.Fatal(err)
	}

	credential, err := (FileProvider{path: path}).Credential(context.Background(), config.Instance{Name: "work"})
	if err != nil {
		t.Fatal(err)
	}
	if credential.Value() != "new-value" {
		t.Errorf("credential value = %q, want the replacement value", credential.Value())
	}
}

func TestFileProviderReportsUnavailableCredentials(t *testing.T) {
	tests := []struct {
		name         string
		contents     string
		instanceName string
		writeFile    bool
	}{
		{name: "missing file", instanceName: "work"},
		{name: "unknown instance", contents: "[instances.\"personal\"]\ntoken = \"stored-value\"\n", instanceName: "work", writeFile: true},
		{name: "empty token", contents: "[instances.\"work\"]\ntoken = \"\"\n", instanceName: "work", writeFile: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := credentialsPathIn(t)
			if test.writeFile {
				writeCredentials(t, path, test.contents)
			}
			_, err := (FileProvider{path: path}).Credential(context.Background(), config.Instance{Name: test.instanceName})
			if !errors.Is(err, applicationauth.ErrCredentialUnavailable) {
				t.Fatalf("Credential() error = %v, want ErrCredentialUnavailable", err)
			}
		})
	}
}

func TestFileProviderReportsMalformedFileWithoutTokenMaterial(t *testing.T) {
	path := credentialsPathIn(t)
	writeCredentials(t, path, "[instances.\"work\"]\ntoken = \"stored-value\"\n[")

	_, err := (FileProvider{path: path}).Credential(context.Background(), config.Instance{Name: "work"})
	if err == nil || errors.Is(err, applicationauth.ErrCredentialUnavailable) {
		t.Fatalf("Credential() error = %v, want a real failure", err)
	}
	expected, _ := path()
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Credential() error = %q, want it to name the path", err)
	}
	if strings.Contains(err.Error(), "stored-value") {
		t.Errorf("Credential() error exposes token material: %q", err)
	}
}

func writeCredentials(t *testing.T, path func() (string, error), contents string) {
	t.Helper()
	resolved, err := path()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(resolved), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(resolved, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
