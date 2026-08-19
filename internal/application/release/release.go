package release

import "context"

type Release struct {
	ID         int64
	TagName    string
	Title      string
	Draft      bool
	Prerelease bool
}

type ListRequest struct {
	Owner string
	Name  string
	Page  int
	Limit int
}

type Lister interface {
	List(context.Context, ListRequest) ([]Release, error)
}
