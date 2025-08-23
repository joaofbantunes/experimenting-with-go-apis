package orders

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus int

const (
	OrderStatusRegistered        = OrderStatus(10)
	OrderStatusInPreparation     = OrderStatus(20)
	OrderStatusWaitingForCourier = OrderStatus(30)
	OrderStatusOutForDelivery    = OrderStatus(40)
	OrderStatusCompleted         = OrderStatus(50)
	OrderStatusCancelled         = OrderStatus(60)
)

type Order struct {
	ID           uuid.UUID
	Items        []OrderItem
	Status       OrderStatus
	RegisteredAt time.Time
}

func RegisterOrder(items []OrderItem, now time.Time) (Order, OrderRegistered) {
	order := Order{
		ID:           uuid.New(),
		Items:        items,
		Status:       OrderStatusRegistered,
		RegisteredAt: now,
	}

	event := OrderRegistered{
		OrderId:    order.ID,
		OccurredAt: now,
		Items:      make([]OrderRegisteredItem, len(items)),
	}

	for i, item := range items {
		event.Items[i] = OrderRegisteredItem{
			DishId:   item.DishId,
			Quantity: item.Quantity,
		}
	}

	return order, event
}

type OrderItem struct {
	DishId   uuid.UUID
	Quantity uint8
}

func NewOrderItem(dishId uuid.UUID, quantity uint8) OrderItem {
	return OrderItem{
		DishId:   dishId,
		Quantity: quantity,
	}
}

type DishRef struct {
	ID         int64
	ExternalID uuid.UUID
}
