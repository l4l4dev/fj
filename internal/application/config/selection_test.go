package config

import "testing"

func TestConfigurationSelectInstance(t *testing.T) {
	work := Instance{Name: "work", Endpoint: "https://work.example", Credential: "work-token"}
	personal := Instance{Name: "personal", Endpoint: "https://personal.example", Credential: "personal-token"}

	tests := []struct {
		name          string
		configuration Configuration
		explicitName  string
		want          Instance
		wantError     string
	}{
		{name: "explicit profile", configuration: Configuration{Instances: []Instance{work, personal}}, explicitName: "personal", want: personal},
		{name: "sole profile", configuration: Configuration{Instances: []Instance{work}}, want: work},
		{name: "missing profile", configuration: Configuration{Instances: []Instance{work}}, explicitName: "missing", wantError: "instance profile \"missing\" not found"},
		{name: "ambiguous profiles", configuration: Configuration{Instances: []Instance{work, personal}}, wantError: "instance selection is ambiguous: specify one of 2 configured profiles"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			instance, err := test.configuration.SelectInstance(test.explicitName)
			if test.wantError != "" {
				if err == nil || err.Error() != test.wantError {
					t.Fatalf("SelectInstance() error = %v, want %q", err, test.wantError)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if instance != test.want {
				t.Errorf("SelectInstance() = %#v, want %#v", instance, test.want)
			}
		})
	}
}
