package internal

import (
	"log/slog"

	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
)

type CompositionRoot struct {
	Config         *Config
	LoggerProvider func(string) *slog.Logger
	problemEncoder shared.ProblemEncoder
	InitApp        func()
}

func NewCompositionRoot(config *Config) *CompositionRoot {
	return &CompositionRoot{
		Config:         config,
		LoggerProvider: CreateLogger,
		problemEncoder: shared.NewProblemEncoder(CreateLogger),
		InitApp:        Init,
	}
}
