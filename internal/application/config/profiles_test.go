package config

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestConfigurationListProfiles(t *testing.T) {
	configuration := Configuration{Instances: []Instance{
		{Name: "work", Endpoint: "https://work.example", Credential: "work-secret"},
		{Name: "personal", Endpoint: "https://personal.example", Credential: "personal-secret"},
	}}

	profiles, err := configuration.ListProfiles()
	if err != nil {
		t.Fatal(err)
	}
	want := []Profile{
		{Name: "work", Endpoint: "https://work.example"},
		{Name: "personal", Endpoint: "https://personal.example"},
	}
	if !reflect.DeepEqual(profiles, want) {
		t.Errorf("ListProfiles() = %#v, want %#v", profiles, want)
	}
	if output := fmt.Sprint(profiles); strings.Contains(output, "secret") {
		t.Errorf("ListProfiles() output exposes a secret: %q", output)
	}
}

func TestConfigurationInspectProfile(t *testing.T) {
	configuration := Configuration{Instances: []Instance{
		{Name: "work", Endpoint: "https://work.example", Credential: "work-secret"},
	}}

	profile, err := configuration.InspectProfile("work")
	if err != nil {
		t.Fatal(err)
	}
	want := Profile{Name: "work", Endpoint: "https://work.example"}
	if profile != want {
		t.Errorf("InspectProfile() = %#v, want %#v", profile, want)
	}
	if output := fmt.Sprint(profile); strings.Contains(output, "secret") {
		t.Errorf("InspectProfile() output exposes a secret: %q", output)
	}
}

func TestConfigurationInspectProfileNotFound(t *testing.T) {
	configuration := Configuration{Instances: []Instance{
		{Name: "work", Endpoint: "https://work.example", Credential: "work-secret"},
	}}

	_, err := configuration.InspectProfile("missing")
	if err == nil || err.Error() != "profile \"missing\" not found" {
		t.Fatalf("InspectProfile() error = %v", err)
	}
}

func TestConfigurationListProfilesDoesNotExposeUnsafeValues(t *testing.T) {
	configuration := Configuration{Instances: []Instance{
		{Name: "work", Endpoint: "https://user:endpoint-secret@work.example", Credential: "credential-secret"},
	}}

	_, err := configuration.ListProfiles()
	if err == nil {
		t.Fatal("ListProfiles() error = nil")
	}
	if strings.Contains(err.Error(), "secret") {
		t.Errorf("ListProfiles() error exposes a secret: %q", err)
	}
}
