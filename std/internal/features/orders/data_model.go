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
	Id           int64
	ExternalId   uuid.UUID
	Items        []OrderItem
	Status       OrderStatus
	RegisteredAt time.Time
}

func NewOrder(items []OrderItem, now time.Time) Order {
	return Order{
		ExternalId:   uuid.New(),
		Items:        items,
		Status:       OrderStatusRegistered,
		RegisteredAt: now,
	}
}

type OrderItem struct {
	OrderId  int64
	DishId   int64
	Quantity uint8
}

func NewOrderItem(dishId int64, quantity uint8) OrderItem {
	return OrderItem{
		DishId:   dishId,
		Quantity: quantity,
	}
}

type DishRef struct {
	ID         int64
	ExternalID uuid.UUID
}
