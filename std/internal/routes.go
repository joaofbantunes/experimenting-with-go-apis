package internal

import (
	"net/http"

	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/hello_world"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders/cancel_order"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders/get_order_details"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders/register_order"
)

func addRoutes(mux *http.ServeMux, root *CompositionRoot) {
	mux.HandleFunc("GET /hello", hello_world.NewHelloEndpoint(root.LoggerProvider))

	mux.HandleFunc("POST /api/v1/orders", register_order.NewRegisterOrderEndpoint(
		root.OrdersDataAccess,
		root.ProblemEncoder,
		root.LoggerProvider,
		root.TimeProvider))
	mux.HandleFunc("POST /api/v1/orders/{orderId}/cancel", cancel_order.NewCancelOrderEndpoint(
		root.OrdersDataAccess,
		root.ProblemEncoder,
		root.LoggerProvider,
		root.TimeProvider))
	mux.HandleFunc("GET /api/v1/orders/{orderId}", get_order_details.NewGetOrderDetailsEndpoint(
		root.OrdersDataAccess,
		root.ProblemEncoder,
		root.LoggerProvider))
}
