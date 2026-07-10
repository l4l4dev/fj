package repository

import (
	"context"
	"strconv"
)

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

type RemoteError struct {
	operation  string
	statusCode int
}

func NewRemoteError(operation string, statusCode int) RemoteError {
	return RemoteError{operation: operation, statusCode: statusCode}
}

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
