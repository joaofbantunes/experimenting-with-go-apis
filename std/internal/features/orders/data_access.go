package orders

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
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

func (db *DataAccess) RegisterOrder(ctx context.Context, order Order, event OrderRegistered) error {
	// Start a transaction
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			db.logger.ErrorContext(ctx, "failed to rollback transaction", slog.Any("error", err))
		}
	}(tx, ctx)

	dishIds := make([]uuid.UUID, len(order.Items))
	quantities := make([]uint8, len(order.Items))
	for i, item := range order.Items {
		dishIds[i] = item.DishId
		quantities[i] = item.Quantity
	}

	_, err = tx.Conn().Exec(ctx,
		// language=postgresql
		`WITH inserted_order AS (
			INSERT INTO orders (external_id, status, registered_at)
			VALUES ($1, $2, $3)
			RETURNING id
		)
		INSERT INTO order_items (order_id, dish_id, quantity)
		SELECT inserted_order.id, d.id, t.q
		FROM 
		    inserted_order,
		    UNNEST($4::uuid[], $5::smallint[]) AS t(d, q)
				INNER JOIN dishes d ON d.external_id = t.d;`,
		order.ID,
		order.Status,
		order.RegisteredAt,
		dishIds,
		quantities)

	if err != nil {
		return err
	}

	msg, err := db.mapRegisteredMsg(event)
	if err != nil {
		return err
	}

	err = shared.InsertOutboxMessage(ctx, tx.Conn(), msg)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// in this case we don't need info from the aggregate, otherwise we would pass it as a parameter
func (db *DataAccess) mapRegisteredMsg(event OrderRegistered) (shared.OutboxMessage, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return shared.OutboxMessage{}, err
	}

	return shared.NewOutboxMessage(
		"OrderRegistered",
		payload,
		event.OccurredAt,
		nil, // TODO: trace context
	), nil
}
