package target

import (
	"fmt"
	"strings"

	"github.com/l4l4dev/fj/internal/application/config"
)

type Repository struct {
	Owner string
	Name  string
}

type Request struct {
	InstanceName       string
	ExplicitRepository *Repository
	DetectedRepository *Repository
}

type Context struct {
	Instance   config.Profile
	Repository Repository
}

func Resolve(configuration config.Configuration, request Request) (Context, error) {
	instance, err := configuration.SelectInstance(request.InstanceName)
	if err != nil {
		return Context{}, err
	}

	repository, err := resolveRepository(request.ExplicitRepository, request.DetectedRepository)
	if err != nil {
		return Context{}, err
	}

	return Context{
		Instance:   config.Profile{Name: instance.Name, Endpoint: instance.Endpoint},
		Repository: repository,
	}, nil
}

func resolveRepository(explicit, detected *Repository) (Repository, error) {
	if explicit == nil && detected == nil {
		return Repository{}, fmt.Errorf("repository context is required")
	}
	if explicit != nil && detected != nil && *explicit != *detected {
		return Repository{}, fmt.Errorf("repository context conflicts: explicit and detected repositories differ")
	}

	repository := detected
	if explicit != nil {
		repository = explicit
	}
	if strings.TrimSpace(repository.Owner) == "" {
		return Repository{}, fmt.Errorf("repository owner is required")
	}
	if strings.TrimSpace(repository.Name) == "" {
		return Repository{}, fmt.Errorf("repository name is required")
	}

	return *repository, nil
}
