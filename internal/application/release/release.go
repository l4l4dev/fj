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

type Asset struct {
	ID   int64
	Name string
	Size int64
}

type ReleaseDetail struct {
	ID         int64
	TagName    string
	Title      string
	Draft      bool
	Prerelease bool
	Notes      string
	Assets     []Asset
}

type InspectRequest struct {
	Owner string
	Name  string
	Tag   string
}

type Inspector interface {
	Inspect(context.Context, InspectRequest) (ReleaseDetail, error)
}

type CreateRequest struct {
	Owner      string
	Name       string
	Tag        string
	Title      string
	Notes      string
	Prerelease bool
}

type Creator interface {
	Create(context.Context, CreateRequest) (ReleaseDetail, error)
}
