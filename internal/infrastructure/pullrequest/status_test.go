package pullrequest

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
)

type statusTransport struct {
	pullBody    string
	reviewsBody string
	checksBody  string
	reviewsErr  error
	checksErr   error
	paths       []string
}

func (transport *statusTransport) Do(_ context.Context, _ string, path string, _ url.Values) (*http.Response, error) {
	transport.paths = append(transport.paths, path)
	switch {
	case strings.HasSuffix(path, "/reviews"):
		if transport.reviewsErr != nil {
			return nil, transport.reviewsErr
		}
		return statusResponse(transport.reviewsBody), nil
	case strings.HasSuffix(path, "/status"):
		if transport.checksErr != nil {
			return nil, transport.checksErr
		}
		return statusResponse(transport.checksBody), nil
	default:
		return statusResponse(transport.pullBody), nil
	}
}

func statusResponse(body string) *http.Response {
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}
}

func TestRESTAdapterViewStatus(t *testing.T) {
	transport := statusTransport{
		pullBody:    `{"number":12,"mergeable":true,"head":{"sha":"abc123"},"requested_reviewers":[]}`,
		reviewsBody: `[{"id":7,"state":"APPROVED","dismissed":false,"stale":false,"user":{"id":42}}]`,
		checksBody:  `{"statuses":[{"status":"success"}]}`,
	}
	result, err := NewRESTAdapter(&transport).ViewStatus(context.Background(), applicationpullrequest.StatusRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || result.Number != 12 || !result.ReviewsAvailable || len(result.Reviews) != 1 || result.Reviews[0].ID != 7 || result.Reviews[0].ReviewerID != 42 || result.Reviews[0].State != "APPROVED" || !result.ChecksAvailable || len(result.Checks) != 1 || result.Checks[0] != "success" || result.Mergeable != applicationpullrequest.MergeableYes {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
	wantPaths := []string{"/api/v1/repos/alice/project/pulls/12", "/api/v1/repos/alice/project/pulls/12/reviews", "/api/v1/repos/alice/project/commits/abc123/status"}
	if len(transport.paths) != len(wantPaths) {
		t.Fatalf("unexpected paths: %v", transport.paths)
	}
	for index := range wantPaths {
		if transport.paths[index] != wantPaths[index] {
			t.Fatalf("path[%d] = %q, want %q", index, transport.paths[index], wantPaths[index])
		}
	}
}

func TestRESTAdapterViewStatusUsesUnavailableForMissingComponents(t *testing.T) {
	transport := statusTransport{
		pullBody:   `{"number":12,"head":{}}`,
		reviewsErr: statusError(http.StatusNotFound),
	}
	result, err := NewRESTAdapter(&transport).ViewStatus(context.Background(), applicationpullrequest.StatusRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || result.ReviewsAvailable || result.ChecksAvailable || result.Mergeable != applicationpullrequest.MergeableUnavailable {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
}

func TestRESTAdapterViewStatusUsesUnavailableForMissingCheckAPI(t *testing.T) {
	transport := statusTransport{
		pullBody:    `{"number":12,"mergeable":false,"head":{"sha":"abc123"}}`,
		reviewsBody: `[]`,
		checksErr:   statusError(http.StatusNotFound),
	}
	result, err := NewRESTAdapter(&transport).ViewStatus(context.Background(), applicationpullrequest.StatusRequest{Owner: "alice", Name: "project", Number: 12})
	if err != nil || result.ChecksAvailable || result.Mergeable != applicationpullrequest.MergeableNo {
		t.Fatalf("unexpected result: %+v err=%v", result, err)
	}
}
