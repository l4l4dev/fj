package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	applicationrepository "github.com/l4l4dev/fj/internal/application/repository"
)

type repositoryService struct {
	result  []applicationrepository.Repository
	request applicationrepository.ListRequest
	err     error
}

func (service *repositoryService) Get(_ context.Context, request applicationrepository.GetRequest) (applicationrepository.Repository, error) {
	if service.err != nil {
		return applicationrepository.Repository{}, service.err
	}
	return applicationrepository.Repository{Owner: request.Owner, Name: request.Name, Description: "example", DefaultBranch: "main"}, nil
}

func (service *repositoryService) Create(_ context.Context, request applicationrepository.CreateRequest) (applicationrepository.Repository, error) {
	if service.err != nil {
		return applicationrepository.Repository{}, service.err
	}
	return applicationrepository.Repository{Owner: "alice", Name: request.Name, Description: request.Description, Private: request.Private}, nil
}

func (service *repositoryService) Update(_ context.Context, request applicationrepository.UpdateRequest) (applicationrepository.Repository, error) {
	if service.err != nil {
		return applicationrepository.Repository{}, service.err
	}
	result := applicationrepository.Repository{Owner: request.Owner, Name: request.Name, Description: "", Archived: false, DefaultBranch: "main"}
	if request.Description != nil {
		result.Description = *request.Description
	}
	if request.Private != nil {
		result.Private = *request.Private
	}
	return result, nil
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

func TestRepositoryInspectCommandPrintsDetails(t *testing.T) {
	var output bytes.Buffer
	command := NewRootCommandWithRepositoryService(&repositoryService{})
	command.SetOut(&output)
	command.SetArgs([]string{"repo", "inspect", "alice/project"})

	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	want := "Repository: alice/project\nDescription: example\nPrivate: false\nArchived: false\nDefault branch: main\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestRepositoryInspectCommandRejectsInvalidTarget(t *testing.T) {
	for _, target := range []string{"project", "alice/", "/project", "alice/project/extra"} {
		command := NewRootCommandWithRepositoryService(&repositoryService{})
		command.SetArgs([]string{"repo", "inspect", target})
		if err := command.Execute(); err == nil {
			t.Fatalf("target %q was accepted", target)
		}
	}
}

func TestRepositoryInspectCommandMapsNotFoundAndAuthenticationErrors(t *testing.T) {
	for _, test := range []struct {
		status   int
		category errorCategory
	}{
		{status: 404, category: categoryRemote},
		{status: 401, category: categoryAuthentication},
		{status: 403, category: categoryAuthentication},
		{status: 500, category: categoryRemote},
	} {
		command := NewRootCommandWithRepositoryService(&repositoryService{err: applicationrepository.NewRemoteError("inspect repository", test.status)})
		command.SetArgs([]string{"repo", "inspect", "alice/project"})
		err := command.Execute()
		var classified commandError
		if !errors.As(err, &classified) || classified.category != test.category {
			t.Fatalf("status %d: error = %v, category = %v", test.status, err, classified.category)
		}
		if test.status == 404 && classified.Error() != "inspect repository: repository not found" {
			t.Fatalf("status 404 leaked detail: %q", classified.Error())
		}
	}
}

func TestRepositoryCreateCommandSendsVisibilityAndDescription(t *testing.T) {
	for _, test := range []struct {
		visibility string
		private    bool
	}{{"public", false}, {"private", true}} {
		var output bytes.Buffer
		command := NewRootCommandWithRepositoryService(&repositoryService{})
		command.SetOut(&output)
		command.SetArgs([]string{"repo", "create", "project", "--description", "demo", "--visibility", test.visibility})
		if err := command.Execute(); err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(output.Bytes(), []byte("Repository: alice/project\nDescription: demo\nPrivate: "+map[bool]string{true: "true", false: "false"}[test.private])) {
			t.Fatalf("output = %q", output.String())
		}
	}
}

func TestRepositoryCreateCommandRejectsInvalidInput(t *testing.T) {
	for _, args := range [][]string{{"repo", "create", "   "}, {"repo", "create", "project", "--visibility", "unknown"}} {
		command := NewRootCommandWithRepositoryService(&repositoryService{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("args %v were accepted", args)
		}
	}
}

func TestRepositoryCreateCommandMapsErrorsSafely(t *testing.T) {
	for _, test := range []struct {
		status int
		want   string
	}{{401, "authentication failed"}, {403, "authentication failed"}, {409, "repository already exists"}, {500, "remote operation failed"}} {
		command := NewRootCommandWithRepositoryService(&repositoryService{err: applicationrepository.NewRemoteError("create repository", test.status)})
		command.SetArgs([]string{"repo", "create", "project"})
		err := command.Execute()
		if err == nil || !strings.Contains(err.Error(), test.want) || strings.Contains(err.Error(), "secret") {
			t.Fatalf("status %d: %v", test.status, err)
		}
	}
}

func TestRepositoryUpdateCommandPrintsChangedFields(t *testing.T) {
	var output bytes.Buffer
	command := NewRootCommandWithRepositoryService(&repositoryService{})
	command.SetOut(&output)
	command.SetArgs([]string{"repo", "update", "alice/project", "--description", "", "--visibility", "public"})
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	want := "Repository: alice/project\nChanged fields: description, visibility\nDescription: -\nPrivate: false\nArchived: false\nDefault branch: main\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestRepositoryUpdateCommandRejectsInvalidInput(t *testing.T) {
	for _, args := range [][]string{{"repo", "update", "alice/project"}, {"repo", "update", "bad"}, {"repo", "update", "alice/project", "--visibility", "other"}} {
		command := NewRootCommandWithRepositoryService(&repositoryService{})
		command.SetArgs(args)
		if err := command.Execute(); err == nil {
			t.Fatalf("args %v accepted", args)
		}
	}
}

func TestRepositoryUpdateCommandMapsErrors(t *testing.T) {
	for _, test := range []struct {
		status int
		want   string
	}{{401, "authentication failed"}, {403, "authentication failed"}, {404, "repository not found"}, {409, "repository update conflict"}, {500, "remote operation failed"}} {
		service := &repositoryService{err: applicationrepository.NewRemoteError("update repository", test.status)}
		command := NewRootCommandWithRepositoryService(service)
		command.SetArgs([]string{"repo", "update", "alice/project", "--visibility", "public"})
		err := command.Execute()
		if err == nil || !strings.Contains(err.Error(), test.want) {
			t.Fatalf("status %d: %v", test.status, err)
		}
	}
}
