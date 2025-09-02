package register_order

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared"
	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/domain"
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
		order := createOrder(r, dishes, now)
		event := createEvent(order, now)
		err = saveOrder(err, db, order, event)
		if err != nil {
			return err
		}
		return c.Status(http.StatusCreated).JSON(response{ID: order.ExternalID})
	}
}

func createOrder(r request, dishes map[uuid.UUID]domain.Dish, now time.Time) *domain.Order {
	order := &domain.Order{
		ExternalID:   uuid.New(),
		Status:       domain.OrderStatusRegistered,
		Items:        make([]domain.OrderItem, len(r.Body.Items)),
		RegisteredAt: now,
	}
	for i, item := range r.Body.Items {
		dish := dishes[item.DishID]
		order.Items[i] = domain.OrderItem{
			DishID:   dish.ID,
			Dish:     dish,
			Quantity: uint8(item.Quantity),
		}
	}
	return order
}

func createEvent(order *domain.Order, now time.Time) *domain.OrderRegistered {
	event := &domain.OrderRegistered{
		OrderId:    order.ExternalID,
		OccurredAt: now,
		Items:      make([]domain.OrderRegisteredItem, len(order.Items)),
	}
	for i, item := range order.Items {
		event.Items[i] = domain.OrderRegisteredItem{
			DishId:   item.Dish.ExternalID,
			Quantity: item.Quantity,
		}
	}
	return event
}

func saveOrder(err error, db *gorm.DB, order *domain.Order, event *domain.OrderRegistered) error {
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func queryDishesByID(ctx context.Context, db *gorm.DB, dishIds []uuid.UUID) (map[uuid.UUID]domain.Dish, error) {
	var dishes []domain.Dish
	dishes, err := gorm.G[domain.Dish](db).Where("external_id IN ?", dishIds).Find(ctx)
	if err != nil {
		return nil, err
	}
	dishMap := make(map[uuid.UUID]domain.Dish, len(dishes))
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

func createUnknownDishesError(ctx context.Context, dishIds []uuid.UUID, dishes map[uuid.UUID]domain.Dish) error {
	missingIds := make([]uuid.UUID, 0)
	for _, id := range dishIds {
		if _, ok := dishes[id]; !ok {
			missingIds = append(missingIds, id)
		}
	}
	return problems.NewProblemError(
		ctx,
		problems.ProblemOrdersUnknownDishes,
		http.StatusUnprocessableEntity,
		"Some dishes are not known",
		"Some dishes are not known",
		unknownDishesError{DishIds: missingIds},
	)
}

func decodeAndValidate(c *fiber.Ctx) (request, error) {
	var b body
	err := c.BodyParser(&b)
	if err != nil {
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
	errors := make([]problems.ValidationError, 0)
	if len(b.Items) == 0 {
		errors = append(errors, problems.ValidationError{
			Description: "At least one item is required",
			Pointer:     shared.JsonPointerForSegments([]string{"items"}).String(),
		})
	}

	for i, item := range b.Items {
		if item.DishID == uuid.Nil {
			errors = append(errors, problems.ValidationError{
				Description: "Order item missing dish is required",
				Pointer:     shared.JsonPointerForSegments([]string{"items", strconv.Itoa(i), "dishId"}).String(),
			})
		}
		if item.Quantity <= 0 {
			errors = append(errors, problems.ValidationError{
				Description: "Order item quantity must be greater than zero",
				Pointer:     shared.JsonPointerForSegments([]string{"items", strconv.Itoa(i), "quantity"}).String(),
			})
		}
		if item.Quantity > 100 {
			errors = append(errors, problems.ValidationError{
				Description: "Order item quantity must be less than or equal to 100",
				Pointer:     shared.JsonPointerForSegments([]string{"items", strconv.Itoa(i), "quantity"}).String(),
			})
		}
	}

	if len(errors) > 0 {
		return request{}, problems.NewValidationProblemError(
			c.UserContext(),
			"Invalid request",
			"Invalid request",
			errors,
		)
	}

	return request{Body: b}, nil
}
