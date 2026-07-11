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
