package cancel_order

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
)

type request struct {
	OrderID uuid.UUID
}

func NewCancelOrderEndpoint(
	db *orders.DataAccess,
	pe shared.ProblemEncoder,
	loggerProvider func(name string) *slog.Logger,
	tp shared.TimeProvider) func(w http.ResponseWriter, r *http.Request) {
	logger := loggerProvider("cancel_order_endpoint")
	return func(w http.ResponseWriter, r *http.Request) {
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
			pe.EncodeProblem(r.Context(), w, shared.NewProblem(
				r.Context(),
				shared.ProblemGeneralNotFound,
				http.StatusNotFound,
				"Order not found",
				"Order not found",
				nil,
			))
			return
		}

		o, event, err := o.Cancel(tp.Now())

		if err != nil {
			if errors.Is(err, orders.ErrOrderNoLongerCancellable{}) {
				pe.EncodeProblem(r.Context(), w, shared.NewProblem(
					r.Context(),
					shared.ProblemOrdersNoLongerCancellable,
					http.StatusUnprocessableEntity,
					"Order cannot be cancelled",
					"Only orders with 'registered' status can be cancelled",
					nil,
				))
				return
			}
			shared.EncodeInternalServerError(r.Context(), w, logger, err)
			return
		}

		err = db.UpdateOrder(r.Context(), o, event)

		w.WriteHeader(http.StatusNoContent)
	}
}
func decodeAndValidate(r *http.Request) (request, *shared.ValidationProblem) {
	orderId, err := uuid.Parse(r.PathValue("orderId"))
	if err != nil {
		return request{}, shared.NewValidationProblem(
			r.Context(),
			"Invalid request",
			"Invalid request",
			[]shared.ValidationError{
				{
					Description: "Invalid order id",
					Parameter:   "orderId",
				},
			},
		)
	}
	return request{OrderID: orderId}, nil
}
