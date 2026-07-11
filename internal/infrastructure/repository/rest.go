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

type jsonTransport interface {
	DoJSON(context.Context, string, string, url.Values, []byte) (*http.Response, error)
}

type RESTAdapter struct {
	transport transport
}

type forgejoRepository struct {
	Name          string       `json:"name"`
	Owner         forgejoOwner `json:"owner"`
	Description   string       `json:"description"`
	Private       bool         `json:"private"`
	Archived      bool         `json:"archived"`
	DefaultBranch string       `json:"default_branch"`
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
		return nil, translateRemoteError(err, "list repositories")
	}
	defer response.Body.Close()

	var decoded []forgejoRepository
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, applicationrepository.NewRemoteError("decode repositories", 0)
	}

	result := make([]applicationrepository.Repository, 0, len(decoded))
	for _, repository := range decoded {
		result = append(result, toApplicationRepository(repository))
	}
	return result, nil
}

func (adapter *RESTAdapter) Get(ctx context.Context, request applicationrepository.GetRequest) (applicationrepository.Repository, error) {
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name)
	response, err := adapter.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return applicationrepository.Repository{}, translateRemoteError(err, "inspect repository")
	}
	defer response.Body.Close()

	var decoded forgejoRepository
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationrepository.Repository{}, applicationrepository.NewRemoteError("inspect repository", 0)
	}
	return toApplicationRepository(decoded), nil
}

func (adapter *RESTAdapter) Create(ctx context.Context, request applicationrepository.CreateRequest) (applicationrepository.Repository, error) {
	transport, ok := adapter.transport.(jsonTransport)
	if !ok {
		return applicationrepository.Repository{}, applicationrepository.NewRemoteError("create repository", 0)
	}
	body, err := json.Marshal(struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		Private     bool   `json:"private"`
	}{Name: request.Name, Description: request.Description, Private: request.Private})
	if err != nil {
		return applicationrepository.Repository{}, applicationrepository.NewRemoteError("create repository", 0)
	}
	response, err := transport.DoJSON(ctx, http.MethodPost, "/api/v1/user/repos", nil, body)
	if err != nil {
		return applicationrepository.Repository{}, translateRemoteError(err, "create repository")
	}
	defer response.Body.Close()
	var decoded forgejoRepository
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationrepository.Repository{}, applicationrepository.NewRemoteError("create repository", 0)
	}
	return toApplicationRepository(decoded), nil
}

func (adapter *RESTAdapter) Update(ctx context.Context, request applicationrepository.UpdateRequest) (applicationrepository.Repository, error) {
	transport, ok := adapter.transport.(jsonTransport)
	if !ok {
		return applicationrepository.Repository{}, applicationrepository.NewRemoteError("update repository", 0)
	}
	body := make(map[string]interface{})
	if request.Description != nil {
		body["description"] = *request.Description
	}
	if request.Private != nil {
		body["private"] = *request.Private
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return applicationrepository.Repository{}, applicationrepository.NewRemoteError("update repository", 0)
	}
	response, err := transport.DoJSON(ctx, http.MethodPatch, "/api/v1/repos/"+url.PathEscape(request.Owner)+"/"+url.PathEscape(request.Name), nil, encoded)
	if err != nil {
		return applicationrepository.Repository{}, translateRemoteError(err, "update repository")
	}
	defer response.Body.Close()
	var decoded forgejoRepository
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationrepository.Repository{}, applicationrepository.NewRemoteError("update repository", 0)
	}
	return toApplicationRepository(decoded), nil
}

func (adapter *RESTAdapter) SetArchived(ctx context.Context, request applicationrepository.ArchiveRequest) (applicationrepository.Repository, error) {
	transport, ok := adapter.transport.(jsonTransport)
	if !ok {
		return applicationrepository.Repository{}, applicationrepository.NewRemoteError("archive repository", 0)
	}
	body, err := json.Marshal(struct {
		Archived bool `json:"archived"`
	}{Archived: request.Archived})
	if err != nil {
		return applicationrepository.Repository{}, applicationrepository.NewRemoteError("archive repository", 0)
	}
	response, err := transport.DoJSON(ctx, http.MethodPatch, "/api/v1/repos/"+url.PathEscape(request.Owner)+"/"+url.PathEscape(request.Name), nil, body)
	if err != nil {
		op := "archive repository"
		if !request.Archived {
			op = "restore repository"
		}
		return applicationrepository.Repository{}, translateRemoteError(err, op)
	}
	defer response.Body.Close()
	var decoded forgejoRepository
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		op := "archive repository"
		if !request.Archived {
			op = "restore repository"
		}
		return applicationrepository.Repository{}, applicationrepository.NewRemoteError(op, 0)
	}
	return toApplicationRepository(decoded), nil
}

func toApplicationRepository(repository forgejoRepository) applicationrepository.Repository {
	return applicationrepository.Repository{
		Owner:         repository.Owner.Login,
		Name:          repository.Name,
		Description:   repository.Description,
		Private:       repository.Private,
		Archived:      repository.Archived,
		DefaultBranch: repository.DefaultBranch,
	}
}

func translateRemoteError(err error, operation string) error {
	var safeError interface {
		Operation() string
		StatusCode() int
	}
	if errors.As(err, &safeError) {
		return applicationrepository.NewRemoteError(safeError.Operation(), safeError.StatusCode())
	}
	return applicationrepository.NewRemoteError(operation, 0)
}

var _ applicationrepository.Service = (*RESTAdapter)(nil)
var _ applicationrepository.Getter = (*RESTAdapter)(nil)
var _ applicationrepository.Creator = (*RESTAdapter)(nil)
var _ applicationrepository.Updater = (*RESTAdapter)(nil)
var _ applicationrepository.Archiver = (*RESTAdapter)(nil)
