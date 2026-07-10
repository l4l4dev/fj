package cli

import (
	"bytes"
	"context"
	"errors"
	"testing"

	applicationrepository "github.com/l4l4dev/fj/internal/application/repository"
)

type repositoryService struct {
	result  []applicationrepository.Repository
	request applicationrepository.ListRequest
	err     error
}

func (service *repositoryService) List(_ context.Context, request applicationrepository.ListRequest) ([]applicationrepository.Repository, error) {
	service.request = request
	return service.result, service.err
}

func TestRepositoryListCommandPrintsRepositories(t *testing.T) {
	service := &repositoryService{result: []applicationrepository.Repository{{Owner: "alice", Name: "project"}}}
	var output bytes.Buffer
	command := NewRootCommandWithRepositoryService(service)
	command.SetOut(&output)
	command.SetArgs([]string{"repo", "list", "--page", "2", "--limit", "4"})

	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if output.String() != "alice/project\n" {
		t.Fatalf("output = %q", output.String())
	}
	if service.request != (applicationrepository.ListRequest{Page: 2, Limit: 4}) {
		t.Fatalf("request = %+v", service.request)
	}
}

func TestRepositoryListCommandPrintsEmptyResult(t *testing.T) {
	service := &repositoryService{}
	var output bytes.Buffer
	command := NewRootCommandWithRepositoryService(service)
	command.SetOut(&output)
	command.SetArgs([]string{"repo", "list"})

	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	if output.String() != "No repositories found.\n" {
		t.Fatalf("output = %q", output.String())
	}
	if service.request != (applicationrepository.ListRequest{Page: 1, Limit: 30}) {
		t.Fatalf("request = %+v", service.request)
	}
}

func TestRepositoryListCommandRejectsNonPositiveFlags(t *testing.T) {
	for _, args := range [][]string{{"repo", "list", "--page", "0"}, {"repo", "list", "--limit", "0"}} {
		command := NewRootCommandWithRepositoryService(&repositoryService{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("args %v were accepted", args)
		}
	}
}

func TestRepositoryListCommandMapsRemoteError(t *testing.T) {
	for _, test := range []struct {
		name     string
		status   int
		category errorCategory
	}{
		{name: "unauthorized", status: 401, category: categoryAuthentication},
		{name: "forbidden", status: 403, category: categoryAuthentication},
		{name: "other remote status", status: 503, category: categoryRemote},
	} {
		t.Run(test.name, func(t *testing.T) {
			service := &repositoryService{err: applicationrepository.NewRemoteError("list repositories", test.status)}
			command := NewRootCommandWithRepositoryService(service)
			command.SetArgs([]string{"repo", "list"})

			err := command.Execute()
			var classified commandError
			if !errors.As(err, &classified) || classified.category != test.category {
				t.Fatalf("error = %v, category = %v", err, classified.category)
			}
		})
	}
}
