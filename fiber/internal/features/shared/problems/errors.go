package problems

import "context"

type ProblemError struct {
	Problem *Problem
}

func (pe *ProblemError) Error() string {
	return pe.Problem.Detail
}

func NewProblemError(ctx context.Context, typ string, status int, title string, detail string, objectDetail any) *ProblemError {
	return &ProblemError{
		Problem: NewProblem(ctx, typ, status, title, detail, objectDetail),
	}
}

type ValidationProblemError struct {
	ValidationProblem *ValidationProblem
}

func (vpe *ValidationProblemError) Error() string {
	return vpe.ValidationProblem.Detail
}

func NewValidationProblemError(ctx context.Context, title string, detail string, errors []ValidationError) *ValidationProblemError {
	return &ValidationProblemError{
		ValidationProblem: NewValidationProblem(ctx, title, detail, errors),
	}
}
