package orders

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusRegistered        = OrderStatus("registered")
	OrderStatusInPreparation     = OrderStatus("in_preparation")
	OrderStatusWaitingForCourier = OrderStatus("waiting_for_courier")
	OrderStatusOutForDelivery    = OrderStatus("out_for_delivery")
	OrderStatusCompleted         = OrderStatus("completed")
	OrderStatusCancelled         = OrderStatus("cancelled")
)

type Order struct {
	ID           uuid.UUID
	Items        []OrderItem
	Status       OrderStatus
	RegisteredAt time.Time
}

func RegisterOrder(items []OrderItem, now time.Time) (*Order, *OrderRegistered) {
	order := &Order{
		ID:           uuid.New(),
		Items:        items,
		Status:       OrderStatusRegistered,
		RegisteredAt: now,
	}

	event := &OrderRegistered{
		OrderId:    order.ID,
		OccurredAt: now,
		Items:      make([]OrderRegisteredItem, len(items)),
	}

	for i, item := range items {
		event.Items[i] = OrderRegisteredItem{
			DishId:   item.DishID,
			Quantity: item.Quantity,
		}
	}

	return order, event
}

func (o *Order) Cancel(now time.Time) (*Order, *OrderCancelled, error) {
	if o.Status != OrderStatusRegistered {
		return o, nil, ErrOrderNoLongerCancellable{}
	}

	o.Status = OrderStatusCancelled

	event := &OrderCancelled{
		OrderID:    o.ID,
		OccurredAt: now,
	}

	return o, event, nil
}

type OrderItem struct {
	DishID   uuid.UUID
	Quantity uint8
}

func NewOrderItem(dishID uuid.UUID, quantity uint8) OrderItem {
	return OrderItem{
		DishID:   dishID,
		Quantity: quantity,
	}
}

type DishRef struct {
	ID         int64
	ExternalID uuid.UUID
}

type ErrOrderNoLongerCancellable struct{}

func (e ErrOrderNoLongerCancellable) Error() string {
	return "order no longer cancellable"
}
