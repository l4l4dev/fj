package target

import (
	"testing"

	"github.com/l4l4dev/fj/internal/application/config"
)

func TestResolve(t *testing.T) {
	work := config.Instance{Name: "work", Endpoint: "https://work.example", Credential: "work-token"}
	personal := config.Instance{Name: "personal", Endpoint: "https://personal.example", Credential: "personal-token"}
	repository := Repository{Owner: "octo", Name: "project"}
	otherRepository := Repository{Owner: "octo", Name: "other"}

	tests := []struct {
		name          string
		configuration config.Configuration
		request       Request
		want          Context
		wantError     string
	}{
		{
			name:          "explicit instance and repository",
			configuration: config.Configuration{Instances: []config.Instance{work, personal}},
			request:       Request{InstanceName: "personal", ExplicitRepository: &repository, DetectedRepository: &repository},
			want:          Context{Instance: config.Profile{Name: "personal", Endpoint: personal.Endpoint}, Repository: repository},
		},
		{
			name:          "sole instance and detected repository",
			configuration: config.Configuration{Instances: []config.Instance{work}},
			request:       Request{DetectedRepository: &repository},
			want:          Context{Instance: config.Profile{Name: "work", Endpoint: work.Endpoint}, Repository: repository},
		},
		{
			name:          "missing repository",
			configuration: config.Configuration{Instances: []config.Instance{work}},
			wantError:     "repository context is required",
		},
		{
			name:          "conflicting repositories",
			configuration: config.Configuration{Instances: []config.Instance{work}},
			request:       Request{ExplicitRepository: &repository, DetectedRepository: &otherRepository},
			wantError:     "repository context conflicts: explicit and detected repositories differ",
		},
		{
			name:          "missing owner",
			configuration: config.Configuration{Instances: []config.Instance{work}},
			request:       Request{ExplicitRepository: &Repository{Name: "project"}},
			wantError:     "repository owner is required",
		},
		{
			name:          "missing name",
			configuration: config.Configuration{Instances: []config.Instance{work}},
			request:       Request{ExplicitRepository: &Repository{Owner: "octo"}},
			wantError:     "repository name is required",
		},
		{
			name:          "ambiguous instance",
			configuration: config.Configuration{Instances: []config.Instance{work, personal}},
			request:       Request{ExplicitRepository: &repository},
			wantError:     "instance selection is ambiguous: specify one of 2 configured profiles",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, err := Resolve(test.configuration, test.request)
			if test.wantError != "" {
				if err == nil || err.Error() != test.wantError {
					t.Fatalf("Resolve() error = %v, want %q", err, test.wantError)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if context != test.want {
				t.Errorf("Resolve() = %#v, want %#v", context, test.want)
			}
		})
	}
}
