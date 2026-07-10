package config

import (
	"strings"
	"testing"
)

func TestConfigurationValidate(t *testing.T) {
	validInstance := Instance{
		Name:       "work",
		Endpoint:   "https://forgejo.example",
		Credential: "super-secret-reference",
	}
	tests := []struct {
		name          string
		configuration Configuration
		wantError     string
	}{
		{name: "valid", configuration: Configuration{Instances: []Instance{validInstance}}},
		{name: "no instances", configuration: Configuration{}, wantError: "configuration must contain at least one instance"},
		{name: "missing name", configuration: Configuration{Instances: []Instance{{Endpoint: validInstance.Endpoint, Credential: validInstance.Credential}}}, wantError: "instance 1: name is required"},
		{name: "duplicate name", configuration: Configuration{Instances: []Instance{validInstance, validInstance}}, wantError: "instance \"work\": name must be unique"},
		{name: "missing endpoint", configuration: Configuration{Instances: []Instance{{Name: "work", Credential: validInstance.Credential}}}, wantError: "instance \"work\": endpoint is required"},
		{name: "invalid endpoint", configuration: Configuration{Instances: []Instance{{Name: "work", Endpoint: "forgejo.example", Credential: validInstance.Credential}}}, wantError: "instance \"work\": endpoint must be an absolute HTTP or HTTPS URL"},
		{name: "endpoint credentials", configuration: Configuration{Instances: []Instance{{Name: "work", Endpoint: "https://user:secret@forgejo.example", Credential: validInstance.Credential}}}, wantError: "instance \"work\": endpoint must not contain credentials"},
		{name: "missing credential reference", configuration: Configuration{Instances: []Instance{{Name: "work", Endpoint: validInstance.Endpoint}}}, wantError: "instance \"work\": credential reference is required"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.configuration.Validate()
			if test.wantError == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v", err)
				}
				return
			}
			if err == nil || err.Error() != test.wantError {
				t.Fatalf("Validate() error = %v, want %q", err, test.wantError)
			}
			for _, sensitiveValue := range []string{"secret", string(validInstance.Credential)} {
				if strings.Contains(err.Error(), sensitiveValue) {
					t.Errorf("Validate() error exposes sensitive value %q", sensitiveValue)
				}
			}
		})
	}
}
