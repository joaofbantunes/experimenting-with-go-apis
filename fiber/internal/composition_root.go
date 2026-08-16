package internal

import (
	"context"
	"io"
	"log/slog"

	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type CompositionRoot struct {
	Config       *Config
	TimeProvider shared.TimeProvider
	O11y         *O11yContext
	stdout       io.Writer
	DB           *gorm.DB
}

func (root *CompositionRoot) CreateLogger(name string) *slog.Logger {
	consoleHandler := slog.NewTextHandler(root.stdout, &slog.HandlerOptions{Level: slog.LevelDebug}).WithAttrs([]slog.Attr{slog.String("logger_name", name)})
	otlpHandler := otelslog.NewHandler(name)
	return slog.New(slogmulti.Fanout(consoleHandler, otlpHandler))
}

func (root *CompositionRoot) CreateTracer(name string) trace.Tracer {
	return root.O11y.TracerProvider.Tracer(name)
}

func (root *CompositionRoot) CreateMeter(name string) metric.Meter {
	return root.O11y.MeterProvider.Meter(name)
}

func (root *CompositionRoot) Shutdown(ctx context.Context) error {
	return root.O11y.Shutdown(ctx)
}

func NewCompositionRoot(
	ctx context.Context,
	config *Config,
	getenv func(string) string,
	stdout io.Writer) (*CompositionRoot, error) {
	o11y, err := SetupOTelSDK(ctx, getenv, stdout)
	if err != nil {
		return nil, err
	}

	db, err := CreateDb(ctx, config, o11y.TracerProvider)
	if err != nil {
		return nil, err
	}

	root := &CompositionRoot{
		Config:       config,
		O11y:         o11y,
		TimeProvider: shared.NewSystemTimeProvider(),
		stdout:       stdout,
		DB:           db,
	}

	err = MigrateDB(ctx, db, root.CreateLogger("migrations"), root.CreateTracer("migrations"))
	if err != nil {
		return nil, err
	}

	return root, nil
}
