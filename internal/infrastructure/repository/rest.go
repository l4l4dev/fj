package repository

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	applicationrepository "github.com/l4l4dev/fj/internal/application/repository"
)

type transport interface {
	Do(context.Context, string, string, url.Values) (*http.Response, error)
}

type RESTAdapter struct {
	transport transport
}

type forgejoRepository struct {
	Name  string       `json:"name"`
	Owner forgejoOwner `json:"owner"`
}

type forgejoOwner struct {
	Login string `json:"login"`
}

func NewRESTAdapter(transport transport) *RESTAdapter {
	return &RESTAdapter{transport: transport}
}

func (adapter *RESTAdapter) List(ctx context.Context, request applicationrepository.ListRequest) ([]applicationrepository.Repository, error) {
	query := url.Values{}
	query.Set("page", strconv.Itoa(request.Page))
	query.Set("limit", strconv.Itoa(request.Limit))

	response, err := adapter.transport.Do(ctx, http.MethodGet, "/api/v1/user/repos", query)
	if err != nil {
		return nil, translateRemoteError(err)
	}
	defer response.Body.Close()

	var decoded []forgejoRepository
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, applicationrepository.NewRemoteError("decode repositories", 0)
	}

	result := make([]applicationrepository.Repository, 0, len(decoded))
	for _, repository := range decoded {
		result = append(result, applicationrepository.Repository{
			Owner: repository.Owner.Login,
			Name:  repository.Name,
		})
	}
	return result, nil
}

func translateRemoteError(err error) error {
	var safeError interface {
		Operation() string
		StatusCode() int
	}
	if errors.As(err, &safeError) {
		return applicationrepository.NewRemoteError(safeError.Operation(), safeError.StatusCode())
	}
	return applicationrepository.NewRemoteError("list repositories", 0)
}

var _ applicationrepository.Service = (*RESTAdapter)(nil)
