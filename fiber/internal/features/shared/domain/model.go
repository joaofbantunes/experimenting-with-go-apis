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
	ID           int64        `gorm:"primaryKey"`
	ExternalID   uuid.UUID    `gorm:"type:uuid;index:idx_orders_external_id,unique"`
	Items        []*OrderItem `gorm:"foreignKey:OrderID"`
	Status       OrderStatus
	RegisteredAt time.Time
}

func RegisterOrder(items []*OrderItem, now time.Time) (*Order, *OrderRegistered) {
	order := &Order{
		ExternalID:   uuid.New(),
		Items:        items,
		Status:       OrderStatusRegistered,
		RegisteredAt: now,
	}

	for i := range order.Items {
		order.Items[i].Order = order
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
	OrderID  int64  `gorm:"primaryKey"`
	Order    *Order `gorm:"constraint:OnDelete:CASCADE"`
	DishID   int64  `gorm:"primaryKey"`
	Dish     *Dish  `gorm:"constraint:OnDelete:RESTRICT"`
	Quantity uint8
}

func NewOrderItem(dish *Dish, quantity uint8) *OrderItem {
	return &OrderItem{
		Dish:     dish,
		Quantity: quantity,
	}
}

type Dish struct {
	ID         int64     `gorm:"primaryKey"`
	ExternalID uuid.UUID `gorm:"type:uuid;index:idx_dishes_external_id,unique"`
	Name       string
}
