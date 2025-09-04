package outbox

import (
	"time"
)

type OutboxMessage struct {
	ID           int64 `gorm:"primaryKey"`
	Type         string
	Payload      []byte
	CreatedAt    time.Time
	TraceContext []byte
}

func NewOutboxMessage(messageType string, payload []byte, createdAt time.Time, traceContext []byte) *OutboxMessage {
	return &OutboxMessage{
		Type:         messageType,
		Payload:      payload,
		CreatedAt:    createdAt,
		TraceContext: traceContext,
	}
}
