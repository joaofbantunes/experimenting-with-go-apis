package shared

import (
	"context"
	"io"
	"log/slog"
)

// Close calls the Close method of the given io.Closer and logs any error using the provided logger.
// If the io.Closer implements CloserWithContext, it will call Close with the provided context.

func Close(ctx context.Context, c io.Closer, logger *slog.Logger) {
	err := c.Close()
	if err != nil {
		logger.ErrorContext(ctx, "error closing resource", slog.Any("error", err))
	}
}
