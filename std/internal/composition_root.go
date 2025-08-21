package internal

import "log/slog"

type CompositionRoot struct {
	Config         *Config
	LoggerProvider func(string) *slog.Logger
}

func NewCompositionRoot(config *Config) *CompositionRoot {
	return &CompositionRoot{
		Config:         config,
		LoggerProvider: CreateLogger,
	}
}
