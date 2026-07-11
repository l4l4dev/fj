package pullrequest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationpullrequest "github.com/l4l4dev/fj/internal/application/pullrequest"
)

type transport interface {
	Do(context.Context, string, string, url.Values) (*http.Response, error)
}

type RESTAdapter struct{ transport transport }

type forgejoPullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Head   struct {
		Ref string `json:"ref"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
}

func NewRESTAdapter(t transport) *RESTAdapter { return &RESTAdapter{transport: t} }

func (a *RESTAdapter) List(ctx context.Context, request applicationpullrequest.ListRequest) ([]applicationpullrequest.PullRequest, error) {
	query := url.Values{}
	query.Set("page", strconv.Itoa(request.Page))
	query.Set("limit", strconv.Itoa(request.Limit))
	query.Set("state", string(request.State))
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/pulls"
	response, err := a.transport.Do(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, translateError(err)
	}
	defer response.Body.Close()
	var decoded []forgejoPullRequest
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, apperror.New(apperror.Remote, "list pull requests", "")
	}
	result := make([]applicationpullrequest.PullRequest, 0, len(decoded))
	for _, item := range decoded {
		state := applicationpullrequest.State(item.State)
		result = append(result, applicationpullrequest.PullRequest{Number: item.Number, Title: item.Title, State: state, HeadBranch: item.Head.Ref, BaseBranch: item.Base.Ref})
	}
	return result, nil
}

func translateError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) && (status.StatusCode() == 401 || status.StatusCode() == 403) {
		return apperror.New(apperror.Authentication, "list pull requests", "")
	}
	return apperror.New(apperror.Remote, "list pull requests", "")
}

var _ applicationpullrequest.PullRequestLister = (*RESTAdapter)(nil)
