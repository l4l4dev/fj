package release

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type PublishUseCase struct {
	inspector Inspector
	publisher Publisher
}

func NewPublishUseCase(inspector Inspector, publisher Publisher) PublishUseCase {
	return PublishUseCase{inspector: inspector, publisher: publisher}
}

func (u PublishUseCase) Execute(ctx context.Context, request InspectRequest) (ReleaseDetail, error) {
	if strings.TrimSpace(request.Owner) == "" || strings.TrimSpace(request.Name) == "" {
		return ReleaseDetail{}, apperror.NewValidation("publish release", "OWNER/NAME owner and name are required")
	}
	if strings.TrimSpace(request.Tag) == "" {
		return ReleaseDetail{}, apperror.NewValidation("publish release", "release tag is required")
	}
	if u.inspector == nil || u.publisher == nil {
		return ReleaseDetail{}, apperror.New(apperror.Internal, "publish release", "")
	}
	detail, err := u.inspector.Inspect(ctx, request)
	if err != nil {
		return ReleaseDetail{}, err
	}
	if !detail.Draft {
		return ReleaseDetail{}, apperror.New(apperror.Conflict, "publish release", "release is already published")
	}
	return u.publisher.Publish(ctx, PublishRequest{Owner: request.Owner, Name: request.Name, ID: detail.ID})
}
