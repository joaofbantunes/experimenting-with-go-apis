package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Problem struct {
	Type         string `json:"type"`
	Status       int    `json:"status"`
	Title        string `json:"title"`
	Detail       string `json:"detail,omitempty"`
	TraceId      string `json:"traceId,omitempty"`
	ObjectDetail any    `json:"objectDetail,omitempty"`
}

func NewProblem(ctx context.Context, typ string, status int, title string, detail string, objectDetail any) *Problem {
	return &Problem{
		Type:         typ,
		Status:       status,
		Title:        title,
		Detail:       detail,
		TraceId:      "", // TODO: grab from context
		ObjectDetail: objectDetail,
	}
}

type ValidationProblem struct {
	Type    string            `json:"type"`
	Status  int               `json:"status"`
	Title   string            `json:"title"`
	Detail  string            `json:"detail,omitempty"`
	TraceId string            `json:"traceId,omitempty"`
	Errors  []ValidationError `json:"errors"`
}

type ValidationError struct {
	Description string `json:"description"`
	Pointer     string `json:"pointer,omitempty"`
	Parameter   string `json:"parameter,omitempty"`
	Header      string `json:"header,omitempty"`
}

func NewValidationProblem(ctx context.Context, title string, detail string, errors []ValidationError) *ValidationProblem {
	return &ValidationProblem{
		Type:    Validation,
		Status:  http.StatusBadRequest,
		Title:   title,
		Detail:  detail,
		TraceId: "", // TODO: grab from context
		Errors:  errors,
	}
}

type ProblemEncoder interface {
	EncodeProblem(ctx context.Context, w http.ResponseWriter, status int, problem *Problem)
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

func (pe *problemEncoder) EncodeProblem(ctx context.Context, w http.ResponseWriter, status int, problem *Problem) {
	err := encodeProblem(w, status, problem)
	if err != nil {
		InternalServerError(ctx, w, pe.logger, err)
		return
	}
}

func (pe *problemEncoder) EncodeValidationProblem(ctx context.Context, w http.ResponseWriter, validationProblem *ValidationProblem) {
	err := encodeProblem(w, http.StatusBadRequest, validationProblem)
	if err != nil {
		InternalServerError(ctx, w, pe.logger, err)
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
