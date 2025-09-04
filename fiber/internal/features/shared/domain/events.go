package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderRegisteredItem struct {
	DishID   uuid.UUID
	Quantity uint8
}

type OrderRegistered struct {
	OrderID    uuid.UUID
	Items      []OrderRegisteredItem
	OccurredAt time.Time
}

type OrderCancelled struct {
	OrderId    uuid.UUID
	OccurredAt time.Time
}
