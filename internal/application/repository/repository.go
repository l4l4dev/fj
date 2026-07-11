package repository

import (
	"context"
	"github.com/l4l4dev/fj/internal/application/apperror"
	"strconv"
)

type ValidationError = apperror.ValidationError

type Repository struct {
	Owner         string
	Name          string
	Description   string
	Private       bool
	Archived      bool
	DefaultBranch string
}

type ListRequest struct {
	Page  int
	Limit int
}

type Service interface {
	List(context.Context, ListRequest) ([]Repository, error)
}

type GetRequest struct {
	Owner string
	Name  string
}

type Getter interface {
	Get(context.Context, GetRequest) (Repository, error)
}

type CreateRequest struct {
	Name        string
	Description string
	Private     bool
}

type Creator interface {
	Create(context.Context, CreateRequest) (Repository, error)
}

type UpdateRequest struct {
	Owner       string
	Name        string
	Description *string
	Private     *bool
}

type Updater interface {
	Update(context.Context, UpdateRequest) (Repository, error)
}

type ArchiveRequest struct {
	Owner    string
	Name     string
	Archived bool
}

type Archiver interface {
	SetArchived(context.Context, ArchiveRequest) (Repository, error)
}

type Permission string

const (
	PermissionRead  Permission = "read"
	PermissionWrite Permission = "write"
	PermissionAdmin Permission = "admin"
)

type Collaborator struct {
	Username   string
	Permission Permission
}
type RepositoryAccess struct {
	Owner         string
	Name          string
	Collaborators []Collaborator
}
type AccessRequest struct {
	Owner string
	Name  string
}

type AccessViewer interface {
	ViewAccess(context.Context, AccessRequest) (RepositoryAccess, error)
}

type RemoteError struct {
	operation  string
	statusCode int
	category   apperror.Category
}

func NewRemoteError(operation string, statusCode int) RemoteError {
	category := apperror.Remote
	switch statusCode {
	case 401, 403:
		category = apperror.Authentication
	case 404:
		category = apperror.NotFound
	case 409:
		category = apperror.Conflict
	}
	return RemoteError{operation: operation, statusCode: statusCode, category: category}
}

func (err RemoteError) Category() apperror.Category { return err.category }

func (err RemoteError) Error() string {
	if err.statusCode != 0 {
		return "remote operation failed: " + err.operation + " (status " + formatStatus(err.statusCode) + ")"
	}
	return "remote operation failed: " + err.operation
}

func (err RemoteError) Operation() string {
	return err.operation
}

func (err RemoteError) StatusCode() int {
	return err.statusCode
}

func formatStatus(statusCode int) string {
	if statusCode < 0 {
		return "unknown"
	}
	return strconv.Itoa(statusCode)
}
