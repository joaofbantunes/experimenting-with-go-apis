package internal

import (
	"net/http"

	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/hello_world"
	"github.com/joaofbantunes/experimenting-with-go-apis/std/internal/features/orders/register_order"
)

func addRoutes(mux *http.ServeMux, root *CompositionRoot) {
	mux.HandleFunc("GET /hello", hello_world.NewHelloEndpoint(root.LoggerProvider))
	mux.HandleFunc("POST /orders", register_order.NewRegisterOrderEndpoint(
		root.OrdersDataAccess,
		root.ProblemEncoder,
		root.LoggerProvider,
		root.TimeProvider))
}
