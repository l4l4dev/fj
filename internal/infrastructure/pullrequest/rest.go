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

type jsonTransport interface {
	DoJSON(context.Context, string, string, url.Values, []byte) (*http.Response, error)
}

type RESTAdapter struct{ transport transport }

type forgejoPullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Head   struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
}

type forgejoPullRequestDetail struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	State     string `json:"state"`
	Body      string `json:"body"`
	Mergeable *bool  `json:"mergeable"`
	Head      struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
	RequestedReviewers []struct{} `json:"requested_reviewers"`
}

type forgejoPullReview struct {
	ID        int64  `json:"id"`
	State     string `json:"state"`
	Dismissed bool   `json:"dismissed"`
	Stale     bool   `json:"stale"`
	User      struct {
		ID int64 `json:"id"`
	} `json:"user"`
}

type forgejoCombinedStatus struct {
	Statuses []struct {
		Status string `json:"status"`
	} `json:"statuses"`
}

func NewRESTAdapter(t transport) *RESTAdapter { return &RESTAdapter{transport: t} }

func (a *RESTAdapter) SubmitReview(ctx context.Context, submission applicationpullrequest.ReviewSubmission) (applicationpullrequest.SubmittedReview, error) {
	jsonClient, ok := a.transport.(jsonTransport)
	if !ok {
		return applicationpullrequest.SubmittedReview{}, apperror.New(apperror.Remote, "submit pull request review", "")
	}
	body, err := json.Marshal(struct {
		Event applicationpullrequest.ReviewEvent `json:"event"`
		Body  string                             `json:"body,omitempty"`
	}{Event: submission.Event, Body: submission.Body})
	if err != nil {
		return applicationpullrequest.SubmittedReview{}, apperror.New(apperror.Remote, "submit pull request review", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(submission.Owner) + "/" + url.PathEscape(submission.Name) + "/pulls/" + strconv.Itoa(submission.Number) + "/reviews"
	response, err := jsonClient.DoJSON(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return applicationpullrequest.SubmittedReview{}, translateSubmitReviewError(err)
	}
	defer response.Body.Close()
	var decoded forgejoPullReview
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationpullrequest.SubmittedReview{}, apperror.New(apperror.Remote, "submit pull request review", "")
	}
	return applicationpullrequest.SubmittedReview{State: decoded.State}, nil
}

func (a *RESTAdapter) ViewStatus(ctx context.Context, request applicationpullrequest.StatusRequest) (applicationpullrequest.StatusData, error) {
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/pulls/" + strconv.Itoa(request.Number)
	response, err := a.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return applicationpullrequest.StatusData{}, translateStatusError(err)
	}
	defer response.Body.Close()
	var pull forgejoPullRequestDetail
	if err := json.NewDecoder(response.Body).Decode(&pull); err != nil {
		return applicationpullrequest.StatusData{}, apperror.New(apperror.Remote, "view pull request status", "")
	}

	result := applicationpullrequest.StatusData{
		Number:    pull.Number,
		Mergeable: applicationpullrequest.MergeableUnavailable,
	}
	if pull.Mergeable != nil {
		result.Mergeable = applicationpullrequest.MergeableNo
		if *pull.Mergeable {
			result.Mergeable = applicationpullrequest.MergeableYes
		}
	}

	reviews, available, err := a.reviews(ctx, path+"/reviews")
	if err != nil {
		return applicationpullrequest.StatusData{}, err
	}
	result.Reviews = reviews
	result.ReviewsAvailable = available
	result.RequestedReviewers = len(pull.RequestedReviewers)

	if pull.Head.SHA != "" {
		checkPath := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/commits/" + url.PathEscape(pull.Head.SHA) + "/status"
		checks, available, err := a.checks(ctx, checkPath)
		if err != nil {
			return applicationpullrequest.StatusData{}, err
		}
		result.Checks = checks
		result.ChecksAvailable = available
	}
	return result, nil
}

func (a *RESTAdapter) reviews(ctx context.Context, path string) ([]applicationpullrequest.Review, bool, error) {
	const pageSize = 50
	reviews := make([]forgejoPullReview, 0)
	for page := 1; ; page++ {
		query := url.Values{"page": {strconv.Itoa(page)}, "limit": {strconv.Itoa(pageSize)}}
		response, err := a.transport.Do(ctx, http.MethodGet, path, query)
		if err != nil {
			if componentUnavailable(err) {
				return nil, false, nil
			}
			return nil, false, translateStatusError(err)
		}
		var pageReviews []forgejoPullReview
		decodeErr := json.NewDecoder(response.Body).Decode(&pageReviews)
		response.Body.Close()
		if decodeErr != nil {
			return nil, false, apperror.New(apperror.Remote, "view pull request status", "")
		}
		reviews = append(reviews, pageReviews...)
		if len(pageReviews) < pageSize {
			break
		}
	}
	result := make([]applicationpullrequest.Review, 0, len(reviews))
	for _, review := range reviews {
		result = append(result, applicationpullrequest.Review{ID: review.ID, ReviewerID: review.User.ID, State: review.State, Dismissed: review.Dismissed, Stale: review.Stale})
	}
	return result, true, nil
}

func (a *RESTAdapter) checks(ctx context.Context, path string) ([]string, bool, error) {
	response, err := a.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		if componentUnavailable(err) {
			return nil, false, nil
		}
		return nil, false, translateStatusError(err)
	}
	defer response.Body.Close()
	var combined forgejoCombinedStatus
	if err := json.NewDecoder(response.Body).Decode(&combined); err != nil {
		return nil, false, apperror.New(apperror.Remote, "view pull request status", "")
	}
	result := make([]string, 0, len(combined.Statuses))
	for _, status := range combined.Statuses {
		result = append(result, status.Status)
	}
	return result, true, nil
}

func componentUnavailable(err error) bool {
	var status interface{ StatusCode() int }
	return errors.As(err, &status) && status.StatusCode() == http.StatusNotFound
}

func (a *RESTAdapter) Create(ctx context.Context, request applicationpullrequest.CreateRequest) (applicationpullrequest.PullRequestDetail, error) {
	jsonClient, ok := a.transport.(jsonTransport)
	if !ok {
		return applicationpullrequest.PullRequestDetail{}, apperror.New(apperror.Remote, "create pull request", "")
	}
	body, err := json.Marshal(struct {
		Title string `json:"title"`
		Head  string `json:"head"`
		Base  string `json:"base"`
	}{Title: request.Title, Head: request.HeadBranch, Base: request.BaseBranch})
	if err != nil {
		return applicationpullrequest.PullRequestDetail{}, apperror.New(apperror.Remote, "create pull request", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/pulls"
	response, err := jsonClient.DoJSON(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return applicationpullrequest.PullRequestDetail{}, translateCreateError(err)
	}
	defer response.Body.Close()
	var decoded forgejoPullRequestDetail
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationpullrequest.PullRequestDetail{}, apperror.New(apperror.Remote, "create pull request", "")
	}
	return applicationpullrequest.PullRequestDetail{Number: decoded.Number, Title: decoded.Title, State: applicationpullrequest.State(decoded.State), HeadBranch: decoded.Head.Ref, BaseBranch: decoded.Base.Ref}, nil
}

func (a *RESTAdapter) Update(ctx context.Context, request applicationpullrequest.UpdateRequest) (applicationpullrequest.PullRequestDetail, error) {
	jsonClient, ok := a.transport.(jsonTransport)
	if !ok {
		return applicationpullrequest.PullRequestDetail{}, apperror.New(apperror.Remote, "update pull request", "")
	}
	payload := make(map[string]string)
	if request.Title != nil {
		payload["title"] = *request.Title
	}
	if request.Body != nil {
		payload["body"] = *request.Body
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return applicationpullrequest.PullRequestDetail{}, apperror.New(apperror.Remote, "update pull request", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/pulls/" + strconv.Itoa(request.Number)
	response, err := jsonClient.DoJSON(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return applicationpullrequest.PullRequestDetail{}, translateUpdateError(err)
	}
	defer response.Body.Close()
	var decoded forgejoPullRequestDetail
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationpullrequest.PullRequestDetail{}, apperror.New(apperror.Remote, "update pull request", "")
	}
	return applicationpullrequest.PullRequestDetail{Number: decoded.Number, Title: decoded.Title, State: applicationpullrequest.State(decoded.State), HeadBranch: decoded.Head.Ref, BaseBranch: decoded.Base.Ref, Body: decoded.Body}, nil
}

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

func (a *RESTAdapter) Inspect(ctx context.Context, request applicationpullrequest.InspectRequest) (applicationpullrequest.PullRequestDetail, error) {
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/pulls/" + strconv.Itoa(request.Number)
	response, err := a.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return applicationpullrequest.PullRequestDetail{}, translateInspectError(err)
	}
	defer response.Body.Close()
	var decoded forgejoPullRequestDetail
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationpullrequest.PullRequestDetail{}, apperror.New(apperror.Remote, "inspect pull request", "")
	}
	return applicationpullrequest.PullRequestDetail{Number: decoded.Number, Title: decoded.Title, State: applicationpullrequest.State(decoded.State), HeadBranch: decoded.Head.Ref, BaseBranch: decoded.Base.Ref, Body: decoded.Body}, nil
}

func translateError(err error) error {
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
		return apperror.New(category, "list pull requests", message)
	}
	return apperror.New(apperror.Remote, "list pull requests", "")
}

func translateInspectError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, "inspect pull request", "")
		case 404:
			return apperror.New(apperror.NotFound, "inspect pull request", "pull request not found")
		}
	}
	return apperror.New(apperror.Remote, "inspect pull request", "")
}

func translateCreateError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, "create pull request", "")
		case 404:
			return apperror.New(apperror.NotFound, "create pull request", "repository or branch not found")
		case 409, 422:
			return apperror.New(apperror.Conflict, "create pull request", "pull request branches are invalid or conflict with an existing pull request")
		}
	}
	return apperror.New(apperror.Remote, "create pull request", "")
}

func translateUpdateError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, "update pull request", "permission denied or authentication failed")
		case 404:
			return apperror.New(apperror.NotFound, "update pull request", "pull request not found")
		case 409, 422:
			return apperror.New(apperror.Validation, "update pull request", "pull request fields were rejected by the remote service")
		}
	}
	return apperror.New(apperror.Remote, "update pull request", "")
}

func translateStatusError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, "view pull request status", "")
		case 404:
			return apperror.New(apperror.NotFound, "view pull request status", "pull request not found")
		}
	}
	return apperror.New(apperror.Remote, "view pull request status", "")
}

func translateSubmitReviewError(err error) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		switch status.StatusCode() {
		case 401, 403:
			return apperror.New(apperror.Authentication, "submit pull request review", "permission denied or authentication failed")
		case 404:
			return apperror.New(apperror.NotFound, "submit pull request review", "pull request not found")
		case 409:
			return apperror.New(apperror.Conflict, "submit pull request review", "review could not be submitted in the pull request's current state")
		case 422:
			return apperror.New(apperror.Validation, "submit pull request review", "review outcome or body was rejected by the remote service")
		}
	}
	return apperror.New(apperror.Remote, "submit pull request review", "")
}

var _ applicationpullrequest.PullRequestLister = (*RESTAdapter)(nil)
var _ applicationpullrequest.PullRequestInspector = (*RESTAdapter)(nil)
var _ applicationpullrequest.PullRequestCreator = (*RESTAdapter)(nil)
var _ applicationpullrequest.Updater = (*RESTAdapter)(nil)
var _ applicationpullrequest.StatusViewer = (*RESTAdapter)(nil)
var _ applicationpullrequest.ReviewSubmitter = (*RESTAdapter)(nil)
