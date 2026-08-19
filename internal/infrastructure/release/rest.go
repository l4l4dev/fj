package release

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationrelease "github.com/l4l4dev/fj/internal/application/release"
)

type transport interface {
	Do(context.Context, string, string, url.Values) (*http.Response, error)
}

type jsonTransport interface {
	DoJSON(context.Context, string, string, url.Values, []byte) (*http.Response, error)
}

type rawTransport interface {
	DoRaw(context.Context, string, string, url.Values, []byte, string) (*http.Response, error)
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

func (a *RESTAdapter) Create(ctx context.Context, request applicationrelease.CreateRequest) (applicationrelease.ReleaseDetail, error) {
	jsonClient, ok := a.transport.(jsonTransport)
	if !ok {
		return applicationrelease.ReleaseDetail{}, apperror.New(apperror.Remote, "create release", "")
	}
	body, err := json.Marshal(struct {
		TagName    string `json:"tag_name"`
		Name       string `json:"name"`
		Body       string `json:"body"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
	}{TagName: request.Tag, Name: request.Title, Body: request.Notes, Draft: true, Prerelease: request.Prerelease})
	if err != nil {
		return applicationrelease.ReleaseDetail{}, apperror.New(apperror.Remote, "create release", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/releases"
	response, err := jsonClient.DoJSON(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return applicationrelease.ReleaseDetail{}, translateCreateError(err)
	}
	defer response.Body.Close()
	var decoded forgejoReleaseDetail
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationrelease.ReleaseDetail{}, apperror.New(apperror.Remote, "create release", "")
	}
	assets := make([]applicationrelease.Asset, 0, len(decoded.Assets))
	for _, asset := range decoded.Assets {
		assets = append(assets, applicationrelease.Asset{ID: asset.ID, Name: asset.Name, Size: asset.Size})
	}
	return applicationrelease.ReleaseDetail{ID: decoded.ID, TagName: decoded.TagName, Title: decoded.Name, Draft: decoded.Draft, Prerelease: decoded.Prerelease, Notes: decoded.Body, Assets: assets}, nil
}

func (a *RESTAdapter) Update(ctx context.Context, request applicationrelease.UpdateRequest) (applicationrelease.ReleaseDetail, error) {
	jsonClient, ok := a.transport.(jsonTransport)
	if !ok {
		return applicationrelease.ReleaseDetail{}, apperror.New(apperror.Remote, "update release", "")
	}
	tagPath := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/releases/tags/" + url.PathEscape(request.Tag)
	tagResponse, err := a.transport.Do(ctx, http.MethodGet, tagPath, nil)
	if err != nil {
		return applicationrelease.ReleaseDetail{}, translateUpdateError(err)
	}
	var existing forgejoReleaseDetail
	decodeErr := json.NewDecoder(tagResponse.Body).Decode(&existing)
	tagResponse.Body.Close()
	if decodeErr != nil {
		return applicationrelease.ReleaseDetail{}, apperror.New(apperror.Remote, "update release", "")
	}

	payload := make(map[string]any)
	if request.Title != nil {
		payload["name"] = *request.Title
	}
	if request.Notes != nil {
		payload["body"] = *request.Notes
	}
	if request.Prerelease != nil {
		payload["prerelease"] = *request.Prerelease
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return applicationrelease.ReleaseDetail{}, apperror.New(apperror.Remote, "update release", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/releases/" + strconv.FormatInt(existing.ID, 10)
	response, err := jsonClient.DoJSON(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return applicationrelease.ReleaseDetail{}, translateUpdateError(err)
	}
	defer response.Body.Close()
	var decoded forgejoReleaseDetail
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationrelease.ReleaseDetail{}, apperror.New(apperror.Remote, "update release", "")
	}
	assets := make([]applicationrelease.Asset, 0, len(decoded.Assets))
	for _, asset := range decoded.Assets {
		assets = append(assets, applicationrelease.Asset{ID: asset.ID, Name: asset.Name, Size: asset.Size})
	}
	return applicationrelease.ReleaseDetail{ID: decoded.ID, TagName: decoded.TagName, Title: decoded.Name, Draft: decoded.Draft, Prerelease: decoded.Prerelease, Notes: decoded.Body, Assets: assets}, nil
}

func (a *RESTAdapter) Publish(ctx context.Context, request applicationrelease.PublishRequest) (applicationrelease.ReleaseDetail, error) {
	jsonClient, ok := a.transport.(jsonTransport)
	if !ok {
		return applicationrelease.ReleaseDetail{}, apperror.New(apperror.Remote, "publish release", "")
	}
	body, err := json.Marshal(struct {
		Draft bool `json:"draft"`
	}{Draft: false})
	if err != nil {
		return applicationrelease.ReleaseDetail{}, apperror.New(apperror.Remote, "publish release", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/releases/" + strconv.FormatInt(request.ID, 10)
	response, err := jsonClient.DoJSON(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return applicationrelease.ReleaseDetail{}, translatePublishError(err)
	}
	defer response.Body.Close()
	var decoded forgejoReleaseDetail
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationrelease.ReleaseDetail{}, apperror.New(apperror.Remote, "publish release", "")
	}
	assets := make([]applicationrelease.Asset, 0, len(decoded.Assets))
	for _, asset := range decoded.Assets {
		assets = append(assets, applicationrelease.Asset{ID: asset.ID, Name: asset.Name, Size: asset.Size})
	}
	return applicationrelease.ReleaseDetail{ID: decoded.ID, TagName: decoded.TagName, Title: decoded.Name, Draft: decoded.Draft, Prerelease: decoded.Prerelease, Notes: decoded.Body, Assets: assets}, nil
}

func (a *RESTAdapter) Delete(ctx context.Context, request applicationrelease.DeleteRequest) error {
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/releases/" + strconv.FormatInt(request.ID, 10)
	response, err := a.transport.Do(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return translateDeleteError(err)
	}
	response.Body.Close()
	return nil
}

func (a *RESTAdapter) UploadAsset(ctx context.Context, request applicationrelease.UploadAssetRequest) (applicationrelease.Asset, error) {
	rawClient, ok := a.transport.(rawTransport)
	if !ok {
		return applicationrelease.Asset{}, apperror.New(apperror.Remote, "upload release asset", "")
	}
	existing, err := a.resolveReleaseByTag(ctx, request.Owner, request.Name, request.Tag, "upload release asset")
	if err != nil {
		return applicationrelease.Asset{}, err
	}
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	part, err := writer.CreateFormFile("attachment", request.AssetName)
	if err != nil {
		return applicationrelease.Asset{}, apperror.New(apperror.Remote, "upload release asset", "")
	}
	if _, err := part.Write(request.Content); err != nil {
		return applicationrelease.Asset{}, apperror.New(apperror.Remote, "upload release asset", "")
	}
	if err := writer.Close(); err != nil {
		return applicationrelease.Asset{}, apperror.New(apperror.Remote, "upload release asset", "")
	}
	query := url.Values{}
	query.Set("name", request.AssetName)
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/releases/" + strconv.FormatInt(existing.ID, 10) + "/assets"
	response, err := rawClient.DoRaw(ctx, http.MethodPost, path, query, buffer.Bytes(), writer.FormDataContentType())
	if err != nil {
		return applicationrelease.Asset{}, translateAssetError(err, "upload release asset")
	}
	defer response.Body.Close()
	var decoded struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Size int64  `json:"size"`
	}
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationrelease.Asset{}, apperror.New(apperror.Remote, "upload release asset", "")
	}
	return applicationrelease.Asset{ID: decoded.ID, Name: decoded.Name, Size: decoded.Size}, nil
}

func (a *RESTAdapter) DeleteAsset(ctx context.Context, request applicationrelease.DeleteAssetRequest) error {
	existing, err := a.resolveReleaseByTag(ctx, request.Owner, request.Name, request.Tag, "delete release asset")
	if err != nil {
		return err
	}
	assetID := int64(0)
	found := false
	for _, asset := range existing.Assets {
		if asset.Name == request.AssetName {
			assetID, found = asset.ID, true
			break
		}
	}
	if !found {
		return apperror.New(apperror.NotFound, "delete release asset", "asset not found")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/releases/" + strconv.FormatInt(existing.ID, 10) + "/assets/" + strconv.FormatInt(assetID, 10)
	response, err := a.transport.Do(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return translateAssetError(err, "delete release asset")
	}
	response.Body.Close()
	return nil
}

func (a *RESTAdapter) resolveReleaseByTag(ctx context.Context, owner, name, tag, operation string) (forgejoReleaseDetail, error) {
	path := "/api/v1/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(name) + "/releases/tags/" + url.PathEscape(tag)
	response, err := a.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return forgejoReleaseDetail{}, translateAssetError(err, operation)
	}
	var decoded forgejoReleaseDetail
	decodeErr := json.NewDecoder(response.Body).Decode(&decoded)
	response.Body.Close()
	if decodeErr != nil {
		return forgejoReleaseDetail{}, apperror.New(apperror.Remote, operation, "")
	}
	return decoded, nil
}

func translateDeleteError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, "delete release", "permission denied or authentication failed")
		case 404:
			return apperror.New(apperror.NotFound, "delete release", "release not found")
		case 405, 409, 422:
			return apperror.New(apperror.Conflict, "delete release", "release could not be deleted in its current state")
		}
	}
	return apperror.New(apperror.Remote, "delete release", "")
}

func translateAssetError(err error, operation string) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, operation, "permission denied or authentication failed")
		case 404:
			return apperror.New(apperror.NotFound, operation, "release not found")
		case 409:
			return apperror.New(apperror.Conflict, operation, "an asset with this name already exists")
		case 413, 422:
			return apperror.New(apperror.Validation, operation, "asset was rejected by the remote service")
		}
	}
	return apperror.New(apperror.Remote, operation, "")
}

func translatePublishError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, "publish release", "permission denied or authentication failed")
		case 404:
			return apperror.New(apperror.NotFound, "publish release", "release not found")
		case 409, 422:
			return apperror.New(apperror.Conflict, "publish release", "release could not be published in its current state")
		}
	}
	return apperror.New(apperror.Remote, "publish release", "")
}

func translateUpdateError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, "update release", "permission denied or authentication failed")
		case 404:
			return apperror.New(apperror.NotFound, "update release", "release not found")
		case 409, 422:
			return apperror.New(apperror.Validation, "update release", "release fields were rejected by the remote service")
		}
	}
	return apperror.New(apperror.Remote, "update release", "")
}

func translateCreateError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, "create release", "permission denied or authentication failed")
		case 404:
			return apperror.New(apperror.NotFound, "create release", "repository not found")
		case 409:
			return apperror.New(apperror.Conflict, "create release", "a release for this tag already exists")
		case 422:
			return apperror.New(apperror.Validation, "create release", "release tag or metadata was rejected by the remote service")
		}
	}
	return apperror.New(apperror.Remote, "create release", "")
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
var _ applicationrelease.Creator = (*RESTAdapter)(nil)
var _ applicationrelease.Updater = (*RESTAdapter)(nil)
var _ applicationrelease.Publisher = (*RESTAdapter)(nil)
var _ applicationrelease.AssetUploader = (*RESTAdapter)(nil)
var _ applicationrelease.AssetDeleter = (*RESTAdapter)(nil)
var _ applicationrelease.Deleter = (*RESTAdapter)(nil)
