package forgejo

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	applicationauth "github.com/l4l4dev/fj/internal/application/auth"
	"github.com/l4l4dev/fj/internal/application/config"
)

const defaultHTTPTimeout = 30 * time.Second

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

type client struct {
	endpoint   string
	credential applicationauth.Credential
	version    string
	httpClient doer
}

func NewClient(instance config.Instance, credential applicationauth.Credential, version string, httpClient doer) client {
	if httpClient == nil {
		httpClient = newHTTPClient()
	}
	return client{
		endpoint:   strings.TrimRight(strings.TrimSpace(string(instance.Endpoint)), "/"),
		credential: credential,
		version:    version,
		httpClient: httpClient,
	}
}

func (forgejoClient client) Do(ctx context.Context, method, apiPath string, query url.Values) (*http.Response, error) {
	return forgejoClient.do(ctx, method, apiPath, query, nil)
}

func (forgejoClient client) DoJSON(ctx context.Context, method, apiPath string, query url.Values, body []byte) (*http.Response, error) {
	return forgejoClient.do(ctx, method, apiPath, query, body)
}

func (forgejoClient client) do(ctx context.Context, method, apiPath string, query url.Values, body []byte) (*http.Response, error) {
	target, err := url.JoinPath(forgejoClient.endpoint, apiPath)
	if err != nil {
		return nil, newRemoteError("build request", 0)
	}
	requestURL, err := url.Parse(target)
	if err != nil {
		return nil, newRemoteError("build request", 0)
	}
	requestURL.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, method, requestURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, newRemoteError("build request", 0)
	}
	request.Header.Set("Authorization", "token "+forgejoClient.credential.Value())
	request.Header.Set("User-Agent", "fj/"+forgejoClient.version)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := forgejoClient.httpClient.Do(request)
	if err != nil {
		return nil, newRemoteError("request", 0)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		response.Body.Close()
		return nil, newRemoteError("request", response.StatusCode)
	}
	return response, nil
}

type RemoteError struct {
	operation  string
	statusCode int
}

func newRemoteError(operation string, statusCode int) RemoteError {
	return RemoteError{operation: operation, statusCode: statusCode}
}

func (err RemoteError) Error() string {
	if err.statusCode != 0 {
		return fmt.Sprintf("%s: remote request failed (status %d)", err.operation, err.statusCode)
	}
	return fmt.Sprintf("%s: remote request failed", err.operation)
}

func (err RemoteError) Operation() string {
	return err.operation
}

func (err RemoteError) StatusCode() int {
	return err.statusCode
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultHTTPTimeout,
		CheckRedirect: func(request *http.Request, previous []*http.Request) error {
			if len(previous) == 0 {
				return nil
			}
			last := previous[len(previous)-1].URL
			if request.URL.Scheme != last.Scheme || request.URL.Host != last.Host {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
}
