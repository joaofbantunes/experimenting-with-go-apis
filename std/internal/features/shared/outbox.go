package shared

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type OutboxMessage struct {
	ID           int64
	Type         string
	Payload      []byte
	CreatedAt    time.Time
	TraceContext []byte
}

func NewOutboxMessage(messageType string, payload []byte, createdAt time.Time, traceContext []byte) OutboxMessage {
	return OutboxMessage{
		Type:         messageType,
		Payload:      payload,
		CreatedAt:    createdAt,
		TraceContext: traceContext,
	}
}

func InsertOutboxMessage(ctx context.Context, conn *pgx.Conn, message OutboxMessage) error {
	_, err := conn.Exec(
		ctx,
		// language=postgresql
		`INSERT INTO outbox_messages (type, payload, created_at, trace_context) VALUES ($1, $2, $3, $4)`,
		message.Type,
		message.Payload,
		message.CreatedAt,
		message.TraceContext,
	)
	return err
}
