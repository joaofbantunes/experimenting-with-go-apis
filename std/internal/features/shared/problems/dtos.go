package problems

import (
	"context"
	"net/http"

	trace "go.opentelemetry.io/otel/trace"
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
		TraceId:      getTraceIDFromContext(ctx),
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
		Type:    ProblemGeneralValidation,
		Status:  http.StatusBadRequest,
		Title:   title,
		Detail:  detail,
		TraceId: getTraceIDFromContext(ctx),
		Errors:  errors,
	}
}

func getTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}
