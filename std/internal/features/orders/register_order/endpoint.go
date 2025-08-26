package register_order

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
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
	db *orders.DataAccess,
	pe shared.ProblemEncoder,
	loggerProvider func(name string) *slog.Logger,
	tp shared.TimeProvider) http.Handler {
	logger := loggerProvider("register_order_endpoint")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, prob := decodeAndValidate(r)

		if prob != nil {
			pe.EncodeValidationProblem(r.Context(), w, prob)
			return
		}

		dishIds := getDishIds(req.Body.Items)
		dishRefs, err := db.GetDishesByID(r.Context(), dishIds)
		if err != nil {
			shared.EncodeInternalServerError(r.Context(), w, logger, err)
			return
		}
		if len(dishRefs) != len(dishIds) {
			encodeUnknownDishes(w, r, dishIds, dishRefs, pe)
			return
		}

		order, event := orders.RegisterOrder(mapOrderItems(req.Body.Items), tp.Now())

		err = db.RegisterOrder(r.Context(), order, event)

		if err != nil {
			shared.EncodeInternalServerError(r.Context(), w, logger, err)
			return
		}

		err = shared.Encode(
			w,
			http.StatusCreated,
			response{
				ID: order.ID,
			})

		if err != nil {
			shared.EncodeInternalServerError(r.Context(), w, logger, err)
			return
		}
	})
}

func encodeUnknownDishes(w http.ResponseWriter, r *http.Request, dishIds []uuid.UUID, dishRefs map[uuid.UUID]orders.DishRef, pe shared.ProblemEncoder) {
	missingIds := make([]uuid.UUID, 0)
	for _, id := range dishIds {
		if _, ok := dishRefs[id]; !ok {
			missingIds = append(missingIds, id)
		}
	}
	pe.EncodeProblem(
		r.Context(),
		w,
		shared.NewProblem(
			r.Context(),
			shared.ProblemOrdersUnknownDishes,
			http.StatusUnprocessableEntity,
			"Some dishes are not known",
			"Some dishes are not known",
			unknownDishesError{DishIds: missingIds},
		),
	)
}

func getDishIds(items []orderItem) []uuid.UUID {
	ids := make([]uuid.UUID, len(items))
	for i, item := range items {
		ids[i] = item.DishID
	}
	return ids
}

func mapOrderItems(items []orderItem) []orders.OrderItem {
	result := make([]orders.OrderItem, len(items))
	for i, item := range items {
		result[i] = orders.OrderItem{
			DishID:   item.DishID,
			Quantity: uint8(item.Quantity),
		}
	}
	return result
}

func decodeAndValidate(r *http.Request) (request, *shared.ValidationProblem) {
	b, err := shared.Decode[body](r)
	if err != nil {
		return request{}, shared.NewValidationProblem(
			r.Context(),
			"Invalid request",
			"Invalid request",
			[]shared.ValidationError{
				{
					Description: "Invalid request body",
					Pointer:     shared.RootPointer.String(),
				},
			},
		)
	}
	errors := make([]shared.ValidationError, 0)
	if len(b.Items) == 0 {
		errors = append(errors, shared.ValidationError{
			Description: "At least one item is required",
			Pointer:     shared.JsonPointerForSegments([]string{"items"}).String(),
		})
	}

	for i, item := range b.Items {
		if item.DishID == uuid.Nil {
			errors = append(errors, shared.ValidationError{
				Description: "Order item missing dish is required",
				Pointer:     shared.JsonPointerForSegments([]string{"items", strconv.Itoa(i), "dishId"}).String(),
			})
		}
		if item.Quantity <= 0 {
			errors = append(errors, shared.ValidationError{
				Description: "Order item quantity must be greater than zero",
				Pointer:     shared.JsonPointerForSegments([]string{"items", strconv.Itoa(i), "quantity"}).String(),
			})
		}
		if item.Quantity > 100 {
			errors = append(errors, shared.ValidationError{
				Description: "Order item quantity must be less than or equal to 100",
				Pointer:     shared.JsonPointerForSegments([]string{"items", strconv.Itoa(i), "quantity"}).String(),
			})
		}
	}

	if len(errors) > 0 {
		return request{}, shared.NewValidationProblem(
			r.Context(),
			"Invalid request",
			"Invalid request",
			errors,
		)
	}

	return request{Body: b}, nil
}
