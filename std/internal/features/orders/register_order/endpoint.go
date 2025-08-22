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
	Quantity int       `json:"quantity"`
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

func NewRegisterOrderEndpoint(
	db *orders.DataAccess,
	pe shared.ProblemEncoder,
	loggerProvider func(name string) *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	logger := loggerProvider("register_order_endpoint")
	return func(w http.ResponseWriter, r *http.Request) {
		req, problem := decodeAndValidate(r)

		if problem != nil {
			pe.EncodeValidationProblem(r.Context(), w, problem)
			return
		}

		err := shared.Encode(
			w,
			http.StatusOK,
			response{
				Greeting: "Hello, World!",
			})

		if err != nil {
			shared.InternalServerError(r.Context(), w, logger, err)
			return
		}
	}

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
				Description: "Dish ID is required",
				Pointer:     shared.JsonPointerForSegments([]string{"items", strconv.Itoa(i), "dishId"}).String(),
			})
		}
		if item.Quantity <= 0 {
			errors = append(errors, shared.ValidationError{
				Description: "Quantity must be greater than zero",
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
