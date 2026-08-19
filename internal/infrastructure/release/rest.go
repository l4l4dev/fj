package release

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationrelease "github.com/l4l4dev/fj/internal/application/release"
)

type transport interface {
	Do(context.Context, string, string, url.Values) (*http.Response, error)
}

type RESTAdapter struct{ transport transport }

type forgejoRelease struct {
	ID         int64  `json:"id"`
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

type forgejoReleaseDetail struct {
	ID         int64  `json:"id"`
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
	Assets     []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Size int64  `json:"size"`
	} `json:"assets"`
}

func NewRESTAdapter(t transport) *RESTAdapter { return &RESTAdapter{transport: t} }

func (a *RESTAdapter) List(ctx context.Context, request applicationrelease.ListRequest) ([]applicationrelease.Release, error) {
	query := url.Values{}
	query.Set("page", strconv.Itoa(request.Page))
	query.Set("limit", strconv.Itoa(request.Limit))
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/releases"
	response, err := a.transport.Do(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, translateListError(err)
	}
	defer response.Body.Close()
	var decoded []forgejoRelease
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, apperror.New(apperror.Remote, "list releases", "")
	}
	result := make([]applicationrelease.Release, 0, len(decoded))
	for _, item := range decoded {
		result = append(result, applicationrelease.Release{ID: item.ID, TagName: item.TagName, Title: item.Name, Draft: item.Draft, Prerelease: item.Prerelease})
	}
	return result, nil
}

func (a *RESTAdapter) Inspect(ctx context.Context, request applicationrelease.InspectRequest) (applicationrelease.ReleaseDetail, error) {
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/releases/tags/" + url.PathEscape(request.Tag)
	response, err := a.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return applicationrelease.ReleaseDetail{}, translateInspectError(err)
	}
	defer response.Body.Close()
	var decoded forgejoReleaseDetail
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationrelease.ReleaseDetail{}, apperror.New(apperror.Remote, "inspect release", "")
	}
	assets := make([]applicationrelease.Asset, 0, len(decoded.Assets))
	for _, asset := range decoded.Assets {
		assets = append(assets, applicationrelease.Asset{ID: asset.ID, Name: asset.Name, Size: asset.Size})
	}
	return applicationrelease.ReleaseDetail{ID: decoded.ID, TagName: decoded.TagName, Title: decoded.Name, Draft: decoded.Draft, Prerelease: decoded.Prerelease, Notes: decoded.Body, Assets: assets}, nil
}

func translateInspectError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, "inspect release", "permission denied or authentication failed")
		case 404:
			return apperror.New(apperror.NotFound, "inspect release", "release not found")
		}
	}
	return apperror.New(apperror.Remote, "inspect release", "")
}

func translateListError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, "list releases", "")
		case 404:
			return apperror.New(apperror.NotFound, "list releases", "repository not found")
		}
	}
	return apperror.New(apperror.Remote, "list releases", "")
}

var _ applicationrelease.Lister = (*RESTAdapter)(nil)
var _ applicationrelease.Inspector = (*RESTAdapter)(nil)
