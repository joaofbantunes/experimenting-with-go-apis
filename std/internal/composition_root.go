package internal

import (
	"context"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
)

type CompositionRoot struct {
	Config           *Config
	LoggerProvider   func(string) *slog.Logger
	OrdersDataAccess *orders.DataAccess
	ProblemEncoder   shared.ProblemEncoder
	TimeProvider     shared.TimeProvider
	pool             *pgxpool.Pool
	once             sync.Once
}

func NewCompositionRoot(config *Config) (*CompositionRoot, error) {
	root := &CompositionRoot{
		Config:         config,
		LoggerProvider: CreateLogger,
		ProblemEncoder: shared.NewProblemEncoder(CreateLogger),
		TimeProvider:   shared.NewSystemTimeProvider(),
	}

	return root, nil
}

func (root *CompositionRoot) InitApp(ctx context.Context) error {
	var err error
	root.once.Do(func() {
		root.pool, err = CreatePool(ctx, root.Config)
		if err != nil {
			return
		}
		err = MigrateDB(ctx, root.pool, root.LoggerProvider("migrations"))
		if err != nil {
			return
		}
		root.OrdersDataAccess = orders.NewDataAccess(root.pool, root.LoggerProvider)
	})
	return err
}
