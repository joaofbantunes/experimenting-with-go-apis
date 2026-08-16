package shared

import (
	"context"

	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/problems"
)

type Validator struct {
	ctx    context.Context
	errors []problems.ValidationError
}

func NewValidator(ctx context.Context) *Validator {
	return &Validator{ctx: ctx, errors: make([]problems.ValidationError, 0)}
}

func (v *Validator) Required(pointer JSONPointer, ok bool, description string) {
	if !ok {
		v.add(pointer.String(), description)
	}
}

func (v *Validator) InRange(pointer JSONPointer, value, min, max int, description string) {
	if value < min || value > max {
		v.add(pointer.String(), description)
	}
}

func (v *Validator) add(pointer, description string) {
	v.errors = append(v.errors, problems.ValidationError{
		Description: description,
		Pointer:     pointer,
	})
}

func (v *Validator) ToError(title string) error {
	if len(v.errors) == 0 {
		return nil
	}
	return problems.NewValidationProblemError(v.ctx, title, title, v.errors)
}
