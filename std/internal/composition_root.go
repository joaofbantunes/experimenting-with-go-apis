package internal

import (
	"context"
	"log/slog"

	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type CompositionRoot struct {
	Config           *Config
	OrdersDataAccess *orders.DataAccess
	ProblemEncoder   shared.ProblemEncoder
	TimeProvider     shared.TimeProvider
	o11y             *O11yContext
}

func (root *CompositionRoot) CreateLogger(name string) *slog.Logger {
	return root.o11y.LoggerProvider(name)
}

func (root *CompositionRoot) CreateTracer(name string) trace.Tracer {
	return root.o11y.TracerProvider(name)
}

func (root *CompositionRoot) CreateMeter(name string) metric.Meter {
	return root.o11y.MeterProvider(name)
}

func (root *CompositionRoot) Shutdown(ctx context.Context) error {
	return root.o11y.Shutdown(ctx)
}

func NewCompositionRoot(ctx context.Context, config *Config) (*CompositionRoot, error) {
	o11y, err := SetupOTelSDK(ctx)
	if err != nil {
		return nil, err
	}

	pool, err := CreatePool(ctx, config)
	if err != nil {
		return nil, err
	}

	root := &CompositionRoot{
		Config:       config,
		o11y:         o11y,
		TimeProvider: shared.NewSystemTimeProvider(),
	}
	root.ProblemEncoder = shared.NewProblemEncoder(root.CreateLogger)
	root.OrdersDataAccess = orders.NewDataAccess(pool, root.CreateLogger)

	err = MigrateDB(ctx, pool, root.CreateLogger("migrations"), root.CreateTracer("migrations"))
	if err != nil {
		return nil, err
	}

	return root, nil
}
