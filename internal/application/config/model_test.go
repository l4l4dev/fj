package config

import (
	"reflect"
	"testing"
)

func TestConfigurationRepresentsMultipleInstances(t *testing.T) {
	want := []Instance{
		{Name: "work", Endpoint: Endpoint("https://forgejo.work.example"), Credential: "work-token"},
		{Name: "personal", Endpoint: Endpoint("https://forgejo.example"), Credential: "personal-token"},
	}
	configuration := Configuration{Instances: want}

	if !reflect.DeepEqual(configuration.Instances, want) {
		t.Errorf("Instances = %#v, want %#v", configuration.Instances, want)
	}
}
