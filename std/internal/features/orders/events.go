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
