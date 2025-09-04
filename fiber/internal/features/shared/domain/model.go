package domain

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
	ID           int64       `gorm:"primaryKey"`
	ExternalID   uuid.UUID   `gorm:"index:idx_order_external_id,unique"`
	Items        []OrderItem `form:"many2many:order_items;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status       OrderStatus
	RegisteredAt time.Time
}

func RegisterOrder(items []OrderItem, now time.Time) (*Order, *OrderRegistered) {
	order := &Order{
		ExternalID:   uuid.New(),
		Items:        items,
		Status:       OrderStatusRegistered,
		RegisteredAt: now,
	}

	event := &OrderRegistered{
		OrderID:    order.ExternalID,
		OccurredAt: now,
		Items:      make([]OrderRegisteredItem, len(items)),
	}

	for i, item := range items {
		event.Items[i] = OrderRegisteredItem{
			DishID:   item.Dish.ExternalID,
			Quantity: item.Quantity,
		}
	}

	return order, event
}

type OrderItem struct {
	OrderID  int64 `gorm:"primaryKey"`
	Order    Order
	DishID   int64 `gorm:"primaryKey"`
	Dish     Dish
	Quantity uint8
}

func NewOrderItem(dish Dish, quantity uint8) OrderItem {
	return OrderItem{
		Dish:     dish,
		Quantity: quantity,
	}
}

type Dish struct {
	ID         int64
	ExternalID uuid.UUID
	Name       string
}
