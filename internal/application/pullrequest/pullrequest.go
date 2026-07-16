package pullrequest

import "context"

type State string

const (
	StateOpen   State = "open"
	StateClosed State = "closed"
	StateAll    State = "all"
)

type PullRequest struct {
	Number     int
	Title      string
	State      State
	HeadBranch string
	BaseBranch string
}

type PullRequestDetail struct {
	Number     int
	Title      string
	State      State
	HeadBranch string
	BaseBranch string
	Body       string
}

type ListRequest struct {
	Owner string
	Name  string
	Page  int
	Limit int
	State State
}

type PullRequestLister interface {
	List(context.Context, ListRequest) ([]PullRequest, error)
}

type InspectRequest struct {
	Owner  string
	Name   string
	Number int
}

type PullRequestInspector interface {
	Inspect(context.Context, InspectRequest) (PullRequestDetail, error)
}

type CreateRequest struct {
	Owner      string
	Name       string
	Title      string
	HeadBranch string
	BaseBranch string
}

type PullRequestCreator interface {
	Create(context.Context, CreateRequest) (PullRequestDetail, error)
}

type AggregateState string

const (
	AggregatePending     AggregateState = "pending"
	AggregateSuccess     AggregateState = "success"
	AggregateFailed      AggregateState = "failed"
	AggregateUnavailable AggregateState = "unavailable"
)

type MergeableState string

const (
	MergeableYes         MergeableState = "yes"
	MergeableNo          MergeableState = "no"
	MergeableUnavailable MergeableState = "unavailable"
)

type StatusRequest struct {
	Owner  string
	Name   string
	Number int
}

type PullRequestStatus struct {
	Number    int
	Review    AggregateState
	Check     AggregateState
	Mergeable MergeableState
}

type Review struct {
	ID         int64
	ReviewerID int64
	State      string
	Dismissed  bool
	Stale      bool
}

type StatusData struct {
	Number             int
	Reviews            []Review
	ReviewsAvailable   bool
	RequestedReviewers int
	Checks             []string
	ChecksAvailable    bool
	Mergeable          MergeableState
}

type StatusViewer interface {
	ViewStatus(context.Context, StatusRequest) (StatusData, error)
}
