package get_order_details

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared/problems"
)

type request struct {
	OrderID uuid.UUID
}

type orderItem struct {
	DishID   uuid.UUID `json:"dishId"`
	Quantity uint8     `json:"quantity"`
}

type response struct {
	OrderID uuid.UUID          `json:"orderId"`
	Status  orders.OrderStatus `json:"status"`
	Items   []orderItem        `json:"items"`
}

func NewGetOrderDetailsEndpoint(
	db *orders.DataAccess,
	pe problems.ProblemEncoder,
	loggerProvider func(name string) *slog.Logger,
) http.Handler {
	logger := loggerProvider("get_order_details_endpoint")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, prob := decodeAndValidate(r)

		if prob != nil {
			pe.EncodeValidationProblem(r.Context(), w, prob)
			return
		}

		o, present, err := db.GetOrderByID(r.Context(), req.OrderID)

		if err != nil {
			shared.EncodeInternalServerError(r.Context(), w, logger, err)
			return
		}

		if !present {
			pe.EncodeProblem(r.Context(), w, problems.NewProblem(
				r.Context(),
				problems.ProblemGeneralNotFound,
				http.StatusNotFound,
				"Order not found",
				"Order not found",
				nil,
			))
			return
		}

		resp := response{
			OrderID: o.ID,
			Status:  o.Status,
			Items:   make([]orderItem, len(o.Items)),
		}
		for i, item := range o.Items {
			resp.Items[i] = orderItem{
				DishID:   item.DishID,
				Quantity: item.Quantity,
			}
		}

		err = shared.Encode(
			w,
			http.StatusOK,
			resp)

		if err != nil {
			shared.EncodeInternalServerError(r.Context(), w, logger, err)
			return
		}
	})
}

func decodeAndValidate(r *http.Request) (request, *problems.ValidationProblem) {
	orderId, err := uuid.Parse(r.PathValue("orderId"))
	if err != nil {
		return request{}, problems.NewValidationProblem(
			r.Context(),
			"Invalid request",
			"Invalid request",
			[]problems.ValidationError{
				{
					Description: "Invalid order id",
					Parameter:   "orderId",
				},
			},
		)
	}
	return request{OrderID: orderId}, nil
}
