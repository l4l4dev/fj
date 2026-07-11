package issue

import (
	"context"
	"strings"

	"github.com/l4l4dev/fj/internal/application/apperror"
)

type AddLabelUseCase struct{ adder LabelAdder }

func NewAddLabelUseCase(adder LabelAdder) AddLabelUseCase { return AddLabelUseCase{adder: adder} }

func (useCase AddLabelUseCase) Execute(ctx context.Context, request AddLabelRequest) (Label, error) {
	if err := validateLabelRequest(request.Owner, request.Name, request.Number, request.Label, "add issue label"); err != nil {
		return Label{}, err
	}
	if useCase.adder == nil {
		return Label{}, apperror.New(apperror.Internal, "add issue label", "")
	}
	return useCase.adder.AddLabel(ctx, request)
}

type RemoveLabelUseCase struct{ remover LabelRemover }

func NewRemoveLabelUseCase(remover LabelRemover) RemoveLabelUseCase {
	return RemoveLabelUseCase{remover: remover}
}

func (useCase RemoveLabelUseCase) Execute(ctx context.Context, request RemoveLabelRequest) (Label, error) {
	if err := validateLabelRequest(request.Owner, request.Name, request.Number, request.Label, "remove issue label"); err != nil {
		return Label{}, err
	}
	if useCase.remover == nil {
		return Label{}, apperror.New(apperror.Internal, "remove issue label", "")
	}
	return useCase.remover.RemoveLabel(ctx, request)
}

func validateLabelRequest(owner, name string, number int, label, operation string) error {
	if strings.TrimSpace(owner) == "" || strings.TrimSpace(name) == "" {
		return apperror.NewValidation(operation, "OWNER/NAME owner and name are required")
	}
	if number < 1 {
		return apperror.NewValidation(operation, "issue number must be a positive integer")
	}
	if strings.TrimSpace(label) == "" {
		return apperror.NewValidation(operation, "label is required")
	}
	return nil
}
