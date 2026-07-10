package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigurationPath(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    string
		wantErr string
	}{
		{
			name: "absolute XDG config home",
			env:  map[string]string{"XDG_CONFIG_HOME": "/tmp/config", "HOME": "/tmp/home"},
			want: "/tmp/config/fj/config.toml",
		},
		{
			name: "home fallback",
			env:  map[string]string{"HOME": "/tmp/home"},
			want: "/tmp/home/.config/fj/config.toml",
		},
		{
			name:    "relative XDG config home",
			env:     map[string]string{"XDG_CONFIG_HOME": "config", "HOME": "/tmp/home"},
			wantErr: "XDG_CONFIG_HOME must be an absolute path",
		},
		{
			name:    "home is unset",
			env:     map[string]string{},
			wantErr: "HOME must be set when XDG_CONFIG_HOME is unset",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := configurationPath(func(name string) (string, bool) {
				value, ok := test.env[name]
				return value, ok
			})
			if test.wantErr != "" {
				if err == nil || err.Error() != test.wantErr {
					t.Fatalf("configurationPath() error = %v, want %q", err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("configurationPath() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestLoadUsesXDGConfiguration(t *testing.T) {
	configHome := t.TempDir()
	home := t.TempDir()
	writeConfiguration(t, filepath.Join(configHome, "fj", "config.toml"), "work", "https://xdg.example", "WORK_TOKEN")
	writeConfiguration(t, filepath.Join(home, ".config", "fj", "config.toml"), "home", "https://home.example", "HOME_TOKEN")

	configuration, err := load(func(name string) (string, bool) {
		switch name {
		case "XDG_CONFIG_HOME":
			return configHome, true
		case "HOME":
			return home, true
		default:
			return "", false
		}
	}, os.ReadFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(configuration.Instances) != 1 || configuration.Instances[0].Name != "work" {
		t.Fatalf("Load() = %#v, want XDG configuration", configuration)
	}
}

func TestLoadConfigurationErrorsAreActionableAndSafe(t *testing.T) {
	home := t.TempDir()
	tests := []struct {
		name        string
		contents    string
		wantError   string
		sensitive   string
		writeConfig bool
	}{
		{
			name:        "missing file",
			wantError:   "configuration file not found",
			writeConfig: false,
		},
		{
			name:        "malformed TOML",
			contents:    "credential = \"secret-token\"\n[",
			wantError:   "configuration file contains malformed TOML",
			sensitive:   "secret-token",
			writeConfig: true,
		},
		{
			name:        "invalid configuration",
			contents:    "[[instances]]\nname = \"work\"\nendpoint = \"not-an-url\"\ncredential = \"secret-token\"\n",
			wantError:   "configuration is invalid: instance \"work\": endpoint must be an absolute HTTP or HTTPS URL",
			sensitive:   "secret-token",
			writeConfig: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configHome := t.TempDir()
			path := filepath.Join(configHome, "fj", "config.toml")
			if test.writeConfig {
				writeRawConfiguration(t, path, test.contents)
			}
			_, err := load(func(name string) (string, bool) {
				if name == "XDG_CONFIG_HOME" {
					return configHome, true
				}
				return home, true
			}, os.ReadFile)
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("Load() error = %v, want text %q", err, test.wantError)
			}
			if test.sensitive != "" && strings.Contains(err.Error(), test.sensitive) {
				t.Errorf("Load() error exposes sensitive value: %q", err)
			}
		})
	}
}

func writeConfiguration(t *testing.T, path, name, endpoint, credential string) {
	t.Helper()
	writeRawConfiguration(t, path, "[[instances]]\nname = \""+name+"\"\nendpoint = \""+endpoint+"\"\ncredential = \""+credential+"\"\n")
}

func writeRawConfiguration(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
