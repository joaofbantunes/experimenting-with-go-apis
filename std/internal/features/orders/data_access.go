package orders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (db *DataAccess) GetOrderByID(ctx context.Context, id uuid.UUID) (*Order, bool, error) {
	batch := &pgx.Batch{}
	batch.Queue(
		// language=postgresql
		`SELECT external_id, status, registered_at
		FROM orders
		WHERE external_id = $1`, id)
	batch.Queue(
		// language=postgresql
		`SELECT d.external_id, oi.quantity
		FROM order_items oi
		INNER JOIN orders o ON o.id = oi.order_id
		INNER JOIN dishes d ON d.id = oi.dish_id
		WHERE o.external_id = $1`, id)

	br := db.pool.SendBatch(ctx, batch)
	defer shared.Close(ctx, br, db.logger)

	var order Order
	err := br.QueryRow().Scan(&order.ID, &order.Status, &order.RegisteredAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	rows, err := br.Query()
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	order.Items = make([]OrderItem, 0)
	for rows.Next() {
		var item OrderItem
		err := rows.Scan(&item.DishID, &item.Quantity)
		if err != nil {
			return nil, false, err
		}
		order.Items = append(order.Items, item)
	}

	if rows.Err() != nil {
		return nil, false, rows.Err()
	}

	return &order, true, nil
}

func (db *DataAccess) RegisterOrder(ctx context.Context, o *Order, event *OrderRegistered) error {
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

	dishIds := make([]uuid.UUID, len(o.Items))
	quantities := make([]uint8, len(o.Items))
	for i, item := range o.Items {
		dishIds[i] = item.DishID
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
		o.ID,
		o.Status,
		o.RegisteredAt,
		dishIds,
		quantities)

	if err != nil {
		return err
	}

	msg, err := db.mapRegisteredMsg(o, event)
	if err != nil {
		return err
	}

	err = shared.InsertOutboxMessage(ctx, tx.Conn(), msg)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (db *DataAccess) UpdateOrder(ctx context.Context, o *Order, event any) error {
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

	_, err = tx.Conn().Exec(ctx,
		// language=postgresql
		`UPDATE orders
		SET status = $1
		WHERE external_id = $2`,
		o.Status,
		o.ID)

	if err != nil {
		return err
	}

	msg, err := db.mapMsg(o, event)
	if err != nil {
		return err
	}

	err = shared.InsertOutboxMessage(ctx, tx.Conn(), msg)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (db *DataAccess) mapMsg(o *Order, event any) (shared.OutboxMessage, error) {
	switch e := event.(type) {
	case *OrderRegistered:
		return db.mapRegisteredMsg(o, e)
	case *OrderCancelled:
		return db.mapCancelledMsg(o, e)
	default:
		return shared.OutboxMessage{}, errors.New(fmt.Sprintf("unsupported message type: %T", e))
	}
}

// in this case we don't need info from the aggregate, but in other cases we might, to populate the integration event
func (db *DataAccess) mapRegisteredMsg(o *Order, event *OrderRegistered) (shared.OutboxMessage, error) {
	// serializing the domain event directly, but normally, as the integration event might be different, we would map to a different struct
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

// in this case we don't need info from the aggregate, but in other cases we might, to populate the integration event
func (db *DataAccess) mapCancelledMsg(o *Order, event *OrderCancelled) (shared.OutboxMessage, error) {
	// serializing the domain event directly, but normally, as the integration event might be different, we would map to a different struct
	payload, err := json.Marshal(event)
	if err != nil {
		return shared.OutboxMessage{}, err
	}

	return shared.NewOutboxMessage(
		"OrderCancelled",
		payload,
		event.OccurredAt,
		nil, // TODO: trace context
	), nil
}
