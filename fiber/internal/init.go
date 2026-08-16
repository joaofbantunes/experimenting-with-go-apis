package internal

import (
	"context"
	"log/slog"
	"strings"

	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/domain"
	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/outbox"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func CreateDb(ctx context.Context, config *Config, tracerProvider trace.TracerProvider) (*gorm.DB, error) {
	replacer := strings.NewReplacer(
		"{user}", config.Database.User,
		"{password}", config.Database.Password,
	)
	db, err := gorm.Open(
		postgres.Open(replacer.Replace(config.Database.BaseConnStr)),
		&gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, db.Use(tracing.NewPlugin(tracing.WithTracerProvider(tracerProvider)))
}

func MigrateDB(ctx context.Context, db *gorm.DB, logger *slog.Logger, tracer trace.Tracer) error {
	ctx, span := tracer.Start(ctx, "migrate db")
	logger.DebugContext(ctx, "migrate db")
	defer span.End()

	// using auto migrate just because I'm testing things out, wouldn't use it for a real application

	// adjacent, but it seems gorm doesn't support stuff like postgres enums, at least not in the migrations,
	// so would need to investigate how to work around this
	// for now, letting it use the default string

	err := db.AutoMigrate(
		&domain.Order{},
		&domain.Dish{},
		&domain.OrderItem{},
		&outbox.OutboxMessage{},
	)

	return err
}
