package internal

import (
	"net/http"

	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/hello_world"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders/cancel_order"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders/get_order_details"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders/register_order"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/shared/auth"
)

func addRoutes(mux *http.ServeMux, root *CompositionRoot) {
	mux.Handle("GET /hello", hello_world.NewHelloEndpoint(root.CreateLogger))

	authenticateMw := auth.Authenticate(root.CreateLogger)
	requireAuthMw := auth.RequireAuthentication()

	mux.Handle("POST /api/v1/orders",
		shared.Chain(
			authenticateMw,
			requireAuthMw,
			auth.RequirePermission("orders.register", root.CreateLogger))(
			register_order.NewRegisterOrderEndpoint(
				root.OrdersDataAccess,
				root.ProblemEncoder,
				root.CreateLogger,
				root.TimeProvider)))
	mux.Handle("POST /api/v1/orders/{orderId}/cancel",
		shared.Chain(
			authenticateMw,
			requireAuthMw,
			auth.RequirePermission("orders.cancel", root.CreateLogger))(
			cancel_order.NewCancelOrderEndpoint(
				root.OrdersDataAccess,
				root.ProblemEncoder,
				root.CreateLogger,
				root.TimeProvider)))
	mux.Handle("GET /api/v1/orders/{orderId}",
		shared.Chain(
			authenticateMw,
			requireAuthMw,
			auth.RequirePermission("orders.read", root.CreateLogger))(
			get_order_details.NewGetOrderDetailsEndpoint(
				root.OrdersDataAccess,
				root.ProblemEncoder,
				root.CreateLogger)))
}
