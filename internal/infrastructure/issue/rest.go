package issue

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/l4l4dev/fj/internal/application/apperror"
	applicationissue "github.com/l4l4dev/fj/internal/application/issue"
)

type transport interface {
	Do(context.Context, string, string, url.Values) (*http.Response, error)
}

type jsonTransport interface {
	DoJSON(context.Context, string, string, url.Values, []byte) (*http.Response, error)
}

type RESTAdapter struct{ transport transport }

type forgejoIssue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Body   string `json:"body"`
}

type forgejoComment struct {
	ID   int64  `json:"id"`
	Body string `json:"body"`
}

type forgejoLabel struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type forgejoMilestone struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type forgejoIssueMilestone struct {
	Number    int64             `json:"number"`
	Milestone *forgejoMilestone `json:"milestone"`
}

type forgejoIssueAssignee struct {
	Username string `json:"username"`
}

type forgejoIssueAssignment struct {
	Assignee *forgejoIssueAssignee `json:"assignee"`
}

func (adapter *RESTAdapter) currentIssueAssignee(ctx context.Context, owner, name string, number int) (*applicationissue.Assignment, error) {
	path := "/api/v1/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(name) + "/issues/" + strconv.Itoa(number)
	response, err := adapter.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, translateAssignmentError(err, "get issue assignee")
	}
	defer response.Body.Close()
	var decoded forgejoIssueAssignment
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, apperror.New(apperror.Remote, "get issue assignee", "")
	}
	if decoded.Assignee == nil {
		return nil, nil
	}
	return &applicationissue.Assignment{Username: decoded.Assignee.Username}, nil
}

func (adapter *RESTAdapter) Assign(ctx context.Context, request applicationissue.AssignRequest) (applicationissue.Assignment, error) {
	current, err := adapter.currentIssueAssignee(ctx, request.Owner, request.Name, request.Number)
	if err != nil {
		return applicationissue.Assignment{}, err
	}
	if current != nil && current.Username == request.Username {
		return *current, nil
	}
	jsonClient, ok := adapter.transport.(jsonTransport)
	if !ok {
		return applicationissue.Assignment{}, apperror.New(apperror.Remote, "assign issue", "")
	}
	body, _ := json.Marshal(struct {
		Assignee string `json:"assignee"`
	}{request.Username})
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues/" + strconv.Itoa(request.Number)
	response, err := jsonClient.DoJSON(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return applicationissue.Assignment{}, translateAssignmentError(err, "assign issue")
	}
	defer response.Body.Close()
	var decoded forgejoIssueAssignment
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil || decoded.Assignee == nil {
		return applicationissue.Assignment{}, apperror.New(apperror.Remote, "assign issue", "")
	}
	return applicationissue.Assignment{Username: decoded.Assignee.Username}, nil
}

func (adapter *RESTAdapter) Unassign(ctx context.Context, request applicationissue.UnassignRequest) error {
	current, err := adapter.currentIssueAssignee(ctx, request.Owner, request.Name, request.Number)
	if err != nil {
		return err
	}
	if current == nil {
		return nil
	}
	jsonClient, ok := adapter.transport.(jsonTransport)
	if !ok {
		return apperror.New(apperror.Remote, "unassign issue", "")
	}
	body := []byte(`{"assignee":null}`)
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues/" + strconv.Itoa(request.Number)
	response, err := jsonClient.DoJSON(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return translateAssignmentError(err, "unassign issue")
	}
	response.Body.Close()
	return nil
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

func (adapter *RESTAdapter) Create(ctx context.Context, request applicationissue.CreateRequest) (applicationissue.IssueDetail, error) {
	jsonClient, ok := adapter.transport.(jsonTransport)
	if !ok {
		return applicationissue.IssueDetail{}, apperror.New(apperror.Remote, "create issue", "")
	}
	body, err := json.Marshal(struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}{Title: request.Title, Body: request.Body})
	if err != nil {
		return applicationissue.IssueDetail{}, apperror.New(apperror.Remote, "create issue", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues"
	response, err := jsonClient.DoJSON(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return applicationissue.IssueDetail{}, translateCreateError(err)
	}
	defer response.Body.Close()
	var decoded forgejoIssue
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationissue.IssueDetail{}, apperror.New(apperror.Remote, "create issue", "")
	}
	state := applicationissue.StateClosed
	if decoded.State == string(applicationissue.StateOpen) {
		state = applicationissue.StateOpen
	}
	return applicationissue.IssueDetail{Number: decoded.Number, Title: decoded.Title, State: state, Body: decoded.Body}, nil
}

func (adapter *RESTAdapter) Update(ctx context.Context, request applicationissue.UpdateRequest) (applicationissue.IssueDetail, error) {
	jsonClient, ok := adapter.transport.(jsonTransport)
	if !ok {
		return applicationissue.IssueDetail{}, apperror.New(apperror.Remote, "update issue", "")
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
		return applicationissue.IssueDetail{}, apperror.New(apperror.Remote, "update issue", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues/" + strconv.Itoa(request.Number)
	response, err := jsonClient.DoJSON(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return applicationissue.IssueDetail{}, translateUpdateError(err)
	}
	defer response.Body.Close()
	var decoded forgejoIssue
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationissue.IssueDetail{}, apperror.New(apperror.Remote, "update issue", "")
	}
	state := applicationissue.StateClosed
	if decoded.State == string(applicationissue.StateOpen) {
		state = applicationissue.StateOpen
	}
	return applicationissue.IssueDetail{Number: decoded.Number, Title: decoded.Title, State: state, Body: decoded.Body}, nil
}

func (adapter *RESTAdapter) ChangeState(ctx context.Context, request applicationissue.ChangeStateRequest) (applicationissue.IssueDetail, error) {
	jsonClient, ok := adapter.transport.(jsonTransport)
	if !ok {
		return applicationissue.IssueDetail{}, apperror.New(apperror.Remote, "change issue state", "")
	}
	body, err := json.Marshal(struct {
		State applicationissue.State `json:"state"`
	}{State: request.State})
	if err != nil {
		return applicationissue.IssueDetail{}, apperror.New(apperror.Remote, "change issue state", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues/" + strconv.Itoa(request.Number)
	response, err := jsonClient.DoJSON(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return applicationissue.IssueDetail{}, translateStateError(err)
	}
	defer response.Body.Close()
	var decoded forgejoIssue
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationissue.IssueDetail{}, apperror.New(apperror.Remote, "change issue state", "")
	}
	state := applicationissue.StateClosed
	if decoded.State == string(applicationissue.StateOpen) {
		state = applicationissue.StateOpen
	}
	return applicationissue.IssueDetail{Number: decoded.Number, Title: decoded.Title, State: state, Body: decoded.Body}, nil
}

func (adapter *RESTAdapter) ListComments(ctx context.Context, request applicationissue.ListCommentsRequest) ([]applicationissue.Comment, error) {
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues/" + strconv.Itoa(request.Number) + "/comments"
	response, err := adapter.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, translateCommentError(err, "list issue comments")
	}
	defer response.Body.Close()
	var decoded []forgejoComment
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, apperror.New(apperror.Remote, "list issue comments", "")
	}
	result := make([]applicationissue.Comment, 0, len(decoded))
	for _, comment := range decoded {
		result = append(result, applicationissue.Comment{ID: comment.ID, Body: comment.Body})
	}
	return result, nil
}

func (adapter *RESTAdapter) AddComment(ctx context.Context, request applicationissue.AddCommentRequest) (applicationissue.Comment, error) {
	jsonClient, ok := adapter.transport.(jsonTransport)
	if !ok {
		return applicationissue.Comment{}, apperror.New(apperror.Remote, "add issue comment", "")
	}
	body, err := json.Marshal(struct {
		Body string `json:"body"`
	}{Body: request.Body})
	if err != nil {
		return applicationissue.Comment{}, apperror.New(apperror.Remote, "add issue comment", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues/" + strconv.Itoa(request.Number) + "/comments"
	response, err := jsonClient.DoJSON(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return applicationissue.Comment{}, translateCommentError(err, "add issue comment")
	}
	defer response.Body.Close()
	var decoded forgejoComment
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationissue.Comment{}, apperror.New(apperror.Remote, "add issue comment", "")
	}
	return applicationissue.Comment{ID: decoded.ID, Body: decoded.Body}, nil
}

func (adapter *RESTAdapter) AddLabel(ctx context.Context, request applicationissue.AddLabelRequest) (applicationissue.Label, error) {
	labels, err := adapter.issueLabels(ctx, request.Owner, request.Name, request.Number)
	if err != nil {
		return applicationissue.Label{}, err
	}
	for _, label := range labels {
		if label.Name == request.Label {
			return label, nil
		}
	}
	jsonClient, ok := adapter.transport.(jsonTransport)
	if !ok {
		return applicationissue.Label{}, apperror.New(apperror.Remote, "add issue label", "")
	}
	body, err := json.Marshal(struct {
		Labels []string `json:"labels"`
	}{Labels: []string{request.Label}})
	if err != nil {
		return applicationissue.Label{}, apperror.New(apperror.Remote, "add issue label", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues/" + strconv.Itoa(request.Number) + "/labels"
	response, err := jsonClient.DoJSON(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return applicationissue.Label{}, translateLabelError(err, "add issue label")
	}
	defer response.Body.Close()
	var decoded []forgejoLabel
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return applicationissue.Label{}, apperror.New(apperror.Remote, "add issue label", "")
	}
	for _, label := range decoded {
		if label.Name == request.Label {
			return applicationissue.Label{ID: label.ID, Name: label.Name}, nil
		}
	}
	return applicationissue.Label{}, apperror.New(apperror.Remote, "add issue label", "")
}

func (adapter *RESTAdapter) RemoveLabel(ctx context.Context, request applicationissue.RemoveLabelRequest) (applicationissue.Label, error) {
	labels, err := adapter.issueLabels(ctx, request.Owner, request.Name, request.Number)
	if err != nil {
		return applicationissue.Label{}, err
	}
	var target applicationissue.Label
	for _, label := range labels {
		if label.Name == request.Label {
			target = label
			break
		}
	}
	if target.Name == "" {
		return applicationissue.Label{Name: request.Label}, nil
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues/" + strconv.Itoa(request.Number) + "/labels/" + strconv.FormatInt(target.ID, 10)
	response, err := adapter.transport.Do(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return applicationissue.Label{}, translateLabelError(err, "remove issue label")
	}
	defer response.Body.Close()
	return target, nil
}

func (adapter *RESTAdapter) SetMilestone(ctx context.Context, request applicationissue.SetMilestoneRequest) (applicationissue.Milestone, error) {
	milestones, err := adapter.resolveMilestones(ctx, request.Owner, request.Name)
	if err != nil {
		return applicationissue.Milestone{}, err
	}
	var target applicationissue.Milestone
	for _, milestone := range milestones {
		if milestone.Title == request.Milestone {
			target = milestone
			break
		}
	}
	if target.Title == "" {
		return applicationissue.Milestone{}, apperror.New(apperror.Remote, "resolve milestone", "milestone not found")
	}
	current, err := adapter.currentIssueMilestone(ctx, request.Owner, request.Name, request.Number)
	if err != nil {
		return applicationissue.Milestone{}, err
	}
	if current != nil && current.ID == target.ID {
		return target, nil
	}
	jsonClient, ok := adapter.transport.(jsonTransport)
	if !ok {
		return applicationissue.Milestone{}, apperror.New(apperror.Remote, "set issue milestone", "")
	}
	body, err := json.Marshal(struct {
		Milestone int64 `json:"milestone"`
	}{Milestone: target.ID})
	if err != nil {
		return applicationissue.Milestone{}, apperror.New(apperror.Remote, "set issue milestone", "")
	}
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues/" + strconv.Itoa(request.Number)
	response, err := jsonClient.DoJSON(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return applicationissue.Milestone{}, translateMilestoneError(err, "set issue milestone")
	}
	defer response.Body.Close()
	var decoded forgejoIssueMilestone
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil || decoded.Milestone == nil {
		return applicationissue.Milestone{}, apperror.New(apperror.Remote, "set issue milestone", "")
	}
	return applicationissue.Milestone{ID: decoded.Milestone.ID, Title: decoded.Milestone.Title}, nil
}

func (adapter *RESTAdapter) ClearMilestone(ctx context.Context, request applicationissue.ClearMilestoneRequest) error {
	current, err := adapter.currentIssueMilestone(ctx, request.Owner, request.Name, request.Number)
	if err != nil {
		return err
	}
	if current == nil {
		return nil
	}
	jsonClient, ok := adapter.transport.(jsonTransport)
	if !ok {
		return apperror.New(apperror.Remote, "clear issue milestone", "")
	}
	body := []byte(`{"milestone":null}`)
	path := "/api/v1/repos/" + url.PathEscape(request.Owner) + "/" + url.PathEscape(request.Name) + "/issues/" + strconv.Itoa(request.Number)
	response, err := jsonClient.DoJSON(ctx, http.MethodPatch, path, nil, body)
	if err != nil {
		return translateMilestoneError(err, "clear issue milestone")
	}
	defer response.Body.Close()
	var decoded forgejoIssueMilestone
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil && err != io.EOF {
		return apperror.New(apperror.Remote, "clear issue milestone", "")
	}
	return nil
}

func (adapter *RESTAdapter) resolveMilestones(ctx context.Context, owner, name string) ([]applicationissue.Milestone, error) {
	path := "/api/v1/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(name) + "/milestones"
	response, err := adapter.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, translateMilestoneError(err, "resolve milestone")
	}
	defer response.Body.Close()
	var decoded []forgejoMilestone
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, apperror.New(apperror.Remote, "resolve milestone", "")
	}
	result := make([]applicationissue.Milestone, 0, len(decoded))
	for _, milestone := range decoded {
		result = append(result, applicationissue.Milestone{ID: milestone.ID, Title: milestone.Title})
	}
	return result, nil
}

func (adapter *RESTAdapter) currentIssueMilestone(ctx context.Context, owner, name string, number int) (*applicationissue.Milestone, error) {
	path := "/api/v1/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(name) + "/issues/" + strconv.Itoa(number)
	response, err := adapter.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, translateMilestoneError(err, "get issue milestone")
	}
	defer response.Body.Close()
	var decoded forgejoIssueMilestone
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, apperror.New(apperror.Remote, "get issue milestone", "")
	}
	if decoded.Milestone == nil {
		return nil, nil
	}
	return &applicationissue.Milestone{ID: decoded.Milestone.ID, Title: decoded.Milestone.Title}, nil
}

func (adapter *RESTAdapter) issueLabels(ctx context.Context, owner, name string, number int) ([]applicationissue.Label, error) {
	path := "/api/v1/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(name) + "/issues/" + strconv.Itoa(number) + "/labels"
	response, err := adapter.transport.Do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, translateLabelError(err, "resolve issue labels")
	}
	defer response.Body.Close()
	var decoded []forgejoLabel
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, apperror.New(apperror.Remote, "resolve issue labels", "")
	}
	labels := make([]applicationissue.Label, 0, len(decoded))
	for _, label := range decoded {
		labels = append(labels, applicationissue.Label{ID: label.ID, Name: label.Name})
	}
	return labels, nil
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

func translateCreateError(err error) error {
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
		case 409:
			category = apperror.Conflict
			message = "issue creation conflict"
		}
		return apperror.New(category, "create issue", message)
	}
	return apperror.New(apperror.Remote, "create issue", "")
}

func translateUpdateError(err error) error {
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
		case 409:
			category = apperror.Conflict
			message = "issue update conflict"
		}
		return apperror.New(category, "update issue", message)
	}
	return apperror.New(apperror.Remote, "update issue", "")
}

func translateStateError(err error) error {
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
		case 409:
			category = apperror.Conflict
			message = "issue state conflict"
		}
		return apperror.New(category, "change issue state", message)
	}
	return apperror.New(apperror.Remote, "change issue state", "")
}

func translateCommentError(err error, operation string) error {
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
		return apperror.New(category, operation, message)
	}
	return apperror.New(apperror.Remote, operation, "")
}

func translateLabelError(err error, operation string) error {
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
		return apperror.New(category, operation, message)
	}
	return apperror.New(apperror.Remote, operation, "")
}

func translateMilestoneError(err error, operation string) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) {
		category := apperror.Remote
		if status.StatusCode() == 401 || status.StatusCode() == 403 {
			category = apperror.Authentication
		}
		return apperror.New(category, operation, "")
	}
	return apperror.New(apperror.Remote, operation, "")
}

func translateAssignmentError(err error, operation string) error {
	var status interface{ StatusCode() int }
	if errors.As(err, &status) && (status.StatusCode() == 401 || status.StatusCode() == 403) {
		return apperror.New(apperror.Authentication, operation, "")
	}
	return apperror.New(apperror.Remote, operation, "")
}

var _ applicationissue.Lister = (*RESTAdapter)(nil)
var _ applicationissue.Inspector = (*RESTAdapter)(nil)
var _ applicationissue.Creator = (*RESTAdapter)(nil)
var _ applicationissue.Updater = (*RESTAdapter)(nil)
var _ applicationissue.StateChanger = (*RESTAdapter)(nil)
var _ applicationissue.CommentViewer = (*RESTAdapter)(nil)
var _ applicationissue.CommentCreator = (*RESTAdapter)(nil)
var _ applicationissue.LabelAdder = (*RESTAdapter)(nil)
var _ applicationissue.LabelRemover = (*RESTAdapter)(nil)
var _ applicationissue.MilestoneSetter = (*RESTAdapter)(nil)
var _ applicationissue.MilestoneClearer = (*RESTAdapter)(nil)
var _ applicationissue.Assigner = (*RESTAdapter)(nil)
var _ applicationissue.Unassigner = (*RESTAdapter)(nil)
