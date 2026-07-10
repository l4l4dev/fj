package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type serviceFunc func(context.Context, ListRequest) ([]Repository, error)

func (service serviceFunc) List(ctx context.Context, request ListRequest) ([]Repository, error) {
	return service(ctx, request)
}

func TestServiceContractPassesContextAndPaginationRequest(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{}{}, "value")
	wantRequest := ListRequest{Page: 2, Limit: 25}
	service := serviceFunc(func(received context.Context, request ListRequest) ([]Repository, error) {
		if received != ctx {
			t.Errorf("context was not passed through")
		}
		if request != wantRequest {
			t.Errorf("request = %#v, want %#v", request, wantRequest)
		}
		return []Repository{{Owner: "octo", Name: "project"}}, nil
	})

	result, err := service.List(ctx, wantRequest)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0] != (Repository{Owner: "octo", Name: "project"}) {
		t.Errorf("result = %#v", result)
	}
}

func TestRemoteErrorIsSafeAndClassifiable(t *testing.T) {
	const secret = "secret-token"
	remoteError := NewRemoteError("list repositories", 503)
	var classified RemoteError
	if !errors.As(remoteError, &classified) {
		t.Fatal("RemoteError cannot be classified with errors.As")
	}
	if classified.Operation() != "list repositories" || classified.StatusCode() != 503 {
		t.Errorf("classification = %q/%d", classified.Operation(), classified.StatusCode())
	}
	for _, diagnostic := range []string{remoteError.Error(), fmt.Sprint(remoteError), fmt.Sprintf("%#v", remoteError)} {
		if strings.Contains(diagnostic, secret) || strings.Contains(diagnostic, "https://") {
			t.Errorf("RemoteError diagnostic exposes sensitive data: %q", diagnostic)
		}
	}
}

var _ Service = serviceFunc(nil)
