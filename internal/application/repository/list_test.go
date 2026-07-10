package repository

import (
	"context"
	"testing"
)

type listService struct {
	request ListRequest
}

func (service *listService) List(_ context.Context, request ListRequest) ([]Repository, error) {
	service.request = request
	return []Repository{{Owner: "alice", Name: "project"}}, nil
}

func TestListUseCaseValidatesPaginationAndDelegates(t *testing.T) {
	service := &listService{}
	result, err := NewListUseCase(service).Execute(context.Background(), ListRequest{Page: 2, Limit: 4})
	if err != nil {
		t.Fatal(err)
	}
	if service.request != (ListRequest{Page: 2, Limit: 4}) {
		t.Fatalf("request = %+v", service.request)
	}
	if len(result) != 1 || result[0].Owner != "alice" || result[0].Name != "project" {
		t.Fatalf("result = %+v", result)
	}
}

func TestListUseCaseRejectsNonPositivePagination(t *testing.T) {
	for _, request := range []ListRequest{{Page: 0, Limit: 1}, {Page: 1, Limit: 0}} {
		if _, err := NewListUseCase(&listService{}).Execute(context.Background(), request); err == nil {
			t.Fatalf("request %+v was accepted", request)
		}
	}
}
