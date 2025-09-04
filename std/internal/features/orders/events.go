package orders

import (
	"time"

	"github.com/google/uuid"
)

type OrderRegisteredItem struct {
	DishId   uuid.UUID
	Quantity uint8
}

type OrderRegistered struct {
	OrderId    uuid.UUID
	Items      []OrderRegisteredItem
	OccurredAt time.Time
}

type OrderCancelled struct {
	OrderID    uuid.UUID
	OccurredAt time.Time
}
