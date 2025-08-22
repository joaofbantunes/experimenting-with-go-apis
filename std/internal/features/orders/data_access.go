package orders

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DataAccess struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewDataAccess(pool *pgxpool.Pool, loggerProvider func(name string) *slog.Logger) *DataAccess {
	return &DataAccess{
		pool:   pool,
		logger: loggerProvider("orders_data_access"),
	}
}

func (db *DataAccess) GetDishesByID(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]DishRef, error) {
	rows, err := db.pool.Query(ctx,
		// language=postgresql
		`SELECT id, external_id
		FROM dishes
		WHERE external_id = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishes = make(map[uuid.UUID]DishRef)
	for rows.Next() {
		var dish DishRef
		err := rows.Scan(&dish.ID, &dish.ExternalID)
		if err != nil {
			return nil, err
		}
		dishes[dish.ExternalID] = dish
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return dishes, nil
}

func (db *DataAccess) RegisterOrder(ctx context.Context, order Order) error {
	// Start a transaction
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			db.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("error", err))
		}
	}(tx, ctx)

	dishIds := make([]int64, len(order.Items))
	quantities := make([]uint8, len(order.Items))
	for i, item := range order.Items {
		dishIds[i] = item.DishId
		quantities[i] = item.Quantity
	}

	_, err = db.pool.Exec(ctx,
		// language=postgresql
		`WITH inserted_order AS (
			INSERT INTO orders (external_id, status, registered_at)
			VALUES ($1, $2, $3)
			RETURNING id
		)
		INSERT INTO order_items (order_id, dish_id, quantity)
		SELECT inserted_order.id, d, q
		FROM inserted_order, UNNEST($4::bigint[], $5::smallint[]) AS t(d, q)`,
		order.ExternalId, order.Status, order.RegisteredAt, dishIds, quantities)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
