package cli

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationconfig "github.com/l4l4dev/fj/internal/application/config"
)

type tokenVerifierStub struct {
	instance applicationconfig.Instance
	token    string
	username string
	err      error
	calls    int
}

func (stub *tokenVerifierStub) Verify(_ context.Context, instance applicationconfig.Instance, token string) (string, error) {
	stub.calls++
	stub.instance = instance
	stub.token = token
	return stub.username, stub.err
}

type credentialStoreStub struct {
	instanceName string
	token        string
	path         string
	calls        int
}

func (stub *credentialStoreStub) Save(instanceName, token string) (string, error) {
	stub.calls++
	stub.instanceName = instanceName
	stub.token = token
	return stub.path, nil
}

func withConfiguredInstance(t *testing.T) {
	t.Helper()
	configHome := t.TempDir()
	path := filepath.Join(configHome, "fj", "config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	contents := "[[instances]]\nname = \"playground\"\nendpoint = \"https://forgejo.example\"\n"
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", configHome)
}

func TestAuthLoginStoresTokenFromStandardInput(t *testing.T) {
	withConfiguredInstance(t)
	verifier := &tokenVerifierStub{username: "example-owner"}
	store := &credentialStoreStub{path: "/tmp/config/fj/credentials.toml"}
	command := newAuthLoginCommand(verifier, store)
	command.SetArgs([]string{"--token-stdin"})
	command.SetIn(strings.NewReader("supplied-value\n"))
	var output, errorOutput strings.Builder
	command.SetOut(&output)
	command.SetErr(&errorOutput)

	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	want := "Logged in to https://forgejo.example as example-owner\nToken stored in /tmp/config/fj/credentials.toml\n"
	if output.String() != want {
		t.Fatalf("unexpected output: %q", output.String())
	}
	if verifier.token != "supplied-value" || verifier.instance.Name != "playground" {
		t.Errorf("verifier received instance %q", verifier.instance.Name)
	}
	if store.instanceName != "playground" || store.token != "supplied-value" {
		t.Errorf("store received instance %q", store.instanceName)
	}
}

func TestAuthLoginRejectsBlankToken(t *testing.T) {
	withConfiguredInstance(t)
	verifier := &tokenVerifierStub{username: "example-owner"}
	store := &credentialStoreStub{}
	command := newAuthLoginCommand(verifier, store)
	command.SetArgs([]string{"--token-stdin"})
	command.SetIn(strings.NewReader("   \n"))
	var output strings.Builder
	command.SetOut(&output)
	command.SetErr(&strings.Builder{})

	err := command.Execute()
	if err == nil || err.Error() != "auth login: token is required" {
		t.Fatalf("Execute() error = %v, want a validation error", err)
	}
	if verifier.calls != 0 || store.calls != 0 {
		t.Error("Execute() used its dependencies for a blank token")
	}
	if strings.Contains(output.String(), "Logged in to") {
		t.Errorf("Execute() printed success output: %q", output.String())
	}
}

func TestAuthLoginRequiresTokenStdinWithoutTerminal(t *testing.T) {
	withConfiguredInstance(t)
	verifier := &tokenVerifierStub{username: "example-owner"}
	store := &credentialStoreStub{}
	command := newAuthLoginCommand(verifier, store)
	command.SetArgs(nil)
	command.SetIn(strings.NewReader("supplied-value\n"))
	command.SetOut(&strings.Builder{})
	command.SetErr(&strings.Builder{})

	err := command.Execute()
	if err == nil || !strings.Contains(err.Error(), "invalid input") {
		t.Fatalf("Execute() error = %v, want a validation error", err)
	}
	if !strings.Contains(errors.Unwrap(err).Error(), "use --token-stdin") {
		t.Errorf("Execute() cause = %v, want the --token-stdin hint", errors.Unwrap(err))
	}
	if verifier.calls != 0 || store.calls != 0 {
		t.Error("Execute() used its dependencies without a token source")
	}
}

func TestReadMaskedTokenEchoesMaskPerByteAndTrims(t *testing.T) {
	var out strings.Builder
	token, err := readMaskedToken(bufio.NewReader(strings.NewReader("ab\n")), &out)
	if err != nil {
		t.Fatal(err)
	}
	if token != "ab" {
		t.Errorf("token = %q, want %q", token, "ab")
	}
	if out.String() != "●●\r\n" {
		t.Errorf("echoed = %q, want %q", out.String(), "●●\r\n")
	}
}

func TestReadMaskedTokenBackspaceErasesLastByte(t *testing.T) {
	var out strings.Builder
	token, err := readMaskedToken(bufio.NewReader(strings.NewReader("ab\x7fc\n")), &out)
	if err != nil {
		t.Fatal(err)
	}
	if token != "ac" {
		t.Errorf("token = %q, want %q", token, "ac")
	}
	if out.String() != "●●\b \b●\r\n" {
		t.Errorf("echoed = %q, want %q", out.String(), "●●\b \b●\r\n")
	}
}

func TestReadMaskedTokenBackspaceOnEmptyBufferIsNoop(t *testing.T) {
	var out strings.Builder
	token, err := readMaskedToken(bufio.NewReader(strings.NewReader("\x7fa\n")), &out)
	if err != nil {
		t.Fatal(err)
	}
	if token != "a" {
		t.Errorf("token = %q, want %q", token, "a")
	}
}

func TestReadMaskedTokenCtrlCInterrupts(t *testing.T) {
	var out strings.Builder
	_, err := readMaskedToken(bufio.NewReader(strings.NewReader("ab\x03cd\n")), &out)
	if !errors.Is(err, errTokenEntryInterrupted) {
		t.Fatalf("err = %v, want errTokenEntryInterrupted", err)
	}
}

func TestAuthLoginDoesNotStoreRejectedToken(t *testing.T) {
	withConfiguredInstance(t)
	verifier := &tokenVerifierStub{err: apperror.New(apperror.Authentication, "auth login", "token was rejected by the instance")}
	store := &credentialStoreStub{}
	command := newAuthLoginCommand(verifier, store)
	command.SetArgs([]string{"--token-stdin"})
	command.SetIn(strings.NewReader("supplied-value\n"))
	var output strings.Builder
	command.SetOut(&output)
	command.SetErr(&strings.Builder{})

	err := command.Execute()
	if err == nil || err.Error() != "auth login: token was rejected by the instance" {
		t.Fatalf("Execute() error = %v, want an authentication error", err)
	}
	if store.calls != 0 {
		t.Error("Execute() stored a token that failed verification")
	}
	if strings.Contains(output.String(), "Logged in to") {
		t.Errorf("Execute() printed success output: %q", output.String())
	}
	if strings.Contains(err.Error(), "supplied-value") {
		t.Errorf("Execute() error exposes the token: %q", err)
	}
}
