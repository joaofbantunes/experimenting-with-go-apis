package problems

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
)

type ProblemEncoder interface {
	EncodeProblem(ctx context.Context, w http.ResponseWriter, problem *Problem)
	EncodeValidationProblem(ctx context.Context, w http.ResponseWriter, validationProblem *ValidationProblem)
}

type problemEncoder struct {
	logger *slog.Logger
}

func NewProblemEncoder(loggerProvider func(string) *slog.Logger) ProblemEncoder {
	return &problemEncoder{
		logger: loggerProvider("problem_encoder"),
	}
}

func (pe *problemEncoder) EncodeProblem(ctx context.Context, w http.ResponseWriter, problem *Problem) {
	err := encodeProblem(w, problem.Status, problem)
	if err != nil {
		shared.EncodeInternalServerError(ctx, w, pe.logger, err)
		return
	}
}

func (pe *problemEncoder) EncodeValidationProblem(ctx context.Context, w http.ResponseWriter, validationProblem *ValidationProblem) {
	err := encodeProblem(w, http.StatusBadRequest, validationProblem)
	if err != nil {
		shared.EncodeInternalServerError(ctx, w, pe.logger, err)
		return
	}
}

func encodeProblem[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode problem: %w", err)
	}
	return nil
}
