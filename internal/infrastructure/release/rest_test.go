package release

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationrelease "github.com/l4l4dev/fj/internal/application/release"
)

type stubTransport struct {
	path  string
	query url.Values
	body  string
}

func (s *stubTransport) Do(_ context.Context, _ string, path string, query url.Values) (*http.Response, error) {
	s.path, s.query = path, query
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(s.body))}, nil
}

type statusError int

func (err statusError) Error() string   { return "remote failure" }
func (err statusError) StatusCode() int { return int(err) }

type errorTransport struct{ err error }

func (transport errorTransport) Do(context.Context, string, string, url.Values) (*http.Response, error) {
	return nil, transport.err
}

func (transport errorTransport) DoJSON(context.Context, string, string, url.Values, []byte) (*http.Response, error) {
	return nil, transport.err
}

type jsonStubTransport struct {
	method   string
	path     string
	body     []byte
	err      error
	response string
}

func (stub *jsonStubTransport) Do(context.Context, string, string, url.Values) (*http.Response, error) {
	return nil, errors.New("unexpected Do call")
}

func (stub *jsonStubTransport) DoJSON(_ context.Context, method, path string, _ url.Values, body []byte) (*http.Response, error) {
	stub.method, stub.path, stub.body = method, path, body
	if stub.err != nil {
		return nil, stub.err
	}
	response := stub.response
	if response == "" {
		response = `{"id":1,"tag_name":"v1.0.0","name":"First release","draft":true,"prerelease":false}`
	}
	return &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(strings.NewReader(response))}, nil
}

func TestRESTAdapterList(t *testing.T) {
	transport := &stubTransport{body: `[{"id":1,"tag_name":"v1.0.0","name":"First release","draft":false,"prerelease":false},{"id":2,"tag_name":"v1.1.0-rc1","name":"Release candidate","draft":false,"prerelease":true},{"id":3,"tag_name":"v0.1.0","name":"Draft","draft":true,"prerelease":false}]`}
	result, err := NewRESTAdapter(transport).List(context.Background(), applicationrelease.ListRequest{Owner: "alice", Name: "project", Page: 2, Limit: 20})
	if err != nil || len(result) != 3 {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if result[0].TagName != "v1.0.0" || result[0].Title != "First release" || result[0].Draft || result[0].Prerelease {
		t.Fatalf("unexpected release 0: %+v", result[0])
	}
	if result[1].Prerelease != true {
		t.Fatalf("unexpected release 1: %+v", result[1])
	}
	if result[2].Draft != true {
		t.Fatalf("unexpected release 2: %+v", result[2])
	}
	if transport.path != "/api/v1/repos/alice/project/releases" || transport.query.Get("page") != "2" || transport.query.Get("limit") != "20" {
		t.Fatalf("unexpected request: path=%s query=%v", transport.path, transport.query)
	}
}

func TestRESTAdapterListEmpty(t *testing.T) {
	result, err := NewRESTAdapter(&stubTransport{body: `[]`}).List(context.Background(), applicationrelease.ListRequest{Owner: "alice", Name: "project", Page: 1, Limit: 20})
	if err != nil || result == nil || len(result) != 0 {
		t.Fatalf("unexpected empty result: %#v err=%v", result, err)
	}
}

func TestRESTAdapterListMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{http.StatusUnauthorized, apperror.Authentication},
		{http.StatusForbidden, apperror.Authentication},
		{http.StatusNotFound, apperror.NotFound},
		{http.StatusInternalServerError, apperror.Remote},
	}
	for _, test := range tests {
		_, err := NewRESTAdapter(errorTransport{err: statusError(test.status)}).List(context.Background(), applicationrelease.ListRequest{Owner: "alice", Name: "project", Page: 1, Limit: 20})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %#v", test.status, err)
		}
	}
}

func TestRESTAdapterListMapsNotFoundSafely(t *testing.T) {
	_, err := NewRESTAdapter(errorTransport{err: statusError(http.StatusNotFound)}).List(context.Background(), applicationrelease.ListRequest{Owner: "example-owner", Name: "example-repository", Page: 1, Limit: 20})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "repository not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
	}
}

func TestRESTAdapterInspect(t *testing.T) {
	transport := &stubTransport{body: `{"id":1,"tag_name":"v1.0.0","name":"First release","body":"Release notes","draft":false,"prerelease":false,"assets":[{"id":11,"name":"fj_darwin_arm64.tar.gz","size":1024}]}`}
	result, err := NewRESTAdapter(transport).Inspect(context.Background(), applicationrelease.InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
	if err != nil || result.TagName != "v1.0.0" || result.Title != "First release" || result.Notes != "Release notes" || result.Draft || result.Prerelease {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if len(result.Assets) != 1 || result.Assets[0].ID != 11 || result.Assets[0].Name != "fj_darwin_arm64.tar.gz" || result.Assets[0].Size != 1024 {
		t.Fatalf("unexpected assets: %+v", result.Assets)
	}
	if transport.path != "/api/v1/repos/alice/project/releases/tags/v1.0.0" {
		t.Fatalf("unexpected path: %s", transport.path)
	}
}

func TestRESTAdapterInspectEscapesTag(t *testing.T) {
	transport := &stubTransport{body: `{}`}
	_, err := NewRESTAdapter(transport).Inspect(context.Background(), applicationrelease.InspectRequest{Owner: "alice", Name: "project", Tag: "v1/rc1"})
	if err != nil {
		t.Fatal(err)
	}
	if transport.path != "/api/v1/repos/alice/project/releases/tags/v1%2Frc1" {
		t.Fatalf("unexpected path: %s", transport.path)
	}
}

func TestRESTAdapterInspectMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{http.StatusUnauthorized, apperror.Authentication},
		{http.StatusForbidden, apperror.Authentication},
		{http.StatusNotFound, apperror.NotFound},
		{http.StatusInternalServerError, apperror.Remote},
	}
	for _, test := range tests {
		_, err := NewRESTAdapter(errorTransport{err: statusError(test.status)}).Inspect(context.Background(), applicationrelease.InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %#v", test.status, err)
		}
	}
}

func TestRESTAdapterInspectMapsNotFoundSafely(t *testing.T) {
	_, err := NewRESTAdapter(errorTransport{err: statusError(http.StatusNotFound)}).Inspect(context.Background(), applicationrelease.InspectRequest{Owner: "alice", Name: "project", Tag: "v1.0.0"})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "release not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
	}
}

func TestRESTAdapterCreateSendsDraftPayload(t *testing.T) {
	transport := &jsonStubTransport{}
	request := applicationrelease.CreateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: "First release", Notes: "Notes", Prerelease: true}
	result, err := NewRESTAdapter(transport).Create(context.Background(), request)
	if err != nil || result.TagName != "v1.0.0" || result.Title != "First release" {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if transport.method != http.MethodPost || transport.path != "/api/v1/repos/alice/project/releases" {
		t.Fatalf("unexpected request: method=%s path=%s", transport.method, transport.path)
	}
	if string(transport.body) != `{"tag_name":"v1.0.0","name":"First release","body":"Notes","draft":true,"prerelease":true}` {
		t.Fatalf("unexpected body: %s", transport.body)
	}
}

func TestRESTAdapterCreateMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{http.StatusUnauthorized, apperror.Authentication},
		{http.StatusForbidden, apperror.Authentication},
		{http.StatusNotFound, apperror.NotFound},
		{http.StatusConflict, apperror.Conflict},
		{http.StatusUnprocessableEntity, apperror.Validation},
		{http.StatusInternalServerError, apperror.Remote},
	}
	for _, test := range tests {
		_, err := NewRESTAdapter(&jsonStubTransport{err: statusError(test.status)}).Create(context.Background(), applicationrelease.CreateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: "First release"})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %#v", test.status, err)
		}
	}
}

func TestRESTAdapterCreateMapsNotFoundSafely(t *testing.T) {
	_, err := NewRESTAdapter(&jsonStubTransport{err: statusError(http.StatusNotFound)}).Create(context.Background(), applicationrelease.CreateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: "First release"})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "repository not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
	}
}

type updateStubTransport struct {
	getPath       string
	getBody       string
	getErr        error
	patchMethod   string
	patchPath     string
	patchBody     []byte
	patchErr      error
	patchResponse string
}

func (stub *updateStubTransport) Do(_ context.Context, _ string, path string, _ url.Values) (*http.Response, error) {
	stub.getPath = path
	if stub.getErr != nil {
		return nil, stub.getErr
	}
	body := stub.getBody
	if body == "" {
		body = `{"id":7,"tag_name":"v1.0.0","name":"First release","draft":true,"prerelease":false}`
	}
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}, nil
}

func (stub *updateStubTransport) DoJSON(_ context.Context, method, path string, _ url.Values, body []byte) (*http.Response, error) {
	stub.patchMethod, stub.patchPath, stub.patchBody = method, path, body
	if stub.patchErr != nil {
		return nil, stub.patchErr
	}
	response := stub.patchResponse
	if response == "" {
		response = `{"id":7,"tag_name":"v1.0.0","name":"Updated title","draft":true,"prerelease":true}`
	}
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(response))}, nil
}

func TestRESTAdapterUpdateSendsOnlySuppliedFields(t *testing.T) {
	transport := &updateStubTransport{}
	title := "Updated title"
	prerelease := true
	result, err := NewRESTAdapter(transport).Update(context.Background(), applicationrelease.UpdateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: &title, Prerelease: &prerelease})
	if err != nil || result.Title != "Updated title" || !result.Prerelease {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if transport.getPath != "/api/v1/repos/alice/project/releases/tags/v1.0.0" {
		t.Fatalf("unexpected resolve path: %s", transport.getPath)
	}
	if transport.patchMethod != http.MethodPatch || transport.patchPath != "/api/v1/repos/alice/project/releases/7" {
		t.Fatalf("unexpected patch request: method=%s path=%s", transport.patchMethod, transport.patchPath)
	}
	if string(transport.patchBody) != `{"name":"Updated title","prerelease":true}` {
		t.Fatalf("unexpected payload: %s", transport.patchBody)
	}
}

func TestRESTAdapterUpdateMapsResolveHTTPError(t *testing.T) {
	title := "Updated title"
	transport := &updateStubTransport{getErr: statusError(http.StatusNotFound)}
	_, err := NewRESTAdapter(transport).Update(context.Background(), applicationrelease.UpdateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: &title})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "release not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
	}
}

func TestRESTAdapterUpdateMapsPatchHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{http.StatusUnauthorized, apperror.Authentication},
		{http.StatusForbidden, apperror.Authentication},
		{http.StatusConflict, apperror.Validation},
		{http.StatusUnprocessableEntity, apperror.Validation},
		{http.StatusInternalServerError, apperror.Remote},
	}
	title := "Updated title"
	for _, test := range tests {
		transport := &updateStubTransport{patchErr: statusError(test.status)}
		_, err := NewRESTAdapter(transport).Update(context.Background(), applicationrelease.UpdateRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", Title: &title})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %#v", test.status, err)
		}
	}
}

func TestRESTAdapterPublishSendsDraftFalsePayload(t *testing.T) {
	transport := &jsonStubTransport{response: `{"id":7,"tag_name":"v1.0.0","name":"First release","draft":false,"prerelease":false}`}
	result, err := NewRESTAdapter(transport).Publish(context.Background(), applicationrelease.PublishRequest{Owner: "alice", Name: "project", ID: 7})
	if err != nil || result.TagName != "v1.0.0" || result.Draft {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if transport.method != http.MethodPatch || transport.path != "/api/v1/repos/alice/project/releases/7" {
		t.Fatalf("unexpected request: method=%s path=%s", transport.method, transport.path)
	}
	if string(transport.body) != `{"draft":false}` {
		t.Fatalf("unexpected body: %s", transport.body)
	}
}

func TestRESTAdapterPublishMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{http.StatusUnauthorized, apperror.Authentication},
		{http.StatusForbidden, apperror.Authentication},
		{http.StatusNotFound, apperror.NotFound},
		{http.StatusConflict, apperror.Conflict},
		{http.StatusUnprocessableEntity, apperror.Conflict},
		{http.StatusInternalServerError, apperror.Remote},
	}
	for _, test := range tests {
		_, err := NewRESTAdapter(&jsonStubTransport{err: statusError(test.status)}).Publish(context.Background(), applicationrelease.PublishRequest{Owner: "alice", Name: "project", ID: 7})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %#v", test.status, err)
		}
	}
}

func TestRESTAdapterPublishMapsNotFoundSafely(t *testing.T) {
	_, err := NewRESTAdapter(&jsonStubTransport{err: statusError(http.StatusNotFound)}).Publish(context.Background(), applicationrelease.PublishRequest{Owner: "alice", Name: "project", ID: 7})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "release not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
	}
}

type assetStubTransport struct {
	getPath     string
	getBody     string
	getErr      error
	doMethod    string
	doPath      string
	doCalls     int
	doErr       error
	rawMethod   string
	rawPath     string
	rawQuery    url.Values
	rawBody     []byte
	rawType     string
	rawErr      error
	rawResponse string
}

func (stub *assetStubTransport) Do(_ context.Context, method, path string, _ url.Values) (*http.Response, error) {
	if method == http.MethodGet {
		stub.getPath = path
		if stub.getErr != nil {
			return nil, stub.getErr
		}
		body := stub.getBody
		if body == "" {
			body = `{"id":7,"tag_name":"v1.0.0","assets":[{"id":11,"name":"fj.tar.gz","size":7}]}`
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}, nil
	}
	stub.doMethod, stub.doPath, stub.doCalls = method, path, stub.doCalls+1
	if stub.doErr != nil {
		return nil, stub.doErr
	}
	return &http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(strings.NewReader(""))}, nil
}

func (stub *assetStubTransport) DoRaw(_ context.Context, method, path string, query url.Values, body []byte, contentType string) (*http.Response, error) {
	stub.rawMethod, stub.rawPath, stub.rawQuery, stub.rawBody, stub.rawType = method, path, query, body, contentType
	if stub.rawErr != nil {
		return nil, stub.rawErr
	}
	response := stub.rawResponse
	if response == "" {
		response = `{"id":11,"name":"fj.tar.gz","size":7}`
	}
	return &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(strings.NewReader(response))}, nil
}

func TestRESTAdapterUploadAssetPostsMultipartToResolvedRelease(t *testing.T) {
	transport := &assetStubTransport{}
	result, err := NewRESTAdapter(transport).UploadAsset(context.Background(), applicationrelease.UploadAssetRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz", Content: []byte("payload")})
	if err != nil || result.ID != 11 || result.Name != "fj.tar.gz" || result.Size != 7 {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	if transport.getPath != "/api/v1/repos/alice/project/releases/tags/v1.0.0" {
		t.Fatalf("unexpected resolve path: %s", transport.getPath)
	}
	if transport.rawMethod != http.MethodPost || transport.rawPath != "/api/v1/repos/alice/project/releases/7/assets" {
		t.Fatalf("unexpected upload request: method=%s path=%s", transport.rawMethod, transport.rawPath)
	}
	if transport.rawQuery.Get("name") != "fj.tar.gz" {
		t.Fatalf("unexpected query: %v", transport.rawQuery)
	}
	if !strings.HasPrefix(transport.rawType, "multipart/form-data; boundary=") {
		t.Fatalf("unexpected content type: %s", transport.rawType)
	}
	boundary := strings.TrimPrefix(transport.rawType, "multipart/form-data; boundary=")
	reader := multipart.NewReader(bytes.NewReader(transport.rawBody), boundary)
	part, err := reader.NextPart()
	if err != nil {
		t.Fatal(err)
	}
	if part.FormName() != "attachment" || part.FileName() != "fj.tar.gz" {
		t.Fatalf("unexpected part: name=%s filename=%s", part.FormName(), part.FileName())
	}
	content, err := io.ReadAll(part)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "payload" {
		t.Fatalf("unexpected part content: %q", content)
	}
}

func TestRESTAdapterUploadAssetMapsResolveHTTPError(t *testing.T) {
	transport := &assetStubTransport{getErr: statusError(http.StatusNotFound)}
	_, err := NewRESTAdapter(transport).UploadAsset(context.Background(), applicationrelease.UploadAssetRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz", Content: []byte("payload")})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "release not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
	}
	if transport.rawMethod != "" {
		t.Fatal("upload must not be attempted when the release cannot be resolved")
	}
}

func TestRESTAdapterUploadAssetMapsUploadHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
		message  string
	}{
		{http.StatusUnauthorized, apperror.Authentication, "permission denied or authentication failed"},
		{http.StatusForbidden, apperror.Authentication, "permission denied or authentication failed"},
		{http.StatusConflict, apperror.Conflict, "an asset with this name already exists"},
		{http.StatusRequestEntityTooLarge, apperror.Validation, "asset was rejected by the remote service"},
		{http.StatusUnprocessableEntity, apperror.Validation, "asset was rejected by the remote service"},
		{http.StatusInternalServerError, apperror.Remote, ""},
	}
	for _, test := range tests {
		transport := &assetStubTransport{rawErr: statusError(test.status)}
		_, err := NewRESTAdapter(transport).UploadAsset(context.Background(), applicationrelease.UploadAssetRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz", Content: []byte("payload")})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category || appErr.Message != test.message {
			t.Fatalf("status %d: unexpected error %#v", test.status, err)
		}
	}
}

func TestRESTAdapterDeleteAssetUsesResolvedAssetID(t *testing.T) {
	transport := &assetStubTransport{}
	if err := NewRESTAdapter(transport).DeleteAsset(context.Background(), applicationrelease.DeleteAssetRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz"}); err != nil {
		t.Fatal(err)
	}
	if transport.getPath != "/api/v1/repos/alice/project/releases/tags/v1.0.0" {
		t.Fatalf("unexpected resolve path: %s", transport.getPath)
	}
	if transport.doMethod != http.MethodDelete || transport.doPath != "/api/v1/repos/alice/project/releases/7/assets/11" {
		t.Fatalf("unexpected delete request: method=%s path=%s", transport.doMethod, transport.doPath)
	}
}

func TestRESTAdapterDeleteAssetReportsMissingAsset(t *testing.T) {
	transport := &assetStubTransport{getBody: `{"id":7,"tag_name":"v1.0.0","assets":[{"id":11,"name":"other.tar.gz","size":7}]}`}
	err := NewRESTAdapter(transport).DeleteAsset(context.Background(), applicationrelease.DeleteAssetRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz"})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "asset not found" {
		t.Fatalf("error = %#v, want asset not-found application error", err)
	}
	if transport.doCalls != 0 {
		t.Fatal("delete must not be attempted when the asset is absent")
	}
}

func TestRESTAdapterDeleteAssetMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{http.StatusUnauthorized, apperror.Authentication},
		{http.StatusForbidden, apperror.Authentication},
		{http.StatusNotFound, apperror.NotFound},
		{http.StatusInternalServerError, apperror.Remote},
	}
	for _, test := range tests {
		transport := &assetStubTransport{doErr: statusError(test.status)}
		err := NewRESTAdapter(transport).DeleteAsset(context.Background(), applicationrelease.DeleteAssetRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz"})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %#v", test.status, err)
		}
	}
}

func TestRESTAdapterDeleteAssetMapsResolveHTTPError(t *testing.T) {
	transport := &assetStubTransport{getErr: statusError(http.StatusNotFound)}
	err := NewRESTAdapter(transport).DeleteAsset(context.Background(), applicationrelease.DeleteAssetRequest{Owner: "alice", Name: "project", Tag: "v1.0.0", AssetName: "fj.tar.gz"})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "release not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
	}
}

func TestRESTAdapterDeleteUsesResolvedID(t *testing.T) {
	transport := &assetStubTransport{}
	if err := NewRESTAdapter(transport).Delete(context.Background(), applicationrelease.DeleteRequest{Owner: "alice", Name: "project", ID: 7}); err != nil {
		t.Fatal(err)
	}
	if transport.doMethod != http.MethodDelete || transport.doPath != "/api/v1/repos/alice/project/releases/7" {
		t.Fatalf("unexpected delete request: method=%s path=%s", transport.doMethod, transport.doPath)
	}
}

func TestRESTAdapterDeleteMapsHTTPError(t *testing.T) {
	tests := []struct {
		status   int
		category apperror.Category
	}{
		{http.StatusUnauthorized, apperror.Authentication},
		{http.StatusForbidden, apperror.Authentication},
		{http.StatusNotFound, apperror.NotFound},
		{http.StatusMethodNotAllowed, apperror.Conflict},
		{http.StatusConflict, apperror.Conflict},
		{http.StatusUnprocessableEntity, apperror.Conflict},
		{http.StatusInternalServerError, apperror.Remote},
	}
	for _, test := range tests {
		transport := &assetStubTransport{doErr: statusError(test.status)}
		err := NewRESTAdapter(transport).Delete(context.Background(), applicationrelease.DeleteRequest{Owner: "alice", Name: "project", ID: 7})
		var appErr apperror.Error
		if !errors.As(err, &appErr) || appErr.Category != test.category {
			t.Fatalf("status %d: unexpected error %#v", test.status, err)
		}
	}
}

func TestRESTAdapterDeleteMapsNotFoundSafely(t *testing.T) {
	transport := &assetStubTransport{doErr: statusError(http.StatusNotFound)}
	err := NewRESTAdapter(transport).Delete(context.Background(), applicationrelease.DeleteRequest{Owner: "alice", Name: "project", ID: 7})
	var appErr apperror.Error
	if !errors.As(err, &appErr) || appErr.Category != apperror.NotFound || appErr.Message != "release not found" {
		t.Fatalf("error = %#v, want safe not-found application error", err)
	}
}
