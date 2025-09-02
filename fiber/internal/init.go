package internal

import (
	"context"
	"log/slog"
	"strings"

	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/domain"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDb(ctx context.Context, config *Config) (*gorm.DB, error) {
	replacer := strings.NewReplacer(
		"{user}", config.Database.User,
		"{password}", config.Database.Password,
	)
	return gorm.Open(
		postgres.Open(replacer.Replace(config.Database.BaseConnStr)),
		&gorm.Config{})
}

func MigrateDB(ctx context.Context, db *gorm.DB, logger *slog.Logger, tracer trace.Tracer) error {
	ctx, span := tracer.Start(ctx, "migrate db")
	logger.DebugContext(ctx, "migrate db")
	defer span.End()

	// using auto migrate just because I'm testing things out, wouldn't use it for a real application

	// adjacent, but it seems gorm doesn't support stuff like postgres enums, at least not in the migrations,
	// so would need to investigate how to work around this
	// for now, letting it use the default string

	return db.AutoMigrate(
		&domain.Order{},
		&domain.Dish{},
		&domain.OrderItem{},
	)
}
