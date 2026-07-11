package issue

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationissue "github.com/l4l4dev/fj/internal/application/issue"
)

type transport interface {
	Do(context.Context, string, string, url.Values) (*http.Response, error)
}

type RESTAdapter struct{ transport transport }

type forgejoIssue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Body   string `json:"body"`
}

func (adapter *RESTAdapter) Inspect(ctx context.Context, request applicationissue.InspectRequest) (applicationissue.IssueDetail, error) {
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues/" + strconv.Itoa(request.Number)
	response, err := adapter.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return applicationissue.IssueDetail{}, translateInspectError(err)
	}
	defer response.Body.Close()
	var decoded forgejoIssue
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationissue.IssueDetail{}, apperror.New(apperror.Remote, "inspect issue", "")
	}
	state := applicationissue.StateClosed
	if decoded.State == string(applicationissue.StateOpen) {
		state = applicationissue.StateOpen
	}
	return applicationissue.IssueDetail{Number: decoded.Number, Title: decoded.Title, State: state, Body: decoded.Body}, nil
}

func NewRESTAdapter(transport transport) *RESTAdapter { return &RESTAdapter{transport: transport} }

func (adapter *RESTAdapter) List(ctx context.Context, request applicationissue.ListRequest) (applicationissue.Page, error) {
	query := url.Values{}
	query.Set("page", strconv.Itoa(request.Page))
	query.Set("limit", strconv.Itoa(request.Limit))
	query.Set("state", string(request.State))
	if request.Filter.Assignee != "" {
		query.Set("assignee", request.Filter.Assignee)
	}
	if request.Filter.Label != "" {
		query.Set("labels", request.Filter.Label)
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues"
	response, err := adapter.transport.Do(ctx, http.MethodGet, path, query)
	if err != nil {
		return applicationissue.Page{}, translateRemoteError(err)
	}
	defer response.Body.Close()
	var decoded []forgejoIssue
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationissue.Page{}, apperror.New(apperror.Remote, "list issues", "")
	}
	result := applicationissue.Page{Issues: make([]applicationissue.Issue, 0, len(decoded)), Page: request.Page, Limit: request.Limit, MorePages: len(decoded) == request.Limit}
	for _, item := range decoded {
		state := applicationissue.StateClosed
		if item.State == string(applicationissue.StateOpen) {
			state = applicationissue.StateOpen
		}
		result.Issues = append(result.Issues, applicationissue.Issue{Number: item.Number, Title: item.Title, State: state})
	}
	return result, nil
}

func translateRemoteError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		category := apperror.Remote
		message := ""
		switch status.StatusCode() {
		case 401, 403:
			category = apperror.Authentication
		case 404:
			category = apperror.NotFound
			message = "repository not found"
		}
		return apperror.New(category, "list issues", message)
	}
	return apperror.New(apperror.Remote, "list issues", "")
}

func translateInspectError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		category := apperror.Remote
		message := ""
		switch status.StatusCode() {
		case 401, 403:
			category = apperror.Authentication
		case 404:
			category = apperror.NotFound
			message = "issue not found"
		}
		return apperror.New(category, "inspect issue", message)
	}
	return apperror.New(apperror.Remote, "inspect issue", "")
}

var _ applicationissue.Lister = (*RESTAdapter)(nil)
var _ applicationissue.Inspector = (*RESTAdapter)(nil)
