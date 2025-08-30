package internal

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	apiMetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	apiTrace "go.opentelemetry.io/otel/trace"
)

type O11yContext struct {
	Shutdown       func(ctx context.Context) error
	LoggerProvider func(string) *slog.Logger
	TracerProvider func(string) apiTrace.Tracer
	MeterProvider  func(string) apiMetric.Meter
}

// SetupOTelSDK sets up OpenTelemetry SDK with trace, metric, and log providers.
// It returns a struct with providers for the signals, as well as a shutdown function that should be called to clean up resources.
func SetupOTelSDK(ctx context.Context) (o11y *O11yContext, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	environment := os.Getenv("APP_ENV")

	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceInstanceIDKey.String(hostname),
		semconv.TelemetrySDKLanguageGo,
		attribute.String("environment", environment),
	)

	// Set up trace provider.
	tracerProvider, err := newTracerProvider(ctx, r)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider(ctx, r)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider(ctx, r)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return &O11yContext{
		Shutdown: shutdown,
		LoggerProvider: func(name string) *slog.Logger {
			consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}).WithAttrs([]slog.Attr{slog.String("logger_name", name)})
			otlpHandler := otelslog.NewHandler(name)
			return slog.New(slogmulti.Fanout(consoleHandler, otlpHandler))
		},
		TracerProvider: func(name string) apiTrace.Tracer {
			return tracerProvider.Tracer(name)
		},
		MeterProvider: func(name string) apiMetric.Meter {
			return meterProvider.Meter(name)
		},
	}, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(ctx context.Context, r *resource.Resource) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(r),
		trace.WithBatcher(traceExporter),
	)
	return tracerProvider, nil
}

func newMeterProvider(ctx context.Context, r *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(r),
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithView(metric.NewView(
			metric.Instrument{
				// we can also use wildcards, but now just wanted to apply to this specific instrument
				Name: "http.server.duration",
				// could use Kind, if we wanted to apply this to any histogram, not just rpc.server.duration
				// Kind: metric.InstrumentKindHistogram,
			},
			metric.Stream{
				Aggregation: metric.AggregationExplicitBucketHistogram{
					Boundaries: []float64{0, 0.5, 1, 2.5, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000},
				}},
		)),
	)
	return meterProvider, nil
}

func newLoggerProvider(ctx context.Context, r *resource.Resource) (*log.LoggerProvider, error) {
	// fanning out through slog, because otel log exporter only exports json
	//consoleLogExporter, err := stdoutlog.New()
	//if err != nil {
	//	return nil, err
	//}
	grpcLogExporter, err := otlploggrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithResource(r),
		// log.WithProcessor(log.NewBatchProcessor(consoleLogExporter)),
		log.WithProcessor(log.NewBatchProcessor(grpcLogExporter)),
	)
	return loggerProvider, nil
}
