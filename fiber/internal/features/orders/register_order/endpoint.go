package register_order

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared"
	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/domain"
	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/outbox"
	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/problems"
	"gorm.io/gorm"
)

type orderItem struct {
	DishID   uuid.UUID `json:"dishId"`
	Quantity int       `json:"quantity"` // using the int type here to capture negative and large values for validation
}
type body struct {
	Items []orderItem `json:"items"`
}
type request struct {
	Body body
}
type response struct {
	ID uuid.UUID `json:"id"`
}

type unknownDishesError struct {
	DishIds []uuid.UUID `json:"dishIds"`
}

func NewRegisterOrderEndpoint(
	loggerProvider func(name string) *slog.Logger,
	db *gorm.DB,
	tp shared.TimeProvider) fiber.Handler {
	// logger := loggerProvider("register_order_endpoint")
	return func(c *fiber.Ctx) error {
		r, err := decodeAndValidate(c)
		if err != nil {
			return err
		}

		dishIds := getDishIds(r.Body.Items)
		dishes, err := queryDishesByID(c.UserContext(), db, dishIds)
		if err != nil {
			return err
		}
		if len(dishes) != len(dishIds) {
			return createUnknownDishesError(c.UserContext(), dishIds, dishes)
		}
		now := tp.Now()
		order, event := domain.RegisterOrder(mapOrderItems(r.Body.Items, dishes), now)
		err = saveOrder(c.UserContext(), db, order, event)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(response{ID: order.ExternalID})
	}
}

func mapOrderItems(items []orderItem, dishes map[uuid.UUID]*domain.Dish) []*domain.OrderItem {
	orderItems := make([]*domain.OrderItem, len(items))
	for i, item := range items {
		dish := dishes[item.DishID]
		orderItems[i] = domain.NewOrderItem(dish, uint8(item.Quantity))
	}
	return orderItems
}

func saveOrder(ctx context.Context, db *gorm.DB, order *domain.Order, event *domain.OrderRegistered) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[domain.Order](tx).Create(ctx, order); err != nil {
			return err
		}
		msg, err := outbox.MapMsg(order, event)
		if err != nil {
			return err
		}
		if err := gorm.G[outbox.OutboxMessage](tx).Create(ctx, msg); err != nil {
			return err
		}
		return nil
	})
}

func queryDishesByID(ctx context.Context, db *gorm.DB, dishIds []uuid.UUID) (map[uuid.UUID]*domain.Dish, error) {
	var dishes []*domain.Dish
	dishes, err := gorm.G[*domain.Dish](db).Where("external_id IN ?", dishIds).Find(ctx)
	if err != nil {
		return nil, err
	}
	dishMap := make(map[uuid.UUID]*domain.Dish, len(dishes))
	for _, dish := range dishes {
		dishMap[dish.ExternalID] = dish
	}
	return dishMap, nil
}

func getDishIds(items []orderItem) []uuid.UUID {
	ids := make([]uuid.UUID, len(items))
	for i, item := range items {
		ids[i] = item.DishID
	}
	return ids
}

func createUnknownDishesError(ctx context.Context, dishIds []uuid.UUID, dishes map[uuid.UUID]*domain.Dish) error {
	missingIds := make([]uuid.UUID, 0)
	for _, id := range dishIds {
		if _, ok := dishes[id]; !ok {
			missingIds = append(missingIds, id)
		}
	}
	return problems.NewProblemError(
		ctx,
		problems.ProblemOrdersUnknownDishes,
		fiber.StatusUnprocessableEntity,
		"Some dishes are not known",
		"Some dishes are not known",
		unknownDishesError{DishIds: missingIds},
	)
}

func decodeAndValidate(c *fiber.Ctx) (request, error) {
	var b body
	if err := c.BodyParser(&b); err != nil {
		return request{}, problems.NewValidationProblemError(
			c.UserContext(),
			"Invalid request",
			"Invalid request",
			[]problems.ValidationError{
				{
					Description: "Invalid request body",
					Pointer:     shared.RootPointer.String(),
				},
			},
		)
	}

	v := shared.NewValidator(c.UserContext())
	v.Required(shared.NewJSONPointerBuilder().Key("items").Build(), len(b.Items) > 0, "At least one item is required")
	for i, item := range b.Items {
		v.Required(
			shared.NewJSONPointerBuilder().Key("items").Index(i).Key("dishId").Build(),
			item.DishID != uuid.Nil,
			"Order item missing dish is required",
		)
		v.InRange(
			shared.NewJSONPointerBuilder().Key("items").Index(i).Key("quantity").Build(),
			item.Quantity,
			1,
			100,
			"Order item quantity must be between 1 and 100",
		)
	}
	if err := v.ToError("Invalid request"); err != nil {
		return request{}, err
	}

	return request{Body: b}, nil
}
