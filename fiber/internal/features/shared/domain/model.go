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
	ID           int64
	ExternalID   uuid.UUID
	Items        []OrderItem
	Status       OrderStatus
	RegisteredAt time.Time
}

type OrderItem struct {
	OrderID  int64 `gorm:"primaryKey"`
	Order    Order
	DishID   int64 `gorm:"primaryKey"`
	Dish     Dish
	Quantity uint8
}

type Dish struct {
	ID         int64
	ExternalID uuid.UUID
	Name       string
}
